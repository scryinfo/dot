export let jsonParse = function (jsonStr:any) {
    let parseJson = (json:any) => {
        let result:any = [];
        let keys = Object.keys(json);
        keys.forEach((k, index) => {
            let val = json[k];
            let parsedVal = val;
            if (getType(val) == "object") {
                parsedVal = parseJson(val);

            } else if (getType(val) == "array") {
                parsedVal = parseArray(val);
            }

            let opt:any = {
                name: k,
                type: getType(val)
            };

            if (opt.type == "array" || opt.type == "object") {
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
    let parseArray = (arrayObj:any) => {
        let result = [];
        for (let i = 0; i < arrayObj.length; ++i) {
            let val = arrayObj[i];
            let parsedVal = val;
            if (getType(val) == "object") {
                parsedVal = parseJson(val);

            } else if (getType(val) == "array") {
                parsedVal = parseArray(val);
            }

            let opt:any = {
                name: null,
                type: getType(val)
            };

            if (opt.type == "array" || opt.type == "object") {
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
    let parseBody = (json:any) => {
        let r = parseJson(json);
        return r;
    };

    return parseBody(jsonStr);
};

export let getType = function (obj:any) {
    switch (Object.prototype.toString.call(obj)) {
        case "[object Array]":
            return "array";
            break;
        case "[object Object]":
            return "object";
            break;
        default:
            return typeof obj;
            break;
    }
};

export let makeJson = function (dataArr:any) {
    let revertWithObj = function(data:any) {
        let r:any = {};
        for (let i = 0; i < data.length; ++i) {
            let el = data[i];
            let key, val;
            key = el.name;
            if (el.type == "array") {
                val = revertWithArray(el.childParams);
            } else if (el.type == "object") {
                val = revertWithObj(el.childParams);
            } else {
                val = el.remark;
            }

            r[key] = val;
        }
        return r;
    };

    let revertWithArray:any = function(data:any) {
        let arr = [];
        for (let i = 0; i < data.length; ++i) {
            let el = data[i];
            let r;
            if (el.type == "array") {
                r = revertWithArray(el.childParams);
            } else if (el.type == "object") {
                r = revertWithObj(el.childParams);
            } else {
                r = el.remark;
            }

            arr.push(r);
        }
        return arr;
    };

    let revertMain = function(data:any) {
        let r = revertWithObj(data);
        return r;
    };

    return revertMain(dataArr);
}
export function jsonParseRely(Json:any) {
    if (Json === null){
        Json = {}
    }
    let result:any = [];
    let keys = Object.keys(Json);
    keys.forEach((k, index) => {
        let val = Json[k];
        let newObject = {name: k, remark: val};
        result.push(newObject)
    });
    return result;
};
export function makeJsonRely(ParData:any) {
    let Revert:any = {};
    for (let i = 0; i < ParData.length; ++i) {
        let el = ParData[i];
        let key:string, val:string;
        key = el.name;
        val = el.remark;
        Revert[key] = val;
    }
    return Revert;
}
