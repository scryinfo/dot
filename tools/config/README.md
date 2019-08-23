
## tool/config
---
### 组件配置工具
---

* finddot  在目录找出有“func () *dot.TypeLives”或“func () []dot.TypeLives”定义函数的go文件，获得其返回值（组件信息），并加载到dotlist  
* dotlist　你可以在这里修改组件信息或者将其添加到配置　　

* config　你可以在这里完成组件的配置信息，支持逐个填写或json文本输入

* import　你可以在这里导入组件信息(json)和配置信息(json,yaml,toml)

* export　你可以在这里进行组件和配置的导出工作

---

### 使用方式

```
cd dot/tools/config/ui/apps

npm install

npm run build

cd ../../

./config.sh  or   config.bat

```
注：若浏览器未成功弹出，请手动打开浏览器并访问http://localhost:9090  
提示：本工具将会占用5012,6868以及9090端口，关闭终端时释放所有资源