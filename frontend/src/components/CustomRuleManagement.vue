<template>
  <div style="padding: 24px; min-height: calc(100vh - 64px);">
    <!-- 第一行：规则管理标题 -->
    <n-card size="large" style="margin-bottom: 24px;">
      <n-space direction="vertical" size="medium" align="center">
        <n-avatar size="huge" color="#18a058">
          <n-icon size="48">
            <shield-outline />
          </n-icon>
        </n-avatar>
        
        <n-space direction="vertical" size="small" align="center">
          <n-h1 style="margin: 0; font-size: 2.5rem;">
            <n-gradient-text type="success">规则管理</n-gradient-text>
          </n-h1>
          <n-text depth="3" style="font-size: 16px;">
            创建和管理WAF防护规则，构建多层安全防线
          </n-text>
        </n-space>
      </n-space>
    </n-card>

    <!-- 第二行：规则统计概览 -->
    <n-card title="规则统计概览" size="large" style="margin-bottom: 24px;">
      <template #header-extra>
        <n-button type="primary" @click="fetchRules" :loading="refreshing">
          <template #icon>
            <n-icon><refresh-outline /></n-icon>
          </template>
          刷新数据
        </n-button>
      </template>
      
      <n-grid :cols="4" :x-gap="24" responsive="screen">
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="规则总数" :value="totalCount" suffix="条">
              <template #prefix>
                <n-icon color="#18a058" size="24">
                  <list-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="启用规则" :value="activeRulesCount" suffix="条">
              <template #prefix>
                <n-icon color="#2080f0" size="24">
                  <checkmark-circle-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="高危规则" :value="highRiskRulesCount" suffix="条">
              <template #prefix>
                <n-icon color="#d03050" size="24">
                  <warning-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="今日触发" :value="todayTriggered" suffix="次">
              <template #prefix>
                <n-icon color="#f0a020" size="24">
                  <flash-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-card>

    <!-- 第三行：快速操作 -->
    <n-card title="快速操作" size="large" style="margin-bottom: 24px;">
      <n-space size="large" justify="center">
        <n-button-group size="large">
          <n-button type="primary" @click="showAddRuleModal = true">
            <template #icon>
              <n-icon size="20"><add-outline /></n-icon>
            </template>
            添加新规则
          </n-button>
          
          <n-button @click="importRules">
            <template #icon>
              <n-icon size="20"><cloud-upload-outline /></n-icon>
            </template>
            导入规则
          </n-button>
          
          <n-button @click="exportRules">
            <template #icon>
              <n-icon size="20"><cloud-download-outline /></n-icon>
            </template>
            导出规则
          </n-button>
          
          <n-button @click="batchOperation">
            <template #icon>
              <n-icon size="20"><options-outline /></n-icon>
            </template>
            批量操作
          </n-button>
        </n-button-group>
      </n-space>
    </n-card>

    <!-- 第四行：规则列表 -->
    <n-card title="防护规则列表" size="large" style="margin-bottom: 24px;">
      <template #header-extra>
        <n-space>
          <n-input
            v-model:value="searchQuery"
            placeholder="搜索规则名称或描述"
            size="medium"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <n-icon><search-outline /></n-icon>
            </template>
          </n-input>
          
          <n-select
            v-model:value="statusFilter"
            placeholder="状态筛选"
            size="medium"
            style="width: 120px;"
            @update:value="handleFilter"
          >
            <n-option value="" label="全部" />
            <n-option value="active" label="启用" />
            <n-option value="inactive" label="禁用" />
          </n-select>
          
          <n-select
            v-model:value="typeFilter"
            placeholder="类型筛选"
            size="medium"
            style="width: 140px;"
            @update:value="handleFilter"
          >
            <n-option value="" label="全部类型" />
            <n-option value="sql_injection" label="SQL注入" />
            <n-option value="xss" label="XSS攻击" />
            <n-option value="path_traversal" label="路径遍历" />
            <n-option value="command_injection" label="命令注入" />
          </n-select>
        </n-space>
      </template>
      
      <!-- 规则表格 -->
      <n-data-table
        :columns="columns"
        :data="displayRules"
        :loading="loading"
        :pagination="{
          page: currentPage,
          pageSize: pageSize,
          itemCount: totalCount,
          showSizePicker: true,
          pageSizes: [10, 20, 50],
          onUpdatePage: changePage,
          onUpdatePageSize: changePageSize
        }"
        :row-key="row => row.id"
        size="large"
      >
        <template #empty>
          <n-empty description="暂无规则数据">
            <template #icon>
              <n-icon size="48" color="#d0d0d0">
                <shield-outline />
              </n-icon>
            </template>
            <template #extra>
              <n-button type="primary" @click="showAddRuleModal = true">
                添加第一个规则
              </n-button>
            </template>
          </n-empty>
        </template>
      </n-data-table>
    </n-card>

    <!-- 添加规则模态框 -->
    <n-modal v-model:show="showAddRuleModal" preset="card" title="添加新防护规则" size="large" style="width: 800px;">
      <n-form
        ref="formRef"
        :model="newRule"
        :rules="formRules"
        label-placement="left"
        label-width="auto"
        size="large"
      >
        <n-grid :cols="2" :x-gap="24">
          <n-grid-item>
            <n-form-item label="规则名称" path="name">
              <n-input v-model:value="newRule.name" placeholder="请输入规则名称" />
            </n-form-item>
          </n-grid-item>
          
          <n-grid-item>
            <n-form-item label="规则类型" path="type">
              <n-select v-model:value="newRule.type" placeholder="选择规则类型">
                <n-option value="sql_injection" label="SQL注入" />
                <n-option value="xss" label="XSS攻击" />
                <n-option value="path_traversal" label="路径遍历" />
                <n-option value="command_injection" label="命令注入" />
                <n-option value="file_upload" label="文件上传" />
                <n-option value="custom" label="自定义" />
              </n-select>
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <n-form-item label="规则描述" path="description">
          <n-input
            v-model:value="newRule.description"
            type="textarea"
            placeholder="请输入规则描述"
            :autosize="{ minRows: 2, maxRows: 4 }"
          />
        </n-form-item>

        <n-form-item label="匹配模式" path="pattern">
          <n-input
            v-model:value="newRule.pattern"
            placeholder="请输入正则表达式或匹配模式"
          />
        </n-form-item>

        <n-grid :cols="3" :x-gap="24">
          <n-grid-item>
            <n-form-item label="威胁等级" path="risk_level">
              <n-select v-model:value="newRule.risk_level" placeholder="选择威胁等级">
                <n-option value="low" label="低危" />
                <n-option value="medium" label="中危" />
                <n-option value="high" label="高危" />
                <n-option value="critical" label="极危" />
              </n-select>
            </n-form-item>
          </n-grid-item>
          
          <n-grid-item>
            <n-form-item label="处理动作" path="action">
              <n-select v-model:value="newRule.action" placeholder="选择处理动作">
                <n-option value="block" label="阻断" />
                <n-option value="monitor" label="监控" />
                <n-option value="redirect" label="重定向" />
              </n-select>
            </n-form-item>
          </n-grid-item>
          
          <n-grid-item>
            <n-form-item label="规则状态">
              <n-switch v-model:value="newRule.enabled">
                <template #checked>启用</template>
                <template #unchecked>禁用</template>
              </n-switch>
            </n-form-item>
          </n-grid-item>
        </n-grid>
      </n-form>
      
      <template #footer>
        <n-space justify="end">
          <n-button @click="showAddRuleModal = false">取消</n-button>
          <n-button type="primary" @click="handleAddRule" :loading="submitting">
            添加规则
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script>
import { ref, computed, onMounted, h } from 'vue';
import { useMessage } from 'naive-ui';
import api from '../api/axiosInstance';
import {
  ShieldOutline,
  RefreshOutline,
  ListOutline,
  CheckmarkCircleOutline,
  WarningOutline,
  FlashOutline,
  AddOutline,
  CloudUploadOutline,
  CloudDownloadOutline,
  OptionsOutline,
  SearchOutline,
  CreateOutline,
  TrashOutline,
  ToggleOutline,
  EyeOutline
} from '@vicons/ionicons5';

export default {
  name: 'CustomRuleManagement',
  components: {
    ShieldOutline,
    RefreshOutline,
    ListOutline,
    CheckmarkCircleOutline,
    WarningOutline,
    FlashOutline,
    AddOutline,
    CloudUploadOutline,
    CloudDownloadOutline,
    OptionsOutline,
    SearchOutline,
    CreateOutline,
    TrashOutline,
    ToggleOutline,
    EyeOutline
  },
  setup() {
    const message = useMessage();
    const rules = ref([]);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalCount = ref(24);
    const loading = ref(false);
    const refreshing = ref(false);
    
    // 搜索和筛选
    const searchQuery = ref('');
    const statusFilter = ref('');
    const typeFilter = ref('');
    
    // 添加规则模态框
    const showAddRuleModal = ref(false);
    const submitting = ref(false);
    const formRef = ref(null);
    
    // 新规则表单数据
    const newRule = ref({
      name: '',
      type: '',
      description: '',
      pattern: '',
      risk_level: 'medium',
      action: 'block',
      enabled: true
    });

    // 表单验证规则
    const formRules = {
      name: {
        required: true,
        message: '请输入规则名称',
        trigger: 'blur'
      },
      type: {
        required: true,
        message: '请选择规则类型',
        trigger: 'change'
      },
      description: {
        required: true,
        message: '请输入规则描述',
        trigger: 'blur'
      },
      pattern: {
        required: true,
        message: '请输入匹配模式',
        trigger: 'blur'
      }
    };

    // 模拟规则数据

    // 计算属性
    const activeRulesCount = computed(() => 
      rules.value.filter(rule => rule.enabled).length
    );
    
    const highRiskRulesCount = computed(() => 
      rules.value.filter(rule => rule.risk_level === 'high' || rule.risk_level === 'critical').length
    );
    
    const todayTriggered = computed(() => 
      rules.value.reduce((sum, rule) => sum + (rule.triggered_count || 0), 0)
    );

    const displayRules = computed(() => {
      let filtered = rules.value;
      
      // 搜索过滤
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase();
        filtered = filtered.filter(rule => 
          rule.name.toLowerCase().includes(query) ||
          rule.description.toLowerCase().includes(query)
        );
      }
      
      // 状态过滤
      if (statusFilter.value) {
        if (statusFilter.value === 'active') {
          filtered = filtered.filter(rule => rule.enabled);
        } else if (statusFilter.value === 'inactive') {
          filtered = filtered.filter(rule => !rule.enabled);
        }
      }
      
      // 类型过滤
      if (typeFilter.value) {
        filtered = filtered.filter(rule => rule.type === typeFilter.value);
      }
      
      return filtered;
    });

    // 表格列定义
    const columns = [
      {
        title: '规则名称',
        key: 'name',
        width: 200,
        ellipsis: {
          tooltip: true
        }
      },
      {
        title: '类型',
        key: 'type',
        width: 120,
        render: (row) => {
          const typeMap = {
            'sql_injection': 'SQL注入',
            'xss': 'XSS攻击',
            'path_traversal': '路径遍历',
            'command_injection': '命令注入',
            'file_upload': '文件上传',
            'custom': '自定义'
          };
          return typeMap[row.type] || row.type;
        }
      },
      {
        title: '威胁等级',
        key: 'risk_level',
        width: 100,
        render: (row) => {
          const levelMap = {
            'low': { type: 'success', text: '低危' },
            'medium': { type: 'warning', text: '中危' },
            'high': { type: 'error', text: '高危' },
            'critical': { type: 'error', text: '极危' }
          };
          const level = levelMap[row.risk_level] || { type: 'default', text: '未知' };
          return h('n-tag', { type: level.type, size: 'small' }, level.text);
        }
      },
      {
        title: '状态',
        key: 'enabled',
        width: 80,
        render: (row) => {
          return h('n-tag', { 
            type: row.enabled ? 'success' : 'default',
            size: 'small'
          }, row.enabled ? '启用' : '禁用');
        }
      },
      {
        title: '触发次数',
        key: 'triggered_count',
        width: 100,
        render: (row) => row.triggered_count || 0
      },
      {
        title: '创建时间',
        key: 'created_at',
        width: 120
      },
      {
        title: '操作',
        key: 'actions',
        width: 200,
        render: (row) => {
          return h('n-space', [
            h('n-button', {
              size: 'small',
              onClick: () => editRule(row)
            }, {
              default: () => '编辑',
              icon: () => h('n-icon', null, { default: () => h(CreateOutline) })
            }),
            h('n-button', {
              size: 'small',
              type: row.enabled ? 'warning' : 'success',
              onClick: () => toggleRule(row)
            }, {
              default: () => row.enabled ? '禁用' : '启用',
              icon: () => h('n-icon', null, { default: () => h(ToggleOutline) })
            }),
            h('n-button', {
              size: 'small',
              type: 'error',
              onClick: () => deleteRule(row)
            }, {
              default: () => '删除',
              icon: () => h('n-icon', null, { default: () => h(TrashOutline) })
            })
          ]);
        }
      }
    ];

    // 方法
    const fetchRules = async () => {
      refreshing.value = true;
      loading.value = true;
      try {
        const response = await api.get('/interception-rules', {
          params: { page: currentPage.value, pageSize: pageSize.value }
        });
        rules.value = response.data.rules || [];
        totalCount.value = response.data.totalCount || rules.value.length;
      } catch (error) {
        console.error('获取规则失败:', error);
        message.error('获取规则失败，请稍后重试');
        rules.value = [];
        totalCount.value = 0;
      } finally {
        loading.value = false;
        refreshing.value = false;
      }
    };

    const handleAddRule = async () => {
      if (!formRef.value) return;
      
      try {
        await formRef.value.validate();
        submitting.value = true;
        
        await api.post('/interception-rules', newRule.value);
        await fetchRules();
        message.success('规则添加成功');
        showAddRuleModal.value = false;
        resetForm();
      } catch (error) {
        if (error.length) return; // 表单验证失败
        console.error('添加规则失败:', error);
        message.error('添加规则失败');
      } finally {
        submitting.value = false;
      }
    };

    const editRule = (rule) => {
      message.info(`编辑规则: ${rule.name}`);
    };

    const toggleRule = async (rule) => {
      try {
        const newStatus = !rule.enabled;
        await api.patch(`/interception-rules/${rule.id}`, { enabled: newStatus });
        rule.enabled = newStatus;
        message.success(`规则已${newStatus ? '启用' : '禁用'}`);
      } catch (error) {
        console.error('切换规则状态失败:', error);
        message.error('操作失败');
      }
    };

    const deleteRule = async (rule) => {
      try {
        await api.delete(`/interception-rules/${rule.id}`);
        await fetchRules();
        message.success(`规则 "${rule.name}" 已删除`);
      } catch (error) {
        console.error('删除规则失败:', error);
        message.error(`删除规则失败`);
      }
    };

    const resetForm = () => {
      newRule.value = {
        name: '',
        type: '',
        description: '',
        pattern: '',
        risk_level: 'medium',
        action: 'block',
        enabled: true
      };
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

    const handleSearch = () => {
      // 搜索逻辑在计算属性中处理
    };

    const handleFilter = () => {
      // 过滤逻辑在计算属性中处理
    };

    const importRules = () => {
      message.info('规则导入功能开发中');
    };

    const exportRules = () => {
      message.info('规则导出功能开发中');
    };

    const batchOperation = () => {
      message.info('批量操作功能开发中');
    };

    onMounted(() => {
      fetchRules();
    });

    return {
      // 数据
      rules,
      displayRules,
      currentPage,
      pageSize,
      totalCount,
      loading,
      refreshing,
      searchQuery,
      statusFilter,
      typeFilter,
      showAddRuleModal,
      submitting,
      formRef,
      newRule,
      formRules,
      columns,
      
      // 计算属性
      activeRulesCount,
      highRiskRulesCount,
      todayTriggered,
      
      // 方法
      fetchRules,
      handleAddRule,
      editRule,
      toggleRule,
      deleteRule,
      changePage,
      changePageSize,
      handleSearch,
      handleFilter,
      importRules,
      exportRules,
      batchOperation
    };
  }
};
</script>

<style scoped>
</style>
