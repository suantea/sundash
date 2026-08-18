<template>
  <div class="login-page" :class="{ 'has-login-bg': appStore.loginBgUrl }">
    <!-- Custom background image -->
    <div 
      v-if="appStore.loginBgUrl" 
      class="login-bg-image" 
      :style="{ backgroundImage: `url(${appStore.loginBgUrl})` }"
    ></div>
    <!-- Animated background -->
    <div class="login-bg" :class="{ 'has-image': appStore.loginBgUrl }">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
      <div class="bg-orb bg-orb-3"></div>
    </div>

    <div class="login-container">
      <!-- Logo & Title -->
      <div class="login-header">
        <div class="logo-mark">
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <rect x="4" y="4" width="18" height="18" rx="5" fill="#007AFF"/>
            <rect x="26" y="4" width="18" height="18" rx="5" fill="#007AFF" opacity="0.5"/>
            <rect x="4" y="26" width="18" height="18" rx="5" fill="#007AFF" opacity="0.5"/>
            <rect x="26" y="26" width="18" height="18" rx="5" fill="#007AFF" opacity="0.2"/>
          </svg>
        </div>
        <h1>Asuan</h1>
        <p class="subtitle">自托管导航面板</p>
      </div>

      <!-- Login Card -->
      <div class="login-card">
        <!-- When registration is disabled, show login form directly without tabs -->
        <template v-if="!appStore.allowRegistration">
          <n-form ref="loginFormRef" :model="loginForm" :rules="loginRules" size="large">
            <n-form-item path="username" :show-label="false">
              <n-input v-model:value="loginForm.username" placeholder="用户名" round>
                <template #prefix>
                  <Icon icon="mdi:account-outline" :size="18" color="#8E8E93" />
                </template>
              </n-input>
            </n-form-item>
            <n-form-item path="password" :show-label="false">
              <n-input v-model:value="loginForm.password" type="password" placeholder="密码" round
                show-password-on="click" @keyup.enter="handleLogin">
                <template #prefix>
                  <Icon icon="mdi:lock-outline" :size="18" color="#8E8E93" />
                </template>
              </n-input>
            </n-form-item>
            <n-button type="primary" block size="large" round :loading="loading" @click="handleLogin"
              class="login-btn">
              登录
            </n-button>
          </n-form>
        </template>

        <!-- When registration is enabled, show tabs for login/register -->
        <n-tabs v-else type="segment" v-model:value="activeTab" animated>
          <n-tab-pane name="login" tab="登录">
            <n-form ref="loginFormRef" :model="loginForm" :rules="loginRules" size="large">
              <n-form-item path="username" :show-label="false">
                <n-input v-model:value="loginForm.username" placeholder="用户名" round>
                  <template #prefix>
                    <Icon icon="mdi:account-outline" :size="18" color="#8E8E93" />
                  </template>
                </n-input>
              </n-form-item>
              <n-form-item path="password" :show-label="false">
                <n-input v-model:value="loginForm.password" type="password" placeholder="密码" round
                  show-password-on="click" @keyup.enter="handleLogin">
                  <template #prefix>
                    <Icon icon="mdi:lock-outline" :size="18" color="#8E8E93" />
                  </template>
                </n-input>
              </n-form-item>
              <n-button type="primary" block size="large" round :loading="loading" @click="handleLogin"
                class="login-btn">
                登录
              </n-button>
            </n-form>
          </n-tab-pane>

          <n-tab-pane v-if="appStore.allowRegistration" name="register" tab="注册">
            <n-form ref="registerFormRef" :model="registerForm" :rules="registerRules" size="large">
              <n-form-item path="username" :show-label="false">
                <n-input v-model:value="registerForm.username" placeholder="用户名" round>
                  <template #prefix>
                    <Icon icon="mdi:account-outline" :size="18" color="#8E8E93" />
                  </template>
                </n-input>
              </n-form-item>
              <n-form-item path="password" :show-label="false">
                <n-input v-model:value="registerForm.password" type="password" placeholder="密码（至少6位）" round
                  show-password-on="click">
                  <template #prefix>
                    <Icon icon="mdi:lock-outline" :size="18" color="#8E8E93" />
                  </template>
                </n-input>
              </n-form-item>
              <n-form-item path="displayName" :show-label="false">
                <n-input v-model:value="registerForm.displayName" placeholder="显示名称（可选）" round>
                  <template #prefix>
                    <Icon icon="mdi:card-account-details-outline" :size="18" color="#8E8E93" />
                  </template>
                </n-input>
              </n-form-item>
              <n-button type="primary" block size="large" round :loading="loading" @click="handleRegister"
                class="login-btn">
                注册
              </n-button>
            </n-form>
          </n-tab-pane>
        </n-tabs>
      </div>

      <!-- Footer -->
      <div class="login-footer">
        <button class="theme-toggle" @click="toggleTheme" :title="isDark ? '浅色模式' : '深色模式'">
          <Icon :icon="isDark ? 'mdi:weather-sunny' : 'mdi:weather-night'" :size="18" />
        </button>
        <div class="copyright-notice">Asuan · Open Source · MIT License</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { NForm, NFormItem, useMessage, type FormInst, type FormRules, NInput, NButton, NTabs, NTabPane } from 'naive-ui'
import { useAppStore } from '../stores/app'
import { useUserStore } from '../stores/user'

const appStore = useAppStore()
const userStore = useUserStore()
const message = useMessage()
const loading = ref(false)
const activeTab = ref('login')
const isDark = computed(() => appStore.isDark)

onMounted(() => {
  appStore.fetchAuthSettings()
})

const loginFormRef = ref<FormInst | null>(null)
const registerFormRef = ref<FormInst | null>(null)

const loginForm = ref({ username: '', password: '' })
const registerForm = ref({ username: '', password: '', displayName: '' })

const loginRules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
}

const registerRules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur', min: 3, max: 32 },
  password: { required: true, message: '请输入密码', trigger: 'blur', min: 6 },
}

async function handleLogin() {
  try {
    await loginFormRef.value?.validate()
  } catch { return }

  loading.value = true
  try {
    await userStore.login(loginForm.value.username, loginForm.value.password)
    message.success('登录成功')
  } catch (e: any) {
    message.error(e.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  try {
    await registerFormRef.value?.validate()
  } catch { return }

  loading.value = true
  try {
    await userStore.register(registerForm.value.username, registerForm.value.password, registerForm.value.displayName)
    if (appStore.requireApproval) {
      message.success('注册成功，请等待管理员审批后登录')
      activeTab.value = 'login'
    } else {
      message.success('注册成功')
    }
  } catch (e: any) {
    message.error(e.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
}

function toggleTheme() {
  appStore.setTheme(appStore.isDark ? 'light' : 'dark')
}
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--sd-bg-root);
  position: relative;
  overflow: hidden;
}

.login-bg-image {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  z-index: 0;
}

/* Animated gradient background orbs */
.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  transition: opacity 0.3s ease;
}

.login-bg.has-image {
  opacity: 0.3;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float 20s ease-in-out infinite;
}

[data-theme="dark"] .bg-orb {
  opacity: 0.15;
}

.bg-orb-1 {
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, #007AFF, #5856D6);
  top: -200px;
  right: -100px;
  animation-delay: 0s;
}

.bg-orb-2 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #5856D6, #FF2D55);
  bottom: -200px;
  left: -100px;
  animation-delay: -7s;
}

.bg-orb-3 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #34C759, #007AFF);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -14s;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(30px, -40px) scale(1.05); }
  50% { transform: translate(-20px, 20px) scale(0.95); }
  75% { transform: translate(40px, 30px) scale(1.02); }
}

.login-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 400px;
  padding: var(--sd-space-6);
  animation: fadeInUp 0.6s var(--sd-ease-spring) both;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.login-header {
  text-align: center;
  margin-bottom: var(--sd-space-8);
}

.logo-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 72px;
  height: 72px;
  margin-bottom: var(--sd-space-5);
  background: var(--sd-bg-surface);
  border-radius: var(--sd-radius-xl);
  box-shadow: var(--sd-shadow-lg);
}

.login-header h1 {
  font-size: var(--sd-text-3xl);
  font-weight: var(--sd-weight-bold);
  color: var(--sd-text-primary);
  letter-spacing: -0.5px;
  margin: 0 0 var(--sd-space-2);
}

.subtitle {
  font-size: var(--sd-text-base);
  color: var(--sd-text-secondary);
  margin: 0;
}

.login-card {
  background: var(--sd-glass-bg);
  backdrop-filter: var(--sd-blur);
  -webkit-backdrop-filter: var(--sd-blur);
  border: 1px solid var(--sd-glass-border);
  border-radius: var(--sd-radius-2xl);
  padding: var(--sd-space-6);
  box-shadow: var(--sd-shadow-xl);
}

.login-card :deep(.n-tabs) {
  --n-tab-border-radius: 10px;
}

.login-card :deep(.n-tabs-tab) {
  font-weight: 500;
}

.login-card :deep(.n-form-item) {
  margin-bottom: 18px;
}

.login-card :deep(.n-form-item:last-of-type) {
  margin-bottom: 24px;
}

.login-btn {
  height: 44px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  letter-spacing: 0.3px;
}

.login-footer {
  text-align: center;
  margin-top: var(--sd-space-8);
}

.theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  padding: 0;
  border: none;
  background: transparent;
  color: var(--sd-text-secondary);
  border-radius: 8px;
  cursor: pointer;
  transition: all var(--sd-duration-fast) var(--sd-ease);
}

.theme-toggle:hover {
  background: var(--sd-bg-surface-secondary);
  color: var(--sd-text-primary);
}

.copyright-notice {
  margin-top: 16px;
  font-size: 12px;
  color: var(--sd-text-tertiary);
  opacity: 0.6;
}
</style>
