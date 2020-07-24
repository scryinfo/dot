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
        <div>
            <h1>server stream</h1>
            <input placeholder="Please input data" v-model="serverStreamData" @keyup.enter="onServerStream">
            <div v-text="serverStreamReturn"></div>
        </div>
        <div>
            <h1>client stream</h1>
            <input placeholder="Please input data" v-model="clientStreamData" @keyup.enter="onClientStream">
            <div v-text="clientStreamReturn"></div>
        </div>
        <div>
            <h1>both sides</h1>
            <input placeholder="Please input data" v-model="bothData" @keyup.enter="onBoth">
            <div v-text="bothReturn"></div>
        </div>
    </div>
</template>

<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator';
    import Hi from "@/api_wrapper/hi_client";
    import {HiRes, WriteRes} from "@/api/hi_pb";

    @Component
    export default class Home extends Vue {
        private hiName = ''

        private data = '';
        private dataReturn = '';

        private serverStreamData = '';
        private serverStreamReturn = '';

        private clientStreamData = '';
        private clientStreamReturn = '';

        private bothData = '';
        private bothReturn = '';


        private onHi(): boolean {
            Hi.hi("hi dot").then( res  => {
               this.hiName = (res as HiRes).getName();
            });
            return false;
        }

        private onWrite(): boolean {
            Hi.write(this.data).then(res => {
                this.dataReturn = (res as WriteRes).getData();
            });
            return false;
        }
        private onServerStream(): boolean {
            Hi.serverStream(this.serverStreamData).then(res => {
                this.serverStreamReturn = res.getReply();
            });
            return false;
        }
        private onClientStream(): boolean {
            Hi.clientStream(this.clientStreamData).then(res => {
                this.clientStreamReturn = res.getReply();
            });
            return false;
        }
        private onBoth(): boolean {
            Hi.bothSides(this.bothData).then(res => {
                this.bothReturn = res.getReply();
            });
            return false;
        }
    }
</script>
