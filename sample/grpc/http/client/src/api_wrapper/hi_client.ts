import {HiDotClient} from "@/api/hi_pb_service";
import {HiReq, WriteReq} from "@/api/hi_pb";

export class HiWrapper {
    private hiDot = new HiDotClient(process.env.VUE_APP_BASE_API);
    public hi(name: string) {
        const req = new HiReq();
        req.setName(name);
        return new Promise( (resolve, reject) => {
            this.hiDot.hi(req, (err, res) => {
                if(err != null){
                    reject(res);
                }else{
                    resolve(res);
                }
            })
        });
    }

    public write(data: string) {
        const req = new WriteReq();
        req.setData(data);
        return new Promise((resolve, reject) => {
            this.hiDot.write(req, (err, res) =>{
                if (err != null) {
                    reject(res);
                }else{
                    resolve(res);
                }
            })
        })
    }
}
const Hi = new HiWrapper();
export default Hi;

