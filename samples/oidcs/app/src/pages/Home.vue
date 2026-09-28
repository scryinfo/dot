<script setup lang="ts">
import { config } from '@/config'
import { authService } from '@/api_impl/client'
import { ref } from 'vue'

const redirectUrl = ref('')
;(() => {
  const t = new URL(config.OIDC_LOGIN)
  t.searchParams.set('_redirect_uri', window.location.href + 'logined')
  redirectUrl.value = t.toString()
})()

async function loginApi() {
  try {
    const res = await authService.login({})
    console.log(res)
  } catch (error) {
    console.error(error)
  }
}
</script>

<template>
  <div style="display: flex; flex-direction: column; gap: 5px">
    <div>
      <a @click="loginApi">Login api</a>
    </div>
    <div>
      <a :href="redirectUrl">Login</a>
    </div>
    <div>
      <a :href="config.OIDC_LOGOUT">Logout</a>
    </div>
    <div></div>
  </div>
</template>
