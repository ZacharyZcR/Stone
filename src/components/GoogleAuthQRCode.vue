<template>
  <div class="bg-gray-900 flex items-center justify-center min-h-screen">
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-md text-center transform transition-all duration-700 ease-in-out hover:scale-105 opacity-0 translate-x-full animate-fade-in-right">
      <h2 class="text-2xl font-bold mb-6 text-white">Google Authenticator 设置 🔐</h2>
      <div v-if="qrCodeUrl" class="mb-6">
        <img :src="qrCodeUrl" alt="Google Authenticator QR Code" class="mx-auto">
      </div>
      <p v-else class="text-white mb-4">正在加载二维码...</p>
      <p class="text-gray-300 mb-4">
        请使用 Google Authenticator 应用扫描此二维码以设置双因素认证。
      </p>
      <button
          @click="refreshQRCode"
          class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700 transform hover:scale-105 transition duration-300"
      >
        刷新二维码 🔄
      </button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';

export default {
  name: 'GoogleAuthQRCode',
  setup() {
    const qrCodeUrl = ref('');

    const fetchQRCode = async () => {
      try {
        const response = await axios.get('http://172.20.2.226:8081/auth/qrcode', { responseType: 'blob' });
        qrCodeUrl.value = URL.createObjectURL(response.data);
      } catch (error) {
        console.error('获取二维码失败:', error);
      }
    };

    const refreshQRCode = () => {
      qrCodeUrl.value = '';
      fetchQRCode();
    };

    onMounted(() => {
      fetchQRCode();
    });

    return {
      qrCodeUrl,
      refreshQRCode
    };
  }
};
</script>

<style scoped>
@keyframes fade-in-right {
  0% {
    opacity: 0;
    transform: translateX(100%);
  }
  100% {
    opacity: 1;
    transform: translateX(0);
  }
}

.animate-fade-in-right {
  animation: fade-in-right 1s forwards;
}
</style>
