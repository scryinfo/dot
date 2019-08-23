import Vue from 'vue';
import App from './App.vue';
import router from './router';
import './registerServiceWorker';
import Element from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import ExtendConfigEditor from './components/initial';

Vue.use(Element)
Vue.config.productionTip = false;
Vue.use(ExtendConfigEditor);

new Vue({
  router,

  data: () => {
    return{
      Dots: [],
      Configs: [],
    };
  },
  render: (h) => h(App),
}).$mount('#app');
