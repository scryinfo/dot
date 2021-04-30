export let jsonParse = function (jsonStr: any) {
    const parseJson = (json: any) => {
        const result: any = [];
        const keys = Object.keys(json);
        keys.forEach((k, index) => {
            const val = json[k];
            let parsedVal = val;
            if (getType(val) == 'object') {
                parsedVal = parseJson(val);

            } else if (getType(val) == 'array') {
                parsedVal = parseArray(val);
            }

            const opt: any = {
                name: k,
                type: getType(val),
            };

            if (opt.type == 'array' || opt.type == 'object') {
                opt.childParams = parsedVal;
                opt.remark = null;
            } else {
                opt.childParams = null;
                opt.remark = parsedVal;
            }

            result.push(opt);
        });
        return result;
    };

    //
    const parseArray = (arrayObj: any) => {
        const result = [];
        for (let i = 0; i < arrayObj.length; ++i) {
            const val = arrayObj[i];
            let parsedVal = val;
            if (getType(val) == 'object') {
                parsedVal = parseJson(val);

            } else if (getType(val) == 'array') {
                parsedVal = parseArray(val);
            }

            const opt: any = {
                name: null,
                type: getType(val),
            };

            if (opt.type == 'array' || opt.type == 'object') {
                opt.childParams = parsedVal;
                opt.remark = null;
            } else {
                opt.childParams = null;
                opt.remark = parsedVal;
            }

            result.push(opt);
        }
        return result;
    };

    // --
    const parseBody = (json: any) => {
        const r = parseJson(json);
        return r;
    };

    return parseBody(jsonStr);
};

export let getType = function (obj: any) {
    switch (Object.prototype.toString.call(obj)) {
        case '[object Array]':
            return 'array';
        case '[object Object]':
            return 'object';
        default:
            return typeof obj;
    }
};

export let makeJson = function (dataArr: any) {
    const revertWithObj = function (data: any) {
        const r: any = {};
        for (let i = 0; i < data.length; ++i) {
            const el = data[i];
            let key, val;
            key = el.name;
            if (el.type == 'array') {
                val = revertWithArray(el.childParams);
            } else if (el.type == 'object') {
                val = revertWithObj(el.childParams);
            } else {
                val = el.remark;
            }

            r[key] = val;
        }
        return r;
    };

    const revertWithArray: any = function (data: any) {
        const arr = [];
        for (let i = 0; i < data.length; ++i) {
            const el = data[i];
            let r;
            if (el.type == 'array') {
                r = revertWithArray(el.childParams);
            } else if (el.type == 'object') {
                r = revertWithObj(el.childParams);
            } else {
                r = el.remark;
            }

            arr.push(r);
        }
        return arr;
    };

    const revertMain = function (data: any) {
        const r = revertWithObj(data);
        return r;
    };

    return revertMain(dataArr);
};

export function jsonParseRely(Json: any) {
    if (Json === null || Json === undefined) {
        Json = {};
    }
    const result: any = [];
    const keys = Object.keys(Json);
    keys.forEach((k, index) => {
        const val = Json[k];
        const newObject = {name: k, remark: val};
        result.push(newObject);
    });
    return result;
}

export function makeJsonRely(ParData: any) {
    const Revert: any = {};
    for (let i = 0; i < ParData.length; ++i) {
        const el = ParData[i];
        let key: string, val: string;
        key = el.name;
        val = el.remark;
        Revert[key] = val;
    }
    return Revert;
}
