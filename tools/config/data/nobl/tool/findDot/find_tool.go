package findDot

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

//保存返回的组件通用信息
type packageInfo struct {
	absDir      string //该目录\包的绝对路径
	packageName string //正在扫描的包名，需要绝对路径
	astFiles    []*ast.File
	ImportDir   string //导入路径
	Alias       string //别名

	Funcs   []*FuncOfDot //返回组件普通信息的函数
	isExist bool         //是否存在返回dot.TypeLives的函数

	ConfigFuncNames []string //返回组件特有配置信息的函数
	IsExistConfig   bool     //是否存在返回config的函数
}
type FuncOfDot struct {
	FuncName string    //函数名
	IsSlice  bool      //返回值是否是数组
	location token.Pos //函数位置
}

//dirs用户传入的数据
//返回组件数据、不存在的目录信息包括mod依赖目录、错误信息
func FindDots(dirs []string) (data []byte, notExistDir []string, err error) {

	var paths []string //保存有效的目录数据

	//处理目录问题
	{
		//处理掉用户输入的不存在目录
		yesDir, noDir := splitDir(dirs)
		if len(yesDir) == 0 {
			err = errors.New("Please enter a valid directory")
			return
		}
		paths = append(paths, yesDir...)
		notExistDir = append(notExistDir, noDir...)
	}

	//得到绝对路径
	{
		for ix := range paths {
			if !filepath.IsAbs(paths[ix]) {
				paths[ix], err = filepath.Abs(paths[ix])
				if err != nil {
					fmt.Println("获取绝对路径出错:", err)
					return
				}
			}
		}
	}

	//处理重复的目录，原因是多个mod中可能会重复依赖多个项目
	// 或者a\b\c和a\b这类目录也要处理
	{
		paths = RemoveRepByMap(paths)
	}

	// 对于每个子目录通过package.config获取包名，以及当前位置的所有ast.file
	var allInfo []*packageInfo
	{
		//获取一个项目下所有的子目录
		for i := range paths {
			//遍历获取子目录
			dirs, err = getAllSonDirs(paths[i])
			if err != nil {
				fmt.Println("获取一个项目下所有的子目录出错：", err)
				return
			}
			//加载路径获取配置信息
			cfg := &packages.Config{
				Mode: packages.LoadSyntax, //不包含依赖,尝试下面这个
				Dir:  paths[i],            //设置当前目录
			}
			pkginfos, errs := packages.Load(cfg, dirs...)
			err = errs
			if err != nil {
				fmt.Println("packages.Load err:", err)
				return
			}
			for ix := range dirs {
				pkg := pkginfos[ix]
				astFiles := make([]*ast.File, len(pkg.Syntax))
				for i, file := range pkg.Syntax {
					astFiles[i] = file
				}
				p := packageInfo{
					absDir:      dirs[ix],
					packageName: pkg.Name,
					astFiles:    astFiles,
					ImportDir:   pkg.PkgPath,
				}
				allInfo = append(allInfo, &p)
			}
		}
	}

	//查找[]*packageInfo的每一个对象，根据ast,node判断
	{
		for _, p := range allInfo {
			p.findFuncNodeOnAst(false) //查找*dot.TypeLives和[]*dot.TypeLives
			p.findFuncNodeOnAst(true)  //查找*dot.ConfigTypeLives
		}
	}

	//将满足普通条件的包筛选出来
	var exitFuncInfos []*packageInfo
	{
		//为isExist字段赋值
		{
			for _, p := range allInfo {
				if len(p.Funcs) == 0 {
					//这个目录下没有满足条件的函数
					p.isExist = false
				} else {
					p.isExist = true
				}
			}
		}
		//赋值
		{
			for _, p := range allInfo {
				if p.isExist {
					exitFuncInfos = append(exitFuncInfos, p)
				}
			}
		}

		//检测有没有结果
		{
			if len(exitFuncInfos) == 0 {
				fmt.Println("没有找到符合条件的函数")
				err = errors.New("没有找到符合条件的组件")
				return
			}
		}
	}

	//解决同名包问题-别名
	{
		//存储每个包出现的次数
		map1 := make(map[string]int)
		//赋值
		{
			for _, p := range exitFuncInfos {
				//该包是否已经放入
				if _, ok := map1[p.packageName]; ok {
					//已有
					map1[p.packageName]++
				} else {
					//没有
					map1[p.packageName] = 1
				}
			}
		}
		//利用map1构建别名，只出现一次别名默认为包名
		{
			for _, p := range exitFuncInfos {
				if v, ok := map1[p.packageName]; ok {
					if v == 1 {
						p.Alias = p.packageName
					} else {
						p.Alias = p.packageName + "_" + strconv.Itoa(v)
						map1[p.packageName]--
					}
				}
			}
		}
	}

	//为isExistConfig字段赋值
	{
		for _, p := range exitFuncInfos {
			if len(p.ConfigFuncNames) == 0 {
				//这个目录下没有满足条件的函数
				p.IsExistConfig = false
			} else {
				p.IsExistConfig = true
			}

		}
	}

	//生成代码文件
	{
		buildCodeFromTemplate(exitFuncInfos)

	}

	//调用执行callMethod生成序列化的json文件
	{
		//运行生成的代码文件
		{
			cmd := exec.Command(getGOROOTBin(), "run", "./run_out/callMethod.go")
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error %v executing command!", err)
				return
			}
		}
		//读取组件信息
		{
			data, err = ioutil.ReadFile("./run_out/result.json")
			if err != nil {
				fmt.Println("File reading error", err)
			}
			return
		}
	}
}

//将目录集合分为两部分，分别是存在的目录和不存在的目录
func splitDir(dirs []string) (existDir, notExistDir []string) {
	for _, dir := range dirs {
		if isDirectory(dir) {
			//该目录已找到
			existDir = append(existDir, dir)
		} else {
			//该目录没找到
			notExistDir = append(notExistDir, dir)
		}
	}
	return
}

//是目录还是文件
func isDirectory(name string) bool {
	if isDirExist(name) {
		info, err := os.Stat(name) //return fileinfo
		if err != nil {
			return false
		}
		return info.IsDir() //true or false
	}
	return false
}

//存不存在
func isDirExist(paths string) bool {
	_, err := os.Stat(paths)
	return err == nil || os.IsExist(err)
}

//解决重复目录
func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for i := range slc {
		l := len(tempMap)
		tempMap[slc[i]] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, slc[i])
		}
	}
	return result
}

//获取指定目录下的所有子目录
func getAllSonDirs(dirpath string) ([]string, error) {
	var dir_list []string
	dir_err := filepath.Walk(dirpath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() && isTrueDir(path) {
				dir_list = append(dir_list, path)
				return nil
			}
			return nil
		})
	return dir_list, dir_err
}
func isTrueDir(path string) bool {
	if strings.Index(path, "node_modules") == -1 {
		if strings.Index(path, ".git") == -1 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//查找满足条件的函数节点
func (p *packageInfo) findFuncNodeOnAst(isConfig bool) {
	///var FuncNames []string
	for _, astFile := range p.astFiles {
		if astFile == nil {
			break
		}
		ast.Inspect(astFile, func(node ast.Node) bool {
			for {
				//must be a func
				funcNode, ok := node.(*ast.FuncDecl)
				if !ok {
					break //不是函数
				}
				if funcNode.Recv != nil {
					break //排除有接收者的函数
				}
				if !(funcNode.Type.Params.List == nil) {
					break //排除需要参数的函数
				}
				result := funcNode.Type.Results
				if result == nil {
					break //排除没有返回值的函数
				}
				resultList := result.List
				if len(resultList) != 1 {
					break //排除返回值不值一个的函数
				}
				if isConfig {
					b1, _ := returnValueJudgmentOfConfig(resultList[0])
					if !b1 {
						break //排除返回值类型不匹配的函数
					}
					//保存函数信息
					p.ConfigFuncNames = append(p.ConfigFuncNames, funcNode.Name.Name)
				} else {
					b1, b2 := returnValueJudgment(resultList[0])
					if !b1 {
						break //排除返回值类型不匹配的函数
					}
					//保存函数信息
					funcInfo := &FuncOfDot{
						FuncName: funcNode.Name.Name, //函数名
						IsSlice:  b2,
						location: funcNode.Type.Func,
					}
					p.Funcs = append(p.Funcs, funcInfo)
				}
				return true
			}
			return true
		})
	}
}

//查找通用组件信息
//第一个代表函数的返回值是否符合条件
//第二个代表返回值是否是数组
func returnValueJudgment(ret *ast.Field) (bool, bool) {
	retype, ok := (ret.Type).(*ast.StarExpr) //找到*
	if ok {                                  //是一个指针
		x, ok1 := (retype.X).(*ast.SelectorExpr) //有选择器的表达式  a.b
		if !ok1 {
			return false, false //类似*xxx
		}
		xx := x.X.(*ast.Ident)
		xsel := x.Sel.Name
		if xx.Name == "dot" {
			if xsel == "TypeLives" {
				return true, false //返回值是*dot.Typelives
			}
		}
		return false, false //指针指向的结构错误
	}

	retype2, ok := (ret.Type).(*ast.ArrayType)
	if ok { //是一个切片
		elt, ok := (retype2.Elt).(*ast.StarExpr)
		if ok { //切片存放的指针数据
			x, ok1 := (elt.X).(*ast.SelectorExpr) //有选择器的表达式  a.b
			if !ok1 {
				return false, false //类似*xxx
			}
			xx := x.X.(*ast.Ident)
			xsel := x.Sel.Name
			if xx.Name == "dot" {
				if xsel == "TypeLives" {
					return true, true //返回值是[]*dot.Typelives
				}
			}
			return false, false //指针指向的结构错误
		}
		return false, false //切片存放的数据不是指针
	}
	return false, false //返回值类型错误
}

//查找特有组件信息
//第一个代表函数的返回值是否符合条件
//第二个代表返回值是否是数组
func returnValueJudgmentOfConfig(ret *ast.Field) (bool, bool) {
	retype, ok := (ret.Type).(*ast.StarExpr) //找到*
	if ok {                                  //是一个指针
		x, ok1 := (retype.X).(*ast.SelectorExpr) //有选择器的表达式  a.b
		if !ok1 {
			return false, false //类似*xxx
		}
		xx := x.X.(*ast.Ident)
		xsel := x.Sel.Name
		if xx.Name == "dot" {
			if xsel == "ConfigTypeLives" {
				return true, false //返回值是*dot.ConfigTypelives
			}
		}
		return false, false //指针指向的结构错误
	}
	return false, false //返回值类型错误
}

//
func getGOPATHsrc() string {
	gopath := os.Getenv("GOPATH")
	switch runtime.GOOS {
	case "windows":
		gopath = gopath + "\\src\\"
	case "linux":
		gopath = gopath + "/src/"
	default:
		log.Fatal("无法识别的操作系统")
	}
	return gopath
}
func getGOROOTBin() string {
	gopath := os.Getenv("GOROOT")
	switch runtime.GOOS {
	case "windows":
		gopath = gopath + "\\bin\\go.exe"
	case "linux":
		gopath = gopath + "/bin/go"
	default:
		log.Fatal("无法识别的操作系统")
	}
	return gopath
}

//模板生成
func buildCodeFromTemplate(e []*packageInfo) {
	buf := bytes.Buffer{}
	//使用模板
	var filepaths = "./nobl/tool/findDot/file1.tmpl"
	filepaths = filepath.FromSlash(filepaths)
	t, err := template.ParseFiles(filepaths)
	if err != nil {
		fmt.Println("parseFileErr err", err)
		return
	}
	err = t.Execute(&buf, e)
	if err != nil {
		fmt.Println("executing template err", err)
		return
	}
	//file
	baseName := "callMethod.go"
	baseName = filepath.Join("./run_out", baseName)
	err = ioutil.WriteFile(baseName, buf.Bytes(), 0777)
	if err != nil {
		fmt.Println("writing callMethod.go err")
		return
	}
}
