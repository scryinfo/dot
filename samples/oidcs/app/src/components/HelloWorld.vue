<script setup lang="ts">
import { ref } from 'vue'

const count = ref(0)
</script>

<template>
  <div>
    <button type="button" class="counter" @click="count++">
      Count is {{ count }}
    </button>
  </div>
</template>
