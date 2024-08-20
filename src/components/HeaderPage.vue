<template>
  <nav class="bg-gradient-to-r from-gray-700 via-gray-800 to-gray-900 p-4 shadow-md fixed w-full z-10 transition-all duration-500">
    <div class="container mx-auto flex justify-between items-center">
      <div class="text-2xl font-bold text-white">🔒 Stone ⛰️ WAF 管理面板</div>
      <div class="space-x-4">
        <router-link v-if="!isAuthenticated" to="/login" v-slot="{ navigate }">
          <button
              @click="navigate"
              class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
          >
            登录 🔐
          </button>
        </router-link>
        <template v-else>
          <router-link to="/" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              仪表盘 📊
            </button>
          </router-link>
          <router-link to="/attacker-profile" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              攻击者画像追踪 🕵️
            </button>
          </router-link>
          <router-link to="/custom-rule" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              规则管理 ⚙️
            </button>
          </router-link>
          <router-link to="/log-analysis" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              日志分析 📜
            </button>
          </router-link>
          <router-link to="/system-configuration" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              系统配置 ⚙️
            </button>
          </router-link>
          <router-link to="/user-management" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              用户管理 👥
            </button>
          </router-link>
          <router-link to="/dashboard" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              仪表盘 📊
            </button>
          </router-link>
          <button
              @click="logout"
              class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
          >
            登出 🚪
          </button>
        </template>
      </div>
    </div>
  </nav>
</template>

<script>
import { mapState, mapActions } from 'vuex'

export default {
  name: 'HeaderPage',
  computed: {
    ...mapState(['isAuthenticated'])
  },
  methods: {
    ...mapActions(['logout'])
  }
}
</script>
