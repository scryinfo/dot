<template>
    <div>
        <el-container>
            <el-header>
                Dot
            </el-header>
            <el-main>
                <div>
                    <span>dotPath :</span>
                    <el-input type="text" v-model="dotPath"/>
                    <el-button @click="delDotPath()">Clear</el-button>
                </div>
                <el-row>
                    <el-button @click="importDot()">Import Dot
                    </el-button>

                </el-row>
            </el-main>
        </el-container>
        <el-container>
            <el-header>
                Config
            </el-header>
            <el-main>
                <div>
                    <span>confPath:</span>
                    <el-input type="text" v-model="confPath"/>
                    <el-button @click="delConfPath()">Clear</el-button>
                    <!--          <el-button id="removeC" @click="del(index)" style="margin-left: 2%">remove</el-button>-->
                </div>
                <el-row>
                    <el-button @click="importConf()">
                        ImportConfig
                    </el-button>
                </el-row>
            </el-main>
        </el-container>
    </div>


</template>

<script lang="ts">
    import {checkType} from "../components/utils/checkType";
    import DotWrapper from "@/rpc_face/data_wrapper";
    import {Component, Vue} from "vue-property-decorator";
    import {Dot, state} from '@/views/home/store';

    @Component
    export default class Import extends Vue {

        private dotPath = '';
        private confPath = '';

        private importDot() {
            if (this.dotPath !== '') {
                DotWrapper.importByDot(this.dotPath).then(data => {
                    if (data.getError() !== '') {
                        let err = data.getError();
                        this.$message({
                            type: 'warning',
                            message: err
                        })
                    } else {
                        let dots = [];
                        dots = JSON.parse(data.getJson());
                        //todo 判断类型是数组而且含有metaData.typeId字段
                        if (!(Array.isArray(dots) && dots[0].metaData.typeId != null)) {
                            alert("请选择正确的组件文件．");
                            return false
                        }
                        for (let i = 0; i < dots.length; i++) {
                            let bo = true;
                            for (let j = 0, len = state.Dots.length; j < len; j++) {
                                if (dots[i].metaData.typeId === state.Dots[j].metaData.typeId) {
                                    bo = false;
                                    break
                                }
                            }
                            if (bo) {
                                state.Dots.push(dots[i]);
                            }

                        }
                        checkType(state.Dots, state.Configs);
                        this.$message({
                            type: "success",
                            message: "Import Dot success!"
                        });
                    }
                }).catch(err => {
                    this.$message({
                        type: 'warning',
                        message: err.toString()
                    })
                });
            } else {
                this.$message({
                    type: 'warning',
                    message: 'Please Input DotPath!'
                })
            }
        }

        private importConf() {
            if (this.confPath !== '') {
                DotWrapper.importByConfig(this.confPath).then(data => {
                    if (data.getError() !== '') {
                        let err = data.getError();
                        this.$message({
                            message: err,
                            type: 'error'
                        })
                    } else {
                        let config;
                        {
                            let ob = JSON.parse(data.getJson());
                            if (ob.dots && Object.prototype.toString.call(ob.dots) === '[object Array]') {
                                config = ob.dots;
                            } else {
                                this.$message({
                                    message: 'Dots is not exist in config files!',
                                    type: 'error'
                                });
                                return;
                            }
                        }
                        for (let i = 0, len = config.length; i < len; i++) {
                            if (config[i].metaData.typeId) {
                                let dot = this.findDot(config[i].metaData.typeId);
                                this.assembleByTypeId(dot, config[i]);
                                this.$message({
                                    message: 'Import success!',
                                    type: 'success'
                                });
                            } else {
                                this.$message({
                                    message: 'typeId is not exist in config files!',
                                    type: 'error'
                                });
                                return;
                            }
                        }
                    }
                }).catch(err => {
                    this.$message({
                        message: err.toString(),
                        type: 'error'
                    })
                });
            } else {
                this.$message({
                    message: 'Please input confPath!',
                    type: 'warning'
                })
            }
        }

        delDotPath() {
            this.dotPath = '';
        }

        delConfPath() {
            this.confPath = '';
        }

        findDot(typeId: string) {
            for (let i = 0, len = state.Dots.length; i < len; i++) {
                if (state.Dots[i].metaData.typeId === typeId) {
                    return JSON.parse(JSON.stringify(state.Dots[i]));
                }
            }
            return null;
        }

        assemble(schema: any, source: any) {
            for (let key in schema) {
                schema[key] = this.isObject(schema[key]) ? this.assemble(schema[key], (source[key] ? source[key] : schema[key])) : (source[key] ? source[key] : schema[key]);
            }
            return schema;
        }

        isObject(o: any) {
            return (typeof o === 'object') && o !== null;
        }

        assembleByTypeId(dot: Dot, config: Dot) {
            let flag = false;
            for (let i = 0, len = state.Configs.length; i < len; i++) {
                if (state.Configs[i].metaData.typeId === config.metaData.typeId) {
                    for (let key in config.lives) {
                        if (!this.equalLiveId(state.Configs[i].lives, config.lives[key])) {
                            state.Configs[i].lives.push(config.lives[key]);
                        }
                    }
                    for (let j = 0, l = state.Configs[i].lives.length; j < l; j++) {
                        if (state.Configs[i].lives[j].liveId === "") {
                            state.Configs[i].lives.splice(j, 1);
                        }
                    }
                    flag = true;
                }
            }
            if (!flag) {
                if (dot) {
                    let dotCopy = JSON.parse(JSON.stringify(dot));
                    dotCopy.lives.length = 0;
                    for (let key in config.lives) {
                        dotCopy.lives.push(this.assemble(JSON.parse(JSON.stringify(dot.lives[0])), config.lives[key]));
                    }
                    state.Configs.push(dotCopy);
                } else {
                    config.metaData.flag = 'not-exist';
                    state.Configs.push(config);
                }
            }
        }

        equalLiveId(confLives: any, live: any) {
            for (let key in confLives) {
                if (confLives[key].liveId === live.liveId) {
                    this.assemble(confLives, live);
                    return true;
                }
            }
            return false;
        }

    }
</script>
<style scoped>

</style>
