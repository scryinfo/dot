import {HiDotPromiseClient} from "@/api/hi_grpc_web_pb";
import {HiReq, WriteReq} from "@/api/hi_pb";

export class HiWrapper {
    private hiDot = new HiDotPromiseClient(process.env.VUE_APP_BASE_API, null, null);
    public hi(name: string) {
        const req = new HiReq();
        req.setName(name);
        return this.hiDot.hi(req);
    }

    public write(data: string) {
        const req = new WriteReq();
        req.setData(data);
        return this.hiDot.write(req);
    }
}
const Hi = new HiWrapper();
export default Hi;

