<template>
<div>
    <button @click="getRoot">$Root</button>
</div>
</template>

<script lang="ts">
    export default {
        name: "AddButton",
        methods: {
            getRoot() {
                console.log(this)
                console.log(this.$root)
            }
        }
    }
</script>

<style scoped>

</style>
