<template>
    <div>
        <el-row><el-col :span="2">relyLives</el-col><el-col :span="2"><el-button :disabled="disable" @click="addObeject()">add</el-button></el-col><el-col :span="3"><json-rely-button :objc="parsedData" v-model="parsedData"></json-rely-button></el-col></el-row>
        <el-row v-for="(ob,index) in parsedData" v-model="parsedData"><el-col :span="6" :offset="2"><el-input v-model="ob.name"></el-input></el-col><el-col :span="2"><el-dropdown trigger="click">
      <span class="el-dropdown-link" style="text-align: center;line-height: 46px;">
            Select<i class="el-icon-arrow-down el-icon--right"></i>
      </span>
            <el-dropdown-menu slot="dropdown">
                <div v-for="config in $root.Configs">
                    <div v-for="live in config.lives">
                        <el-dropdown-item @click.native="changeItem(ob,live)">{{live.name}}:{{live.liveId}}</el-dropdown-item>
                    </div>
                </div>
            </el-dropdown-menu>
        </el-dropdown></el-col><el-col :span="10"><el-input v-model="ob.remark"></el-input></el-col><el-col :span="3"><el-button @click="removeObject(index)">remove</el-button></el-col></el-row>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    import JsonRelyButton from './JsonRelyButton.vue'
    export default Vue.extend({
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
        components: {
          "json-rely-button": JsonRelyButton
        },
        watch: {
            objData: {
                handler(newValue, oldValue) {
                    (this as any).parsedData = this.jsonParse(this.objData.relyLives);
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
            jsonParse(Json:any) {
                if (Json === null){
                    Json = {}
                }
                let result:any = [];
                let keys = Object.keys(Json);
                keys.forEach((k, index) => {
                    let val = Json[k];
                    let newObject = {name: k, remark: val};
                    result.push(newObject)
                });
                return result;
            },
            makeJson(ParData:any) {
                let Revert:any = {};
                for (let i = 0; i < ParData.length; ++i) {
                    let el = ParData[i];
                    let key:string, val:string;
                    key = el.name;
                    val = el.remark;
                    Revert[key] = val;
                }
                return Revert;
            },
            checkKey() {
                for (let i = 0; i < this.parsedData.length; ++i) {
                    if ((this as any).parsedData[i].name === 'default'){
                        this.disable = true;
                        return;
                    }
                }
                this.disable = false;
            },
            addObeject() {
                let obj:object = {name: 'default', remark: 'please change default'};
                (this as any).parsedData.push(obj);
            },
            removeObject(index:number) {
                this.parsedData.splice(index,1);
            },
            changeItem(ob:any,live:any) {
                if (live.name){
                    ob.name=live.name;
                }
                if (live.liveId){
                    ob.remark=live.liveId;
                }
            }
        }
    })
</script>

<style scoped>

</style>
