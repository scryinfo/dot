<template>
    <el-collapse>
        <el-button id="updateLiveNum" @click="updateLvNum()" style="margin-bottom: 5px;">Live Num</el-button>
        <el-button id="removeAllDot" @click="removeAllDots()" style="margin-bottom: 5px;">Remove All</el-button>
        <el-row v-for="(v,index) in this.$store.state.Dots">
            <el-col :span="3">
                <div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{livesNum[index]}}
                </div>
            </el-col>
            <el-col :span="3">
                <div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{v.metaData.name}}
                </div>
            </el-col>
            <el-col :span="10">
                <el-collapse-item v-bind:title="v.metaData.typeId" v-bind:name="index">
                    <el-row v-for="(a,b,c) in v.metaData">
                        <el-col :span="6">
                            <div>{{b}}</div>
                        </el-col>
                        <el-col :span="18">
                            <div>{{a}}</div>
                        </el-col>
                    </el-row>
                </el-collapse-item>
            </el-col>
            <el-col :span="2">
                <div style="margin-left: 6px">
                    <el-button v-on:click="open(index)">Name</el-button>
                </div>
            </el-col>
            <el-col :span="3">
                <div>
                    <el-button v-on:click="addConf(index)">Add Config</el-button>
                </div>
            </el-col>
            <el-col :span="3">
                <el-button id="delDot" @click="delDot(index)" style="margin-left: 2%">remove</el-button>
            </el-col>
        </el-row>
    </el-collapse>
</template>
<script>
    import {checkType, removeAllType} from "../components/utils/checkType";

    export default {
        data() {
            return {
                livesNum: [],
                table: this.store.state.Dots,
                activeNames: ['1']
            };
        },
        methods: {
            updateLvNum() {
                let l = this.store.state.Dots.length;
                for (let i = 0; i < l; i++) {
                    this.livesNum[i] = 0;
                    let tId = this.store.state.Dots[i].metaData.typeId;
                    for (let j = 0; j < this.store.state.Configs.length; j++) {
                        if (this.store.state.Configs[j].metaData.typeId === tId) {
                            let livesLen = this.store.state.Configs[j].lives.length;
                            for (let k = 0; k < livesLen; k++) {
                                if (this.store.state.Configs[j].lives[k].liveId !== '') {
                                    this.livesNum[i]++;
                                }
                            }
                            break;
                        }
                    }
                }
                this.$forceUpdate();
                this.$message({
                    type: "success",
                    message: "updata LivesNumber Success",
                });
            },
            open(index) {
                this.$prompt('输入组件名', 'Dot name', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                }).then(({value}) => {
                    this.$message({
                        type: 'success',
                        message: '该组件名更改为: ' + value,
                    });
                    this.store.state.Dots[index].metaData.name = value;
                    let id = this.store.state.Dots[index].metaData.typeId;
                    for (let i = 0; i < this.store.state.Configs.length; i++) {
                        if (this.store.state.Configs[i].metaData.typeId === id) {
                            this.store.state.Configs[i].metaData.name = value;
                            break;
                        }
                    }
                    ;
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '取消输入'
                    });
                });
            },
            addConf(index) {
                let tId = this.store.state.Dots[index].metaData.typeId;
                let confNum = this.store.state.Configs.length;
                for (let i = 0; i < confNum; i++) {
                    if (this.store.state.Configs[i].metaData.typeId === tId) {
                        this.$message({
                            type: "warning",
                            message: "Add default,Exist!",
                        });
                        return
                    }
                }
                this.store.state.Configs.push(JSON.parse(JSON.stringify(this.store.state.Dots[index])));
                this.$message({
                    type: "success",
                    message: "Add Config success",
                });
            },
            removeAllDots() {
                this.store.state.Dots = [];
                removeAllType(this.store.state.Configs);
            },
            delDot(index) {
                this.store.state.Dots.splice(index, 1);
                checkType(this.store.state.Dots, this.store.state.Configs);
            },
        },
    }

</script>

<style scoped>
    .bg-purple {
        background: #d3dce6;
    }

    .grid-content {
        border-radius: 4px;
        min-height: 46px;
    }
</style>
