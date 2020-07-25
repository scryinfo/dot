import {HiDot, HiDotClient} from "@/api/hi_pb_service";
import {HelloRequest, HelloResponse, HiReq, HiRes, WriteReq, WriteRes} from "@/api/hi_pb";
import {grpc} from "@improbable-eng/grpc-web";

export class HiWrapper {
    private client = grpc.client(HiDot.Hi, {
        host: process.env.VUE_APP_BASE_API,
        transport: grpc.WebsocketTransport(),
    });
    //{transport: grpc.WebsocketTransport()}
    private hiDot = new HiDotClient(process.env.VUE_APP_BASE_API,{transport: grpc.WebsocketTransport()});
    private clientStream_ : grpc.Client<HelloRequest, HelloResponse> | null  = null;
    private serverStream_ : grpc.Client<HelloRequest, HelloResponse> | null  = null;
    private bothStream_ : grpc.Client<HelloRequest, HelloResponse> | null  = null;
    public hi(name: string) {
        const req = new HiReq();
        req.setName(name);

        return new Promise<HiRes | null>( (resolve, reject) => {
            this.hiDot.hi(req,(err, res) =>{
               if (err != null) {
                   reject(err);
               }else{
                   resolve(res);
               }
            });
            // grpc.unary(HiDot.Hi,{
            //     request: req,
            //     host: process.env.VUE_APP_BASE_API,
            //     transport: grpc.WebsocketTransport(),
            //     onEnd: res =>{
            //         const { status, statusMessage, headers, message, trailers } = res;
            //         if (status === grpc.Code.OK && message) {
            //             resolve(message as HiRes);
            //         }else{
            //             reject(statusMessage);
            //         }
            //     }
            // });
        });

        // return new Promise<HiRes | null>( (resolve, reject) => {
        //     this.hiDot.hi(req, (err, res) => {
        //         if(err != null){
        //             reject(res);
        //         }else{
        //             resolve(res);
        //         }
        //     })
        // });
    }

    public write(data: string) {
        const req = new WriteReq();
        req.setData(data);
        return new Promise<WriteRes | null>((resolve, reject) => {
            this.hiDot.write(req, (err, res) =>{
                if (err != null) {
                    reject(res);
                }else{
                    resolve(res);
                }
            })
        })
    }

    public clientStream(data:string) {
        const req = new HelloRequest();
        req.setGreeting(data);
        let client = this.clientStream_ as grpc.Client<HelloRequest, HelloResponse>;
        if (!client) {
            //{host:process.env.VUE_APP_BASE_API, transport: grpc.WebsocketTransport()}
            client = grpc.client<HelloRequest, HelloResponse, grpc.MethodDefinition<HelloRequest, HelloResponse> >(HiDot.ClientStream, {host:process.env.VUE_APP_BASE_API, transport: grpc.WebsocketTransport()});
            this.clientStream_ = client;
            client.start();
        }

        return new Promise<HelloResponse> ((resolve, reject) => {
            client.onMessage(msg =>{
                if (msg.getReply().startsWith("close")) {
                    client.close();
                    this.clientStream_ = null;
                }
                resolve(msg);
            });
            client.onEnd((code,message,tra) =>{
                client.close();
                this.clientStream_ = null;
                if(code != 0) {
                    reject(message);
                }
            });

            client.send(req);
        })
    }

    public serverStream(data: string) {
        const req = new HelloRequest();
        req.setGreeting(data);

        let client = this.serverStream_ as grpc.Client<HelloRequest, HelloResponse>;
        if (!client) {
            //{host:process.env.VUE_APP_BASE_API, transport: grpc.WebsocketTransport()}
            client = grpc.client<HelloRequest, HelloResponse, grpc.MethodDefinition<HelloRequest, HelloResponse> >(HiDot.ServerStream, {host:process.env.VUE_APP_BASE_API, transport: grpc.WebsocketTransport()});
            this.serverStream_ = client; //由于 HiDot.ServerStream.requestStream = false, 所以这个流只能发送一次
            client.start();
        }

        return new Promise<HelloResponse>((resolve, reject)=>{
            client.onMessage(msg =>{
                if (msg.getReply().startsWith("close")) {
                    client.close();
                    this.serverStream_ = null;
                }
                resolve(msg);
            });
            client.onEnd((code,message,tra) =>{
                // client.close();
                // this.serverStream_ = null;
                if(code != 0) {
                    reject(message);
                }
            });
            client.send(req);
        })

        // return new Promise<HelloResponse>((resolve, reject)=>{
        //     const stream = this.hiDot.serverStream(req);
        //     stream.on("data", res => {
        //         stream.cancel(); //it w
        //         resolve(res);
        //     });
        //     stream.on("status", status => {
        //        if (status.code != 0) {
        //            reject(status);
        //        }
        //     });
        //     stream.on("end", status => {
        //         stream.cancel()
        //     })
        // })
    }

    public bothSides(data: string) {
        const req = new HelloRequest();
        req.setGreeting(data);

        let client = this.bothStream_ as grpc.Client<HelloRequest, HelloResponse>;
        if (!client) {
            //{host:process.env.VUE_APP_BASE_API, transport: grpc.WebsocketTransport()}
            //{host:process.env.VUE_APP_BASE_API}
            client = grpc.client<HelloRequest, HelloResponse, grpc.MethodDefinition<HelloRequest, HelloResponse> >(HiDot.BothSides,{host:process.env.VUE_APP_BASE_API, transport: grpc.WebsocketTransport()} );
            this.bothStream_ = client;
            client.start();
        }

        return new Promise<HelloResponse>((resolve, reject) =>{
            client.onMessage(msg => {
                if (msg.getReply().startsWith("close")) {
                    client.close();
                    this.bothStream_ = null;
                }
                resolve(msg);
            });
            client.onEnd((code, message, trailers) =>{
                client.close();
                this.bothStream_ = null;
                if(code != 0) {
                    reject(message);
                }
            });
            client.send(req);
        })
    }
}
const Hi = new HiWrapper();
export default Hi;

