<template>
<div>
    <el-button @click="ShowJsonDialog(objc)">JSON</el-button>
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
                placeholder="Please input json data!"
                v-model="textarea">
        </el-input>
    </el-drawer>
</div>
</template>

<script lang="ts">
import Vue from 'vue';
export default Vue.extend({
    name: 'JsonButton',
    props: {
        objc: {},
    },
    data() {
      return {
          dialog: false,
          schemaObject: {},
          textarea: '',
      };
    },
    methods: {
        ShowJsonDialog(obj: any) {
            this.dialog = true;
            (this as any).objc = obj;
            const jsonSchemaGenerator = require('./schemaGenerator/index.js');
            (this as any).schemaObject = jsonSchemaGenerator.jsonToSchema(obj);
            for (const key in (this as any).schemaObject.properties) {
                if (key === 'relyLives') {
                    (this as any).schemaObject.properties[key].required = [];
                    this.$delete((this as any).schemaObject.properties[key], 'maxProperties');
                    //todo: highest json button relyLives validation
                    break;
                }
            }
            this.textarea = JSON.stringify(obj, null, 4);
        },
        handleClose(done: any) {
            try {
                if (this.textarea) {
                    const objct: any = JSON.parse(this.textarea);
                    const tv4 = require('tv4');
                    if (tv4.validate(objct, this.schemaObject)) {
                        this.$emit('input', objct);
                    } else {
                        (this as any).$message.error('json text input error!');
                    }
                } else {
                    (this as any).$message.error('json text input error!');
                }
            } catch (e) {
                (this as any).$message.error('json text input error!');
            } finally {
                done();
            }
        },
    },
});
</script>

<style scoped>

</style>
