<template>
    <div style="line-height: 50px">
        <el-container id="importDot" style="margin-top: -12px;">
            <el-header id="dotH" style="line-height: 30px;height: 30px;text-align: left;background-color: #D3DCE6;">
                Dot
            </el-header>
            <el-main id="dotM">
                <div>
                    <span>dotPath :</span>
                    <el-input type="text" v-model="dotPath" style="width: 60%;margin-left: 2%;"></el-input>
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
                    <el-input type="text" v-model="confPath" style="width: 60%;margin-left: 2%;"></el-input>
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
    import {checkType} from "../components/changeDataStructure/checkType";

    export default {
        data() {
            return {
                dotPath: '',
                confPath: '',

            }
        },
        methods: {
            importDot() {
                if (this.dotPath != '') {
                    var {rpcimportByDot} = require('../plugins/rpcInterface');
                    rpcimportByDot(this.dotPath, (response) => {
                        if (response.getError() != '') {
                            var err = response.getError();
                            console.log(err);
                            this.$message({
                                type:'warning',
                                 message:err
                    })
                        } else {
                            var dots=[];
                            dots= JSON.parse(response.getJson());
                            console.log(dots);
                            //todo 判断类型是数组而且含有metaData.typeId字段
                            if  (!(Array.isArray(dots) && dots[0].metaData.typeId!=null)){
                                alert("请选择正确的组件文件．");
                                return false
                            }
                            for(var i=0;i<dots.length;i++){
                                var bo=true;
                                for (var j = 0, len = this.$root.Dots.length; j < len; j++) {
                                    if(dots[i].metaData.typeId == this.$root.Dots[j].metaData.typeId){
                                        bo=false;
                                        break
                                    }
                                }
                                if (bo) {
                                    this.$root.Dots.push(dots[i]);
                                    this.$root.DotsTem.push(JSON.parse(JSON.stringify(dots[i])));
                                    this.$root.ExportDots.push(JSON.parse(JSON.stringify(dots[i])))
                                }

                            }
                            checkType(this.$root.Dots,this.$root.Configs);
                            this.$message({
                                type:"success",
                                message:"Import Dot success!"
                            });
                        }
                    });
                } else {
                    this.$message({
                        type:'warning',
                        message:'Please Imput DotPath!'
                    })
                }
            },
            importConf() {
                if (this.confPath != '') {
                    let {rpcimportByConfig} = require('../plugins/rpcInterface');
                    rpcimportByConfig(this.confPath, (response) => {
                        if (response.getError() != '') {
                            let err = response.getError()
                            console.log(err)
                            this.$message({
                                message: err,
                                type: 'error'
                            })
                        } else {
                            let config;
                            {
                                let ob = JSON.parse(response.getJson());
                                if(ob.dots && Object.prototype.toString.call(ob.dots) === '[object Array]'){
                                    config = ob.dots;
                                }else {
                                    this.$message({
                                        message: 'Dots is not exist in config files!',
                                        type: 'error'
                                    }) ;
                                    return;
                                }
                            }
                            for (let i = 0, len = config.length; i < len; i++) {
                                if(config[i].metaData.typeId){
                                    let dot = this.findDot(config[i].metaData.typeId);
                                    this.assembleByTypeId(dot, config[i]);
                                    this.$message({
                                        message: 'Import success!',
                                        type: 'success'
                                    });
                                }else {
                                    this.$message({
                                        message: 'typeId is not exist in config files!',
                                        type: 'error'
                                    }) ;
                                    return;
                                }
                            }
                        }
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
                for (let i = 0, len = this.$root.DotsTem.length; i < len; i++) {
                    if (this.$root.DotsTem[i].metaData.typeId === typeId) {
                        return this.$root.DotsTem[i];
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
                for (let i = 0, len = this.$root.Configs.length; i < len; i++) {
                    if (this.$root.Configs[i].metaData.typeId === config.metaData.typeId) {
                        this.$root.Configs[i] = this.assemble(this.$root.Configs[i], config);
                        flag = true;
                    }
                }
                if (!flag) {
                    if (dot) {
                        this.$root.Configs.push(this.assemble(dot, config));
                    } else {
                        config.metaData.flag = 'not-exist';
                        this.$root.Configs.push(config);
                    }
                }
            }
        }
    }
</script>
<style>

</style>
