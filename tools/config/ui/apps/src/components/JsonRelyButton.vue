<template>
    <div>
        <el-button @click="ShowJsonDialog(objc)">JSON</el-button>
        <el-drawer
                title="JSON textarea!"
                :before-close="handleClose"
                :visible.sync="dialog"
                direction="ltr"
                custom-class="demo-drawer"
                ref="drawer"
        >
            <el-input
                    type="textarea"
                    :autosize="{ minRows: 10, maxRows: 30}"
                    placeholder="请输入内容"
                    v-model="textarea">
            </el-input>
        </el-drawer>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue'
    export default Vue.extend({
        name: "JsonRelyButton",
        props: {
            objc: {}
        },
        data() {
            return {
                dialog: false,
                schemaObject: {},
                textarea: ''
            }
        },
        methods: {
            ShowJsonDialog(obj:any){
                this.dialog = true;
                (this as any).objc = obj;
                let data = this.makeJson(obj);
                this.textarea = JSON.stringify(data,null,4);

            },
            handleClose(done:any){
                try{
                    if(this.textarea){
                        let objct:any = JSON.parse(this.textarea);
                        if(this.inputCheck(objct)){
                            this.$emit('input',this.jsonParse(objct));
                        }else {
                            (this as any).$message.error('json text input error!');
                        }
                    }else{
                        (this as any).$message.error('json text input error!');
                    }
                }catch (e) {
                    (this as any).$message.error('json text input error!');
                }finally {
                    done();
                }
            },
            inputCheck(input:any): Boolean{
                for (let i in input){
                    if (typeof input[i] !== typeof ''){
                        return false;
                    }
                }
                return true;
            },
            jsonParse: function (jsonStr:any) {
                let parseJson = (json:any) => {
                    let result:any = [];
                    let keys = Object.keys(json);
                    keys.forEach((k, index) => {
                        let val = json[k];
                        let parsedVal = val;
                        if (this.getType(val) == "object") {
                            parsedVal = parseJson(val);

                        } else if (this.getType(val) == "array") {
                            parsedVal = parseArray(val);
                        }

                        let opt:any = {
                            name: k,
                            type: this.getType(val)
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
                        if (this.getType(val) == "object") {
                            parsedVal = parseJson(val);

                        } else if (this.getType(val) == "array") {
                            parsedVal = parseArray(val);
                        }

                        let opt:any = {
                            name: null,
                            type: this.getType(val)
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
            },

            getType: function (obj:any) {
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
            },
            makeJson: function (dataArr:any) {
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
        }
    })
</script>

<style scoped>

</style>
