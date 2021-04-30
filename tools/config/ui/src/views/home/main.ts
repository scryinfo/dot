import Vue from "vue";
import App from "@/views/home/app.vue";
import router from "@/views/home/router";
import store from "@/views/home/store";
import  MessageBox  from 'element-ui';
Vue.use(MessageBox);
Vue.prototype.MessageBox = MessageBox;


Vue.config.productionTip = false;

new Vue({
    router,
    store,
    render: h => h(App)
}).$mount("#app");

