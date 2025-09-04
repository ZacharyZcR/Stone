<template>
  <n-config-provider :theme="null">
    <n-layout class="register-layout">
      <n-layout-content>
        <div class="register-container">
          <n-grid :cols="2" :x-gap="0" responsive="screen" class="register-grid">
            <!-- 品牌展示区域 -->
            <n-grid-item class="brand-area">
              <n-card :bordered="false" class="brand-card">
                <n-space direction="vertical" align="center" size="large">
                  <n-space direction="vertical" align="center" size="medium">
                    <n-icon size="80" depth="1">
                      <shield-checkmark-outline />
                    </n-icon>
                    <n-h1 class="brand-title">Stone WAF</n-h1>
                    <n-text type="info" class="brand-subtitle">
                      Web Application Firewall
                    </n-text>
                  </n-space>
                  
                  <n-space direction="vertical" size="medium">
                    <n-space align="center" size="small">
                      <n-icon size="20" color="#18a058">
                        <checkmark-circle-outline />
                      </n-icon>
                      <n-text>实时威胁检测</n-text>
                    </n-space>
                    <n-space align="center" size="small">
                      <n-icon size="20" color="#18a058">
                        <checkmark-circle-outline />
                      </n-icon>
                      <n-text>智能规则引擎</n-text>
                    </n-space>
                    <n-space align="center" size="small">
                      <n-icon size="20" color="#18a058">
                        <checkmark-circle-outline />
                      </n-icon>
                      <n-text>安全认证体系</n-text>
                    </n-space>
                  </n-space>
                </n-space>
              </n-card>
            </n-grid-item>

            <!-- 注册表单区域 -->
            <n-grid-item class="form-area">
              <n-card :bordered="false" class="form-card">
                <n-space direction="vertical" size="large">
                  <n-space direction="vertical" align="center" size="small">
                    <n-h2>创建账户</n-h2>
                    <n-text depth="3">注册后即可使用系统功能</n-text>
                  </n-space>

                  <n-form 
                    @submit.prevent="handleRegister" 
                    :model="formModel" 
                    ref="formRef" 
                    :rules="rules"
                  >
                    <n-form-item path="username" :show-label="false">
                      <n-input
                        v-model:value="formModel.username"
                        placeholder="用户名 (3-20位字符)"
                        size="large"
                        :input-props="{ autocomplete: 'username' }"
                      >
                        <template #prefix>
                          <n-icon>
                            <person-outline />
                          </n-icon>
                        </template>
                      </n-input>
                    </n-form-item>
                    
                    <n-form-item path="password" :show-label="false">
                      <n-input
                        v-model:value="formModel.password"
                        placeholder="密码 (至少8位字符)"
                        type="password"
                        show-password-on="mousedown"
                        size="large"
                        :input-props="{ autocomplete: 'new-password' }"
                      >
                        <template #prefix>
                          <n-icon>
                            <key-outline />
                          </n-icon>
                        </template>
                      </n-input>
                    </n-form-item>

                    <n-form-item path="confirmPassword" :show-label="false">
                      <n-input
                        v-model:value="formModel.confirmPassword"
                        placeholder="确认密码"
                        type="password"
                        show-password-on="mousedown"
                        size="large"
                        :input-props="{ autocomplete: 'new-password' }"
                      >
                        <template #prefix>
                          <n-icon>
                            <lock-closed-outline />
                          </n-icon>
                        </template>
                      </n-input>
                    </n-form-item>
                    
                    <n-form-item :show-label="false">
                      <n-space direction="vertical" size="medium" style="width: 100%">
                        <n-button
                          type="primary"
                          size="large"
                          :loading="loading"
                          @click="handleRegister"
                          block
                        >
                          注册账户
                        </n-button>
                        
                        <n-button
                          secondary
                          size="large"
                          @click="goToLogin"
                          block
                        >
                          已有账户？立即登录
                        </n-button>
                      </n-space>
                    </n-form-item>
                  </n-form>

                  <n-divider>
                    <n-text depth="3" style="font-size: 12px;">安全注册</n-text>
                  </n-divider>
                </n-space>
              </n-card>
            </n-grid-item>
          </n-grid>
        </div>
      </n-layout-content>
    </n-layout>
  </n-config-provider>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import api from '../api/axiosInstance'
import { 
  PersonOutline, 
  KeyOutline, 
  LockClosedOutline,
  ShieldCheckmarkOutline,
  CheckmarkCircleOutline 
} from '@vicons/ionicons5'

export default {
  name: 'RegisterPage',
  components: {
    PersonOutline,
    KeyOutline,
    LockClosedOutline,
    ShieldCheckmarkOutline,
    CheckmarkCircleOutline
  },
  setup() {
    const formRef = ref(null)
    const router = useRouter()
    const message = useMessage()
    const loading = ref(false)

    const formModel = reactive({
      username: '',
      password: '',
      confirmPassword: ''
    })

    const rules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度为3-20个字符', trigger: 'blur' },
        { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 8, message: '密码至少8位字符', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        {
          validator: (rule, value) => {
            return value === formModel.password
          },
          message: '两次输入的密码不一致',
          trigger: 'blur'
        }
      ]
    }

    const handleRegister = async () => {
      try {
        await formRef.value?.validate()
        loading.value = true
        
        const response = await api.post('/auth/register', {
          username: formModel.username,
          password: formModel.password
        })
        
        if (response.data.message === '注册成功') {
          message.success('注册成功！请使用新账户登录')
          setTimeout(() => {
            router.push({ name: 'LoginPage' })
          }, 1000)
        } else {
          throw new Error('注册失败')
        }
      } catch (error) {
        if (error?.errors) {
          return
        }
        
        let errorMsg = '注册失败，请稍后重试'
        if (error.response?.data?.error) {
          errorMsg = error.response.data.error
        }
        
        message.error(errorMsg)
      } finally {
        loading.value = false
      }
    }

    const goToLogin = () => {
      router.push('/login')
    }

    return {
      formRef,
      formModel,
      rules,
      loading,
      handleRegister,
      goToLogin
    }
  }
}
</script>

<style scoped>
.register-layout {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.register-grid {
  max-width: 1000px;
  width: 100%;
}

.brand-area :deep(.n-card) {
  min-height: 500px;
  display: flex;
  align-items: center;
}

.form-area :deep(.n-card) {
  min-height: 500px;
  display: flex;
  align-items: center;
}

.brand-title {
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .register-grid {
    margin: 16px;
  }
  
  .register-grid :deep(.n-grid) {
    grid-template-columns: 1fr !important;
  }
  
  .brand-area :deep(.n-card),
  .form-area :deep(.n-card) {
    min-height: auto;
    padding: 32px 24px;
  }
}
</style>