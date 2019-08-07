/**
 * @fileoverview gRPC-Web generated client stub for go_out
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.go_out = require('./hi_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.go_out.HiDotClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.go_out.HiDotPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqData,
 *   !proto.go_out.ResData>}
 */
const methodDescriptor_HiDot_Hi = new grpc.web.MethodDescriptor(
  '/go_out.HiDot/Hi',
  grpc.web.MethodType.UNARY,
  proto.go_out.ReqData,
  proto.go_out.ResData,
  /** @param {!proto.go_out.ReqData} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResData.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.go_out.ReqData,
 *   !proto.go_out.ResData>}
 */
const methodInfo_HiDot_Hi = new grpc.web.AbstractClientBase.MethodInfo(
  proto.go_out.ResData,
  /** @param {!proto.go_out.ReqData} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResData.deserializeBinary
);


/**
 * @param {!proto.go_out.ReqData} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.go_out.ResData)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.go_out.ResData>|undefined}
 *     The XHR Node Readable Stream
 */
proto.go_out.HiDotClient.prototype.hi =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.HiDot/Hi',
      request,
      metadata || {},
      methodDescriptor_HiDot_Hi,
      callback);
};


/**
 * @param {!proto.go_out.ReqData} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.go_out.ResData>}
 *     A native promise that resolves to the response
 */
proto.go_out.HiDotPromiseClient.prototype.hi =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.HiDot/Hi',
      request,
      metadata || {},
      methodDescriptor_HiDot_Hi);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqDirs,
 *   !proto.go_out.ResDots>}
 */
const methodDescriptor_HiDot_FindDot = new grpc.web.MethodDescriptor(
  '/go_out.HiDot/FindDot',
  grpc.web.MethodType.UNARY,
  proto.go_out.ReqDirs,
  proto.go_out.ResDots,
  /** @param {!proto.go_out.ReqDirs} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResDots.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.go_out.ReqDirs,
 *   !proto.go_out.ResDots>}
 */
const methodInfo_HiDot_FindDot = new grpc.web.AbstractClientBase.MethodInfo(
  proto.go_out.ResDots,
  /** @param {!proto.go_out.ReqDirs} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResDots.deserializeBinary
);


/**
 * @param {!proto.go_out.ReqDirs} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.go_out.ResDots)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.go_out.ResDots>|undefined}
 *     The XHR Node Readable Stream
 */
proto.go_out.HiDotClient.prototype.findDot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.HiDot/FindDot',
      request,
      metadata || {},
      methodDescriptor_HiDot_FindDot,
      callback);
};


/**
 * @param {!proto.go_out.ReqDirs} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.go_out.ResDots>}
 *     A native promise that resolves to the response
 */
proto.go_out.HiDotPromiseClient.prototype.findDot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.HiDot/FindDot',
      request,
      metadata || {},
      methodDescriptor_HiDot_FindDot);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqLoad,
 *   !proto.go_out.ResConfig>}
 */
const methodDescriptor_HiDot_LoadByConfig = new grpc.web.MethodDescriptor(
  '/go_out.HiDot/LoadByConfig',
  grpc.web.MethodType.UNARY,
  proto.go_out.ReqLoad,
  proto.go_out.ResConfig,
  /** @param {!proto.go_out.ReqLoad} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResConfig.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.go_out.ReqLoad,
 *   !proto.go_out.ResConfig>}
 */
const methodInfo_HiDot_LoadByConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.go_out.ResConfig,
  /** @param {!proto.go_out.ReqLoad} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResConfig.deserializeBinary
);


/**
 * @param {!proto.go_out.ReqLoad} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.go_out.ResConfig)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.go_out.ResConfig>|undefined}
 *     The XHR Node Readable Stream
 */
proto.go_out.HiDotClient.prototype.loadByConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.HiDot/LoadByConfig',
      request,
      metadata || {},
      methodDescriptor_HiDot_LoadByConfig,
      callback);
};


/**
 * @param {!proto.go_out.ReqLoad} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.go_out.ResConfig>}
 *     A native promise that resolves to the response
 */
proto.go_out.HiDotPromiseClient.prototype.loadByConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.HiDot/LoadByConfig',
      request,
      metadata || {},
      methodDescriptor_HiDot_LoadByConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqImport,
 *   !proto.go_out.ResImport>}
 */
const methodDescriptor_HiDot_ImportByConfig = new grpc.web.MethodDescriptor(
  '/go_out.HiDot/ImportByConfig',
  grpc.web.MethodType.UNARY,
  proto.go_out.ReqImport,
  proto.go_out.ResImport,
  /** @param {!proto.go_out.ReqImport} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResImport.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.go_out.ReqImport,
 *   !proto.go_out.ResImport>}
 */
const methodInfo_HiDot_ImportByConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.go_out.ResImport,
  /** @param {!proto.go_out.ReqImport} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResImport.deserializeBinary
);


/**
 * @param {!proto.go_out.ReqImport} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.go_out.ResImport)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.go_out.ResImport>|undefined}
 *     The XHR Node Readable Stream
 */
proto.go_out.HiDotClient.prototype.importByConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.HiDot/ImportByConfig',
      request,
      metadata || {},
      methodDescriptor_HiDot_ImportByConfig,
      callback);
};


/**
 * @param {!proto.go_out.ReqImport} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.go_out.ResImport>}
 *     A native promise that resolves to the response
 */
proto.go_out.HiDotPromiseClient.prototype.importByConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.HiDot/ImportByConfig',
      request,
      metadata || {},
      methodDescriptor_HiDot_ImportByConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqExport,
 *   !proto.go_out.ResExport>}
 */
const methodDescriptor_HiDot_ExportConfig = new grpc.web.MethodDescriptor(
  '/go_out.HiDot/ExportConfig',
  grpc.web.MethodType.UNARY,
  proto.go_out.ReqExport,
  proto.go_out.ResExport,
  /** @param {!proto.go_out.ReqExport} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResExport.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.go_out.ReqExport,
 *   !proto.go_out.ResExport>}
 */
const methodInfo_HiDot_ExportConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.go_out.ResExport,
  /** @param {!proto.go_out.ReqExport} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResExport.deserializeBinary
);


/**
 * @param {!proto.go_out.ReqExport} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.go_out.ResExport)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.go_out.ResExport>|undefined}
 *     The XHR Node Readable Stream
 */
proto.go_out.HiDotClient.prototype.exportConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.HiDot/ExportConfig',
      request,
      metadata || {},
      methodDescriptor_HiDot_ExportConfig,
      callback);
};


/**
 * @param {!proto.go_out.ReqExport} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.go_out.ResExport>}
 *     A native promise that resolves to the response
 */
proto.go_out.HiDotPromiseClient.prototype.exportConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.HiDot/ExportConfig',
      request,
      metadata || {},
      methodDescriptor_HiDot_ExportConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqExport,
 *   !proto.go_out.ResExport>}
 */
const methodDescriptor_HiDot_ExportDot = new grpc.web.MethodDescriptor(
  '/go_out.HiDot/ExportDot',
  grpc.web.MethodType.UNARY,
  proto.go_out.ReqExport,
  proto.go_out.ResExport,
  /** @param {!proto.go_out.ReqExport} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResExport.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.go_out.ReqExport,
 *   !proto.go_out.ResExport>}
 */
const methodInfo_HiDot_ExportDot = new grpc.web.AbstractClientBase.MethodInfo(
  proto.go_out.ResExport,
  /** @param {!proto.go_out.ReqExport} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.go_out.ResExport.deserializeBinary
);


/**
 * @param {!proto.go_out.ReqExport} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.go_out.ResExport)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.go_out.ResExport>|undefined}
 *     The XHR Node Readable Stream
 */
proto.go_out.HiDotClient.prototype.exportDot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.HiDot/ExportDot',
      request,
      metadata || {},
      methodDescriptor_HiDot_ExportDot,
      callback);
};


/**
 * @param {!proto.go_out.ReqExport} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.go_out.ResExport>}
 *     A native promise that resolves to the response
 */
proto.go_out.HiDotPromiseClient.prototype.exportDot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.HiDot/ExportDot',
      request,
      metadata || {},
      methodDescriptor_HiDot_ExportDot);
};


module.exports = proto.go_out;

