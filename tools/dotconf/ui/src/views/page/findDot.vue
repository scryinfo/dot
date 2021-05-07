<template>
    <div class="find" v-loading.fullscreen.lock="fullscreenLoading">
        <div class="el-row--flex">
            <el-tooltip content="支持多目录查找">
                <el-button>通过代码目录查找组件</el-button>
            </el-tooltip>
            <el-button @click="add()">Add</el-button>
            <el-button @click="removeAll()">Remove All</el-button>
        </div>

        <div v-for="(item,index) in files" class="el-row--flex dotpath">
            <el-tag>filePath{{index+1}}:</el-tag>
            <el-input type="text" v-model="files[index]"/>
            <el-button id="remove" @click="del(index)">remove</el-button>
        </div>
        <el-button  @click="find()">FindDot</el-button>
        <span>{{notExistT}}</span>
    </div>
</template>

<script lang="ts">
    import {checkType} from "../components/utils/checkType";
    import DotWrapper from "@/rpc_face/data_wrapper";
    import {Component, Vue} from "vue-property-decorator";
    import {state} from "@/views/home/store";

    @Component
    export default class FindDot extends Vue {

        private files = [''];
        private fullscreenLoading = false;
        private notExistT = '';

        private add() {
            let end = this.files.length - 1;
            if (this.files[end] === '') {
                this.$message({
                    type: 'error',
                    message: 'the end filepath is empty!'
                })
            } else {
                this.files.push('');
            }
        }

        private removeAll() {
            this.files = ['']
        }


        private del(index: number) {
            this.files.splice(index, 1);
        }

        private find() {
            let dir = this.files;
            if (dir[0] === '') {
                this.$message({
                    type: 'error',
                    message: 'Please Input FilePath!'
                })
            } else {
                this.fullscreenLoading = true;
                DotWrapper.findDot(dir).then(data => {
                    if (data.getError()) {
                        this.$message({
                            type: 'error',
                            message: data.getError(),
                        });
                    } else {
                        let res = JSON.parse(data.getDotsinfo());
                        for (let i = 0; i < res.length; i++) {
                            let bo = true;
                            for (let j = 0, len = state.Dots.length; j < len; j++) {
                                if (res[i].metaData.typeId === state.Dots[j].metaData.typeId) {
                                    bo = false;
                                    break
                                }
                            }
                            if (bo) {
                                state.Dots.push(res[i]);
                            }
                        }
                        checkType(state.Dots, state.Configs);
                        this.$message({
                            type: 'success',
                            message: 'find dot finish!'
                        });
                    }
                }).catch(err => {
                    this.$message({
                        type: 'error',
                        message: err.getError(),
                    });
                });
                this.fullscreenLoading = false;
            }
        }
    }

</script>
<style scoped lang="scss">
    .find {

        .dotpath {
            margin-top: 16px;
        }

    }
</style>
