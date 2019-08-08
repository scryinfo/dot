import ExtendConfigEditor from './ExtendConfigEditor.vue';
import JsonView from './JsonView.vue';
import ArrayView from './ArrayView.vue';
import BoolView from './BoolView.vue';


const install = (Vue:any) => {
    Vue.component('ExtendConfigEditor',ExtendConfigEditor);
    Vue.component('json-view', JsonView);
    Vue.component('array-view', ArrayView);
    Vue.component('bool-view', BoolView);
};

export default install
