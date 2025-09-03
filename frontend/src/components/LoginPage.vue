<template>
  <div class="login-container">
    <n-card class="login-form" :bordered="false">
      <template #header>
        <div class="login-header">
          <n-h2>🔐 Stone 防火墙</n-h2>
          <n-text depth="3">请输入您的认证信息</n-text>
        </div>
      </template>
      
      <n-form @submit.prevent="handleLogin" :model="formModel" ref="formRef" :rules="rules">
        <n-form-item label="账户" path="account">
          <n-input
            v-model:value="formModel.account"
            placeholder="请输入账户名"
            size="large"
            :input-props="{ autocomplete: 'username' }"
          />
        </n-form-item>
        
        <n-form-item label="验证码" path="code">
          <n-input
            v-model:value="formModel.code"
            placeholder="请输入双因素认证码"
            size="large"
            :input-props="{ autocomplete: 'one-time-code' }"
          />
        </n-form-item>
        
        <div class="button-group">
          <n-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            block
          >
            🚀 登录
          </n-button>
          
          <n-button
            type="info"
            size="large"
            @click="goToSetup2FA"
            block
          >
            🔑 设置双因素认证
          </n-button>
        </div>
      </n-form>
    </n-card>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { useMessage } from 'naive-ui'

export default {
  name: 'LoginPage',
  setup() {
    const formRef = ref(null)
    const router = useRouter()
    const store = useStore()
    const message = useMessage()
    const loading = ref(false)

    const formModel = reactive({
      account: '',
      code: ''
    })

    const rules = {
      account: [
        { required: true, message: '请输入账户名', trigger: 'blur' }
      ],
      code: [
        { required: true, message: '请输入验证码', trigger: 'blur' },
        { min: 6, max: 6, message: '验证码必须是6位数字', trigger: 'blur' }
      ]
    }

    const handleLogin = async () => {
      try {
        await formRef.value?.validate()
        loading.value = true
        
        const success = await store.dispatch('login', {
          account: formModel.account,
          code: formModel.code
        })
        
        if (success) {
          message.success('🎉 登录成功！欢迎回来！')
          setTimeout(() => {
            router.push({ name: 'Home' })
          }, 1000)
        } else {
          throw new Error('登录失败')
        }
      } catch (error) {
        if (error?.errors) {
          // 表单验证错误
          return
        }
        message.error('❌ 登录失败，请检查您的验证码和账户信息')
      } finally {
        loading.value = false
      }
    }

    const goToSetup2FA = () => {
      router.push('/setup-2fa')
    }

    return {
      formRef,
      formModel,
      rules,
      loading,
      handleLogin,
      goToSetup2FA
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-form {
  width: 100%;
  max-width: 400px;
  margin: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 20px;
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 24px;
}
</style>