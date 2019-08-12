const { ReqDirs,ReqExport, ReqImport,ReqLoad} = require('./hi_pb.js');
const {HiDotClient} = require('./hi_grpc_web_pb.js');
var rpcweb = new HiDotClient('http://localhost:6868/root');

function rpcFindDot(dirList) {
    var request = new ReqDirs();
    //var dirList = ['/home/jayce/golangPath/src/dot/dots'];
    request.setDirsList(dirList);
    console.log("finddot:",request);
    rpcweb.findDot(request, {}, (err, response) => {
        if (response) {
            console.log("dotsinfo:",response.getDotsinfo());
            console.log("不存在组件的目录：",response.getNoexistdirsList());
            console.log("error:",response.getError());
            var res={
                dotsInfo:null,
                noExistDirs:null,
                error:null
            };
            res.dotsInfo=response.getDotsinfo();
            res.noExistDirs=response.getNoexistdirsList();
            res.error=response.getError();
            return res;
        } else {
            console.log(err);
        }
    })
}
window.rpcFindDot = rpcFindDot;

//loadbyconfig
function rpcLoadByConfig(dir,typeId) {
    var request = new ReqLoad();
    //var dir = '/home/jayce/golangPath/src/dot/sample/grpc/http/server/server_http.json'; //文件路径
    //var typeId = 'afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0';  //typeId
    request.setDatafilepath(dir);
    request.setTypeid(typeId);
    console.log("loadbyconfig",request);
    rpcweb.loadByConfig(request, {}, (err, response) => {
        if (response) {
            console.log("configjson:",response.getConfigjson());
            console.log("errinfo:",response.getErrinfo());
            var res={
                configJson:null,
                errorInfo:null
            };
            res.configJson=response.getConfigjson();
            res.errorInfo=response.getError();
            return res;
        } else {
            console.log(err);
        }
    })
}
window.rpcLoadByConfig = rpcLoadByConfig;

//importByDot
function rpcimportByDot(filepath){
    var request = new ReqImport();
    //var filepath = '/home/jayce/golangPath/src/github.com/scryinfo/dot/tools/config/data/server/testdot.json'; //文件路径
    request.setFilepath(filepath);
    console.log("importbydot",request);
    rpcweb.importByDot(request,{},(err, response)=>{
        if (response){
            console.log("json:",response.getJson());
            console.log("error:",response.getError());
            var res={
                Json:null,
                error:null
            };
            res.Json=response.getJson();
            res.error=response.getError();
            return res;
        }else {
            console.log(err)
        }
    })
}
window.rpcimportByDot=rpcimportByDot;

//importByConfig
function rpcimportByConfig(filepath){
    var request = new ReqImport();
    //var filepath = '/home/jayce/golangPath/src/dot/sample/grpc/http/server/server_http.json'; //文件路径
    request.setFilepath(filepath);
    console.log("importbyconfig",request);
    rpcweb.importByConfig(request,{},(err, response)=>{
        if (response){
            console.log("json:",response.getJson());
            console.log("error:",response.getError());
            var res={
                Json:null,
                error:null
            };
            res.Json=response.getJson();
            res.error=response.getError();
            return res;
        }else {
            console.log(err)
        }
    })
}
window.rpcimportByConfig=rpcimportByConfig;

//test exportDot
function rpcExportDot(data,filename) {
    var request = new ReqExport();
    /*var data =
        [{
            "Meta": {
                "typeId": "4b8b1751-4799-4578-af46-d9b339cf582f",
                "version": "",
                "name": "",
                "showName": "",
                "single": false,
                "relyTypeIds": null
            },
            "Lives": [
                {
                    "TypeId": "",
                    "LiveId": "",
                    "RelyLives": null,
                    "Dot": null,
                    "json": null,
                    "name": ""
                }
            ]
        }];*/
    const objToStr = JSON.stringify(data);
    //var filename = ["testdot.json"];
    request.setFilenameList(filename);
    request.setDotdata(objToStr);
    console.log("exportDot:",request);
    rpcweb.exportDot(request, {}, (err, response) => {
        if (response) {
            console.log("error:",response.getError());
            return response.getError();
        } else {
            console.log(err);
        }
    })
}
window.rpcExportDot = rpcExportDot;

//test exportConfig
function rpcExportConfig(data,filename) {
    /*var data = {
        "log": {
            "file": "log.log",
            "level": "debug"
        },
        "dots": [
            {
                "metaData": {
                    "typeId": "4943e959-7ad7-42c6-84dd-8b24e9ed30bb",
                    "version": "",
                    "name": "",
                    "showName": "",
                    "single": false,
                    "relyTypeIds": null
                },
                "lives": [
                    {
                        "TypeId": "",
                        "liveId": "4943e959-7ad7-42c6-84dd-8b24e9ed30bb",
                        "RelyLives": null,
                        "Dot": null,
                        "name": "",
                        "json": {
                            "addr": ":8080",
                            "keyFile": "",
                            "pemFile": "",
                            "logSkipPaths": ["/sample/!*"]
                        }
                    }
                ]
            },
            {
                "metaData": {
                    "typeId": "6be39d0b-3f5b-47b4-818c-642c049f3166",
                    "version": "",
                    "name": "",
                    "showName": "",
                    "single": false,
                    "relyTypeIds": null
                },
                "lives": [
                    {

                        "TypeId": "",
                        "liveId": "6be39d0b-3f5b-47b4-818c-642c049f3166",
                        "relyLives": {"GinDot_": "4943e959-7ad7-42c6-84dd-8b24e9ed30bb"},
                        "Dot": null,
                        "name": "",
                        "json": {
                            "relativePath": "/"
                        }
                    }
                ]
            }
        ]
    };*/
    var request = new ReqExport();
    const objToStr = JSON.stringify(data);
    //var filename = ["testconfig.json", "testconfig.toml", "testconfig.yaml"];
    request.setFilenameList(filename);
    request.setConfigdata(objToStr);
    console.log("exportConfig:",request);
    rpcweb.exportConfig(request, {}, (err, response) => {
        if (response) {
            console.log("error:",response.getError());
            return response.getError();
        } else {
            console.log(err);
        }s
    })
}
window.rpcExportConfig = rpcExportConfig;