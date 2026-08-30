<template>
  <div class="profile-page">
    <div class="page-bg">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
    </div>

    <header class="page-header">
      <button class="back-btn" @click="$router.push('/')">
        <Icon icon="mdi:chevron-left" :size="18" />
        <span>{{ $t("common.back") }}</span>
      </button>
      <h1>{{ $t("user.profile") }}</h1>
      <div style="width: 70px;"></div>
    </header>

    <div class="profile-content">
      <!-- Profile Info -->
      <div class="section-header">
        <div>
          <h2>{{ $t("user.profileInfo") }}</h2>
          <p>{{ $t("user.manageDisplayNameAndAvatar") }}</p>
        </div>
      </div>

      <div class="settings-card">
        <div class="setting-row">
          <div class="setting-left">
            <div class="user-avatar-lg">
              {{ (userStore.user?.display_name || userStore.user?.username || 'U')[0].toUpperCase() }}
            </div>
            <div class="user-info">
              <div class="user-name">{{ userStore.user?.username }}</div>
              <div class="user-role">
                <span :class="['role-pill', userStore.user?.role === 'admin' ? 'is-admin' : 'is-user']">
                  {{ userStore.user?.role === 'admin' ? $t('admin.roleAdmin') : $t('admin.roleUser') }}
                </span>
              </div>
            </div>
          </div>
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:card-account-details-outline" :size="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">{{ $t("user.displayName") }}</div>
              <div class="setting-desc">{{ $t("user.displayNameDesc") }}</div>
            </div>
          </div>
          <div class="inline-edit">
            <input v-model="displayName" class="mini-text-input" :placeholder="$t('user.enterDisplayName')" />
            <button class="save-inline-btn" @click="saveDisplayName" :disabled="displayNameSaving">
              {{ displayNameSaving ? $t('common.saving') : $t('common.save') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Change Password -->
      <div class="section-header" style="margin-top: 32px;">
        <div>
          <h2>{{ $t("user.changePassword") }}</h2>
          <p>{{ $t("user.changePasswordDesc") }}</p>
        </div>
      </div>

      <div class="settings-card">
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="form-field">
            <label>{{ $t("user.currentPassword") }}</label>
            <n-input v-model:value="passwordForm.oldPassword" type="password" :placeholder="$t('user.enterCurrentPassword')" show-password-on="click" />
          </div>
          <div class="form-field">
            <label>{{ $t("user.newPassword") }}</label>
            <n-input v-model:value="passwordForm.newPassword" type="password" :placeholder="$t('user.enterNewPasswordAtLeast6')" show-password-on="click" />
          </div>
          <div class="form-field">
            <label>{{ $t("user.confirmPassword") }}</label>
            <n-input v-model:value="passwordForm.confirmPassword" type="password" :placeholder="$t('user.enterNewPasswordAgain')" show-password-on="click" @keyup.enter="handleChangePassword" />
          </div>
          <div style="display: flex; justify-content: flex-end;">
            <n-button type="primary" :loading="passwordSaving" @click="handleChangePassword">
              {{ $t("user.changePassword") }}
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import { useMessage, NInput, NButton } from 'naive-ui'
import { useUserStore } from '../stores/user'
import { api } from '../api'

const { t } = useI18n()

const userStore = useUserStore()
const message = useMessage()

const displayName = ref('')
const displayNameSaving = ref(false)
const passwordSaving = ref(false)
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

onMounted(() => {
  userStore.fetchProfile()
  displayName.value = userStore.user?.display_name || userStore.user?.username || ''
})

async function saveDisplayName() {
  if (!displayName.value.trim()) {
    message.warning($t('user.displayNameCannotBeEmpty'))
    return
  }
  displayNameSaving.value = true
  try {
    await api.put('profile', { display_name: displayName.value.trim() })
    if (userStore.user) {
      userStore.user.display_name = displayName.value.trim()
    }
    message.success($t('user.displayNameUpdated'))
  } catch {
    message.error($t('user.updateFailed'))
  } finally {
    displayNameSaving.value = false
  }
}

async function handleChangePassword() {
  if (!passwordForm.value.oldPassword) {
    message.warning($t('user.pleaseEnterCurrentPassword'))
    return
  }
  if (!passwordForm.value.newPassword || passwordForm.value.newPassword.length < 6) {
    message.warning($t('user.newPasswordAtLeast6'))
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    message.warning($t('user.passwordsDoNotMatch'))
    return
  }
  passwordSaving.value = true
  try {
    await api.put('profile/password', {
      old_password: passwordForm.value.oldPassword,
      new_password: passwordForm.value.newPassword,
    })
    message.success($t('user.passwordChanged'))
    passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  } catch (e: any) {
    message.error(e.response?.data?.error || $t('user.changePasswordFailed'))
  } finally {
    passwordSaving.value = false
  }
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: var(--sd-bg-root);
  position: relative;
}

.page-bg {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
}

.bg-orb-1 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(0,122,255,0.07) 0%, transparent 70%);
  top: -150px;
  right: -100px;
  animation: floatOrb 20s ease-in-out infinite;
}

.bg-orb-2 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(88,86,214,0.05) 0%, transparent 70%);
  bottom: -100px;
  left: -80px;
  animation: floatOrb 25s ease-in-out infinite reverse;
}

@keyframes floatOrb {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(20px, -20px); }
}

.page-header {
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--sd-header-height);
  padding: 0 var(--sd-space-6);
  background: rgba(255,255,255,0.8);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid rgba(0,0,0,0.06);
}

:root[data-theme="dark"] .page-header {
  background: rgba(0,0,0,0.7);
  border-bottom-color: rgba(255,255,255,0.08);
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 2px;
  background: none;
  border: none;
  color: var(--sd-primary);
  font-size: 14px;
  font-family: var(--sd-font);
  font-weight: 500;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: background 0.15s ease;
}

.back-btn:hover { background: rgba(0,122,255,0.08); }

.page-header h1 {
  font-size: 17px;
  font-weight: 600;
  margin: 0;
  color: var(--sd-text-primary);
}

.profile-content {
  position: relative;
  z-index: 1;
  max-width: 600px;
  margin: 0 auto;
  padding: var(--sd-space-6) var(--sd-space-6) var(--sd-space-16);
}

.section-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: var(--sd-space-5);
}

.section-header h2 {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 4px;
  color: var(--sd-text-primary);
}

.section-header p {
  font-size: 14px;
  color: var(--sd-text-secondary);
  margin: 0;
}

.settings-card {
  background: rgba(255,255,255,0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border-radius: 14px;
  border: 1px solid rgba(0,0,0,0.06);
  overflow: hidden;
}

:root[data-theme="dark"] .settings-card {
  background: rgba(28,28,30,0.8);
  border-color: rgba(255,255,255,0.08);
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  min-height: 52px;
}

.setting-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.setting-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.setting-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--sd-text-primary);
  line-height: 1.3;
}

.setting-desc {
  font-size: 12px;
  color: var(--sd-text-tertiary);
  line-height: 1.3;
}

.setting-divider {
  height: 1px;
  background: rgba(0,0,0,0.05);
  margin: 0 16px;
}

:root[data-theme="dark"] .setting-divider {
  background: rgba(255,255,255,0.06);
}

.user-avatar-lg {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0,122,255,0.1);
  color: var(--sd-primary);
  border-radius: 14px;
  font-size: 20px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--sd-text-primary);
}

.role-pill {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
}

.is-admin {
  background: rgba(255,149,0,0.12);
  color: #FF9500;
}

.is-user {
  background: rgba(0,122,255,0.1);
  color: #007AFF;
}

.inline-edit {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mini-text-input {
  width: 200px;
  height: 34px;
  padding: 0 12px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 13px;
  font-family: var(--sd-font);
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  outline: none;
  transition: border-color 0.2s;
}

.mini-text-input:focus {
  border-color: #007AFF;
}

.save-inline-btn {
  padding: 0 16px;
  height: 34px;
  background: var(--sd-primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.save-inline-btn:hover {
  background: var(--sd-primary-hover);
}

.save-inline-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-field label {
  font-size: 13px;
  font-weight: 500;
  color: var(--sd-text-secondary);
}
</style>
