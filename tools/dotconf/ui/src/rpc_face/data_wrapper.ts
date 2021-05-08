import {DotConfigFacePromiseClient} from "@/rpc_face/config_grpc_web_pb";
import {ExportReq, FindReq, ImportReq} from "@/rpc_face/config_pb";
import {Dot} from "@/views/home/store";

export class _DotWrapper {
    private client = new DotConfigFacePromiseClient(
        process.env.VUE_APP_BASE_API as string,
        null,
        null
    );

    public findDot(dirs: string[]) {
        const request = new FindReq();
        request.setDirsList(dirs);
        return this.client.findDot(request);
    }

    public importByConfig(path: string) {
        const request = new ImportReq();
        request.setFilepath(path);
        return this.client.importByConfig(request);
    }

    public importByDot(path: string) {
        const request = new ImportReq();
        request.setFilepath(path);
        return this.client.importByDot(request);
    }
    public initImportDot() {
        const request = new ImportReq();
        return this.client.initImport(request);
    }

    public exportConfig(data: Dot[], fileName: string[]) {
        const request = new ExportReq();
        let result = {
            log: {
                file: "log.log",
                level: "debug"
            },
            dots: {},
        };
        let dots = new Array<dot>();
        for (let i = 0; i < data.length; i++) {
            const item = new dot();
            item.metaData.name = data[i].metaData.name;
            item.metaData.typeId = data[i].metaData.typeId;
            item.lives.length = 0;
            for (let j = 0; j < data[i].lives.length; j++) {
                const itemLive = {
                    liveId: data[i].lives[j].liveId,
                    relyLives: data[i].lives[j].relyLives,
                    json: data[i].lives[j].json
                };
                item.lives.push(itemLive)
            }
            dots.push(item)
        }
        result.dots = dots;
        request.setConfigdata(JSON.stringify(result));
        request.setFilenameList(fileName);
        return this.client.exportConfig(request);
    }

    public exportDot(data: string, fileName: string[]) {
        const request = new ExportReq();
        request.setFilenameList(fileName);
        request.setDotdata(data);
        return this.client.exportDot(request);
    }
}

const DotWrapper = new _DotWrapper();
export default DotWrapper;

class dot {
    metaData = {
        name: "",
        typeId: "",
    };
    lives = [{
        liveId: "",
        relyLives: null,
        json: {}
    }];
}