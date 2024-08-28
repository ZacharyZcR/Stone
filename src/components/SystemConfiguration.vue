<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 配置表单 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md transform transition-all duration-500 hover:shadow-2xl animate-fade-in-up mb-8">
        <h2 class="text-2xl font-bold mb-4">系统配置设置 ⚙️</h2>
        <form>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2 hover:text-teal-400 transition duration-300" for="firewallStatus">防火墙状态</label>
            <select class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-teal-500 transition duration-300" id="firewallStatus">
              <option>启用 ✅</option>
              <option>禁用 ❌</option>
            </select>
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2 hover:text-teal-400 transition duration-300" for="logLevel">日志级别</label>
            <select class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-teal-500 transition duration-300" id="logLevel">
              <option>详细 📋</option>
              <option>普通 📄</option>
              <option>简略 🗒️</option>
            </select>
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2 hover:text-teal-400 transition duration-300" for="alertEmail">警报邮箱</label>
            <input class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-teal-500 transition duration-300" id="alertEmail" type="email" placeholder="输入警报接收邮箱 📧">
          </div>
          <div class="flex items-center justify-between">
            <button class="bg-teal-500 text-white px-4 py-2 rounded hover:bg-teal-700 transform hover:scale-105 transition duration-300" type="button">
              保存配置 💾
            </button>
          </div>
        </form>
      </div>

      <!-- 黑名单和白名单 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- 黑名单 IP 列表 -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-md">
          <h2 class="text-2xl font-bold mb-4">黑名单 IP 列表 🚫</h2>
          <table class="min-w-full bg-gray-800">
            <thead>
            <tr>
              <th class="py-2 px-4 border-b-2 border-gray-700">IP 地址</th>
              <th class="py-2 px-4 border-b-2 border-gray-700">操作</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="ip in blacklist" :key="ip" class="hover:bg-gray-700 transition duration-300">
              <td class="py-2 px-4 border-b border-gray-700">{{ ip }}</td>
              <td class="py-2 px-4 border-b border-gray-700">
                <button class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700">移除 🗑️</button>
              </td>
            </tr>
            </tbody>
          </table>
        </div>

        <!-- 白名单 IP 列表 -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-md">
          <h2 class="text-2xl font-bold mb-4">白名单 IP 列表 ✅</h2>
          <table class="min-w-full bg-gray-800">
            <thead>
            <tr>
              <th class="py-2 px-4 border-b-2 border-gray-700">IP 地址</th>
              <th class="py-2 px-4 border-b-2 border-gray-700">操作</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="ip in whitelist" :key="ip" class="hover:bg-gray-700 transition duration-300">
              <td class="py-2 px-4 border-b border-gray-700">{{ ip }}</td>
              <td class="py-2 px-4 border-b border-gray-700">
                <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">移除 🗑️</button>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <FooterPage />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import HeaderPage from './HeaderPage.vue';
import FooterPage from './FooterPage.vue';

export default {
  name: 'SystemConfiguration',
  components: {
    HeaderPage,
    FooterPage
  },
  setup() {
    const whitelist = ref([]);
    const blacklist = ref([]);

    const fetchData = async () => {
      try {
        const response = await axios.get('http://172.20.2.226:8081/ip-control-rules');
        const data = response.data;
        whitelist.value = data.Whitelist;
        blacklist.value = data.Blacklist;
      } catch (error) {
        console.error('请求失败:', error);
      }
    };

    onMounted(() => {
      fetchData();
    });

    return { whitelist, blacklist };
  }
};
</script>
