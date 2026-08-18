<template>
  <div class="admin-page">
    <div class="page-bg">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
    </div>

    <header class="page-header">
      <button class="back-btn" @click="$router.push('/')">
        <Icon icon="mdi:chevron-left" :size="18" />
        <span>返回</span>
      </button>
      <h1>管理面板</h1>
      <div style="width: 70px;"></div>
    </header>

    <div class="admin-content">
      <!-- ===== 1. System Settings ===== -->
      <div class="section-header">
        <div>
          <h2>系统设置</h2>
          <p>管理注册策略和默认配置</p>
        </div>
      </div>

      <div class="settings-card">
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
              <Icon icon="mdi:account-plus" :size="18" color="#34C759" />
            </div>
            <div>
              <div class="setting-title">允许用户注册</div>
              <div class="setting-desc">开启后主页显示注册选项</div>
            </div>
          </div>
          <n-switch v-model:value="globalSettings.allow_registration" @update:value="saveGlobalSetting('allow_registration', $event)" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
              <Icon icon="mdi:shield-check" :size="18" color="#FF9500" />
            </div>
            <div>
              <div class="setting-title">注册需要审批</div>
              <div class="setting-desc">开启后新用户需管理员审批</div>
            </div>
          </div>
          <n-switch v-model:value="globalSettings.require_approval" @update:value="saveGlobalSetting('require_approval', $event)" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:image" :size="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">默认壁纸类型</div>
              <div class="setting-desc">新用户的默认背景</div>
            </div>
          </div>
          <div class="seg-control">
            <button :class="['seg-btn', { active: globalSettings.default_wallpaper_type === '' }]" @click="saveGlobalSetting('default_wallpaper_type', '')">默认</button>
            <button :class="['seg-btn', { active: globalSettings.default_wallpaper_type === 'bing' }]" @click="saveGlobalSetting('default_wallpaper_type', 'bing')">必应</button>
            <button :class="['seg-btn', { active: globalSettings.default_wallpaper_type === 'custom' }]" @click="saveGlobalSetting('default_wallpaper_type', 'custom')">自定义</button>
          </div>
        </div>
        <div v-if="globalSettings.default_wallpaper_type === 'custom'" class="setting-sub">
          <input v-model="globalSettings.default_wallpaper_url" class="bg-url-input" placeholder="输入背景图片 URL" @change="saveGlobalSetting('default_wallpaper_url', globalSettings.default_wallpaper_url)" />
        </div>
      </div>

      <!-- ===== 2. Site Configuration ===== -->
      <div class="section-header" style="margin-top: 32px;">
        <div>
          <h2>站点配置</h2>
          <p>管理网站标题、图标、CDN、统计代码等</p>
        </div>
        <button class="add-btn" @click="saveAllSiteConfig" :disabled="siteSaving">
          <Icon icon="mdi:content-save" :size="16" />
          {{ siteSaving ? '保存中...' : '保存配置' }}
        </button>
      </div>

      <div class="settings-card">
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:text-short" :size="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">网站标题</div>
              <div class="setting-desc">浏览器标签页显示的标题</div>
            </div>
          </div>
          <input v-model="siteConfig.site_title" class="mini-text-input" placeholder="Asuan" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
              <Icon icon="mdi:icon" :size="18" color="#FF9500" />
            </div>
            <div>
              <div class="setting-title">网站图标</div>
              <div class="setting-desc">ICO 或 SVG 图标 URL</div>
            </div>
          </div>
          <input v-model="siteConfig.site_icon_url" class="mini-text-input wide" placeholder="/favicon.svg 或 https://..." />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
              <Icon icon="mdi:cloud" :size="18" color="#34C759" />
            </div>
            <div>
              <div class="setting-title">CDN 地址</div>
              <div class="setting-desc">静态资源 CDN 基础路径</div>
            </div>
          </div>
          <input v-model="siteConfig.site_cdn_url" class="mini-text-input wide" placeholder="https://cdn.example.com" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
              <Icon icon="mdi:chart-line" :size="18" color="#AF52DE" />
            </div>
            <div>
              <div class="setting-title">统计代码</div>
              <div class="setting-desc">百度统计、Google Analytics 等脚本</div>
            </div>
          </div>
          <textarea v-model="siteConfig.site_analytics_code" class="code-input" placeholder="<script>/* 统计代码 */</script>" rows="3"></textarea>
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:code-tags" :size="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">自定义 Head 代码</div>
              <div class="setting-desc">注入到 &lt;head&gt; 标签中的 CSS/JS</div>
            </div>
          </div>
          <textarea v-model="siteConfig.site_custom_head" class="code-input" placeholder="<link rel=&quot;stylesheet&quot; href=&quot;...&quot;>" rows="3"></textarea>
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
              <Icon icon="mdi:code-tags" :size="18" color="#34C759" />
            </div>
            <div>
              <div class="setting-title">自定义页脚代码</div>
              <div class="setting-desc">注入到 &lt;/body&gt; 前的代码</div>
            </div>
          </div>
          <textarea v-model="siteConfig.site_custom_footer" class="code-input" placeholder="<script>/* 页脚代码 */</script>" rows="3"></textarea>
        </div>
      </div>

      <!-- ===== 3. Pending Users ===== -->
      <div class="section-header" style="margin-top: 32px;">
        <div>
          <h2>待审批用户</h2>
          <p>审批新注册的用户账户</p>
        </div>
        <span v-if="pendingUsers.length" class="badge">{{ pendingUsers.length }}</span>
      </div>

      <div class="table-card">
        <div v-for="user in pendingUsers" :key="user.id" class="table-row">
          <div class="td" style="flex: 2.5;">
            <span class="u-avatar pending">{{ user.username[0].toUpperCase() }}</span>
            <span class="u-name">{{ user.username }}</span>
          </div>
          <div class="td" style="flex: 2; color: var(--sd-text-secondary);">{{ user.display_name || '-' }}</div>
          <div class="td" style="flex: 2; gap: 6px;">
            <button class="row-btn approve" @click="approveUser(user)">
              <Icon icon="mdi:check" :size="14" />
              批准
            </button>
            <button class="row-btn danger" @click="rejectUser(user)">
              <Icon icon="mdi:close" :size="14" />
              拒绝
            </button>
          </div>
        </div>
        <div v-if="pendingUsers.length === 0" class="table-empty">暂无待审批用户</div>
      </div>

      <!-- ===== 3. All Users ===== -->
      <div class="section-header" style="margin-top: 32px;">
        <div>
          <h2>用户管理</h2>
          <p>管理系统中的所有用户账户</p>
        </div>
        <button class="add-btn" @click="showAddUser = true">
          <Icon icon="mdi:plus" :size="16" />
          添加用户
        </button>
      </div>

      <div class="table-card">
        <div class="table-head">
          <div class="th" style="flex: 2.5;">用户</div>
          <div class="th" style="flex: 2;">显示名</div>
          <div class="th" style="flex: 1;">角色</div>
          <div class="th" style="flex: 1;">状态</div>
          <div class="th" style="flex: 2.5;">操作</div>
        </div>
        <div v-for="user in users" :key="user.id" class="table-row">
          <div class="td" style="flex: 2.5;">
            <span class="u-avatar">{{ user.username[0].toUpperCase() }}</span>
            <span class="u-name">{{ user.username }}</span>
          </div>
          <div class="td" style="flex: 2; color: var(--sd-text-secondary);">{{ user.display_name || '-' }}</div>
          <div class="td" style="flex: 1;">
            <span :class="['role-pill', user.role === 'admin' ? 'is-admin' : 'is-user']">
              {{ user.role === 'admin' ? '管理员' : '用户' }}
            </span>
          </div>
          <div class="td" style="flex: 1;">
            <span :class="['status-pill', 'is-' + (user.status || 'approved')]">
              {{ statusLabel(user.status) }}
            </span>
          </div>
          <div class="td td-actions">
            <button class="row-btn" @click="toggleRole(user)">
              {{ user.role === 'admin' ? '设为用户' : '设为管理员' }}
            </button>
            <button class="row-btn" @click="viewUserPanel(user)">
              <Icon icon="mdi:view-grid" :size="14" />
              书签
            </button>
            <button class="row-btn danger" @click="deleteUser(user)">删除</button>
            <button class="row-btn" @click="openEditDisplayName(user)">
              <Icon icon="mdi:rename-box" :size="14" />
              改名
            </button>
            <button class="row-btn" @click="openResetPassword(user)" style="color: #FF9500; border-color: rgba(255,149,0,0.3);">
              <Icon icon="mdi:key-variant" :size="14" />
              重置密码
            </button>
          </div>
        </div>
        <div v-if="users.length === 0" class="table-empty">暂无用户数据</div>
      </div>
    </div>

    <!-- ===== User Bookmarks Modal ===== -->
    <n-modal v-model:show="showUserPanel" preset="card" :title="`${viewingUser?.username || ''} 的书签`" style="max-width: 700px;" :bordered="false">
      <div v-if="viewingGroups.length === 0" class="table-empty">该用户暂无书签</div>
      <div v-for="group in viewingGroups" :key="group.id" class="view-group">
        <div class="view-group-name">
          <Icon icon="mdi:folder" :size="16" color="#007AFF" />
          {{ group.name }}
          <span class="view-card-count">{{ group.cards?.length || 0 }} 个书签</span>
        </div>
        <div v-if="group.cards && group.cards.length" class="view-cards">
          <div v-for="card in group.cards" :key="card.id" class="view-card">
            <span class="view-card-title">{{ card.title }}</span>
            <span class="view-card-url">{{ card.url }}</span>
          </div>
        </div>
      </div>
    </n-modal>

    <!-- ===== Add User Modal ===== -->
    <n-modal v-model:show="showAddUser" preset="dialog" title="添加用户" positive-text="创建" negative-text="取消"
      :loading="modalLoading" @positive-click="handleCreateUser">
      <n-form label-placement="top">
        <n-form-item label="用户名">
          <n-input v-model:value="newUser.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item label="密码">
          <n-input v-model:value="newUser.password" type="password" placeholder="请输入密码" show-password-on="click" />
        </n-form-item>
        <n-form-item label="显示名称">
          <n-input v-model:value="newUser.displayName" placeholder="请输入显示名称（可选）" />
        </n-form-item>
      </n-form>
    </n-modal>

    <!-- ===== Reset Password Modal ===== -->
    <n-modal v-model:show="showResetPassword" preset="dialog" title="重置用户密码" positive-text="重置" negative-text="取消"
      :loading="resetPasswordLoading" @positive-click="handleResetPassword">
      <div style="margin-bottom: 12px; color: var(--sd-text-secondary); font-size: 14px;">
        为用户 <strong>{{ resettingUser?.username }}</strong> 设置新密码
      </div>
      <n-input v-model:value="newResetPassword" type="password" placeholder="输入新密码（至少6位）" show-password-on="click" @keyup.enter="handleResetPassword" />
    </n-modal>

    <!-- ===== Edit Display Name Modal ===== -->
    <n-modal v-model:show="showEditDisplayName" preset="dialog" title="修改显示名" positive-text="保存" negative-text="取消"
      :loading="editDisplayNameLoading" @positive-click="handleEditDisplayName">
      <div style="margin-bottom: 12px; color: var(--sd-text-secondary); font-size: 14px;">
        修改用户 <strong>{{ editingDisplayNameUser?.username }}</strong> 的显示名称
      </div>
      <n-input v-model:value="newDisplayName" placeholder="输入新的显示名称" @keyup.enter="handleEditDisplayName" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useMessage, NSwitch, NModal, NForm, NFormItem, NInput } from 'naive-ui'
import { api } from '../api'
import type { User, PanelGroup } from '../types'

const message = useMessage()
const users = ref<User[]>([])
const pendingUsers = ref<User[]>([])
const loading = ref(false)
const showAddUser = ref(false)
const modalLoading = ref(false)
const newUser = ref({ username: '', password: '', displayName: '' })

// User bookmarks viewer
const showUserPanel = ref(false)
const viewingUser = ref<User | null>(null)
const viewingGroups = ref<PanelGroup[]>([])

// Global settings
const globalSettings = reactive({
  allow_registration: false,
  require_approval: false,
  default_wallpaper_type: '',
  default_wallpaper_url: '',
  default_wallpaper_blur: '0',
  default_wallpaper_opacity: '100',
  default_wallpaper_copyright: 'false',
})

// Site configuration
const siteConfig = reactive({
  site_title: '',
  site_icon_url: '',
  site_cdn_url: '',
  site_analytics_code: '',
  site_custom_head: '',
  site_custom_footer: '',
})
const siteSaving = ref(false)

onMounted(() => {
  fetchUsers()
  fetchSettings()
})

function statusLabel(status?: string) {
  if (status === 'pending') return '待审批'
  if (status === 'rejected') return '已拒绝'
  return '正常'
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await api.get('users')
    const allUsers = res.data || []
    users.value = (Array.isArray(allUsers) ? allUsers : []).filter((u: User) => u.status !== 'pending')
    pendingUsers.value = (Array.isArray(allUsers) ? allUsers : []).filter((u: User) => u.status === 'pending')
  } catch {
    message.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchSettings() {
  try {
    const res = await api.get('admin/settings')
    if (res.data) {
      Object.assign(globalSettings, {
        allow_registration: res.data.allow_registration === 'true',
        require_approval: res.data.require_approval === 'true',
        default_wallpaper_type: res.data.default_wallpaper_type || '',
        default_wallpaper_url: res.data.default_wallpaper_url || '',
        default_wallpaper_blur: res.data.default_wallpaper_blur || '0',
        default_wallpaper_opacity: res.data.default_wallpaper_opacity || '100',
        default_wallpaper_copyright: res.data.default_wallpaper_copyright || 'false',
      })
      siteConfig.site_title = res.data.site_title || ''
      siteConfig.site_icon_url = res.data.site_icon_url || ''
      siteConfig.site_cdn_url = res.data.site_cdn_url || ''
      siteConfig.site_analytics_code = res.data.site_analytics_code || ''
      siteConfig.site_custom_head = res.data.site_custom_head || ''
      siteConfig.site_custom_footer = res.data.site_custom_footer || ''
    }
  } catch {
    // ignore
  }
}

async function saveGlobalSetting(key: string, value: any) {
  const strValue = typeof value === 'boolean' ? (value ? 'true' : 'false') : String(value)
  try {
    await api.put('admin/settings', { settings: { [key]: strValue } })
    message.success('设置已保存')
  } catch {
    message.error('保存失败')
  }
}

async function saveAllSiteConfig() {
  siteSaving.value = true
  try {
    await api.put('admin/settings', {
      settings: {
        site_title: siteConfig.site_title,
        site_icon_url: siteConfig.site_icon_url,
        site_cdn_url: siteConfig.site_cdn_url,
        site_analytics_code: siteConfig.site_analytics_code,
        site_custom_head: siteConfig.site_custom_head,
        site_custom_footer: siteConfig.site_custom_footer,
      }
    })
    message.success('站点配置已保存，刷新页面后生效')
  } catch {
    message.error('保存失败')
  } finally {
    siteSaving.value = false
  }
}

async function handleCreateUser() {
  if (!newUser.value.username || !newUser.value.password) return false
  modalLoading.value = true
  try {
    await api.post('auth/register', {
      username: newUser.value.username,
      password: newUser.value.password,
      display_name: newUser.value.displayName,
    })
    message.success('用户已创建')
    showAddUser.value = false
    newUser.value = { username: '', password: '', displayName: '' }
    await fetchUsers()
  } catch (e: any) {
    message.error(e.response?.data?.error || '创建失败')
  } finally {
    modalLoading.value = false
  }
  return true
}

async function toggleRole(user: User) {
  const newRole = user.role === 'admin' ? 'user' : 'admin'
  try {
    await api.put(`/api/users/${user.id}`, { role: newRole })
    user.role = newRole
    message.success('角色已更新')
  } catch {
    message.error('更新失败')
  }
}

async function deleteUser(user: User) {
  try {
    await api.delete(`/api/users/${user.id}`)
    message.success('用户已删除')
    await fetchUsers()
  } catch {
    message.error('删除失败')
  }
}

async function approveUser(user: User) {
  try {
    await api.post(`/api/users/${user.id}/approve`)
    message.success(`${user.username} 已批准`)
    pendingUsers.value = pendingUsers.value.filter(u => u.id !== user.id)
    await fetchUsers()
  } catch (e: any) {
    message.error(e.response?.data?.error || '审批失败')
  }
}

async function rejectUser(user: User) {
  try {
    await api.post(`/api/users/${user.id}/reject`)
    message.success(`${user.username} 已拒绝`)
    pendingUsers.value = pendingUsers.value.filter(u => u.id !== user.id)
  } catch (e: any) {
    message.error(e.response?.data?.error || '操作失败')
  }
}

async function viewUserPanel(user: User) {
  viewingUser.value = user
  try {
    const res = await api.get(`/api/users/${user.id}/panel`)
    viewingGroups.value = res.data.groups || []
    showUserPanel.value = true
  } catch {
    message.error('获取用户书签失败')
  }
}

// Edit display name
const showEditDisplayName = ref(false)
const editingDisplayNameUser = ref<User | null>(null)
const newDisplayName = ref('')
const editDisplayNameLoading = ref(false)

function openEditDisplayName(user: User) {
  editingDisplayNameUser.value = user
  newDisplayName.value = user.display_name || ''
  showEditDisplayName.value = true
}

async function handleEditDisplayName() {
  if (!editingDisplayNameUser.value) return false
  if (!newDisplayName.value.trim()) {
    message.warning('显示名称不能为空')
    return false
  }
  editDisplayNameLoading.value = true
  try {
    await api.put(`/api/users/${editingDisplayNameUser.value.id}`, {
      display_name: newDisplayName.value.trim(),
    })
    message.success(`${editingDisplayNameUser.value.username} 的显示名已更新`)
    showEditDisplayName.value = false
    await fetchUsers()
  } catch (e: any) {
    message.error(e.response?.data?.error || '更新失败')
    return false
  } finally {
    editDisplayNameLoading.value = false
  }
  return true
}

// Reset password
const showResetPassword = ref(false)
const resettingUser = ref<User | null>(null)
const newResetPassword = ref('')
const resetPasswordLoading = ref(false)

function openResetPassword(user: User) {
  resettingUser.value = user
  newResetPassword.value = ''
  showResetPassword.value = true
}

async function handleResetPassword() {
  if (!resettingUser.value) return false
  if (!newResetPassword.value || newResetPassword.value.length < 6) {
    message.warning('密码至少需要6位')
    return false
  }
  resetPasswordLoading.value = true
  try {
    await api.post(`/api/users/${resettingUser.value.id}/reset-password`, {
      new_password: newResetPassword.value,
    })
    message.success(`${resettingUser.value.username} 的密码已重置`)
    showResetPassword.value = false
  } catch (e: any) {
    message.error(e.response?.data?.error || '重置失败')
    return false
  } finally {
    resetPasswordLoading.value = false
  }
  return true
}
</script>

<style scoped>
.admin-page {
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

/* Header */
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

.back-btn:hover {
  background: rgba(0,122,255,0.08);
}

.page-header h1 {
  font-size: 17px;
  font-weight: 600;
  margin: 0;
  color: var(--sd-text-primary);
}

/* Content */
.admin-content {
  position: relative;
  z-index: 1;
  max-width: 900px;
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

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 8px;
  background: #FF3B30;
  color: white;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
}

.add-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 0 18px;
  height: 38px;
  background: var(--sd-primary);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: all 0.2s ease;
}

.add-btn:hover {
  background: var(--sd-primary-hover);
  transform: translateY(-1px);
}

/* Settings Card */
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
  padding: 12px 14px;
  min-height: 48px;
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
  margin: 0 14px;
}

:root[data-theme="dark"] .setting-divider {
  background: rgba(255,255,255,0.06);
}

.setting-sub {
  padding: 8px 14px 14px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.seg-control {
  display: flex;
  background: rgba(0,0,0,0.04);
  border-radius: 8px;
  padding: 2px;
  gap: 2px;
}

:root[data-theme="dark"] .seg-control {
  background: rgba(255,255,255,0.06);
}

.seg-btn {
  padding: 6px 12px;
  font-size: 12px;
  font-weight: 500;
  font-family: var(--sd-font);
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--sd-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.seg-btn.active {
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
}

:root[data-theme="dark"] .seg-btn.active {
  background: rgba(44,44,46,0.9);
  box-shadow: 0 1px 3px rgba(0,0,0,0.3);
}

.bg-url-input {
  flex: 1;
  height: 32px;
  padding: 0 10px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 12px;
  font-family: var(--sd-font);
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  outline: none;
  transition: border-color 0.2s;
}

.bg-url-input:focus {
  border-color: #007AFF;
}

.mini-text-input {
  width: 160px;
  height: 32px;
  padding: 0 10px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 13px;
  font-family: var(--sd-font);
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  outline: none;
}

.mini-text-input:focus {
  border-color: #007AFF;
}

.mini-text-input.wide {
  width: 260px;
}

.code-input {
  width: 100%;
  padding: 10px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 12px;
  font-family: monospace;
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  resize: vertical;
  outline: none;
}

.code-input:focus {
  border-color: #007AFF;
}

/* Table */
.table-card {
  background: rgba(255,255,255,0.85);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-radius: 14px;
  border: 1px solid rgba(0,0,0,0.06);
  overflow: hidden;
}

:root[data-theme="dark"] .table-card {
  background: rgba(28,28,30,0.8);
  border-color: rgba(255,255,255,0.08);
}

.table-head {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  background: rgba(0,0,0,0.03);
  border-bottom: 1px solid rgba(0,0,0,0.06);
}

:root[data-theme="dark"] .table-head {
  background: rgba(255,255,255,0.03);
  border-bottom-color: rgba(255,255,255,0.06);
}

.th {
  font-size: 12px;
  font-weight: 600;
  color: var(--sd-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.table-row {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(0,0,0,0.04);
  transition: background 0.15s ease;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row:hover {
  background: rgba(0,0,0,0.02);
}

:root[data-theme="dark"] .table-row:hover {
  background: rgba(255,255,255,0.03);
}

.td {
  font-size: 14px;
  color: var(--sd-text-primary);
  display: flex;
  align-items: center;
  gap: 10px;
}

.u-avatar {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0,122,255,0.1);
  color: var(--sd-primary);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.u-avatar.pending {
  background: rgba(255,149,0,0.12);
  color: #FF9500;
}

.u-name {
  font-weight: 500;
}

.role-pill, .status-pill {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 12px;
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

.status-pill.is-approved {
  background: rgba(52,199,89,0.1);
  color: #34C759;
}

.status-pill.is-pending {
  background: rgba(255,149,0,0.12);
  color: #FF9500;
}

.status-pill.is-rejected {
  background: rgba(255,59,48,0.1);
  color: #FF3B30;
}

.row-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  background: transparent;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 6px;
  color: var(--sd-text-secondary);
  font-size: 12px;
  font-family: var(--sd-font);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
}

:root[data-theme="dark"] .row-btn {
  border-color: rgba(255,255,255,0.12);
}

.row-btn:hover {
  background: rgba(0,0,0,0.05);
  color: var(--sd-text-primary);
}

.row-btn.approve {
  color: #34C759;
  border-color: rgba(52,199,89,0.3);
}

.row-btn.approve:hover {
  background: rgba(52,199,89,0.08);
}

.row-btn.danger {
  color: #FF3B30;
  border-color: rgba(255,59,48,0.2);
}

.row-btn.danger:hover {
  background: rgba(255,59,48,0.08);
}

.table-empty {
  padding: 48px;
  text-align: center;
  color: var(--sd-text-tertiary);
  font-size: 14px;
}

/* User Bookmarks Viewer */
.view-group {
  margin-bottom: 16px;
}

.view-group-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--sd-text-primary);
  margin-bottom: 8px;
}

.view-card-count {
  font-size: 12px;
  font-weight: 400;
  color: var(--sd-text-tertiary);
}

.view-cards {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-left: 24px;
}

.view-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 10px;
  border-radius: 6px;
  background: rgba(0,0,0,0.02);
}

:root[data-theme="dark"] .view-card {
  background: rgba(255,255,255,0.03);
}

.view-card-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--sd-text-primary);
  min-width: 80px;
}

.view-card-url {
  font-size: 12px;
  color: var(--sd-text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}


/* === Responsive === */
@media (max-width: 768px) {
  .admin-content {
    padding: var(--sd-space-4) var(--sd-space-3) var(--sd-space-16);
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  /* Hide role/status columns on mobile, use card layout */
  .table-head .th:nth-child(3),
  .table-head .th:nth-child(4),
  .table-row .td:nth-child(3),
  .table-row .td:nth-child(4) {
    display: none;
  }

  .th, .td {
    font-size: 12px;
  }

  .td-actions {
    gap: 3px;
  }

  .row-btn {
    padding: 3px 6px;
    font-size: 10px;
  }

  .row-btn span {
    display: none;
  }

  .setting-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .mini-text-input,
  .mini-text-input.wide {
    width: 100%;
  }
}
</style>