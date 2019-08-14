<template>
    <json-view v-model="parsedData" :parsedData="parsedData"></json-view>
</template>

<script lang="ts">
    import Vue from 'vue';
    import {jsonParse,makeJson} from "@/components/changeDataStructure/chDS";

    export default Vue.extend({
        name: "ExtendConfigEditor",
        props:  {
            objData: {
                type: Object,
                required: true
            },
        },
        data () {
            return {
                parsedData: [],
                lastParsedData: {}
            };
        },
        watch: {
            objData: {
                handler(newValue, oldValue) {
                    (this as any).parsedData = jsonParse(this.objData);
                },
                immediate: true
            },
            parsedData: {
                handler(newValue, oldValue) {
                    if (JSON.stringify(newValue) === JSON.stringify(this.lastParsedData)) {
                        return;
                    }
                    this.lastParsedData = newValue;
                    this.$emit("input", makeJson(this.parsedData));
                },
                deep: true
            }
        }
    })
</script>

<style scoped>

</style>
