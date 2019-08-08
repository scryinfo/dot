<template>
    <div>
    <el-collapse v-model="activeTypes">
        <el-row v-for="(config,index) in $root.Configs">
            <el-col :span="1" ><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{index+1}}</div></el-col>
            <el-col :span="2"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{config.Meta.name}}</div></el-col>
            <el-col :span="16">
                <el-collapse-item v-bind:title="config.Meta.typeId" v-bind:name="index">
                    <el-row v-for="(live,index2) in config.Lives">
                        <el-col :span="2"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{live.name}}</div></el-col>
                        <el-col :span="17">
                            <el-collapse-item v-bind:title="live.LiveId" v-bind:name="index+' '+index2">
                                <el-row><el-col :span="2"><label>name</label></el-col><el-col :span="15"><el-input v-model="live.name" placeholder="Name"></el-input></el-col></el-row>
                                <el-row><el-col :span="2"><label>liveId</label></el-col><el-col :span="15"><el-input v-model="live.LiveId" placeholder="Liveid"></el-input></el-col>
                                <el-col :span="4"><el-button @click="UuidGenerator(live)">Generate Live Id</el-button></el-col>
                                </el-row>

                                <rely-lives-editor
                                        :objData="live"
                                        v-model="live.RelyLives"
                                        @JSONDialog="ShowJsonDialog"
                                ></rely-lives-editor>

                                <el-row v-if="live.json"><el-col :span="20"><el-collapse-item title="Extend Config for live" v-bind:name="index+','+index2">
                                    <extend-config-editor
                                            :objData="live.json"
                                            v-model="live.json"
                                    ></extend-config-editor>
                                </el-collapse-item></el-col><el-col :span="4"><el-button @click="ShowJsonDialog(live.json)">JSON</el-button></el-col></el-row>
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
    import RelyLivesEditor from '../components/RelyLivesEditor.vue'
    import ExtendConfigEditor from '../components/ExtendConfigEditor.vue'
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
                schemaObject: {}
            }
        },
        components: {
          "rely-lives-editor": RelyLivesEditor,
        },
        methods: {
            AddObject(config:any,typeId:string,keyss:string) {
                for (let dot of (this as any).$root.Dots) {
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
                let GenerateSchema = require('generate-schema');
                this.schemaObject = GenerateSchema.json(obj)
                this.textarea = JSON.stringify(obj,null,4);

            },
            handleClose(done:any){
                try{
                    if(this.textarea){
                        let objct:any = JSON.parse(this.textarea);
                        let tv4 = require('tv4');
                        if(tv4.validate(objct,this.schemaObject)){
                            for (let prop in objct) {
                                if (objct.hasOwnProperty(prop)) {
                                    (this as any).objc[prop] = objct[prop];
                                }
                            }
                        }else {
                            (this as any).$message.error('json text input error!');
                        }
                    }else{
                        (this as any).$message.error('json text input error!');
                    }
                }catch (e) {
                    (this as any).$message.error('json text input error!');
                }finally {
                    done();
                }
            },
            AddRelyLives(live:any){
                if(live.RelyLives===null){
                    this.$set(live,'RelyLives',{})
                }
                this.$set(live.RelyLives,'default','please change default');
            },
            RemoveRelyLives(relyLives:any,name:string){
                this.$delete(relyLives,name)
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
    .bg-purple {
        background: #d3dce6;
    }
    .grid-content {
        border-radius: 4px;
        min-height: 46px;
    }
</style>
