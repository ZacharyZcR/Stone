<template>
  <n-layout-footer class="footer-container" bordered>
    <n-text depth="3" class="footer-text">
      &copy; 2024 Stone ⛰️ 防火墙. 保留所有权利. 🌐
    </n-text>
  </n-layout-footer>
</template>

<script>
export default {
  name: 'FooterPage'
}
</script>

<style scoped>
.footer-container {
  text-align: center;
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.footer-text {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}
</style>
