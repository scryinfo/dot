
const {ReqData, ResData} = require('../hi_pb.js');
const {HiDotClient} = require('../hi_grpc_web_pb.js');

var rpcweb = new HiDotClient('http://localhost:6868/root');

var request = new ReqData();
request.setName('http hi client');

rpcweb.hi(request, {}, (err, response) => {

    console.log("http");
    if (response) {
        console.log(response.getName());
    }else{
        console.log(err);
    }
});

var rpcwebs = new HiDotClient('https://localhost:6868/root');
rpcwebs.hi(request, {}, (err, response) => {
    console.log("http");
    if (response) {
        console.log(response.getName());
    }else{
        console.log(err);
    }
});