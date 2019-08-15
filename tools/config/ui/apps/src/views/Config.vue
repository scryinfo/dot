<template>
    <div>
    <el-collapse v-model="activeTypes">
        <el-row v-for="(config,index) in $root.Configs">
            <div v-if="config.metaData.flag=='not-exist'">
                <el-col :span="1" ><el-tooltip effect="dark" content="This typeId is not exist in dots!" placement="bottom-start"><div class="grid-content bg-warning" style="text-align: center;line-height: 46px;">{{index+1}}</div></el-tooltip></el-col>
            </div>
            <div v-else>
                <el-col :span="1" ><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{index+1}}</div></el-col>
            </div>
            <el-col :span="2"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{config.metaData.name}}</div></el-col>
            <el-col :span="16">
                <el-collapse-item v-bind:title="config.metaData.typeId" v-bind:name="index">
                    <el-row v-for="(live,index2) in config.lives">
                        <el-col :span="2"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{live.name}}</div></el-col>
                        <el-col :span="17">
                            <el-collapse-item v-bind:title="live.liveId" v-bind:name="index+' '+index2">
                                <el-row><el-col :span="2"><label>name</label></el-col><el-col :span="15"><el-input type="text" v-model="live.name" placeholder="Name"></el-input></el-col></el-row>
                                <el-row><el-col :span="2"><label>liveId</label></el-col><el-col :span="15"><el-input type="text" v-model="live.liveId" placeholder="LiveId"></el-input></el-col>
                                <el-col :span="4"><el-button @click="UuidGenerator(live)">Generate Live Id</el-button></el-col>
                                </el-row>

                                <rely-lives-editor
                                        :objData="live"
                                        v-model="live.relyLives"
                                ></rely-lives-editor>

                                <el-row v-if="live.json"><el-col :span="20"><el-collapse-item title="Extend Config for live" v-bind:name="index+','+index2">
                                    <extend-config-editor
                                            :objData="live.json"
                                            v-model="live.json"
                                    ></extend-config-editor>
                                </el-collapse-item></el-col><el-col :span="4"><json-button :objc="live.json" v-model="live.json"></json-button></el-col></el-row>
                            </el-collapse-item>
                        </el-col>
                        <el-col :span="2"><json-button :objc="config.lives[index2]" v-model="config.lives[index2]"></json-button></el-col><el-col :span="3"><el-button @click="RemoveObject(config,'lives',index2)">Remove Live</el-button></el-col>
                    </el-row>
                </el-collapse-item>
            </el-col>
            <el-col :span="4"><el-button @click="showDialog(config.metaData.typeId,config.lives)">Load By Config</el-button><el-button @click="AddObject(config,config.metaData.typeId,'lives')">Add Live</el-button></el-col>
        </el-row>
    </el-collapse>
        <el-dialog
        title="load by config"
        :visible.sync="dialogVisible"
        width="40%">
            <span style="text-align: center">
              <el-upload
                      action=""
                      class="upload-demo"
                      :http-request="uploadSectionFile"
                      drag
                      :limit=1
                      >
                <i class="el-icon-upload"></i>
                <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
                <div class="el-upload__tip" slot="tip">只能上传json文件</div>
            </el-upload>
            </span>
            <span slot="footer" class="dialog-footer">
                <el-button @click="dialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="handleConfirm()">确 定</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue'
    import RelyLivesEditor from '../components/RelyLivesEditor.vue'
    import JsonButton from '../components/JsonButton.vue'
    const uuidv1 = require('uuid/v1')
    export default Vue.extend({
        data() {
            return {
                activeTypes:[],
                dialogVisible: false,
                textarea: '',
                keyarea: '',
                objc: null,
                model: '',
                schemaObject: {},
                upLoadFile: '',
                typeId: '',
                lives: []
            }
        },
        components: {
          "rely-lives-editor": RelyLivesEditor,
            "json-button": JsonButton
        },
        methods: {
            AddObject(config:any,typeId:string,keyss:string) {
                for (let dot of (this as any).$root.DotsTem) {
                    if (dot.metaData.typeId === typeId) {
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
                live.liveId=uuidv1()
            },
            shallowCopy(src:any):any {
                let dst:any = {};
                for (let prop in src) {
                    if (src.hasOwnProperty(prop)) {
                     dst[prop] = src[prop];
                    }
                }
                return dst;
            },
            uploadSectionFile(param:any){
                let fileObj = param.file;
                let Blb = fileObj.slice();
                let reader = new FileReader();
                reader.readAsText(fileObj,'utf-8');
                let result:any;
                reader.onload = this.fileOnload;
            },
            fileOnload(e:any){
                this.upLoadFile = e.target.result;
            },
            showDialog(typeId:string,lives:any){
                this.typeId = typeId;
                this.lives = lives;
                this.dialogVisible = true;
            },
            handleConfirm(){
                let configLives = this.findConfigLives(this.typeId);
                let dotLive = this.findDotLive(this.typeId);
                if (configLives && dotLive){
                    for(let i = 0, len = configLives.length; i < len; i++){
                        this.assembleByLiveId(dotLive,configLives[i]);
                        // let target:any = (this as any).assemble(dotLive,configLives[i]);
                        // (this as any).lives.push(target);
                    }
                }
                this.dialogVisible = false;

            },
            findConfigLives(typeId:string):any{
                let config = JSON.parse(this.upLoadFile);
                for(let i = 0, len = config.dots.length; i < len; i++){
                    if(config.dots[i].metaData.typeId === typeId){
                        return config.dots[i].lives;
                    }
                }
                return null;
            },
            findDotLive(typeId:string):any{
                for(let i = 0, len = (this as any).$root.DotsTem.length; i < len; i++){
                    if((this as any).$root.DotsTem[i].metaData.typeId === typeId){
                        return (this as any).$root.DotsTem[i].lives[0];
                    }
                }
                return null;
            },
            assemble(schema:any,source:any):any{
                for(let key in schema){
                    schema[key] = this.isObject(schema[key]) ? this.assemble(schema[key],(source[key] ? source[key] : schema[key])) : (source[key] ? source[key] : schema[key]);
                }
                return schema;
            },
            isObject(o:any):any {
                return (typeof o === 'object') && o !==null;
            },
            assembleByLiveId(dotLive:any, configLive:any){
                let flag = false;
                for(let i = 0, len = this.lives.length; i < len; i++){
                    if((this as any).lives[i].liveId === configLive.liveId){
                        (this as any).lives[i] = this.assemble((this as any).lives[i],configLive);
                        flag = true;
                        break;
                    }
                }
                if (!flag){
                    let target:any = (this as any).assemble(dotLive,configLive);
                    (this as any).lives.push(target);
                }
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
    .bg-warning {
        background: #d6a23c;
    }
</style>
