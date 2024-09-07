<template>
  <div class="bg-gray-900 flex items-center justify-center min-h-screen">
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-md text-center transform transition-all duration-700 ease-in-out hover:scale-105 opacity-0 translate-x-full animate-fade-in-right">
      <h2 class="text-2xl font-bold mb-6 text-white">Google Authenticator 神奇设置 🧙‍♂️✨</h2>

      <div v-if="interfaceClosed">
        <p class="text-red-500 mb-4">哎呀！注册通道已经关闭啦 🚪🔒</p>
        <p class="text-yellow-300">别灰心，可能是在进行魔法维护呢！ 🧹✨</p>
      </div>

      <div v-else>
        <div v-if="qrCodeUrl" class="mb-6">
          <img :src="qrCodeUrl" alt="Google Authenticator QR Code" class="mx-auto">
          <p class="text-green-400 mt-2">瞧！你的专属魔法二维码 🎭✨</p>
        </div>
        <p v-else-if="loading" class="text-white mb-4">正在召唤神奇的二维码... 🔮</p>
        <p v-else class="text-white mb-4">准备好开启你的魔法之旅了吗？点击下方按钮！ 🚀</p>

        <p class="text-gray-300 mb-4">
          用你的 Google Authenticator 魔法棒（呃，我是说 App）扫描这个二维码，开启双重保护咒语！ 🧙‍♀️📱
        </p>

        <div class="flex flex-col space-y-4 mt-6">
          <button
              @click="fetchQRCode"
              class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transform hover:scale-105 transition duration-300 border-2 border-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
              :disabled="loading"
          >
            {{ qrCodeUrl ? '刷新魔法二维码' : '召唤二维码' }} 🔄✨
          </button>

          <button
              @click="goToLogin"
              class="bg-green-500 text-white px-6 py-3 rounded-lg hover:bg-green-700 transform hover:scale-105 transition duration-300 border-2 border-green-400 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50"
          >
            前往登录页面 🚀
          </button>
        </div>

      </div>
    </div>

    <!-- 弹窗通知 -->
    <PopupNotification
        v-if="showNotification"
        :message="notificationMessage"
        :emoji="notificationEmoji"
        :type="notificationType"
        @close="showNotification = false"
    />
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api/axiosInstance';
import PopupNotification from './PopupNotification.vue';

export default {
  name: 'GoogleAuthQRCode',
  components: {
    PopupNotification
  },
  setup() {
    const router = useRouter();
    const qrCodeUrl = ref('');
    const loading = ref(false);
    const interfaceClosed = ref(false);
    const showNotification = ref(false);
    const notificationMessage = ref('');
    const notificationEmoji = ref('');
    const notificationType = ref('success');

    const fetchQRCode = async () => {
      loading.value = true;
      try {
        const response = await api.get('/auth/qrcode', { responseType: 'blob' });
        qrCodeUrl.value = URL.createObjectURL(response.data);
        interfaceClosed.value = false;
        showNotification.value = true;
        notificationMessage.value = '魔法二维码召唤成功！';
        notificationEmoji.value = '🎉';
        notificationType.value = 'success';
      } catch (error) {
        console.error('获取二维码失败:', error);
        interfaceClosed.value = true;
        showNotification.value = true;
        notificationMessage.value = '哎呀，魔法通道暂时关闭了！';
        notificationEmoji.value = '🔒';
        notificationType.value = 'error';

        // 如果错误响应是 Blob 类型，需要读取其内容
        if (error.response && error.response.data instanceof Blob) {
          const text = await error.response.data.text();
          const errorData = JSON.parse(text);
          if (errorData.error === "二维码接口已关闭") {
            notificationMessage.value = '二维码接口已关闭，请稍后再试！';
          }
        }
      } finally {
        loading.value = false;
      }
    };

    const goToLogin = () => {
      router.push('/login');
    };

    return {
      qrCodeUrl,
      loading,
      interfaceClosed,
      fetchQRCode,
      showNotification,
      notificationMessage,
      notificationEmoji,
      notificationType,
      goToLogin
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

/* 可以添加一些额外的按钮样式 */
button {
  font-weight: bold;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.1);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

button:hover {
  box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
}
</style>
