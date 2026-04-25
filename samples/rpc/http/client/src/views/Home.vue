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
async function onServerStream() {
  const stream = hiService.serverStream({ greeting: serverStreamData.value })
  for await (const res of stream) {
    serverStreamReturn.value = serverStreamReturn.value + res.reply + '\n'
  }
}
async function onClientStream() {
  const res = await hiService.clientStream(
    (async function* () {
      yield { greeting: clientStreamData.value + 'first' }
      yield { greeting: clientStreamData.value + 'second' }
    })(),
  )
  clientStreamReturn.value = res.greeting
}
async function onBoth() {
  const stream = hiService.bothStream(
    (async function* () {
      yield { greeting: bothData.value + 'first' }
      yield { greeting: bothData.value + 'second' }
    })(),
  )
  for await (const res of stream) {
    bothReturn.value = bothReturn.value + res.greeting + '\n'
  }
}
</script>
