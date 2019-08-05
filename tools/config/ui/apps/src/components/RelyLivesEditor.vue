<template>
    <div></div>
</template>

<script>
    export default {
        name: "RelyLivesEditor",
        props: {
            objData: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                parsedData: [],
                lastParsedData: {}
            }
        },
        watch: {
            objData: {
                handler(newValue, oldValue) {
                    this.parsedData = this.jsonParse(this.objData);
                },
                immediate: true
            },
            parsedData: {
                handler(newValue, oldValue) {
                    if (JSON.stringify(newValue) === JSON.stringify(this.lastParsedData)){
                        return;
                    }
                    this.lastParsedData = newValue;
                    this.$emit("input", this.makeJson(this.parsedData));
                },
                deep: true
            }
        }
    }
</script>

<style scoped>

</style>
