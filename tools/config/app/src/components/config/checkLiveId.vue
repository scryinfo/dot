<template>
    <div>
        <div v-if="checkResult">
            <el-col :span="2">
                <el-tooltip effect="dark" content="Please perfect liveId in this typeId!" placement="bottom-start">
                    <div class="grid-content bg-Danger" style="text-align: center;line-height: 46px;">{{metaName}}</div>
                </el-tooltip>
            </el-col>
        </div>
        <div v-else>
            <el-col :span="2">
                <div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{metaName}}</div>
            </el-col>
        </div>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue'

    export default Vue.extend({
        name: "checkLiveId",
        props: {
            metaName: {
                type: String,
                required: true
            },
            lives: {
                type: Array,
                required: true
            }
        },
        watch: {
            lives: {
                handler(newValue, oldValue) {
                    let flag = false;
                    for (let i = 0, len = newValue.length; i < len; i++) {
                        if (newValue[i].liveId === '') {
                            this.checkResult = true;
                            flag = true;
                            break;
                        }
                    }
                    if (!flag) {
                        this.checkResult = false;
                    }
                },
                immediate: true,
                deep: true
            }
        },
        data() {
            return {
                checkResult: false
            }
        }
    })
</script>

<style scoped>
    .bg-purple {
        background: #d3dce6;
    }

    .grid-content {
        border-radius: 4px;
        min-height: 46px;
    }

    .bg-Danger {
        background: #F56C6C;
    }
</style>
