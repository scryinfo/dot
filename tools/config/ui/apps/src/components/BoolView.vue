<template>
<div>
    {{name}}
    <el-radio v-model="radio" label='true'>true</el-radio>
    <el-radio v-model="radio" label='false'>false</el-radio>
</div>
</template>

<script lang="ts">
    import Vue from 'vue'
    export default Vue.extend({
        name: "BoolView",
        props: {
            boolValue: {
                type: Boolean,
                required: true
            },
            name: {}
        },
        data() {
            return {
                radio: this.boolValue
            }
        },
        watch: {
            boolValue: {
                handler(newValue,oldValue){
                    this.radio = this.boolValue;
                },
                immediate: true
            },
            radio: {
                handler(newValue, oldValue){
                    this.$emit('input',newValue);
                }
            }
        }
    })
</script>

<style scoped>

</style>
