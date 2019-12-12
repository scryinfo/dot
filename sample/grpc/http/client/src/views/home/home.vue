<template>
    <div>
        <div>
            <h1>Hi Dot</h1>
            <button @click="onHi"> send hi </button>
            <span v-text="hiName"></span>
        </div>
        <div>
            <h1>Write data</h1>
            <input placeholder="Please input data" v-model="data" @keyup.enter="onWrite">
            <div v-text="dataReturn"></div>
        </div>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator';
    import Hi from "@/api_wrapper/hi_client";

    @Component
    export default class Home extends Vue {
        private hiName = ''

        private data = '';
        private dataReturn = '';

        private onHi(): boolean {
            Hi.hi("hi dot").then(data =>{
               this.hiName = data.getName();
            });
            return false;
        }

        private onWrite(): boolean {
            Hi.write(this.data).then(data => {
                this.dataReturn = data.getData();
            });
            return false;
        }


    }
</script>
