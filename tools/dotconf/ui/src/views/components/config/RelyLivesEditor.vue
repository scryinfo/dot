<template>
    <div>
        <el-row>
            <el-col :span="2">relyLives</el-col>
            <el-col :span="2">
                <el-button :disabled="disable" @click="addObeject()">add</el-button>
            </el-col>
            <el-col :span="3">
                <json-rely-button :objc="parsedData" v-model="parsedData"/>
            </el-col>
        </el-row>
        <el-row v-for="(ob,index) in parsedData" v-model="parsedData">
            <el-col :span="6" :offset="2">
                <el-input type="text" v-model="ob.name"/>
            </el-col>
            <el-col :span="2">
                <el-dropdown trigger="click">
      <span class="el-dropdown-link" style="text-align: center;line-height: 46px;">
            Select<i class="el-icon-arrow-down el-icon--right"/>
      </span>
                    <el-dropdown-menu slot="dropdown">
                        <div v-for="config in $store.state.Configs">
                            <div v-for="live in config.lives">
                                <el-dropdown-item @click.native="changeItem(ob,live)">{{live.name}}:{{live.liveId}}
                                </el-dropdown-item>
                            </div>
                        </div>
                    </el-dropdown-menu>
                </el-dropdown>
            </el-col>
            <el-col :span="10">
                <el-input type="text" v-model="ob.remark">
                </el-input>
            </el-col>
            <el-col :span="3">
                <el-button @click="removeObject(index)">remove</el-button>
            </el-col>
        </el-row>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    import JsonRelyButton from './JsonRelyButton.vue';
    import {jsonParseRely, makeJsonRely} from "@/views/components/utils/changeDataStruct";

    export default Vue.extend({
        name: "RelyLivesEditor",
        props: {
            objData: {
                type: Object,
                required: true
            }
        },
        data() {
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
                    (this as any).parsedData = jsonParseRely(this.objData.relyLives);
                },
                immediate: true
            },
            parsedData: {
                handler(newValue, oldValue) {
                    this.checkKey();
                    this.$emit("input", makeJsonRely(newValue));
                },
                deep: true
            }
        },
        methods: {
            checkKey() {
                for (let i = 0; i < this.parsedData.length; ++i) {
                    if ((this as any).parsedData[i].name === 'default') {
                        this.disable = true;
                        return;
                    }
                }
                this.disable = false;
            },
            addObeject() {
                let obj: object = {name: 'default', remark: 'please change default'};
                (this as any).parsedData.push(obj);
            },
            removeObject(index: number) {
                this.parsedData.splice(index, 1);
            },
            changeItem(ob: any, live: any) {
                if (live.name) {
                    ob.name = live.name;
                }
                if (live.liveId) {
                    ob.remark = live.liveId;
                }
            }
        }
    })
</script>

<style scoped>

</style>
