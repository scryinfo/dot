<template>
    <div>
        <el-collapse v-model="activeTypes">
            <el-row v-for="(config,index) in $store.state.Configs">
                <div v-if="config.metaData.flag==='not-exist'">
                    <el-col :span="1">
                        <el-tooltip effect="dark" content="This typeId is not exist in dots!" placement="bottom-start">
                            <div class="grid-content bg-warning" style="text-align: center;line-height: 46px;">
                                {{index+1}}
                            </div>
                        </el-tooltip>
                    </el-col>
                </div>
                <div v-else>
                    <el-col :span="1">
                        <div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{index+1}}
                        </div>
                    </el-col>
                </div>
                <check-live-id :metaName="config.metaData.name" :lives="config.lives"/>
                <el-col :span="18">
                    <el-collapse-item v-bind:title="config.metaData.typeId" v-bind:name="index">
                        <el-row v-for="(live,index2) in config.lives">
                            <el-col :span="2">
                                <div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">
                                    {{live.name}}
                                </div>
                            </el-col>
                            <el-col :span="17">
                                <el-collapse-item v-bind:title="live.liveId" v-bind:name="index+' '+index2">
                                    <el-row>
                                        <el-col :span="2"><label>name</label></el-col>
                                        <el-col :span="15">
                                            <el-input type="text" v-model="live.name" placeholder="Name"/>
                                        </el-col>
                                    </el-row>
                                    <el-row>
                                        <el-col :span="2"><label>liveId</label></el-col>
                                        <el-col :span="15">
                                            <el-input type="text" v-model="live.liveId" placeholder="LiveId"/>
                                        </el-col>
                                        <el-col :span="4">
                                            <el-button @click="uuidGenerator(live)">Generate Live Id</el-button>
                                        </el-col>
                                    </el-row>

                                    <rely-lives-editor :objData="live" v-model="live.relyLives"/>

                                    <el-row v-if="live.json">
                                        <el-col :span="20">
                                            <el-collapse-item title="Extend Config for live"
                                                              v-bind:name="index+','+index2">
                                                <extend-config-editor
                                                        :objData="live.json"
                                                        v-model="live.json"/>
                                            </el-collapse-item>
                                        </el-col>
                                        <el-col :span="4">
                                            <json-button :objc="live.json" v-model="live.json"/>
                                        </el-col>
                                    </el-row>
                                </el-collapse-item>
                            </el-col>
                            <el-col :span="2">
                                <json-button :objc="config.lives[index2]" v-model="config.lives[index2]"/>
                            </el-col>
                            <el-col :span="3">
                                <el-button @click="removeObject(config,index2)" :disabled="config.lives.length <= 1">
                                    Remove Live
                                </el-button>
                            </el-col>
                        </el-row>
                    </el-collapse-item>
                </el-col>
                <el-col :span="3">
                    <el-row>
                        <el-button @click="showDialog(config.metaData.typeId,config.lives)" size="mini">Load By Config
                        </el-button>
                    </el-row>
                    <el-row><span><el-button @click="AddObject(config,config.metaData.typeId)"
                                             size="mini">Add Live</el-button></span><span><el-button
                            @click="removeType(index)" size="mini">remove type</el-button></span></el-row>
                </el-col>
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
                <i class="el-icon-upload"/>
                <div class="el-upload__text">Drag and drop files here, or <em>click to upload</em></div>
                <div class="el-upload__tip" slot="tip">Can only upload json files</div>
            </el-upload>
            </span>
            <span slot="footer" class="dialog-footer">
                <el-button @click="dialogVisible = false">Cancel</el-button>
                <el-button type="primary" @click="handleConfirm()">Ok</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from "vue-property-decorator";
    import RelyLivesEditor from '../components/config/RelyLivesEditor.vue';
    import JsonButton from '../components/config/JsonButton.vue';
    import CheckLiveId from '../components/config/checkLiveId.vue';

    const uuidv1 = require('uuid/v1');

    @Component({
        components: {
            "rely-lives-editor": RelyLivesEditor,
            "json-button": JsonButton,
            "check-live-id": CheckLiveId
        }
    })
    export default class Config extends Vue {
        private activeTypes = [];
        private dialogVisible = false;
        private textarea = '';
        private keyarea = '';
        private objc = null;
        private model = '';
        private schemaObject = {};
        private upLoadFile = '';
        private typeId = '';
        private lives = []

        private removeType(index: number) {
            (this as any).$store.state.Configs.splice(index, 1);
        }

        AddObject(config: any, typeId: string) {
            if (config.metaData.flag) {
                let confcopy: any = this.shallowCopy(config.lives[config.lives.length - 1]);
                config.lives.push(confcopy);
            } else {
                for (let dot of (this as any).$store.state.Dots) {
                    if (dot.metaData.typeId === typeId) {
                        let dotcopy: any = this.shallowCopy(dot.lives[0]);
                        config.lives.push(dotcopy);
                        break;
                    }
                }
            }
        }

        private removeObject(config: any, index: number) {
            config.lives.splice(index, 1);
        }

        private uuidGenerator(live: any) {
            live.liveId = uuidv1()
        }

        private shallowCopy(src: any): any {
            let dst: string;
            dst = JSON.stringify(src);
            return JSON.parse(dst);
        }

        private uploadSectionFile(param: any) {
            let fileObj = param.file;
            let Blb = fileObj.slice();
            let reader = new FileReader();
            reader.readAsText(fileObj, 'utf-8');
            let result: any;
            reader.onload = this.fileOnload;
        }

        private fileOnload(e: any) {
            this.upLoadFile = e.target.result;
        }

        private showDialog(typeId: string, lives: any) {
            this.typeId = typeId;
            this.lives = lives;
            this.dialogVisible = true;
        }

        private handleConfirm() {
            let configLives = this.findConfigLives(this.typeId);
            let dotLive = this.findDotLive(this.typeId);
            if (configLives && dotLive) {
                for (let i = 0, len = configLives.length; i < len; i++) {
                    this.assembleByLiveId(dotLive, configLives[i]);
                }
            }
            this.dialogVisible = false;

        }

        private findConfigLives(typeId: string): any {
            let config = JSON.parse(this.upLoadFile);
            for (let i = 0, len = config.dots.length; i < len; i++) {
                if (config.dots[i].metaData.typeId === typeId) {
                    return config.dots[i].lives;
                }
            }
            return null;
        }

        private findDotLive(typeId: string): any {
            for (let i = 0, len = (this as any).$store.state.Dots.length; i < len; i++) {
                if ((this as any).$store.state.Dots[i].metaData.typeId === typeId) {
                    return (this as any).$store.state.Dots[i].lives[0];
                }
            }
            return null;
        }

        private assemble(schema: any, source: any): any {
            for (let key in schema) {
                schema[key] = this.isObject(schema[key]) ? this.assemble(schema[key], (source[key] ? source[key] : schema[key])) : (source[key] ? source[key] : schema[key]);
            }
            return schema;
        }

        private isObject(o: any): any {
            return (typeof o === 'object') && o !== null;
        }

        private assembleByLiveId(dotLive: any, configLive: any) {
            let flag = false;
            for (let i = 0, len = this.lives.length; i < len; i++) {
                if ((this as any).lives[i].liveId === configLive.liveId) {
                    (this as any).lives[i] = JSON.parse(JSON.stringify(this.assemble((this as any).lives[i], configLive)));
                    flag = true;
                    break;
                }
            }
            if (!flag) {
                let target: any = (this as any).assemble(dotLive, configLive);
                (this as any).lives.push(JSON.parse(JSON.stringify(target)));
            }
        }
    }
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
