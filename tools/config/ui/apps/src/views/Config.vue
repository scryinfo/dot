<template>
    <div>
    <el-collapse v-model="activeTypes">
        <el-row v-for="(config,index) in this.$root.Configs">
            <el-col :span="1" ><el-tag type="info" effect="plain">{{index+1}}</el-tag></el-col>
            <el-col :span="2"><el-tag type="info" effect="plain">{{config.Meta.name}}</el-tag></el-col>
            <el-col :span="16">
                <el-collapse-item v-bind:title="config.Meta.typeId" v-bind:name="index">
                    <el-row v-for="(live,index2) in config.Lives">
                        <el-col :span="2"><el-tag type="info" effect="plain">{{live.name}}</el-tag></el-col>
                        <el-col :span="17">
                            <el-collapse-item v-bind:title="live.LiveId" v-bind:name="index+' '+index2">

                                <el-row><el-col :span="2"><label>name</label></el-col><el-col :span="15"><el-input v-model="live.name" placeholder="Name"></el-input></el-col></el-row>
                                <el-row><el-col :span="2"><label>liveId</label></el-col><el-col :span="15"><el-input v-model="live.LiveId" placeholder="Liveid"></el-input></el-col>
                                <el-col :span="4"><el-button @click="UuidGenerator(live)">Generate Live Id</el-button></el-col>
                                </el-row>

                                <el-row><el-col :span="2"><label>relyLives</label></el-col><el-col :span="5"><el-button @click="AddRelyLives(config.Lives[index2])">add</el-button><el-button @click="ShowJsonDialog(live.RelyLives)">JSON</el-button></el-col></el-row>
                                <el-row v-for="(relyLiveId, liveName, index3) in live.RelyLives"><el-col :span="5" :offset="2"><el-input v-model="liveName"></el-input></el-col>
                                    <el-col :span="2"><el-dropdown>
                                        <el-button type="primary">
                                            更多菜单<i class="el-icon-arrow-down el-icon--right"></i>
                                        </el-button>
                                        <el-dropdown-menu slot="dropdown">
                                            <el-dropdown-item>黄金糕</el-dropdown-item>
                                            <el-dropdown-item>狮子头</el-dropdown-item>
                                            <el-dropdown-item>螺蛳粉</el-dropdown-item>
                                            <el-dropdown-item>双皮奶</el-dropdown-item>
                                            <el-dropdown-item>蚵仔煎</el-dropdown-item>
                                        </el-dropdown-menu>
                                    </el-dropdown></el-col>
                                    <el-col :span="2" style="text-align: center">:</el-col ><el-col :span="8"><el-input v-model="live.RelyLives[liveName]"></el-input></el-col><el-col :span="2"><el-button @click="RemoveRelyLives(live.RelyLives,liveName)">remove</el-button></el-col></el-row>
                            </el-collapse-item>
                        </el-col>
                        <el-col :span="5"><el-button @click="ShowJsonDialog(live)">JSON</el-button><el-button @click="RemoveObject(config,'Lives',index2)">Remove Live</el-button></el-col>
                    </el-row>
                </el-collapse-item>
            </el-col>
            <el-col :span="4"><el-button>Load By Config</el-button><el-button @click="AddObject(config,config.Meta.typeId,'Lives')">Add Live</el-button></el-col>
        </el-row>
    </el-collapse>

        <el-drawer
                title="JSON textarea!"
                :before-close="handleClose"
                :visible.sync="dialog"
                direction="ltr"
                custom-class="demo-drawer"
                ref="drawer"
        >
            <el-input
                    type="textarea"
                    :autosize="{ minRows: 10, maxRows: 30}"
                    placeholder="请输入内容"
                    v-model="textarea">
            </el-input>
        </el-drawer>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue'
    const uuidv1 = require('uuid/v1')
    export default Vue.extend({
        data() {
            return {
                activeTypes:[],
                dialog: false,
                textarea: '',
                keyarea: '',
                objc: null,
                model: '',
            }
        },
        methods: {
            AddObject(config:any,typeId:string,keyss:string) {
                for (let dot of this.$root.Dots) {
                    if (dot.Meta.typeId === typeId) {
                        let dotcopy;
                        eval("dotcopy = this.shallowCopy(dot."+keyss+"[0])");
                        eval("config."+keyss+".push(dotcopy)");
                        break;
                    }
                }
            },
            RemoveObject(config:any, keyss:string, index:number) {
                eval("config."+keyss+".splice("+index+","+1+")")
            },
            UuidGenerator(live:any) {
                live.LiveId=uuidv1()
            },
            ShowJsonDialog(obj:any){
                this.dialog = true;
                this.objc = obj;
                this.textarea = JSON.stringify(obj,null,4);

            },
            handleClose(done:any){
                let objct:any = JSON.parse(this.textarea);
                for (let prop in objct) {
                    if (objct.hasOwnProperty(prop)) {
                        this.objc[prop] = objct[prop];
                    }
                }
                done();
            },
            AddRelyLives(live:any){
                if(live.RelyLives===null){
                    console.log("add rely");
                    this.$set(live,'RelyLives',{})
                    console.log(live.RelyLives);
                }
                this.$set(live.RelyLives,'default','please change default');
            },
            RemoveRelyLives(relyLives:any,name:string){
                this.$delete(relyLives,name)
            },
            KeyInputFocus(){

            },
            KeyInputChanges(){

            },
            shallowCopy(src:any):any {
                let dst:any = {};
                for (let prop in src) {
                    if (src.hasOwnProperty(prop)) {
                     dst[prop] = src[prop];
                    }
                }
                return dst;
            }
        }
    })

</script>

<style scoped>
    .el-dropdown {
        vertical-align: top;
    }
    .el-dropdown + .el-dropdown {
        margin-left: 15px;
    }
    .el-icon-arrow-down {
        font-size: 12px;
    }
</style>
