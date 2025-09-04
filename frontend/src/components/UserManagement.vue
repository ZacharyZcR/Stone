<template>
  <div style="padding: 24px; min-height: calc(100vh - 64px);">
    <!-- 第一行：页面标题 -->
    <n-card size="large" style="margin-bottom: 24px;">
      <n-space direction="vertical" size="medium" align="center">
        <n-avatar size="huge" color="#722ed1">
          <n-icon size="48">
            <people-outline />
          </n-icon>
        </n-avatar>
        
        <n-space direction="vertical" size="small" align="center">
          <n-h1 style="margin: 0; font-size: 2.5rem;">
            <n-gradient-text type="primary">用户管理中心</n-gradient-text>
          </n-h1>
          <n-text depth="3" style="font-size: 16px;">
            系统用户账户管理与权限控制
          </n-text>
        </n-space>
      </n-space>
    </n-card>

    <!-- 第二行：用户统计概览 -->
    <n-card title="用户统计" size="large" style="margin-bottom: 24px;">
      <template #header-extra>
        <n-button type="primary" @click="fetchUsers" :loading="refreshing">
          <template #icon>
            <n-icon><refresh-outline /></n-icon>
          </template>
          刷新数据
        </n-button>
      </template>
      
      <n-grid :cols="5" :x-gap="24" responsive="screen">
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="总用户数" :value="userStats.total" suffix="人">
              <template #prefix>
                <n-icon color="#2080f0" size="24">
                  <people-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="活跃用户" :value="userStats.active" suffix="人">
                <template #prefix>
                  <n-icon color="#18a058" size="24">
                    <checkmark-circle-outline />
                  </n-icon>
                </template>
              </n-statistic>
            </n-card>
          </n-grid-item>
          
          <n-grid-item>
            <n-card hoverable embedded>
              <n-statistic label="管理员" :value="userStats.admin" suffix="人">
                <template #prefix>
                  <n-icon color="#f0a020" size="24">
                    <shield-outline />
                  </n-icon>
                </template>
              </n-statistic>
            </n-card>
          </n-grid-item>

          <n-grid-item>
            <n-card hoverable embedded>
              <n-statistic label="今日登录" :value="userStats.todayLogin" suffix="人">
                <template #prefix>
                  <n-icon color="#d03050" size="24">
                    <log-in-outline />
                  </n-icon>
                </template>
              </n-statistic>
            </n-card>
          </n-grid-item>

          <n-grid-item>
            <n-card hoverable embedded>
              <n-statistic label="在线用户" :value="userStats.online" suffix="人">
                <template #prefix>
                  <n-icon color="#722ed1" size="24">
                    <radio-button-on-outline />
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
            <n-button type="primary" @click="showAddUserModal = true">
              <template #icon>
                <n-icon size="20"><person-add-outline /></n-icon>
              </template>
              添加用户
            </n-button>
            
            <n-button @click="exportUserData" :loading="exportLoading">
              <template #icon>
                <n-icon size="20"><download-outline /></n-icon>
              </template>
              导出数据
            </n-button>
            
            <n-button @click="showBatchModal = true" :disabled="selectedUsers.length === 0">
              <template #icon>
                <n-icon size="20"><layers-outline /></n-icon>
              </template>
              批量操作 ({{ selectedUsers.length }})
            </n-button>
          </n-button-group>
        </n-space>
    </n-card>

    <!-- 第四行：用户管理主界面 -->
    <n-grid :cols="3" :x-gap="24" responsive="screen" style="margin-bottom: 24px;">
      <!-- 左侧：用户列表 -->
      <n-grid-item :span="2">
        <n-card title="用户列表" size="large" style="height: 600px;">
            <template #header-extra>
              <n-space>
                <n-input
                  v-model:value="searchQuery"
                  placeholder="搜索用户名..."
                  clearable
                  size="small"
                >
                  <template #prefix>
                    <n-icon size="16" color="#808080">
                      <search-outline />
                    </n-icon>
                  </template>
                </n-input>
                <n-select
                  v-model:value="roleFilter"
                  placeholder="筛选角色"
                  clearable
                  size="small"
                  :options="roleFilterOptions"
                />
              </n-space>
            </template>
            
            <div style="height: 500px; overflow-y: auto;">
              <n-data-table
                v-model:checked-row-keys="selectedUsers"
                :columns="userColumns"
                :data="filteredUsers"
                :row-key="(row) => row.id"
                :pagination="{ pageSize: 15 }"
                :bordered="false"
                striped
                size="small"
                :scroll-x="900"
              />
            </div>
        </n-card>
      </n-grid-item>

      <!-- 右侧：在线用户监控 -->
      <n-grid-item>
        <n-card title="在线用户监控" size="large" style="height: 600px;">
          <div style="height: 500px; overflow-y: auto; padding-right: 8px;">
            <n-list>
              <n-list-item v-for="user in onlineUsers" :key="user.id">
                <n-space justify="space-between">
                  <n-space>
                    <n-avatar size="small" :style="{ backgroundColor: user.role === 'admin' ? '#f0a020' : '#2080f0' }">
                      <n-icon size="12">
                        <person-outline />
                      </n-icon>
                    </n-avatar>
                    <n-text>{{ user.username }}</n-text>
                  </n-space>
                  <n-space>
                    <n-tag :type="user.role === 'admin' ? 'warning' : 'info'" size="small">
                      {{ user.role === 'admin' ? '管理员' : '用户' }}
                    </n-tag>
                    <n-text depth="3" style="font-size: 12px;">{{ user.loginTime }}</n-text>
                  </n-space>
                </n-space>
              </n-list-item>
            </n-list>
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- 第五行：用户活动时间线 -->
    <n-card title="最近用户活动" size="large">
      <div style="height: 300px; overflow-y: auto; padding-right: 8px;">
        <n-timeline size="large">
          <n-timeline-item 
            v-for="(activity, index) in userActivities" 
            :key="index"
            :type="getActivityType(activity.type)"
            :title="activity.title"
            :content="activity.message"
            :time="activity.time"
          />
        </n-timeline>
      </div>
    </n-card>

    <!-- 添加用户模态框 -->
    <n-modal v-model:show="showAddUserModal">
      <n-card title="添加新用户" :bordered="false" size="huge" closable @close="showAddUserModal = false">
        <n-form ref="addUserFormRef" :model="newUser" :rules="userRules" label-placement="left" label-width="120px">
          <n-form-item label="用户名" path="username">
            <n-input v-model:value="newUser.username" placeholder="请输入用户名" />
          </n-form-item>
          <n-form-item label="密码" path="password">
            <n-input v-model:value="newUser.password" type="password" placeholder="请输入密码" show-password-on="click" />
          </n-form-item>
          <n-form-item label="确认密码" path="confirmPassword">
            <n-input v-model:value="newUser.confirmPassword" type="password" placeholder="请再次输入密码" show-password-on="click" />
          </n-form-item>
          <n-form-item label="角色" path="role">
            <n-select v-model:value="newUser.role" :options="roleOptions" placeholder="选择用户角色" />
          </n-form-item>
        </n-form>
        
        <template #footer>
          <n-space justify="end">
            <n-button @click="showAddUserModal = false">取消</n-button>
            <n-button type="primary" @click="handleAddUser" :loading="addUserLoading">确认添加</n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>

    <!-- 批量操作模态框 -->
    <n-modal v-model:show="showBatchModal">
      <n-card title="批量操作" :bordered="false" size="huge" closable @close="showBatchModal = false">
        <n-space vertical size="large">
          <n-alert type="info">
            已选择 {{ selectedUsers.length }} 个用户
          </n-alert>
          
          <n-space vertical size="medium">
            <n-button block @click="batchChangeRole">
              <template #icon>
                <n-icon><key-outline /></n-icon>
              </template>
              批量更改角色
            </n-button>
            <n-button block @click="batchResetPassword">
              <template #icon>
                <n-icon><lock-closed-outline /></n-icon>
              </template>
              批量重置密码
            </n-button>
            <n-button block type="error" @click="batchDeleteUsers">
              <template #icon>
                <n-icon><trash-outline /></n-icon>
              </template>
              批量删除用户
            </n-button>
          </n-space>
        </n-space>
      </n-card>
    </n-modal>
  </div>
</template>

<script>
import { ref, onMounted, h, computed } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import api from '../api/axiosInstance'
import {
  PeopleOutline,
  CheckmarkCircleOutline,
  ShieldOutline,
  LogInOutline,
  RadioButtonOnOutline,
  RefreshOutline,
  PersonAddOutline,
  DownloadOutline,
  LayersOutline,
  KeyOutline,
  SearchOutline,
  PersonOutline,
  LockClosedOutline,
  TrashOutline
} from '@vicons/ionicons5'

export default {
  name: 'UserManagement',
  components: {
    PeopleOutline,
    CheckmarkCircleOutline,
    ShieldOutline,
    LogInOutline,
    RadioButtonOnOutline,
    RefreshOutline,
    PersonAddOutline,
    DownloadOutline,
    LayersOutline,
    KeyOutline,
    SearchOutline,
    PersonOutline,
    LockClosedOutline,
    TrashOutline
  },
  setup() {
    const message = useMessage()
    const dialog = useDialog()
    
    // 基础数据
    const users = ref([])
    const selectedUsers = ref([])
    const refreshing = ref(false)
    const searchQuery = ref('')
    const roleFilter = ref('')
    
    // 模态框状态
    const showAddUserModal = ref(false)
    const showBatchModal = ref(false)
    const addUserLoading = ref(false)
    const exportLoading = ref(false)
    
    // 用户统计数据
    const userStats = ref({
      total: 0,
      active: 0,
      admin: 0,
      todayLogin: 0,
      online: 0
    })

    // 新用户表单数据
    const newUser = ref({
      username: '',
      password: '',
      confirmPassword: '',
      role: ''
    })

    // 在线用户
    const onlineUsers = ref([])

    // 用户活动日志
    const userActivities = ref([])

    // 角色选项
    const roleOptions = [
      { label: '管理员', value: 'admin' },
      { label: '普通用户', value: 'user' }
    ]

    // 角色筛选选项
    const roleFilterOptions = [
      { label: '全部', value: '' },
      { label: '管理员', value: 'admin' },
      { label: '普通用户', value: 'user' }
    ]

    // 表单验证规则
    const userRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度应为3-20个字符', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 8, message: '密码长度至少8个字符', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        {
          validator: (rule, value) => {
            return value === newUser.value.password
          },
          message: '两次输入的密码不一致',
          trigger: 'blur'
        }
      ],
      role: [
        { required: true, message: '请选择用户角色', trigger: 'change' }
      ]
    }

    // 过滤用户列表
    const filteredUsers = computed(() => {
      let result = users.value
      
      // 搜索过滤
      if (searchQuery.value) {
        result = result.filter(user => 
          user.username.toLowerCase().includes(searchQuery.value.toLowerCase())
        )
      }
      
      // 角色过滤
      if (roleFilter.value) {
        result = result.filter(user => user.role === roleFilter.value)
      }
      
      return result
    })

    // 用户表格列定义
    const userColumns = [
      { type: 'selection', width: 50 },
      {
        title: '用户名',
        key: 'username',
        width: 120,
        ellipsis: { tooltip: true }
      },
      {
        title: '角色',
        key: 'role',
        width: 100,
        render(row) {
          const roleMap = {
            admin: { type: 'error', label: '管理员' },
            user: { type: 'default', label: '用户' }
          }
          const role = roleMap[row.role] || { type: 'default', label: '未知' }
          return h('n-tag', { type: role.type, size: 'small' }, role.label)
        }
      },
      {
        title: '创建时间',
        key: 'created',
        width: 150,
        render(row) {
          return row.created ? new Date(row.created).toLocaleDateString('zh-CN') : '-'
        }
      },
      {
        title: '最后登录',
        key: 'last_login',
        width: 150,
        render(row) {
          return row.last_login ? new Date(row.last_login).toLocaleString('zh-CN') : '从未登录'
        }
      },
      {
        title: '状态',
        key: 'active',
        width: 80,
        render(row) {
          return h(
            'n-tag',
            { type: row.active ? 'success' : 'default', size: 'small' },
            row.active ? '活跃' : '禁用'
          )
        }
      },
      {
        title: '操作',
        key: 'actions',
        width: 150,
        render(row) {
          return h(
            'n-space',
            { size: 'small' },
            {
              default: () => [
                h(
                  'n-button',
                  {
                    type: 'primary',
                    size: 'small',
                    quaternary: true,
                    onClick: () => editUser(row)
                  },
                  '编辑'
                ),
                h(
                  'n-button',
                  {
                    type: 'error',
                    size: 'small',
                    quaternary: true,
                    onClick: () => confirmDelete(row.username),
                    disabled: row.username === 'admin'
                  },
                  '删除'
                )
              ]
            }
          )
        }
      }
    ]


    // 获取用户数据
    const fetchUsers = async () => {
      refreshing.value = true
      try {
        const response = await api.get('/users')
        users.value = response.data || []
        
        // 更新统计数据
        userStats.value = {
          total: users.value.length,
          active: users.value.filter(u => u.active).length,
          admin: users.value.filter(u => u.role === 'admin').length,
          todayLogin: users.value.filter(u => isToday(u.last_login)).length,
          online: onlineUsers.value.length
        }
      } catch (error) {
        console.error('获取用户列表失败:', error)
        message.error('获取用户列表失败')
      } finally {
        refreshing.value = false
      }
    }

    // 判断是否为今天
    const isToday = (dateStr) => {
      if (!dateStr) return false
      const today = new Date().toDateString()
      const date = new Date(dateStr).toDateString()
      return today === date
    }

    // 编辑用户
    const editUser = (user) => {
      message.warning('编辑用户功能暂未实现')
    }

    // 删除确认
    const confirmDelete = (username) => {
      dialog.warning({
        title: '删除用户确认',
        content: `您确定要删除用户 ${username} 吗？此操作不可撤销。`,
        positiveText: '确定删除',
        negativeText: '取消',
        onPositiveClick: () => deleteUser(username)
      })
    }

    // 删除用户
    const deleteUser = async (username) => {
      try {
        await api.delete(`/users/${username}`)
        await fetchUsers()
        message.success(`成功删除用户 ${username}`)
      } catch (error) {
        console.error(`删除用户失败: ${username}`, error)
        message.error(`删除用户 ${username} 失败`)
      }
    }

    // 添加用户
    const handleAddUser = async () => {
      addUserLoading.value = true
      try {
        await api.post('/users', {
          username: newUser.value.username,
          password: newUser.value.password,
          role: newUser.value.role
        })
        
        // 重置表单
        newUser.value = {
          username: '',
          password: '',
          confirmPassword: '',
          role: ''
        }
        
        showAddUserModal.value = false
        message.success('用户添加成功')
        await fetchUsers()
      } catch (error) {
        console.error('添加用户失败:', error)
        message.error('添加用户失败')
      } finally {
        addUserLoading.value = false
      }
    }

    // 导出用户数据
    const exportUserData = () => {
      exportLoading.value = true
      setTimeout(() => {
        const data = users.value.map(user => ({
          用户名: user.username,
          角色: user.role,
          创建时间: user.created ? new Date(user.created).toLocaleDateString('zh-CN') : '',
          最后登录: user.last_login ? new Date(user.last_login).toLocaleString('zh-CN') : '从未登录',
          状态: user.active ? '活跃' : '禁用'
        }))
        
        const csv = [Object.keys(data[0]).join(',')]
          .concat(data.map(row => Object.values(row).join(',')))
          .join('\n')
        
        const blob = new Blob([csv], { type: 'text/csv' })
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `用户数据_${new Date().toLocaleDateString()}.csv`
        a.click()
        
        exportLoading.value = false
        message.success('用户数据导出成功')
      }, 1000)
    }

    // 批量操作
    const batchChangeRole = () => {
      message.warning('批量更改角色功能暂未实现')
    }

    const batchResetPassword = () => {
      message.warning('批量重置密码功能暂未实现')
    }

    const batchDeleteUsers = () => {
      dialog.warning({
        title: '批量删除确认',
        content: `您确定要删除选中的 ${selectedUsers.value.length} 个用户吗？此操作不可撤销。`,
        positiveText: '确定删除',
        negativeText: '取消',
        onPositiveClick: async () => {
          try {
            const deletePromises = selectedUsers.value.map(username => 
              api.delete(`/users/${username}`)
            )
            await Promise.all(deletePromises)
            message.success(`已删除 ${selectedUsers.value.length} 个用户`)
            selectedUsers.value = []
            showBatchModal.value = false
            await fetchUsers()
          } catch (error) {
            console.error('批量删除用户失败:', error)
            message.error('批量删除用户失败')
          }
        }
      })
    }

    // 获取活动类型
    const getActivityType = (type) => {
      const typeMap = {
        success: 'success',
        info: 'info',
        warning: 'warning',
        error: 'error'
      }
      return typeMap[type] || 'default'
    }

    onMounted(() => {
      fetchUsers()
    })

    return {
      // 基础数据
      users,
      filteredUsers,
      selectedUsers,
      refreshing,
      searchQuery,
      roleFilter,
      userStats,
      onlineUsers,
      userActivities,
      
      // 模态框状态
      showAddUserModal,
      showBatchModal,
      addUserLoading,
      exportLoading,
      
      // 表单数据
      newUser,
      roleOptions,
      roleFilterOptions,
      userRules,
      
      // 表格列
      userColumns,
      
      // 方法
      fetchUsers,
      editUser,
      confirmDelete,
      deleteUser,
      handleAddUser,
      exportUserData,
      batchChangeRole,
      batchResetPassword,
      batchDeleteUsers,
      getActivityType
    }
  }
}
</script>