// package: hidot
// file: hi.proto

var hi_pb = require("./hi_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var HiDot = (function () {
  function HiDot() {}
  HiDot.serviceName = "hidot.HiDot";
  return HiDot;
}());

HiDot.Hi = {
  methodName: "Hi",
  service: HiDot,
  requestStream: false,
  responseStream: false,
  requestType: hi_pb.HiReq,
  responseType: hi_pb.HiRes
};

HiDot.Write = {
  methodName: "Write",
  service: HiDot,
  requestStream: false,
  responseStream: false,
  requestType: hi_pb.WriteReq,
  responseType: hi_pb.WriteRes
};

exports.HiDot = HiDot;

function HiDotClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

HiDotClient.prototype.hi = function hi(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(HiDot.Hi, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

HiDotClient.prototype.write = function write(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(HiDot.Write, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.HiDotClient = HiDotClient;

