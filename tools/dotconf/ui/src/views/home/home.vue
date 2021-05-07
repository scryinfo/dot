<template>
    <el-row type="flex">
        <el-col :span="3">
            <ul>
                <li>
                    <router-link :to="{name:'findDot'}">组件查找</router-link>
                </li>
                <li>
                    <router-link :to="{name:'dotList'}">组件列表</router-link>
                </li>
                <li>
                    <router-link :to="{name:'import'}">导入</router-link>
                </li>
                <li>
                    <router-link :to="{name:'export'}">导出</router-link>
                </li>
                <li>
                    <router-link :to="{name:'config'}">配置</router-link>
                </li>
            </ul>
        </el-col>
        <el-col :span="20">
            <router-view/>
        </el-col>
    </el-row>
</template>

<script lang="ts">
    import {Component, Vue} from "vue-property-decorator";
    import {
        Button,
        Col,
        DatePicker,
        Form,
        FormItem,
        Input,
        Loading,
        Menu,
        Option,
        Pagination, Radio,
        Row,
        Select,
        Submenu,
        Table,
        TableColumn,
        Tooltip,
        Upload
    } from "element-ui";
    import 'element-ui/lib/theme-chalk/index.css';

    Vue.use(Row).use(Col).use(Menu).use(Submenu).use(Button).use(Tooltip).use(Select).use(Option).use(Table).use(TableColumn).use(Pagination).use(Form).use(FormItem).use(Input).use(Upload).use(DatePicker).use(Radio).use(Loading);


    @Component
    export default class Home extends Vue {
        private mounted(){}
    }

</script>

<style lang="scss">
    .home_page {
    }
</style>