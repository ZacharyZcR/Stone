<template>
  <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
    <h2 class="text-2xl font-bold mb-4">攻击者信息 🕵️</h2>
    <div class="flex items-center">
      <input v-model="localIpAddress" class="bg-gray-700 text-white px-4 py-2 rounded-l-lg focus:outline-none" placeholder="输入IP地址">
      <button @click="fetchProfile" class="bg-blue-500 text-white px-4 py-2 rounded-r-lg hover:bg-blue-600">分析</button>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue';

export default {
  name: 'AttackerInfo',
  props: {
    ipAddress: String
  },
  emits: ['fetch-profile'],
  setup(props, { emit }) {
    const localIpAddress = ref(props.ipAddress);

    watch(() => props.ipAddress, (newValue) => {
      localIpAddress.value = newValue;
    });

    const fetchProfile = () => {
      emit('fetch-profile', localIpAddress.value);
    };

    return {
      localIpAddress,
      fetchProfile
    };
  }
};
</script>
