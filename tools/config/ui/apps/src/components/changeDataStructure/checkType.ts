import Vue from 'vue';
export function checkType(dots: any, configs: any) {
    for (let i = 0, len = configs.length; i < len; i++) {
        if (compare(dots, configs[i])) {
            if (configs[i].metaData.flag) {
                Vue.set(configs[i].metaData, 'flag', null);
            }
        } else {
            if (configs[i].metaData.flag) {

            } else {
                Vue.set(configs[i].metaData, 'flag', 'not-exist');
            }
        }
    }
}

function compare(dots: any, config: any): Boolean {
    for (let i = 0, len = dots.length; i < len; i++) {
        if (dots[i].metaData.typeId === config.metaData.typeId) {
            return true;
        }
    }
    return false;
}

export function removeAllType(configs: any) {
    for (let i = 0, len = configs.length; i < len; i++) {
        if (configs[i].metaData.flag) {

        } else {
            Vue.set(configs[i].metaData, 'flag', 'not-exist');
        }
    }
}
