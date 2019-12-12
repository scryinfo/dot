import * as grpcWeb from 'grpc-web';

import {
  HiReq,
  HiRes,
  WriteReq,
  WriteRes} from './hi_pb';

export class HiDotClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  hi(
    request: HiReq,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: HiRes) => void
  ): grpcWeb.ClientReadableStream<HiRes>;

  write(
    request: WriteReq,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: WriteRes) => void
  ): grpcWeb.ClientReadableStream<WriteRes>;

}

export class HiDotPromiseClient {
  constructor (hostname: string,
               credentials: null | { [index: string]: string; },
               options: null | { [index: string]: string; });

  hi(
    request: HiReq,
    metadata?: grpcWeb.Metadata
  ): Promise<HiRes>;

  write(
    request: WriteReq,
    metadata?: grpcWeb.Metadata
  ): Promise<WriteRes>;

}

