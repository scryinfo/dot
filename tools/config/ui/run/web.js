const {ReqData, ReqDirs, ReqExport, ReqLoad} = require('./hi_pb.js');
const {HiDotClient} = require('./hi_grpc_web_pb.js');
var rpcweb = new HiDotClient('http://localhost:6868/root');

//hi
function func1() {
    var request = new ReqData();
    request.setName('http hi client');

    rpcweb.hi(request, {}, (err, response) => {
        console.log("HI:::");
        if (response) {
            console.log(response.getTest());
        } else {
            console.log(err);
        }
    });
}
window.func1 = func1;

//Finddot
function rpcFindDot(dirList) {
    var request = new ReqDirs();
    //var dirList = ['C:\\go-work\\src\\github.com\\scryinfo\\dot'];
    request.setDirsList(dirList);
    rpcweb.findDot(request, {}, (err, response) => {
        console.log("finddot:::");
        if (response) {
            console.log(response.getDotsinfo());
            console.log(response.getNoexistdirsList());
            console.log(response.getError());
            return response
        } else {
            console.log(err);
        }
    })
}
window.rpcFindDot = rpcFindDot;


//loadbyconfig
function rpcLoadByConfig(dir,typeId) {
    var request = new ReqLoad();
    //var data2 = 'C:\\go-work\\src\\dot\\sample\\grpc\\http\\server\\server_http.json'; //文件路径
    //var data3 = 'afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0';  //typeId
    request.setDatafilepath(dir);
    request.setTypeid(typeId);
    console.log(request);
    rpcweb.loadByConfig(request, {}, (err, response) => {
        console.log("loadByConfig:::");
        if (response) {
            console.log(response.getConfigjson())
            console.log("err:",response.getErrinfo())
            return response
        } else {
            console.log(err);
        }
    })
}
window.rpcLoadByConfig = rpcLoadByConfig;

//importByDot

//importByConfig

//test exportDot
function rpcExportDot(data,filename) {
    var request = new ReqExport();
   /* var data =
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
   // var filename = ["testdot.json"];
    request.setFilenameList(filename);
    request.setDotdata(objToStr);
    rpcweb.exportDot(request, {}, (err, response) => {
        console.log("exportDot:::");
        if (response) {
            console.log(response.getError());
            return response
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
    rpcweb.exportConfig(request, {}, (err, response) => {
        console.log("exportConfig:::");
        if (response) {
            console.log(response.getError());
        } else {
            console.log(err);
        }
    })
}
window.rpcExportConfig = rpcExportConfig;