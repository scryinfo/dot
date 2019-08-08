<template>
    <div>
        <el-row v-for="(member,index) in flowData">
            <el-col :span="3" v-if="member.type == 'object' || member.type == 'array'"><span  style="text-align: center;line-height: 46px;">{{member.name}}</span></el-col>
            <el-col :span="21">
                <div v-if="member.type !== 'object' && member.type !== 'array'">
                    <el-input type="text"
                              v-model="flowData[index].remark"
                              v-if="member.type == 'string'">
                        <template slot="prepend">{{member.name}}</template>
                    </el-input>
                    <el-input
                            type="number"
                            v-model.number="flowData[index].remark"
                            v-if="member.type == 'number'">
                        <template slot="prepend">{{member.name}}</template>
                    </el-input>
                    <bool-view
                            v-model="flowData[index].remark"
                            :boolValue="flowData[index].remark"
                            v-if="member.type == 'boolean'"
                            :name="member.name"
                    >
                    </bool-view>
                </div>
                <div v-if="member.type == 'object'">
                    <json-view v-model="flowData[index].childParams" :parsedData="flowData[index].childParams"></json-view>
                </div>
                <div v-if="member.type == 'array'">
                    <array-view v-model="flowData[index].childParams" :parsedData="flowData[index]"></array-view>
                </div>
            </el-col>
        </el-row>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    export default Vue.extend({
        name: 'JsonView',
        props: {
            parsedData: {},
        },
        data () {
            return {
                flowData: this.parsedData,
            };
        },
        watch: {
            parsedData: {
                handler(newValue, oldValue) {
                    this.flowData = this.parsedData;
                },
                immediate: true
            },
            flowData: {
                handler(newValue, oldValue) {
                    this.$emit('input', newValue);
                },
                deep: true
            }
        }
    })
</script>

<style scoped>

</style>
