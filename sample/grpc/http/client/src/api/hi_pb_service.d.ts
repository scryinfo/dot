// package: hidot
// file: hi.proto

import * as hi_pb from "./hi_pb";
import {grpc} from "@improbable-eng/grpc-web";

type HiDotHi = {
  readonly methodName: string;
  readonly service: typeof HiDot;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof hi_pb.HiReq;
  readonly responseType: typeof hi_pb.HiRes;
};

type HiDotWrite = {
  readonly methodName: string;
  readonly service: typeof HiDot;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof hi_pb.WriteReq;
  readonly responseType: typeof hi_pb.WriteRes;
};

type HiDotServerStream = {
  readonly methodName: string;
  readonly service: typeof HiDot;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof hi_pb.HelloRequest;
  readonly responseType: typeof hi_pb.HelloResponse;
};

type HiDotClientStream = {
  readonly methodName: string;
  readonly service: typeof HiDot;
  readonly requestStream: true;
  readonly responseStream: false;
  readonly requestType: typeof hi_pb.HelloRequest;
  readonly responseType: typeof hi_pb.HelloResponse;
};

type HiDotBothSides = {
  readonly methodName: string;
  readonly service: typeof HiDot;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof hi_pb.HelloRequest;
  readonly responseType: typeof hi_pb.HelloResponse;
};

export class HiDot {
  static readonly serviceName: string;
  static readonly Hi: HiDotHi;
  static readonly Write: HiDotWrite;
  static readonly ServerStream: HiDotServerStream;
  static readonly ClientStream: HiDotClientStream;
  static readonly BothSides: HiDotBothSides;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class HiDotClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  hi(
    requestMessage: hi_pb.HiReq,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: hi_pb.HiRes|null) => void
  ): UnaryResponse;
  hi(
    requestMessage: hi_pb.HiReq,
    callback: (error: ServiceError|null, responseMessage: hi_pb.HiRes|null) => void
  ): UnaryResponse;
  write(
    requestMessage: hi_pb.WriteReq,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: hi_pb.WriteRes|null) => void
  ): UnaryResponse;
  write(
    requestMessage: hi_pb.WriteReq,
    callback: (error: ServiceError|null, responseMessage: hi_pb.WriteRes|null) => void
  ): UnaryResponse;
  serverStream(requestMessage: hi_pb.HelloRequest, metadata?: grpc.Metadata): ResponseStream<hi_pb.HelloResponse>;
  clientStream(metadata?: grpc.Metadata): RequestStream<hi_pb.HelloRequest>;
  bothSides(metadata?: grpc.Metadata): BidirectionalStream<hi_pb.HelloRequest, hi_pb.HelloResponse>;
}

