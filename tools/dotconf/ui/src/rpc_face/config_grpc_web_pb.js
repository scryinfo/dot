/**
 * @fileoverview gRPC-Web generated client stub for rpc
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.rpc = require('./config_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.rpc.DotConfigFaceClient =
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

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.rpc.DotConfigFacePromiseClient =
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

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rpc.FindReq,
 *   !proto.rpc.FindRes>}
 */
const methodDescriptor_DotConfigFace_FindDot = new grpc.web.MethodDescriptor(
  '/rpc.DotConfigFace/FindDot',
  grpc.web.MethodType.UNARY,
  proto.rpc.FindReq,
  proto.rpc.FindRes,
  /**
   * @param {!proto.rpc.FindReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.FindRes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rpc.FindReq,
 *   !proto.rpc.FindRes>}
 */
const methodInfo_DotConfigFace_FindDot = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rpc.FindRes,
  /**
   * @param {!proto.rpc.FindReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.FindRes.deserializeBinary
);


/**
 * @param {!proto.rpc.FindReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rpc.FindRes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rpc.FindRes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rpc.DotConfigFaceClient.prototype.findDot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rpc.DotConfigFace/FindDot',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_FindDot,
      callback);
};


/**
 * @param {!proto.rpc.FindReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rpc.FindRes>}
 *     A native promise that resolves to the response
 */
proto.rpc.DotConfigFacePromiseClient.prototype.findDot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rpc.DotConfigFace/FindDot',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_FindDot);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rpc.ImportReq,
 *   !proto.rpc.ImportRes>}
 */
const methodDescriptor_DotConfigFace_ImportByConfig = new grpc.web.MethodDescriptor(
  '/rpc.DotConfigFace/ImportByConfig',
  grpc.web.MethodType.UNARY,
  proto.rpc.ImportReq,
  proto.rpc.ImportRes,
  /**
   * @param {!proto.rpc.ImportReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.ImportRes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rpc.ImportReq,
 *   !proto.rpc.ImportRes>}
 */
const methodInfo_DotConfigFace_ImportByConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rpc.ImportRes,
  /**
   * @param {!proto.rpc.ImportReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.ImportRes.deserializeBinary
);


/**
 * @param {!proto.rpc.ImportReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rpc.ImportRes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rpc.ImportRes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rpc.DotConfigFaceClient.prototype.importByConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rpc.DotConfigFace/ImportByConfig',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_ImportByConfig,
      callback);
};


/**
 * @param {!proto.rpc.ImportReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rpc.ImportRes>}
 *     A native promise that resolves to the response
 */
proto.rpc.DotConfigFacePromiseClient.prototype.importByConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rpc.DotConfigFace/ImportByConfig',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_ImportByConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rpc.ImportReq,
 *   !proto.rpc.ImportRes>}
 */
const methodDescriptor_DotConfigFace_ImportByDot = new grpc.web.MethodDescriptor(
  '/rpc.DotConfigFace/ImportByDot',
  grpc.web.MethodType.UNARY,
  proto.rpc.ImportReq,
  proto.rpc.ImportRes,
  /**
   * @param {!proto.rpc.ImportReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.ImportRes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rpc.ImportReq,
 *   !proto.rpc.ImportRes>}
 */
const methodInfo_DotConfigFace_ImportByDot = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rpc.ImportRes,
  /**
   * @param {!proto.rpc.ImportReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.ImportRes.deserializeBinary
);


/**
 * @param {!proto.rpc.ImportReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rpc.ImportRes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rpc.ImportRes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rpc.DotConfigFaceClient.prototype.importByDot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rpc.DotConfigFace/ImportByDot',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_ImportByDot,
      callback);
};


/**
 * @param {!proto.rpc.ImportReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rpc.ImportRes>}
 *     A native promise that resolves to the response
 */
proto.rpc.DotConfigFacePromiseClient.prototype.importByDot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rpc.DotConfigFace/ImportByDot',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_ImportByDot);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rpc.ExportReq,
 *   !proto.rpc.ExportRes>}
 */
const methodDescriptor_DotConfigFace_ExportConfig = new grpc.web.MethodDescriptor(
  '/rpc.DotConfigFace/ExportConfig',
  grpc.web.MethodType.UNARY,
  proto.rpc.ExportReq,
  proto.rpc.ExportRes,
  /**
   * @param {!proto.rpc.ExportReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.ExportRes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rpc.ExportReq,
 *   !proto.rpc.ExportRes>}
 */
const methodInfo_DotConfigFace_ExportConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rpc.ExportRes,
  /**
   * @param {!proto.rpc.ExportReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.ExportRes.deserializeBinary
);


/**
 * @param {!proto.rpc.ExportReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rpc.ExportRes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rpc.ExportRes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rpc.DotConfigFaceClient.prototype.exportConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rpc.DotConfigFace/ExportConfig',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_ExportConfig,
      callback);
};


/**
 * @param {!proto.rpc.ExportReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rpc.ExportRes>}
 *     A native promise that resolves to the response
 */
proto.rpc.DotConfigFacePromiseClient.prototype.exportConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rpc.DotConfigFace/ExportConfig',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_ExportConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rpc.ExportReq,
 *   !proto.rpc.ExportRes>}
 */
const methodDescriptor_DotConfigFace_ExportDot = new grpc.web.MethodDescriptor(
  '/rpc.DotConfigFace/ExportDot',
  grpc.web.MethodType.UNARY,
  proto.rpc.ExportReq,
  proto.rpc.ExportRes,
  /**
   * @param {!proto.rpc.ExportReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.ExportRes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rpc.ExportReq,
 *   !proto.rpc.ExportRes>}
 */
const methodInfo_DotConfigFace_ExportDot = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rpc.ExportRes,
  /**
   * @param {!proto.rpc.ExportReq} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rpc.ExportRes.deserializeBinary
);


/**
 * @param {!proto.rpc.ExportReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rpc.ExportRes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rpc.ExportRes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rpc.DotConfigFaceClient.prototype.exportDot =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rpc.DotConfigFace/ExportDot',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_ExportDot,
      callback);
};


/**
 * @param {!proto.rpc.ExportReq} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rpc.ExportRes>}
 *     A native promise that resolves to the response
 */
proto.rpc.DotConfigFacePromiseClient.prototype.exportDot =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rpc.DotConfigFace/ExportDot',
      request,
      metadata || {},
      methodDescriptor_DotConfigFace_ExportDot);
};


module.exports = proto.rpc;

