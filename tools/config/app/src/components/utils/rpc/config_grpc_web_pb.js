/**
 * @fileoverview gRPC-Web generated client stub for go_out
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.go_out = require('./config_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.go_out.DotConfigClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
proto.go_out.DotConfigPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 *   !proto.go_out.ReqDirs,
 *   !proto.go_out.ResDots>}
 */
const methodDescriptor_DotConfig_FindDot = new grpc.web.MethodDescriptor(
  '/go_out.DotConfig/FindDot',
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
const methodInfo_DotConfig_FindDot = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.go_out.DotConfigClient.prototype.findDot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.DotConfig/FindDot',
      request,
      metadata || {},
      methodDescriptor_DotConfig_FindDot,
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
proto.go_out.DotConfigPromiseClient.prototype.findDot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.DotConfig/FindDot',
      request,
      metadata || {},
      methodDescriptor_DotConfig_FindDot);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqImport,
 *   !proto.go_out.ResImport>}
 */
const methodDescriptor_DotConfig_ImportByConfig = new grpc.web.MethodDescriptor(
  '/go_out.DotConfig/ImportByConfig',
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
const methodInfo_DotConfig_ImportByConfig = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.go_out.DotConfigClient.prototype.importByConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.DotConfig/ImportByConfig',
      request,
      metadata || {},
      methodDescriptor_DotConfig_ImportByConfig,
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
proto.go_out.DotConfigPromiseClient.prototype.importByConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.DotConfig/ImportByConfig',
      request,
      metadata || {},
      methodDescriptor_DotConfig_ImportByConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqImport,
 *   !proto.go_out.ResImport>}
 */
const methodDescriptor_DotConfig_ImportByDot = new grpc.web.MethodDescriptor(
  '/go_out.DotConfig/ImportByDot',
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
const methodInfo_DotConfig_ImportByDot = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.go_out.DotConfigClient.prototype.importByDot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.DotConfig/ImportByDot',
      request,
      metadata || {},
      methodDescriptor_DotConfig_ImportByDot,
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
proto.go_out.DotConfigPromiseClient.prototype.importByDot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.DotConfig/ImportByDot',
      request,
      metadata || {},
      methodDescriptor_DotConfig_ImportByDot);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqExport,
 *   !proto.go_out.ResExport>}
 */
const methodDescriptor_DotConfig_ExportConfig = new grpc.web.MethodDescriptor(
  '/go_out.DotConfig/ExportConfig',
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
const methodInfo_DotConfig_ExportConfig = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.go_out.DotConfigClient.prototype.exportConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.DotConfig/ExportConfig',
      request,
      metadata || {},
      methodDescriptor_DotConfig_ExportConfig,
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
proto.go_out.DotConfigPromiseClient.prototype.exportConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.DotConfig/ExportConfig',
      request,
      metadata || {},
      methodDescriptor_DotConfig_ExportConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.go_out.ReqExport,
 *   !proto.go_out.ResExport>}
 */
const methodDescriptor_DotConfig_ExportDot = new grpc.web.MethodDescriptor(
  '/go_out.DotConfig/ExportDot',
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
const methodInfo_DotConfig_ExportDot = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.go_out.DotConfigClient.prototype.exportDot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/go_out.DotConfig/ExportDot',
      request,
      metadata || {},
      methodDescriptor_DotConfig_ExportDot,
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
proto.go_out.DotConfigPromiseClient.prototype.exportDot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/go_out.DotConfig/ExportDot',
      request,
      metadata || {},
      methodDescriptor_DotConfig_ExportDot);
};


module.exports = proto.go_out;

