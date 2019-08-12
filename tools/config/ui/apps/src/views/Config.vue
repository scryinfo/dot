<template>
    <div>
    <el-collapse v-model="activeTypes">
        <el-row v-for="(config,index) in $root.Configs">
            <el-col :span="1" ><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{index+1}}</div></el-col>
            <el-col :span="2"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{config.metaData.name}}</div></el-col>
            <el-col :span="16">
                <el-collapse-item v-bind:title="config.metaData.typeId" v-bind:name="index">
                    <el-row v-for="(live,index2) in config.lives">
                        <el-col :span="2"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{live.name}}</div></el-col>
                        <el-col :span="17">
                            <el-collapse-item v-bind:title="live.LiveId" v-bind:name="index+' '+index2">
                                <el-row><el-col :span="2"><label>name</label></el-col><el-col :span="15"><el-input v-model="live.name" placeholder="Name"></el-input></el-col></el-row>
                                <el-row><el-col :span="2"><label>liveId</label></el-col><el-col :span="15"><el-input v-model="live.liveId" placeholder="Liveid"></el-input></el-col>
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
            <el-col :span="4"><el-button>Load By Config</el-button><el-button @click="AddObject(config,config.metaData.typeId,'lives')">Add Live</el-button></el-col>
        </el-row>
    </el-collapse>
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
            "json-button": JsonButton
        },
        methods: {
            AddObject(config:any,typeId:string,keyss:string) {
                for (let dot of (this as any).$root.Dots) {
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
