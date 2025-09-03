<template>
  <div class="bg-gray-800 p-6 rounded-lg shadow-md transform transition-all duration-500 hover:shadow-2xl">
    <h2 class="text-2xl font-bold mb-4">添加新规则 ➕</h2>
    <form @submit.prevent="submitRule">
      <div class="mb-4">
        <label class="block text-gray-300 text-sm font-bold mb-2" for="ruleName">规则名称</label>
        <input v-model="newRule.name" class="input-field" id="ruleName" type="text" placeholder="输入规则名称">
      </div>
      <div class="mb-4">
        <label class="block text-gray-300 text-sm font-bold mb-2" for="ruleRegex">规则正则</label>
        <input v-model="newRule.regex" class="input-field" id="ruleRegex" type="text" placeholder="输入规则正则">
      </div>
      <div class="mb-4">
        <label class="block text-gray-300 text-sm font-bold mb-2" for="ruleMethod">HTTP 方法</label>
        <select v-model="newRule.method" class="input-field" id="ruleMethod">
          <option value="GET">GET</option>
          <option value="POST">POST</option>
        </select>
      </div>
      <div class="flex items-center justify-between">
        <button class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-700 transform hover:scale-105 transition duration-300" type="submit">
          添加规则 ➕
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  name: 'AddRuleForm',
  emits: ['add-rule'],
  setup(props, { emit }) {
    const newRule = ref({ name: '', regex: '', method: 'GET' });

    const submitRule = () => {
      emit('add-rule', { ...newRule.value });
      newRule.value = { name: '', regex: '', method: 'GET' };
    };

    return { newRule, submitRule };
  }
};
</script>

<style scoped>
.input-field {
  @apply shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-300;
}
</style>
