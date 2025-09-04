import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'
import MainLayout from '@/components/MainLayout.vue'
import Home from '@/components/HomePage.vue'
import AttackerProfileTracking from '@/components/AttackerProfileTracking.vue'
import CustomRuleManagement from '@/components/CustomRuleManagement.vue'
import LogAnalysis from '@/components/LogAnalysis.vue'
import LoginPage from '@/components/LoginPage.vue'
import RegisterPage from '@/components/RegisterPage.vue'
import SystemConfiguration from '@/components/SystemConfiguration.vue'
import UserManagement from '@/components/UserManagement.vue'
import WAFDashboard from '@/components/WAFDashboard.vue'

const routes = [
    {
        path: '/login',
        name: 'LoginPage',
        component: LoginPage
    },
    {
        path: '/register',
        name: 'RegisterPage',
        component: RegisterPage
    },
    {
        path: '/',
        component: MainLayout,
        children: [
            {
                path: '',
                name: 'Home',
                component: Home
            },
            {
                path: 'dashboard',
                name: 'WAFDashboard',
                component: WAFDashboard
            },
            {
                path: 'attacker-profile',
                name: 'AttackerProfileTracking',
                component: AttackerProfileTracking
            },
            {
                path: 'custom-rule',
                name: 'CustomRuleManagement',
                component: CustomRuleManagement
            },
            {
                path: 'log-analysis',
                name: 'LogAnalysis',
                component: LogAnalysis
            },
            {
                path: 'system-configuration',
                name: 'SystemConfiguration',
                component: SystemConfiguration
            },
            {
                path: 'user-management',
                name: 'UserManagement',
                component: UserManagement
            }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

router.beforeEach(async (to, from, next) => {
    await store.dispatch('checkAuth'); // 确保最新的认证状态

    const isAuthenticated = store.state.isAuthenticated;
    const publicRoutes = ['LoginPage', 'RegisterPage'];
    
    // 未认证用户访问根路径或主页，重定向到登录页
    if (!isAuthenticated && (to.name === 'Home' || to.path === '/')) {
        next({ name: 'LoginPage' });
    }
    // 如果用户未认证且访问受保护的页面，重定向到登录页
    else if (!isAuthenticated && !publicRoutes.includes(to.name)) {
        next({ name: 'LoginPage' });
    } 
    // 如果用户已认证且访问登录页，重定向到仪表盘
    else if (isAuthenticated && to.name === 'LoginPage') {
        next({ name: 'WAFDashboard' });
    }
    // 如果用户已认证且访问根路径(但不是具体的Home路由)，重定向到主页
    else if (isAuthenticated && to.path === '/' && to.name !== 'Home') {
        next({ name: 'Home' });
    }
    else {
        next();
    }
});

export default router