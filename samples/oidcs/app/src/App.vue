<script setup lang="ts">
import { config } from './config'
import { authService } from './api_impl/client'
import { LoginRequest } from './api/v1/auth_pb'

async function loginApi() {
  try {
    const res = await authService.login({ username: '', password: '' })
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
      <a :href="config.OIDC_LOGIN">Login</a>
    </div>
    <div>
      <a :href="config.OIDC_LOGOUT">Logout</a>
    </div>
    <div></div>
  </div>
</template>
