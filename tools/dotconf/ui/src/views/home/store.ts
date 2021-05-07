import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);
export const state = {
    Loading: false,
    Dots: new Array<Dot>(),
    Configs: new Array<Dot>(),
};
export default new Vuex.Store({
    state: state,
    mutations: {},
    actions: {},
    modules: {}
})

export class Dot {
    metaData = {
        typeId: "",
        version: "",
        name: "",
        showName: "",
        single: false,
        relyTypeIds: null,
        flag:""
    };
    lives = [{
        typeId: "",
        liveId: "",
        relyLives: null,
        Dot: null,
        json: {},
        name: ""
    }];
    requiredInfo = null
}
