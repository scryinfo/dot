<template>
  <div>
    <div>
      <h1>Hi Dot</h1>
      <button @click="onHi">send hi</button>
      <span v-text="hiName"></span>
    </div>
    <div>
      <h1>Write data</h1>
      <input placeholder="Please input data" v-model="data" @keyup.enter="onWrite" />
      <div v-text="dataReturn"></div>
    </div>
    <div>
      <h1>server stream</h1>
      <input
        placeholder="Please input data"
        v-model="serverStreamData"
        @keyup.enter="onServerStream"
      />
      <div v-text="serverStreamReturn"></div>
    </div>
    <div>
      <h1>client stream</h1>
      <input
        placeholder="Please input data"
        v-model="clientStreamData"
        @keyup.enter="onClientStream"
      />
      <div v-text="clientStreamReturn"></div>
    </div>
    <div>
      <h1>both sides</h1>
      <input placeholder="Please input data" v-model="bothData" @keyup.enter="onBoth" />
      <div v-text="bothReturn"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { hiService } from '@/api_impl/client'

const hiName = ref('')
const data = ref('')
const dataReturn = ref('')
const serverStreamData = ref('')
const serverStreamReturn = ref('')
const clientStreamData = ref('')
const clientStreamReturn = ref('')
const bothData = ref('')
const bothReturn = ref('')

function onHi(): boolean {
  hiService.hi({ name: 'hi dot' }).then((res) => {
    hiName.value = res.name
  })
  return false
}

function onWrite(): boolean {
  hiService.write({ data: data.value }).then((res) => {
    dataReturn.value = res.data
  })
  return false
}
function onServerStream(): boolean {
  hiService.serverStream({ greeting: serverStreamData.value }).then((res) => {
    serverStreamReturn.value = res.reply
  })
  return false
}
function onClientStream(): boolean {
  hiService.clientStream({ greeting: clientStreamData.value }).then((res) => {
    clientStreamReturn.value = res.greeting
  })
  return false
}
function onBoth(): boolean {
  hiService.bothStream({ greeting: bothData.value }).then((res) => {
    bothReturn.value = res.greeting
  })
  return false
}
</script>
