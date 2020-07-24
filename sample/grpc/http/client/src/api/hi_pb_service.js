// package: hidot
// file: hi.proto

var hi_pb = require("./hi_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

// var exports = (typeof module === 'object' && module != null && module.exports !== undefined) ? module.exports : {};
// var exports = {};

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

HiDot.ServerStream = {
  methodName: "ServerStream",
  service: HiDot,
  requestStream: false,
  responseStream: true,
  requestType: hi_pb.HelloRequest,
  responseType: hi_pb.HelloResponse
};

HiDot.ClientStream = {
  methodName: "ClientStream",
  service: HiDot,
  requestStream: true,
  responseStream: false,
  requestType: hi_pb.HelloRequest,
  responseType: hi_pb.HelloResponse
};

HiDot.BothSides = {
  methodName: "BothSides",
  service: HiDot,
  requestStream: true,
  responseStream: true,
  requestType: hi_pb.HelloRequest,
  responseType: hi_pb.HelloResponse
};

// exports.HiDot = HiDot;

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

HiDotClient.prototype.serverStream = function serverStream(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(HiDot.ServerStream, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners.end.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

HiDotClient.prototype.clientStream = function clientStream(metadata) {
  var listeners = {
    end: [],
    status: []
  };
  var client = grpc.client(HiDot.ClientStream, {
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport
  });
  client.onEnd(function (status, statusMessage, trailers) {
    listeners.status.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners.end.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners = null;
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    write: function (requestMessage) {
      if (!client.started) {
        client.start(metadata);
      }
      client.send(requestMessage);
      return this;
    },
    end: function () {
      client.finishSend();
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

HiDotClient.prototype.bothSides = function bothSides(metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.client(HiDot.BothSides, {
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport
  });
  client.onEnd(function (status, statusMessage, trailers) {
    listeners.status.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners.end.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners = null;
  });
  client.onMessage(function (message) {
    listeners.data.forEach(function (handler) {
      handler(message);
    })
  });
  client.start(metadata);
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    write: function (requestMessage) {
      client.send(requestMessage);
      return this;
    },
    end: function () {
      client.finishSend();
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

// exports.HiDotClient = HiDotClient;

export {HiDot, HiDotClient}