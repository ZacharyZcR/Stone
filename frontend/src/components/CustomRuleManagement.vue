<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <HeaderPage />

    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <RuleList
          :rules="rules"
          @delete-rule="deleteRule"
      />

      <LogPagination
          class="mt-8 mb-8"
      :current-page="currentPage"
      :total-pages="totalPages"
      :total-count="totalCount"
      :page-size="pageSize"
      @page-changed="changePage"
      @page-size-changed="changePageSize"
      />

      <AddRuleForm
          class="mt-8"
      @add-rule="addRule"
      />
    </div>

    <FooterPage />

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
import { ref, computed, onMounted } from 'vue';
import api from '../api/axiosInstance';
import HeaderPage from './HeaderPage.vue';
import FooterPage from './FooterPage.vue';
import PopupNotification from './PopupNotification.vue';
import LogPagination from './LogPagination.vue';
import RuleList from './RuleList.vue';
import AddRuleForm from './AddRuleForm.vue';

export default {
  name: 'CustomRuleManagement',
  components: {
    HeaderPage,
    FooterPage,
    PopupNotification,
    LogPagination,
    RuleList,
    AddRuleForm
  },
  setup() {
    const rules = ref([]);
    const showNotification = ref(false);
    const notificationMessage = ref('');
    const notificationEmoji = ref('');
    const notificationType = ref('success');
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalCount = ref(0);

    const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value));

    const showPopup = (message, emoji, type) => {
      notificationMessage.value = message;
      notificationEmoji.value = emoji;
      notificationType.value = type;
      showNotification.value = true;
    };

    const fetchRules = async () => {
      try {
        const response = await api.get('/interception-rules', {
          params: { page: currentPage.value, pageSize: pageSize.value }
        });
        rules.value = response.data.rules || [];
        totalCount.value = response.data.totalCount;
      } catch (error) {
        console.error('获取规则失败:', error);
        showPopup('获取规则失败', '❌', 'error');
      }
    };

    const addRule = async (newRule) => {
      try {
        await api.post('/interception-rules', newRule);
        await fetchRules();
        showPopup('规则添加成功', '✅', 'success');
      } catch (error) {
        console.error('添加规则失败:', error);
        showPopup('添加规则失败', '❌', 'error');
      }
    };

    const deleteRule = async (rule) => {
      try {
        await api.delete(`/interception-rules/${rule.name}`);
        totalCount.value -= 1;
        if (rules.value.length === 1 && currentPage.value > 1) {
          currentPage.value -= 1;
        }
        await fetchRules();
        showPopup(`规则 "${rule.name}" 已删除`, '🗑️', 'success');
      } catch (error) {
        console.error('删除规则失败:', error);
        showPopup(`删除规则 "${rule.name}" 失败`, '❌', 'error');
      }
    };

    const changePage = (page) => {
      currentPage.value = page;
      fetchRules();
    };

    const changePageSize = (size) => {
      pageSize.value = size;
      currentPage.value = 1;
      fetchRules();
    };

    onMounted(fetchRules);

    return {
      rules,
      addRule,
      deleteRule,
      showNotification,
      notificationMessage,
      notificationEmoji,
      notificationType,
      currentPage,
      pageSize,
      totalCount,
      totalPages,
      changePage,
      changePageSize
    };
  }
};
</script>
