<script setup lang="ts">
import { config } from './config';

</script>

<template>
  <div style="display: flex;flex-direction: column; gap: 5px;">
      <div>
          <a :href="config + '/login'">Login</a>
      </div>
      <div>
          <a :href="config + '/logout'">Logout</a>
      </div>
      <div>

      </div>
  </div>
</template>
