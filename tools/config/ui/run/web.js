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
function func2() {
    var request = new ReqDirs();
    var dirList = ['C:\\go-work\\src\\github.com\\scryinfo\\dot'];
    request.setDirsList(dirList);
    rpcweb.findDot(request, {}, (err, response) => {
        console.log("finddot:::");
        if (response) {
            console.log(response.getDotsinfo());
            console.log(response.getNoexistdirsList());
            console.log(response.getError());
        } else {
            console.log(err);
        }
    })
}
window.func2 = func2;


//loadbyconfig
function func3() {
    var request = new ReqLoad();
    var datacopy = {
        "log": {
            "file": "log.log",
            "level": "debug"
        },
        "dots": [
            {
                "metaData": {
                    "name": "ServerNobl",
                    "typeId": "77a766e7-c288-413f-946b-bc9de6df3d70"
                },
                "lives": [
                    {
                        "liveId": "77a766e7-c288-413f-946b-bc9de6df3d70",
                        "json": {
                            "addrs": ["localhost:5012"]
                        }
                    }
                ]
            },
            {
                "metaData": {
                    "name": "grpc http server",
                    "typeId": "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0"
                },
                "lives": [
                    {
                        "liveId": "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0",
                        "json": {
                            "addr": ":6868",
                            "preUrl": "root"
                        }
                    }
                ]
            },
            {
                "metaData": {
                    "name": "server 1",
                    "typeId": "hiserver"
                },
                "lives": [
                    {
                        "liveId": "hiserver",
                        "json": {
                            "name": "server 1"
                        }
                    }
                ]
            }
        ]
    };  //粘贴数据
    var data2 = 'C:\\go-work\\src\\dot\\sample\\grpc\\http\\server\\server_http.json'; //文件路径
    var data3 = 'afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0';  //typeId
    var datatemplate = {
        "typeId": "",
        "liveId": "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0",
        "relyLives": {
            "ServerNobl": "77a766e7-c288-413f-946b-bc9de6df3d70"
        },
        "Dot": null,
        "json": {
            "preUrl": "",
            "addr": "",
            "tls": {
                "caPem": "",
                "pem": "",
                "key": "",
                "serverNameOverride": ""
            }
        },
        "name": ""
    };
    const objToStr = JSON.stringify(datatemplate);
    request.setDotinfo(objToStr);
    request.setTypeid(data3);
    const objToStr2 = JSON.stringify(datacopy);
    request.setDatacopypaste(objToStr2);
    console.log(request);
    rpcweb.loadByConfig(request, {}, (err, response) => {
        console.log("loadByConfig:::");
        if (response) {
            console.log(response.getConfigjson())
            console.log("err:",response.getErrinfo())
        } else {
            console.log(err);
        }
    })
}
window.func3 = func3;

//importByDot

//importByConfig

//test exportDot
function func6() {
    var request = new ReqExport();
    var data =
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
        }];
    const objToStr = JSON.stringify(data);
    var filename = ["testdot.json"];
    request.setFilenameList(filename);
    request.setDotdata(objToStr);
    console.log(objToStr);
    console.log(request);
    rpcweb.exportDot(request, {}, (err, response) => {
        console.log("exportDot:::");
        if (response) {
            console.log(response.getError());
        } else {
            console.log(err);
        }
    })
}
window.func6 = func6;

//test exportConfig
function func7() {
    var data = {
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
                            "logSkipPaths": ["/sample/*"]
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
    };
    var request = new ReqExport();
    const objToStr = JSON.stringify(data);
    var filename = ["testconfig.json", "testconfig.toml", "testconfig.yaml"];
    request.setFilenameList(filename);
    request.setConfigdata(objToStr);
    console.log(objToStr);
    console.log(request);
    rpcweb.exportConfig(request, {}, (err, response) => {
        console.log("exportConfig:::");
        if (response) {
            console.log(response.getError());
        } else {
            console.log(err);
        }
    })
}
window.func5=func5;

function testImportConfig() {
    let request = new ReqImport();
    let filepath = "";
    request.setFilepath(filepath);
    rpcweb.importByConfig(request,{},(err,response)=>{
        console.log('importConfig:::');
        if (response) {
            if(response.getError()==''){
                console.log(response.getJson());
            }else {
                console.log(response.getError());
            }
        }else {
            console.log(err);
        }
    })
}
window.testImportConfig = testImportConfig;























