<template>
    <div style="line-height: 50px">
        <el-container id="importDot" style="margin-top: -12px;">
            <el-header id="dotH" style="line-height: 30px;height: 30px;text-align: left;background-color: #D3DCE6;">
                Dot
            </el-header>
            <el-main id="dotM">
                <div>
                    <span>dotPath :</span>
                    <el-input type="text" v-model="dotPath" style="width: 60%;margin-left: 2%;"/>
                    <el-button id="removeD" @click="delDotPath()" style="margin-left: 2%">Clear</el-button>
                </div>
                <el-row>
                    <el-button id="findD" @click="importDot()" style="margin-right:30%;margin-left: 10%;">Import Dot
                    </el-button>

                </el-row>
            </el-main>
        </el-container>
        <el-container id="importConf">
            <el-header id="confH" style="height: 30px;line-height: 30px; text-align: left;background-color: #D3DCE6;">
                Config
            </el-header>
            <el-main id="confM">
                <div>
                    <span>confPath:</span>
                    <el-input type="text" v-model="confPath" style="width: 60%;margin-left: 2%;"/>
                    <el-button id="removeC" @click="delConfPath()" style="margin-left: 2%">Clear</el-button>
                    <!--          <el-button id="removeC" @click="del(index)" style="margin-left: 2%">remove</el-button>-->

                </div>
                <el-row>
                    <el-button id="findC" @click="importConf()" style="margin-right:29%;margin-left: 10%;">
                        ImportConfig
                    </el-button>
                </el-row>
            </el-main>
        </el-container>
    </div>


</template>

<script>
    import {checkType} from "../components/utils/checkType";
    import DotWrapper from "@/rpc_face/data_wrapper";

    export default {
        data() {
            return {
                dotPath: '',
                confPath: '',

            }
        },
        methods: {
            importDot() {
                if (this.dotPath !== '') {
                    DotWrapper.importByDot(this.dotPath).then(data => {
                        if (data.getError() !== '') {
                            let err = data.getError();
                            console.log(err);
                            this.$message({
                                type: 'warning',
                                message: err
                            })
                        } else {
                            let dots = [];
                            dots = JSON.parse(data.getJson());
                            console.log(dots);
                            //todo 判断类型是数组而且含有metaData.typeId字段
                            if (!(Array.isArray(dots) && dots[0].metaData.typeId != null)) {
                                alert("请选择正确的组件文件．");
                                return false
                            }
                            for (let i = 0; i < dots.length; i++) {
                                let bo = true;
                                for (let j = 0, len = this.store.state.Dots.length; j < len; j++) {
                                    if (dots[i].metaData.typeId === this.store.state.Dots[j].metaData.typeId) {
                                        bo = false;
                                        break
                                    }
                                }
                                if (bo) {
                                    this.store.state.Dots.push(dots[i]);
                                }

                            }
                            checkType(this.store.state.Dots, this.store.state.Configs);
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
                        message: 'Please Imput DotPath!'
                    })
                }
            },
            importConf() {
                if (this.confPath !== '') {
                    DotWrapper.importByConfig(this.confPath).then(data => {
                        if (data.getError() !== '') {
                            let err = data.getError();
                            console.log(err);
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


            },
            delDotPath() {
                this.dotPath = '';
            },
            delConfPath() {
                this.confPath = '';
            },
            findDot(typeId) {
                for (let i = 0, len = this.store.state.Dots.length; i < len; i++) {
                    if (this.store.state.Dots[i].metaData.typeId === typeId) {
                        return JSON.parse(JSON.stringify(this.store.state.Dots[i]));
                    }
                }
                return null;
            },
            assemble(schema, source) {
                for (let key in schema) {
                    schema[key] = this.isObject(schema[key]) ? this.assemble(schema[key], (source[key] ? source[key] : schema[key])) : (source[key] ? source[key] : schema[key]);
                }
                return schema;
            },
            isObject(o) {
                return (typeof o === 'object') && o !== null;
            },
            assembleByTypeId(dot, config) {
                let flag = false;
                for (let i = 0, len = this.store.state.Configs.length; i < len; i++) {
                    if (this.store.state.Configs[i].metaData.typeId === config.metaData.typeId) {
                        for (let key in config.lives) {
                            if (this.equalLiveId(this.store.state.Configs[i].lives, config.lives[key])) {

                            } else {
                                this.store.state.Configs[i].lives.push(config.lives[key]);
                            }
                        }
                        for (let j = 0, l = this.store.state.Configs[i].lives.length; j < l; j++) {
                            if (this.store.state.Configs[i].lives[j].liveId === "") {
                                this.store.state.Configs[i].lives.splice(j, 1);
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
                        this.store.state.Configs.push(dotCopy);
                    } else {
                        config.metaData.flag = 'not-exist';
                        this.store.state.Configs.push(config);
                    }
                }
            },
            equalLiveId(confLives, live) {
                for (let key in confLives) {
                    if (confLives[key].liveId === live.liveId) {
                        this.assemble(confLives, live);
                        return true;
                    }
                }
                return false;
            }
        }
    }
</script>
<style scoped>

</style>
