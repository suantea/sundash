import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import { api } from '../api'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/Login.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/setup',
      name: 'Setup',
      component: () => import('../views/Setup.vue'),
      meta: { requiresAuth: false, isSetup: true },
    },
    {
      path: '/',
      name: 'Home',
      component: () => import('../views/Home.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('../views/Settings.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('../views/Profile.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'Admin',
      component: () => import('../views/Admin.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/bookmarks',
      name: 'Bookmarks',
      component: () => import('../views/Bookmarks.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// 缓存首次安装状态，避免每次导航都打接口
let needsSetupCache: boolean | null = null

async function checkNeedsSetup(): Promise<boolean> {
  if (needsSetupCache !== null) return needsSetupCache
  try {
    const res = await api.get('auth/setup-status')
    needsSetupCache = res.data.needs_setup === true
  } catch {
    needsSetupCache = false
  }
  return needsSetupCache
}

// 初始化完成后重置缓存，允许后续再次检查（例如：管理员又删光了用户）
export function resetSetupCache() {
  needsSetupCache = null
}

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()

  if (to.name === 'Setup') {
    // 已初始化则不允许访问设置页
    const needs = await checkNeedsSetup()
    if (!needs) next({ name: 'Login' })
    else next()
    return
  }

  if (to.name === 'Login') {
    const needs = await checkNeedsSetup()
    if (needs) next({ name: 'Setup' })
    else next()
    return
  }

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    // 未登录时若系统尚未初始化，引导去首次设置
    const needs = await checkNeedsSetup()
    if (needs) next({ name: 'Setup' })
    else next({ name: 'Login' })
  } else if (to.meta.requiresAdmin && userStore.user?.role !== 'admin') {
    next({ name: 'Home' })
  } else {
    next()
  }
})

export default router
