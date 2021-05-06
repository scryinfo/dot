import {DotConfigFacePromiseClient} from "@/rpc_face/config_grpc_web_pb";
import {ExportReq, FindReq, ImportReq} from "@/rpc_face/config_pb";

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

    public exportConfig(data: string, fileName: string[]) {
        const request = new ExportReq();
        request.setConfigdata(data);
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