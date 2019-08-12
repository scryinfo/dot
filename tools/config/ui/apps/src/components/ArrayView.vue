<template>
    <div>
        <el-row><el-col :span="24"><el-button @click="addEvent()">add</el-button><el-button @click="ShowJsonDialog(parsedData)">JSON</el-button></el-col></el-row>
        <el-row v-for="(member,index) in flowData " v-model="flowData">
            <el-col :span="18">
                <div v-if="member.type !== 'object' && member.type !== 'array'"  class="grid-content bg-purple-light">
                <el-input type="text"
                          v-model="flowData[index].remark"
                          v-if="member.type == 'string'">
                </el-input>
                <el-input
                        type="number"
                        v-model.number="flowData[index].remark"
                        v-if="member.type == 'number'">
                </el-input>
                <bool-view
                        v-model="flowData[index].remark"
                        :boolValue="flowData[index].remark"
                        v-if="member.type == 'boolean'"
                >
                </bool-view>
            </div>
            <div v-else  class="grid-content bg-purple-light">
                <json-view v-model="flowData[index].childParams" :parsedData="flowData[index].childParams"></json-view>
            </div>
            </el-col>
            <el-col :span="4"><el-button :disabled="cantRemove" @click="removeEvent(index)">remove</el-button></el-col></el-row>
        <el-row></el-row>

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
    import Vue from 'vue';
    export default Vue.extend({
        name: "ArrayView",
        props: {
            parsedData: {},
        },
        data () {
            return {
                flowData: (this as any).parsedData.childParams,
                dialog: false,
                objc: [],
                textarea: '',
                cantRemove: true,
                schemaObject: {},
                temp: {},
                nameTemp: ''
            }
        },
        watch: {
            parsedData: {
                handler(newValue, oldValue) {
                    this.flowData = (this as any).parsedData.childParams;
                },
                immediate: true
            },
            flowData: {
                handler(newValue, oldValue) {
                    if(newValue.length > 1){
                        this.cantRemove = false;
                    }
                    if (newValue.length === 1){
                        this.cantRemove = true;
                    }
                    this.$emit('input',newValue);
                },
                deep: true
            }
        },
        methods: {
            ShowJsonDialog(obj:any){
                this.dialog = true;
                (this as any).objc.push(obj);
                let jsonSchemaGenerator = require('./schemaGenerator/index.js');
                let data = {};
                this.temp = this.makeJson(this.objc);
                this.nameTemp = obj.name;
                eval("data = this.temp."+this.nameTemp);
                this.schemaObject = jsonSchemaGenerator.jsonToSchema(this.temp);
                this.textarea = JSON.stringify(data,null,4);
            },
            handleClose(done:any){
                try{
                    if(this.textarea){
                        let data = JSON.parse(this.textarea);
                        eval("this.temp."+this.nameTemp+"= data");
                        let tv4 = require('tv4');
                        if(tv4.validate(this.temp, this.schemaObject)){
                            let objct:any = this.jsonParse(data);
                            this.$emit('input',objct);
                        }else{
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
            },
            addEvent () {
                this.flowData.push(this.shallowCopy(this.flowData[this.flowData.length-1]));
            },
            shallowCopy(src:any):any {
                let dst:any = {};
                for (let prop in src) {
                    if (src.hasOwnProperty(prop)) {
                        dst[prop] = src[prop];
                    }
                }
                return dst;
            },
            removeEvent(index:number) {
                this.flowData.splice(index,1);
            }
        }
    })
</script>

<style scoped>
    .bg-purple-light {
        background: #e5e9f2;
    }
    .grid-content {
        border-radius: 4px;
        min-height: 36px;
    }
</style>
