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
                lastParsedData: {},
            }
        },
        watch: {
            objData: {
                handler(newValue, oldValue) {
                    this.parsedData = this.jsonParse(this.objData.RelyLives);
                },
                immediate: true
            },
            parsedData: {
                handler(newValue, oldValue) {
                    if (JSON.stringify(newValue) === JSON.stringify(this.lastParsedData)){
                        return;
                    }
                    this.lastParsedData = newValue;
                    this.makeJson(this.parsedData);
                },
                deep: true
            }
        },
        methods: {
            jsonParse(Json) {
                if (Json === null){
                    Json = {}
                }
                let result = [];
                let keys = Object.keys(Json);
                keys.forEach((k, index) => {
                    let val = Json[k];
                    let newObject = {name: k, remark: val};
                    result.push(newObject)
                });
                return result;
            },
            makeJson(ParData) {
                let Revert = {};
                for (let i = 0; i < data.length; ++i) {
                    let el = data[i];
                    let key, val;
                    key = el.name;
                    val = el.remark;
                    Revert[key] = val;
                }
                return Revert;
            }
        }
    }
</script>

<style scoped>

</style>
