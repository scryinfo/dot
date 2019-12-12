import Vue from 'vue';
import Vuex from 'vuex';
import {createDirectStore} from 'direct-vuex';
import page, {Page} from './page';

export interface State {
    page: Page;
}

Vue.use(Vuex);

const {store, rootActionContext, moduleActionContext} = createDirectStore({
    modules: {
        page,
    },

    getters: {},
    mutations: {},
    actions: {},
});

export default store;
export {rootActionContext, moduleActionContext};
export type AppStore = typeof store;
declare module 'vuex' {
    interface Store<S> {
        direct: AppStore;
    }
}
