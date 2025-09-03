<template>
  <n-config-provider :theme="theme">
    <n-global-style />
    <n-message-provider>
      <div id="app">
        <router-view />
      </div>
    </n-message-provider>
  </n-config-provider>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'App',
  setup() {
    const theme = ref(null) // 使用浅色主题，如果要暗色主题则设为 darkTheme
    
    return {
      theme
    }
  }
}
</script>
