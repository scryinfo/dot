<template>
    <div>
        <el-row>
            <el-col :span="24">
                <el-button @click="addEvent()">add</el-button>
                <el-button @click="ShowJsonDialog(parsedData)">JSON</el-button>
            </el-col>
        </el-row>
        <el-row v-for="(member,index) in flowData " v-model="flowData">
            <el-col :span="18">
                <div v-if="member.type !== 'object' && member.type !== 'array'" class="grid-content bg-purple-light">
                    <el-input type="text"
                              v-model="flowData[index].remark"
                              v-if="member.type == 'string'">
                    </el-input>
                    <el-input-number
                            v-model.number="flowData[index].remark"
                            v-if="member.type == 'number'">
                    </el-input-number>
                    <select
                            name="value"
                            v-model="flowData[index].remark"
                            class="val-input"
                            v-if="member.type == 'boolean'"
                    >
                        <option :value="true">true</option>
                        <option :value="false">false</option>
                    </select>
                </div>
                <div v-else class="grid-content bg-purple-light">
                    <json-view v-model="flowData[index].childParams"
                               :parsedData="flowData[index].childParams"></json-view>
                </div>
            </el-col>
            <el-col :span="4">
                <el-button :disabled="cantRemove" @click="removeEvent(index)">remove</el-button>
            </el-col>
        </el-row>
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
import {jsonParse, makeJson} from '@/components/changeDataStructure/chDS';
export default Vue.extend({
    name: 'ArrayView',
    props: {
        parsedData: {},
    },
    data() {
        return {
            flowData: (this as any).parsedData.childParams,
            dialog: false,
            objc: [],
            textarea: '',
            cantRemove: true,
            schemaObject: {},
            temp: {},
            nameTemp: '',
        };
    },
    watch: {
        parsedData: {
            handler(newValue, oldValue) {
                this.flowData = (this as any).parsedData.childParams;
            },
            immediate: true,
        },
        flowData: {
            handler(newValue, oldValue) {
                if (newValue.length > 1) {
                    this.cantRemove = false;
                }
                if (newValue.length === 1) {
                    this.cantRemove = true;
                }
                this.$emit('input', newValue);
            },
            deep: true,
        },
    },
    methods: {
        ShowJsonDialog(obj: any) {
            this.dialog = true;
            (this as any).objc.push(obj);
            const jsonSchemaGenerator = require('./schemaGenerator/index.js');
            const data = {};
            this.temp = makeJson(this.objc);
            this.nameTemp = obj.name;
            eval('data = this.temp.' + this.nameTemp);
            this.schemaObject = jsonSchemaGenerator.jsonToSchema(this.temp);
            this.textarea = JSON.stringify(data, null, 4);
        },
        handleClose(done: any) {
            try {
                if (this.textarea) {
                    const data = JSON.parse(this.textarea);
                    eval('this.temp.' + this.nameTemp + '= data');
                    const tv4 = require('tv4');
                    if (tv4.validate(this.temp, this.schemaObject)) {
                        const objct: any = jsonParse(data);
                        this.$emit('input', objct);
                    } else {
                        (this as any).$message.error('json text input error!');
                    }
                } else {
                    (this as any).$message.error('json text input error!');
                }
            } catch (e) {
                (this as any).$message.error('json text input error!');
            } finally {
                done();
            }
        },
        addEvent() {
            this.flowData.push(this.shallowCopy(this.flowData[this.flowData.length - 1]));
        },
        shallowCopy(src: any): any {
            const dst: any = {};
            for (const prop in src) {
                if (src.hasOwnProperty(prop)) {
                    dst[prop] = src[prop];
                }
            }
            return dst;
        },
        removeEvent(index: number) {
            this.flowData.splice(index, 1);
        },
    },
});
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
