<template>
    <div class="fold">
        <el-row>
            <el-col :span="3"><el-button>Check Type</el-button></el-col>
        </el-row>
        <el-collapse v-model="activeNames" @change="handleChange">
            <el-row>
                <el-col :span="20">
                    <el-collapse-item title="1      Name TypeId" name="1">
                        <el-row>
                            <el-col :span="20">
                                <el-collapse-item title="Name Live id" name="1.1">
                                </el-collapse-item>
                            </el-col>
                            <el-col :span="4">
                                <el-button>JSON</el-button>
                                <el-button>Remove Live</el-button>
                            </el-col>
                        </el-row>
                        <el-row>
                            <el-col :span="22">
                                <el-collapse-item title="Extend Config for live" name="1.2">
                                </el-collapse-item>
                            </el-col>
                            <el-col :span="2">
                                <el-button>JSON</el-button>
                            </el-col>
                        </el-row>
                    </el-collapse-item>
                </el-col>
                <el-col :span="4">
                    <el-button>Load By Config?</el-button>
                    <el-button>Add Live</el-button>
                </el-col>
            </el-row>
        </el-collapse>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue'
    export default Vue.extend({
        data () {
            return {
                activeNames: []
            }
        },
        methods: {
            handleChange (val:any) {
                console.log(val)
            }
        }
    })
</script>

<style scoped>

</style>
