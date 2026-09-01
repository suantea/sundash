<template>
  <div class="setup-page">
    <!-- Animated background (same style as login) -->
    <div class="setup-bg">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
      <div class="bg-orb bg-orb-3"></div>
    </div>

    <div class="setup-container">
      <div class="setup-header">
        <div class="logo-mark">
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <rect x="4" y="4" width="18" height="18" rx="5" fill="#007AFF"/>
            <rect x="26" y="4" width="18" height="18" rx="5" fill="#007AFF" opacity="0.5"/>
            <rect x="4" y="26" width="18" height="18" rx="5" fill="#007AFF" opacity="0.5"/>
            <rect x="26" y="26" width="18" height="18" rx="5" fill="#007AFF" opacity="0.2"/>
          </svg>
        </div>
        <h1>{{ $t('setup.title') }}</h1>
        <p class="subtitle">{{ $t('setup.subtitle') }}</p>
      </div>

      <div class="setup-card">
        <n-form ref="setupFormRef" :model="setupForm" :rules="setupRules" size="large">
          <n-form-item path="username" :show-label="false">
            <n-input v-model:value="setupForm.username" :placeholder="$t('login.username')" round>
              <template #prefix>
                <Icon icon="mdi:account-outline" :size="18" color="#8E8E93" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="password" :show-label="false">
            <n-input v-model:value="setupForm.password" type="password"
              :placeholder="$t('login.passwordMinLength')" round show-password-on="click">
              <template #prefix>
                <Icon icon="mdi:lock-outline" :size="18" color="#8E8E93" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="displayName" :show-label="false">
            <n-input v-model:value="setupForm.displayName" :placeholder="$t('login.displayName')" round>
              <template #prefix>
                <Icon icon="mdi:card-account-details-outline" :size="18" color="#8E8E93" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="siteTitle" :show-label="false">
            <n-input v-model:value="setupForm.siteTitle" :placeholder="$t('setup.siteTitlePlaceholder')" round>
              <template #prefix>
                <Icon icon="mdi:web" :size="18" color="#8E8E93" />
              </template>
            </n-input>
          </n-form-item>
          <n-button type="primary" block size="large" round :loading="loading" @click="handleSetup"
            class="setup-btn">
            {{ $t('setup.create') }}
          </n-button>
        </n-form>
        <p class="setup-hint">{{ $t('setup.hint') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import { NForm, NFormItem, useMessage, type FormInst, type FormRules, NInput, NButton } from 'naive-ui'
import { api } from '../api'
import { useUserStore } from '../stores/user'
import { resetSetupCache } from '../router'

const { t } = useI18n()
const userStore = useUserStore()
const message = useMessage()
const loading = ref(false)

const setupFormRef = ref<FormInst | null>(null)
const setupForm = ref({ username: '', password: '', displayName: '', siteTitle: '' })

const setupRules: FormRules = {
  username: { required: true, message: t('login.usernameRequired'), trigger: 'blur', min: 3, max: 32 },
  password: { required: true, message: t('login.passwordRequired'), trigger: 'blur', min: 6 },
}

async function handleSetup() {
  try {
    await setupFormRef.value?.validate()
  } catch { return }

  loading.value = true
  try {
    const res = await api.post('auth/setup', {
      username: setupForm.value.username,
      password: setupForm.value.password,
      display_name: setupForm.value.displayName || undefined,
      site_title: setupForm.value.siteTitle || undefined,
    })
    // 系统已初始化，重置路由守卫缓存，避免退出后再次被引导到设置页
    resetSetupCache()
    // 初始化成功：直接写入 token 并进入面板
    userStore.setToken(res.data.token, res.data.user)
    message.success(t('setup.success'))
  } catch (e: any) {
    message.error(e.response?.data?.error || t('setup.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.setup-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--sd-bg-root);
  position: relative;
  overflow: hidden;
}

.setup-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float 20s ease-in-out infinite;
}

[data-theme="dark"] .bg-orb { opacity: 0.15; }

.bg-orb-1 {
  width: 600px; height: 600px;
  background: linear-gradient(135deg, #007AFF, #5856D6);
  top: -200px; right: -100px;
}
.bg-orb-2 {
  width: 500px; height: 500px;
  background: linear-gradient(135deg, #30D158, #64D2FF);
  bottom: -180px; left: -120px;
  animation-delay: -7s;
}
.bg-orb-3 {
  width: 400px; height: 400px;
  background: linear-gradient(135deg, #FF9F0A, #FF375F);
  top: 40%; left: 30%;
  animation-delay: -12s;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(30px, -30px) scale(1.1); }
}

.setup-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 400px;
  padding: 24px;
}

.setup-header { text-align: center; margin-bottom: 28px; }
.setup-header h1 { font-size: 28px; font-weight: 700; margin: 14px 0 6px; color: var(--sd-text-primary); }
.setup-header .subtitle { color: var(--sd-text-secondary); font-size: 14px; margin: 0; }

.setup-card {
  background: var(--sd-card-bg, rgba(255,255,255,0.85));
  backdrop-filter: blur(20px);
  border: 1px solid var(--sd-border-color, rgba(0,0,0,0.08));
  border-radius: 20px;
  padding: 28px 24px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.08);
}
[data-theme="dark"] .setup-card {
  background: var(--sd-card-bg, rgba(28,28,30,0.9));
  border-color: var(--sd-border-color, rgba(255,255,255,0.1));
}

.setup-btn { margin-top: 4px; }
.setup-hint {
  margin-top: 16px;
  font-size: 12px;
  color: var(--sd-text-secondary);
  text-align: center;
  line-height: 1.6;
}
</style>
