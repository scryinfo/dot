import ExtendConfigEditor from './ExtendConfigEditor.vue';
import JsonView from './JsonView.vue';
import ArrayView from './ArrayView.vue';


const install = (Vue: any) => {
    Vue.component('ExtendConfigEditor', ExtendConfigEditor);
    Vue.component('json-view', JsonView);
    Vue.component('array-view', ArrayView);
};

export default install;
