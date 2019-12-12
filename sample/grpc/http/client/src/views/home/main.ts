import Vue from 'vue';
import App from './app.vue';
import store from './store';
import router from './router';
import Element from 'element-ui';

Vue.config.productionTip = false; // 控制:开发模式

Vue.use(Element);
new Vue({
    router,
    store: store.original,
    render: (h) => h(App),
}).$mount('#app');
