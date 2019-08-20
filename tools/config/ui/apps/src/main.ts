import Vue from 'vue';
import App from './App.vue';
import router from './router';

import './registerServiceWorker';
import './plugins/element.js';
import ExtendConfigEditor from './components/initial';

Vue.config.productionTip = false;
Vue.use(ExtendConfigEditor);

new Vue({
  router,

  data: () => {
    return{
      Dots: [],
      ExportDots: [],
      Configs: [],
      DotsTem: [],
    };
  },
  render: (h) => h(App),
}).$mount('#app');
