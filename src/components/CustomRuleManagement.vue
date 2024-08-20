<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 规则列表 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">现有规则 📜</h2>
        <table class="min-w-full bg-gray-800">
          <thead>
          <tr>
            <th class="py-2 px-4 border-b-2 border-gray-700">规则名称</th>
            <th class="py-2 px-4 border-b-2 border-gray-700">描述</th>
            <th class="py-2 px-4 border-b-2 border-gray-700">操作</th>
          </tr>
          </thead>
          <tbody>
          <tr>
            <td class="py-2 px-4 border-b border-gray-700">禁止 SQL 注入</td>
            <td class="py-2 px-4 border-b border-gray-700">阻止所有 SQL 注入尝试</td>
            <td class="py-2 px-4 border-b border-gray-700">
              <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">编辑 ✏️</button>
              <button class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700">删除 🗑️</button>
            </td>
          </tr>
          <tr>
            <td class="py-2 px-4 border-b border-gray-700">禁止 XSS 攻击</td>
            <td class="py-2 px-4 border-b border-gray-700">阻止所有 XSS 攻击尝试</td>
            <td class="py-2 px-4 border-b border-gray-700">
              <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">编辑 ✏️</button>
              <button class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700">删除 🗑️</button>
            </td>
          </tr>
          </tbody>
        </table>
      </div>

      <!-- 添加规则 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold mb-4">添加新规则 ➕</h2>
        <form>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="ruleName">规则名称</label>
            <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="ruleName" type="text" placeholder="输入规则名称">
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="description">描述</label>
            <textarea class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="description" placeholder="输入规则描述"></textarea>
          </div>
          <div class="flex items-center justify-between">
            <button class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-700" type="button">
              添加规则 ➕
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
  name: 'CustomRuleManagement',
  components: {
    HeaderPage,
    FooterPage
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');
</style>
