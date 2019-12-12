import * as jspb from "google-protobuf"

export class HiReq extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HiReq.AsObject;
  static toObject(includeInstance: boolean, msg: HiReq): HiReq.AsObject;
  static serializeBinaryToWriter(message: HiReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HiReq;
  static deserializeBinaryFromReader(message: HiReq, reader: jspb.BinaryReader): HiReq;
}

export namespace HiReq {
  export type AsObject = {
    name: string,
  }
}

export class HiRes extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HiRes.AsObject;
  static toObject(includeInstance: boolean, msg: HiRes): HiRes.AsObject;
  static serializeBinaryToWriter(message: HiRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HiRes;
  static deserializeBinaryFromReader(message: HiRes, reader: jspb.BinaryReader): HiRes;
}

export namespace HiRes {
  export type AsObject = {
    name: string,
  }
}

export class WriteReq extends jspb.Message {
  getData(): string;
  setData(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WriteReq.AsObject;
  static toObject(includeInstance: boolean, msg: WriteReq): WriteReq.AsObject;
  static serializeBinaryToWriter(message: WriteReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WriteReq;
  static deserializeBinaryFromReader(message: WriteReq, reader: jspb.BinaryReader): WriteReq;
}

export namespace WriteReq {
  export type AsObject = {
    data: string,
  }
}

export class WriteRes extends jspb.Message {
  getData(): string;
  setData(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WriteRes.AsObject;
  static toObject(includeInstance: boolean, msg: WriteRes): WriteRes.AsObject;
  static serializeBinaryToWriter(message: WriteRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WriteRes;
  static deserializeBinaryFromReader(message: WriteRes, reader: jspb.BinaryReader): WriteRes;
}

export namespace WriteRes {
  export type AsObject = {
    data: string,
  }
}

