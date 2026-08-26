import { V7Generator } from 'uuidv7'
import { ReqbaseSchema, type Reqbase } from '../oidcapi/v1/reqsbase_pb'
import { create } from '@bufbuild/protobuf'

const uuidv7 = new V7Generator()
export function newReqbase(): Reqbase {
  const bytes = uuidv7.generate().bytes
  const view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength)
  const reqIdLow = view.getBigUint64(0, false)
  const reqIdHigh = view.getBigUint64(8, false)
  return create(ReqbaseSchema, {
    reqIdLow: reqIdLow,
    reqIdHigh: reqIdHigh,
  })
}
