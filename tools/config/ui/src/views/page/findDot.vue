<template>
    <div id="findDot">
        <el-row>
            <el-button id="add" @click="add()">Add</el-button>
            <el-button id="removeAll" @click="removeAll()">Remove All</el-button>
        </el-row>
        <div v-for="(item,index) in files">
            <span>filePath{{index+1}}:</span>
            <el-input type="text" v-model="files[index]" style="width: 78%;margin-left: 2%;"/>
            <el-button id="remove" @click="del(index)" style="margin-left: 2%">remove</el-button>
        </div>
        <el-button id="find" v-loading="fullscreenLoading" @click="find()" style="margin-right:78%;">FindDot</el-button>
        <span style="margin-right: 70%">{{notExistT}}</span>
    </div>
</template>

<script>
    import {checkType} from "../components/utils/checkType";
    import DotWrapper from "@/rpc_face/data_wrapper";

    export default {
        data() {
            return {
                files: [''],
                fullscreenLoading: false,
                notExistT: ''
            }
        },
        methods: {
            add() {
                let end = this.files.length - 1;
                if (this.files[end] === '') {
                    this.$message({
                        type: 'error',
                        message: 'the end filepath is empty!'
                    })
                } else {
                    this.files.push('')
                }
            },
            removeAll() {
                this.files = ['']
            },

            del(index) {
                this.files.splice(index, 1);

            },
            find() {
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
                                for (let j = 0, len = this.store.state.Dots.length; j < len; j++) {
                                    if (res[i].metaData.typeId === this.store.state.Dots[j].metaData.typeId) {
                                        bo = false;
                                        break
                                    }
                                }
                                if (bo) {
                                    this.store.state.Dots.push(res[i]);

                                }
                            }
                            checkType(this.store.state.Dots, this.store.state.Configs);
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
    }
</script>
<style scoped>
    #findDot {
        text-align: right;
        line-height: 50px;
    }
</style>
