<template>
    <el-collapse>
        <el-button @click="updateLvNum()">Live Num</el-button>
        <el-button @click="removeAllDots()">Remove All</el-button>
        <el-row v-for="(v,index) in this.dots" :key="index">
            <el-col :span="1">
                <div class="grid-content bg-purple el-dialog--center">{{livesNum[index]}}</div>
            </el-col>
            <el-col :span="3">
                <div class="grid-content bg-purple ">{{v.metaData.name}}</div>
            </el-col>
            <el-col :span="10">
                <el-collapse-item :title="v.metaData.typeId" :name="index">
                    <el-row v-for="(a,b) in v.metaData" :key="a">
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
                <el-button @click="open(index)">Name</el-button>
            </el-col>
            <el-col :span="3">
                <el-button @click="addConf(index)">Add Config</el-button>
            </el-col>
            <el-col :span="3">
                <el-button @click="delDot(index)">remove</el-button>
            </el-col>
        </el-row>
    </el-collapse>
</template>
<script lang="ts">
    import {checkType, removeAllType} from "../components/utils/checkType";
    import {Dot, state} from "@/views/home/store";
    import {Component, Vue} from "vue-property-decorator";

    @Component
    export default class DotList extends Vue {

        private livesNum = new Array<number>();
        private activeNames = ['1'];
        private dots = new Array<Dot>();

        private updateLvNum() {
            let l = this.dots.length;
            for (let i = 0; i < l; i++) {
                this.livesNum[i] = 0;
                let tId = this.dots[i].metaData.typeId;
                for (let j = 0; j < state.Configs.length; j++) {
                    if (state.Configs[j].metaData.typeId === tId) {
                        let livesLen = state.Configs[j].lives.length;
                        for (let k = 0; k < livesLen; k++) {
                            if (state.Configs[j].lives[k].liveId !== '') {
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
                message: "update LivesNumber Success",
            });
        }

        private open(index: number) {
            this.$prompt('输入组件名', 'Dot name', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
            }).then((data:any) => {
                this.$message({
                    type: 'success',
                    message: '该组件名更改为: ' + data.value,
                });
                this.dots[index].metaData.name = data.value;
                let id = this.dots[index].metaData.typeId;
                for (let i = 0; i < state.Configs.length; i++) {
                    if (state.Configs[i].metaData.typeId === id) {
                        state.Configs[i].metaData.name = data.value;
                        break;
                    }
                }
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '取消输入'
                });
            });
        }

        private addConf(index: number) {
            let tId = this.dots[index].metaData.typeId;
            let confNum = state.Configs.length;
            for (let i = 0; i < confNum; i++) {
                if (state.Configs[i].metaData.typeId === tId) {
                    this.$message({
                        type: "warning",
                        message: "Add default,Exist!",
                    });
                    return
                }
            }
            state.Configs.push(this.dots[index]);
            this.$message({
                type: "success",
                message: "Add Config success",
            });
        }

        private removeAllDots() {
            this.dots = [];
            removeAllType(state.Configs);
        }

        private delDot(index: number) {
            this.dots.splice(index, 1);
            checkType(this.dots, state.Configs);
        }

        private mounted() {
            this.dots = state.Dots;
            // console.log(this.dots);
            // console.log(state.Configs);
        }
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
