<template>
    <div>
        <el-row><el-col :span="2">relyLives</el-col><el-col :span="3"><el-button :disabled="disable" @click="addObeject()">add</el-button></el-col><el-col :span="3"><el-button @click="jsonButton()">JSON</el-button></el-col></el-row>
        <el-row v-for="(ob,index) in parsedData" v-model="parsedData"><el-col :span="6" :offset="2"><el-input v-model="ob.name"></el-input></el-col><el-col :span="2"><el-dropdown trigger="click">
      <span class="el-dropdown-link" style="text-align: center;line-height: 46px;">
            Select<i class="el-icon-arrow-down el-icon--right"></i>
      </span>
            <el-dropdown-menu slot="dropdown">
                <div v-for="config in $root.Configs">
                    <div v-for="live in config.Lives">
                        <el-dropdown-item @click.native="changeItem(ob,live)">{{live.name}}:{{live.LiveId}}</el-dropdown-item>
                    </div>
                </div>
            </el-dropdown-menu>
        </el-dropdown></el-col><el-col :span="10"><el-input v-model="ob.remark"></el-input></el-col><el-col :span="3"><el-button @click="removeObject(index)">remove</el-button></el-col></el-row>
    </div>
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
                disable: false,
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
                    this.checkKey();
                    this.$emit("input",this.makeJson(newValue));
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
                for (let i = 0; i < ParData.length; ++i) {
                    let el = ParData[i];
                    let key, val;
                    key = el.name;
                    val = el.remark;
                    Revert[key] = val;
                }
                return Revert;
            },
            checkKey() {
                for (let i = 0; i < this.parsedData.length; ++i) {
                    if (this.parsedData[i].name === 'default'){
                        this.disable = true;
                        return;
                    }
                }
                this.disable = false;
            },
            addObeject() {
                let obj = {name: 'default', remark: 'please change default'};
                this.parsedData.push(obj);
            },
            removeObject(index) {
                this.parsedData.splice(index,1);
            },
            jsonButton() {
                this.$emit("JSONDialog",this.objData.RelyLives)
            },
            changeItem(ob,live) {
                if (live.name){
                    ob.name=live.name;
                }
                if (live.LiveId){
                    ob.remark=live.LiveId;
                }
            }
        }
    }
</script>

<style scoped>

</style>
