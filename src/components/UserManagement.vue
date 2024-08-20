<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 用户列表 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">用户列表 📋</h2>
        <table class="min-w-full bg-gray-800">
          <thead>
          <tr>
            <th class="py-2 px-4 border-b-2 border-gray-700">用户名</th>
            <th class="py-2 px-4 border-b-2 border-gray-700">邮箱</th>
            <th class="py-2 px-4 border-b-2 border-gray-700">角色</th>
            <th class="py-2 px-4 border-b-2 border-gray-700">操作</th>
          </tr>
          </thead>
          <tbody>
          <tr>
            <td class="py-2 px-4 border-b border-gray-700">张三</td>
            <td class="py-2 px-4 border-b border-gray-700">zhangsan@example.com</td>
            <td class="py-2 px-4 border-b border-gray-700">管理员</td>
            <td class="py-2 px-4 border-b border-gray-700">
              <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">编辑 ✏️</button>
              <button class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700">删除 🗑️</button>
            </td>
          </tr>
          <tr>
            <td class="py-2 px-4 border-b border-gray-700">李四</td>
            <td class="py-2 px-4 border-b border-gray-700">lisi@example.com</td>
            <td class="py-2 px-4 border-b border-gray-700">用户</td>
            <td class="py-2 px-4 border-b border-gray-700">
              <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">编辑 ✏️</button>
              <button class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700">删除 🗑️</button>
            </td>
          </tr>
          </tbody>
        </table>
      </div>

      <!-- 添加用户 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold mb-4">添加新用户 ➕</h2>
        <form>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="username">用户名</label>
            <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="输入用户名">
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="email">邮箱</label>
            <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" type="email" placeholder="输入邮箱">
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="role">角色</label>
            <select class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="role">
              <option>管理员</option>
              <option>用户</option>
            </select>
          </div>
          <div class="flex items-center justify-between">
            <button class="bg-purple-500 text-white px-4 py-2 rounded hover:bg-purple-700" type="button">
              添加用户 ➕
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 页脚 -->
    <FooterPage />
  </div>
</template>

<script>
import HeaderPage from './HeaderPage.vue'
import FooterPage from './FooterPage.vue'

export default {
  name: 'UserManagement',
  components: {
    HeaderPage,
    FooterPage
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');
</style>
