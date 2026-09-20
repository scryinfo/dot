<script setup lang="ts" vapor>
import { ref } from 'vue'

const submitting = ref(false)
function onLogin(event: SubmitEvent) {
  if (submitting.value) {
    return
  }
  const form = event.currentTarget as HTMLFormElement
  if (!form.hasChildNodes()) {
    return
  }
  form.action = window.location.href
  form.method = 'post'
  try {
    submitting.value = true
    form.submit()
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <form class="w-80 flex flex-col items-center justify-center" @submit.prevent="onLogin">
      <input required id="username" class="w-full mb-4 px-3 py-2" placeholder="user name" />
      <input required id="password" class="w-full mb-4 px-3 py-2" type="password" placeholder="password" />
      <button class="w-full px-3 py-2" type="submit" :disabled="submitting">Login</button>
    </form>
  </div>
</template>
