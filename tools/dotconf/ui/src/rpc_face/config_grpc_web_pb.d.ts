import * as grpcWeb from 'grpc-web';

import {
  ExportReq,
  ExportRes,
  FindReq,
  FindRes,
  ImportReq,
  ImportRes} from './config_pb';

export class DotConfigFaceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  findDot(
    request: FindReq,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: FindRes) => void
  ): grpcWeb.ClientReadableStream<FindRes>;

  importByConfig(
    request: ImportReq,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ImportRes) => void
  ): grpcWeb.ClientReadableStream<ImportRes>;

  importByDot(
    request: ImportReq,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ImportRes) => void
  ): grpcWeb.ClientReadableStream<ImportRes>;

  initImport(
    request: ImportReq,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ImportRes) => void
  ): grpcWeb.ClientReadableStream<ImportRes>;

  exportConfig(
    request: ExportReq,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ExportRes) => void
  ): grpcWeb.ClientReadableStream<ExportRes>;

  exportDot(
    request: ExportReq,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ExportRes) => void
  ): grpcWeb.ClientReadableStream<ExportRes>;

}

export class DotConfigFacePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  findDot(
    request: FindReq,
    metadata?: grpcWeb.Metadata
  ): Promise<FindRes>;

  importByConfig(
    request: ImportReq,
    metadata?: grpcWeb.Metadata
  ): Promise<ImportRes>;

  importByDot(
    request: ImportReq,
    metadata?: grpcWeb.Metadata
  ): Promise<ImportRes>;

  initImport(
    request: ImportReq,
    metadata?: grpcWeb.Metadata
  ): Promise<ImportRes>;

  exportConfig(
    request: ExportReq,
    metadata?: grpcWeb.Metadata
  ): Promise<ExportRes>;

  exportDot(
    request: ExportReq,
    metadata?: grpcWeb.Metadata
  ): Promise<ExportRes>;

}

