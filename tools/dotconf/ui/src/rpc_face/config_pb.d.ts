import * as jspb from "google-protobuf"

export class FindReq extends jspb.Message {
  getDirsList(): Array<string>;
  setDirsList(value: Array<string>): FindReq;
  clearDirsList(): FindReq;
  addDirs(value: string, index?: number): FindReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindReq.AsObject;
  static toObject(includeInstance: boolean, msg: FindReq): FindReq.AsObject;
  static serializeBinaryToWriter(message: FindReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindReq;
  static deserializeBinaryFromReader(message: FindReq, reader: jspb.BinaryReader): FindReq;
}

export namespace FindReq {
  export type AsObject = {
    dirsList: Array<string>,
  }
}

export class FindRes extends jspb.Message {
  getDotsinfo(): string;
  setDotsinfo(value: string): FindRes;

  getNoexistdirsList(): Array<string>;
  setNoexistdirsList(value: Array<string>): FindRes;
  clearNoexistdirsList(): FindRes;
  addNoexistdirs(value: string, index?: number): FindRes;

  getError(): string;
  setError(value: string): FindRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FindRes.AsObject;
  static toObject(includeInstance: boolean, msg: FindRes): FindRes.AsObject;
  static serializeBinaryToWriter(message: FindRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FindRes;
  static deserializeBinaryFromReader(message: FindRes, reader: jspb.BinaryReader): FindRes;
}

export namespace FindRes {
  export type AsObject = {
    dotsinfo: string,
    noexistdirsList: Array<string>,
    error: string,
  }
}

export class ImportReq extends jspb.Message {
  getFilepath(): string;
  setFilepath(value: string): ImportReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportReq.AsObject;
  static toObject(includeInstance: boolean, msg: ImportReq): ImportReq.AsObject;
  static serializeBinaryToWriter(message: ImportReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportReq;
  static deserializeBinaryFromReader(message: ImportReq, reader: jspb.BinaryReader): ImportReq;
}

export namespace ImportReq {
  export type AsObject = {
    filepath: string,
  }
}

export class ImportRes extends jspb.Message {
  getJson(): string;
  setJson(value: string): ImportRes;

  getError(): string;
  setError(value: string): ImportRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportRes.AsObject;
  static toObject(includeInstance: boolean, msg: ImportRes): ImportRes.AsObject;
  static serializeBinaryToWriter(message: ImportRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportRes;
  static deserializeBinaryFromReader(message: ImportRes, reader: jspb.BinaryReader): ImportRes;
}

export namespace ImportRes {
  export type AsObject = {
    json: string,
    error: string,
  }
}

export class ExportReq extends jspb.Message {
  getConfigdata(): string;
  setConfigdata(value: string): ExportReq;

  getFilenameList(): Array<string>;
  setFilenameList(value: Array<string>): ExportReq;
  clearFilenameList(): ExportReq;
  addFilename(value: string, index?: number): ExportReq;

  getDotdata(): string;
  setDotdata(value: string): ExportReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExportReq.AsObject;
  static toObject(includeInstance: boolean, msg: ExportReq): ExportReq.AsObject;
  static serializeBinaryToWriter(message: ExportReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExportReq;
  static deserializeBinaryFromReader(message: ExportReq, reader: jspb.BinaryReader): ExportReq;
}

export namespace ExportReq {
  export type AsObject = {
    configdata: string,
    filenameList: Array<string>,
    dotdata: string,
  }
}

export class ExportRes extends jspb.Message {
  getError(): string;
  setError(value: string): ExportRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExportRes.AsObject;
  static toObject(includeInstance: boolean, msg: ExportRes): ExportRes.AsObject;
  static serializeBinaryToWriter(message: ExportRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExportRes;
  static deserializeBinaryFromReader(message: ExportRes, reader: jspb.BinaryReader): ExportRes;
}

export namespace ExportRes {
  export type AsObject = {
    error: string,
  }
}

