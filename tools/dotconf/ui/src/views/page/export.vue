<template>
    <div>
        <el-container>
            <el-header>Dot</el-header>
            <el-main>
                <div>
                    <span>fileName :</span>
                    <el-input type="text" v-model="dotFileName"/>
                </div>
                <el-button @click="ExportDot">Export Dot</el-button>
            </el-main>
        </el-container>
        <el-container>
            <el-header>Config</el-header>
            <el-main>
                <div>
                    <span>fileName :</span>
                    <el-input type="text" v-model="confFileName"/>
                </div>
                <el-row>
                    <template>
                        <el-checkbox
                                :indeterminate="isIndeterminate"
                                v-model="checkAllC"
                                @change="handleCheckAllChangeC"
                        >全选
                        </el-checkbox>
                        <el-checkbox-group v-model="checkedFormatC" @change="handleCheckedCitiesChangeC">
                            <el-checkbox v-for="city in optionsC" :label="city" :key="city">{{city}}</el-checkbox>
                        </el-checkbox-group>
                    </template>
                </el-row>
                <el-button @click="ExportConf">ExportConfig</el-button>
            </el-main>
        </el-container>
    </div>
</template>

<script lang="ts">
    import DotWrapper from "@/rpc_face/data_wrapper";
    import {Component, Vue} from "vue-property-decorator";
    import {Dot, state} from "@/views/home/store";

    @Component
    export default class Export extends Vue {
        private optionsC = ["json", "toml", "yaml"];
        private dotFileName = "";
        private confFileName = "";
        private checkAllC = false;
        private checkedFormatC = ["json"];
        private isIndeterminate = true;

        private handleCheckAllChangeC(val: []) {
            this.checkedFormatC = val ? this.optionsC : [];
            this.isIndeterminate = false;
        }

        private handleCheckedCitiesChangeC(value: []) {
            let checkedCount = value.length;
            this.checkAllC = checkedCount === this.optionsC.length;
            this.isIndeterminate =
                checkedCount > 0 && checkedCount < this.optionsC.length;
        }

        private ExportDot() {
            if (this.dotFileName === "") {
                alert("请输入文件名");
            } else {
                let filename = [this.dotFileName + ".json"];
                DotWrapper.exportDot(JSON.stringify(state.Dots), filename).then(data => {
                    if (!data.getError()) {
                        alert("导出文件成功，文件位置tools/dotconf/run_out目录下");
                    } else {
                        alert("导出文件失败" + data.getError());
                    }
                }).catch(err => {
                    alert("导出文件失败" + err);
                });
            }
        }

        private ExportConf() {
            //console.log(JSON.stringify(state.Configs, null, 4));
            if (this.confFileName === "") {
                alert("请输入文件名");
            } else {
                let confFileList = [];
                //要生成的文件名
                for (let i = 0; i < this.checkedFormatC.length; i++) {
                    confFileList.push(this.confFileName + "." + this.checkedFormatC[i]);
                }
                //判断liveid
                let conf = state.Configs; //config页面数据
                let resultDot = []; //处理掉空配置
                {
                    let liveIds = [];
                    for (let i = 0; i < conf.length; i++) {
                        if (conf[i].lives.length === 0) {
                            //实例数为０跳过
                            continue;
                        }
                        for (let j = 0; j < conf[i].lives.length; j++) {
                            if (conf[i].lives[j].liveId === "") {
                                alert(conf[i].lives[j] + ":liveId is null");
                                return false;
                            } else {
                                for (let z = 0; z < liveIds.length; z++) {
                                    if (conf[i].lives[j].liveId === liveIds[z]) {
                                        alert(conf[i].lives[j].liveId + "liveid重复．");
                                        return false;
                                    }
                                }
                                liveIds.push(conf[i].lives[j].liveId);
                            }
                        }
                        resultDot.push(JSON.parse(JSON.stringify(conf[i])));
                    }
                }
                if (!this.configRequire(resultDot)) {
                    return false
                }
                let result = {
                    log: {
                        file: "log.log",
                        level: "debug"
                    },
                    dots: {},
                };
                result.dots = resultDot;
                DotWrapper.exportConfig(JSON.stringify(result), confFileList).then(data => {
                    if (!data.getError()) {
                        alert("导出文件成功，文件位置tools/dotconf/run_out目录下");
                    } else {
                        alert("导出文件失败" + data.getError());
                    }
                }).catch(err => {
                    alert("导出文件失败" + err);
                });
            }
        }

        private configRequire(configs: Dot[]) {
            for (let key in configs) {
                let typeId = configs[key].metaData.typeId;
                if (!configs[key].requiredInfo) {
                    continue
                }
                let require = configs[key].requiredInfo;
                for (let i in configs[key].lives) {
                    let config = configs[key].lives[i].json;
                    if (this.configConfirm(typeId, config, require, "lives[" + i + "]")) {
                        this.$delete(configs[key], 'requiredInfo');
                    } else {
                        return false;
                    }
                }
            }
            return true;
        }

        private configConfirm(typeId: string, config: any, require: any, index: any) {
            for (let key in require) {
                if (require[key] === true) {
                    if (config[key] === null || config[key] === '' || config[key].length === 0 || this.baseArrayConfirm(config[key])) {
                        alert(index + key + " in extend config of the typeID: " + typeId + " is not exist");
                        return false;
                    }
                } else if (require[key] === false) {

                } else {
                    index = index + key + '.';
                    if (Array.isArray(config[key])) {
                        for (let i in config[key]) {
                            if (!this.configConfirm(typeId, config[key][i], require[key], index)) {
                                return false;
                            }
                        }
                    } else {
                        if (!this.configConfirm(typeId, config[key], require[key], index)) {
                            return false;
                        }
                    }
                }
            }
            return true;
        }

        private baseArrayConfirm(arr: any) {
            for (let i in arr) {
                if (arr[i] === '') {
                    return true;
                }
            }
            return false;
        }

    }
</script>
<style scoped>
    #eDotH,
    #eConfH {
        text-align: left;
        background-color: #d3dce6;
    }
</style>
