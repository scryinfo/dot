<template>
    <div class="container-fluid">
        <router-view></router-view>
    </div>
</template>

<script lang="ts">
    import './index.scss';
    import {Component, Vue} from 'vue-property-decorator';
    @Component
    export default class App extends Vue {
    }
</script>
