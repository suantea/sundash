<template>
  <div class="bookmarks-page">
    <header class="bm-header">
      <div class="bm-header-left">
        <button class="bm-back" @click="router.push('/')" :title="$t('bookmarks.backToHome')">
          <Icon icon="mdi:arrow-left" :width="18" :height="18" />
        </button>
        <h1 class="bm-title">{{ $t('bookmarks.title') }}</h1>
        <n-tag v-if="!status.configured" type="warning" size="small" round>
          {{ $t('bookmarks.tagUnconfigured') }}
        </n-tag>
        <n-tag v-else-if="status.hasSynced" type="success" size="small" round>
          {{ $t('bookmarks.tagSynced') }} · rev {{ status.rev }}
        </n-tag>
        <n-tag v-else type="info" size="small" round>
          {{ $t('bookmarks.tagNotSynced') }}
        </n-tag>
      </div>
      <div class="bm-header-actions">
        <n-input v-model:value="keyword" :placeholder="$t('bookmarks.searchPlaceholder')" clearable class="bm-search">
          <template #prefix><Icon icon="mdi:magnify" :width="15" :height="15" /></template>
        </n-input>
        <n-button size="medium" @click="refresh">
          <template #icon><Icon icon="mdi:refresh" :width="16" :height="16" /></template>
          {{ $t('bookmarks.sync') }}
        </n-button>
        <n-button size="medium" type="primary" @click="openCreate('folder')">
          <template #icon><Icon icon="mdi:folder-plus" :width="16" :height="16" /></template>
          {{ $t('bookmarks.newFolder') }}
        </n-button>
        <n-button size="medium" type="primary" secondary @click="openCreate('bookmark')">
          <template #icon><Icon icon="mdi:bookmark-plus" :width="16" :height="16" /></template>
          {{ $t('bookmarks.newBookmark') }}
        </n-button>
      </div>
    </header>

    <!-- 未配置引导：管理员 → 直达管理面板；普通用户 → 联系管理员 -->
    <n-empty v-if="!status.configured" :description="isAdmin ? $t('bookmarks.unconfiguredAdmin') : $t('bookmarks.unconfiguredUser')">
      <template #extra>
        <n-button v-if="isAdmin" size="small" type="primary" @click="router.push('/admin')">
          {{ $t('bookmarks.goToAdmin') }}
        </n-button>
      </template>
    </n-empty>

    <!-- 已配置但无数据 -->
    <n-empty v-else-if="loading" :description="$t('bookmarks.syncingDesc')">
      <template #extra><n-spin size="small" /></template>
    </n-empty>
    <n-empty v-else-if="!hasContent" :description="$t('bookmarks.emptyDesc')">
    </n-empty>

    <!-- 书签树 -->
    <div v-else class="bm-tree">
      <BookmarkNode
        v-for="root in visibleRoots"
        :key="root.syncId"
        :node="root"
        :nodes="nodesById"
        :keyword="keyword"
        @edit="openEdit"
        @delete="confirmDelete"
      />
    </div>

    <!-- 新建 / 编辑弹窗 -->
    <n-modal
      v-model:show="showEditor" preset="card" style="width: 460px"
      :title="editing.syncId ? $t('bookmarks.edit') : (form.type === 'folder' ? $t('bookmarks.newFolder') : $t('bookmarks.newBookmark'))"
    >
      <n-form label-placement="top">
        <n-form-item :label="$t('bookmarks.name')">
          <n-input v-model:value="form.title" :placeholder="$t('bookmarks.name')" />
        </n-form-item>
        <n-form-item v-if="form.type === 'bookmark'" :label="$t('bookmarks.link')">
          <n-input v-model:value="form.url" placeholder="https://…" />
        </n-form-item>
        <n-form-item :label="$t('bookmarks.parentFolder')">
          <n-select
            v-model:value="form.parentSyncId"
            :options="folderOptions"
            :placeholder="$t('bookmarks.rootFolder')"
            clearable
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="bm-modal-footer">
          <n-button @click="showEditor = false">{{ $t('common.cancel') }}</n-button>
          <n-button type="primary" :loading="saving" @click="save">{{ $t('bookmarks.save') }}</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMessage, NTag, NButton, NInput, NSelect, NForm, NFormItem, NModal, NEmpty, NSpin } from 'naive-ui'
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import router from '../router'
import { useUserStore } from '../stores/user'
import BookmarkNode from '../components/bookmark/BookmarkNode.vue'
import { getStatus, getTree, pullTree, pushChanges, type SyncNode } from '../api/bmsync'

const { t } = useI18n()
const userStore = useUserStore()
const message = useMessage()

// 仅管理员能进入管理面板完成配置（/admin 路由有权限守卫）
const isAdmin = computed(() => userStore.user?.role === 'admin')

const status = ref({ configured: false, serverUrl: '', hasSynced: false, rev: 0 })
const nodes = ref<SyncNode[]>([])
const loading = ref(true)
const keyword = ref('')

const nodesById = computed(() => {
  const m = new Map<string, SyncNode>()
  for (const n of nodes.value) m.set(n.syncId, n)
  return m
})

// 活跃节点（非墓碑）+ 顶层根
const activeNodes = computed(() => nodes.value.filter((n) => !n.deletedAt))
const rootNodes = computed(() => activeNodes.value.filter((n) => !n.parentSyncId || !nodesById.value.has(n.parentSyncId!) || nodesById.value.get(n.parentSyncId!)?.deletedAt))

const hasContent = computed(() => activeNodes.value.length > 0)

// 搜索过滤后的根节点（子节点是否显示由 BookmarkNode 内部判断）
const visibleRoots = computed(() => {
  if (!keyword.value.trim()) return rootNodes.value
  const kw = keyword.value.trim().toLowerCase()
  const match = (n: SyncNode) =>
    n.title.toLowerCase().includes(kw) || (n.url || '').toLowerCase().includes(kw)
  return rootNodes.value.filter((n) => subtreeMatches(n, match))
})

function subtreeMatches(n: SyncNode, match: (n: SyncNode) => boolean): boolean {
  if (match(n)) return true
  const children = activeNodes.value.filter((c) => c.parentSyncId === n.syncId)
  return children.some((c) => subtreeMatches(c, match))
}

// ── 编辑弹窗 ──────────────────────────────
const showEditor = ref(false)
const saving = ref(false)
const editing = ref<SyncNode>({ syncId: '', type: 'bookmark', title: '', url: '', parentSyncId: '', index: 0, createdAt: '', updatedAt: '' })
const form = ref({ type: 'bookmark' as 'folder' | 'bookmark', title: '', url: '', parentSyncId: '' })

function openCreate(type: 'folder' | 'bookmark') {
  editing.value = { syncId: '', type, title: '', url: '', parentSyncId: '', index: 0, createdAt: '', updatedAt: '' }
  form.value = { type, title: '', url: '', parentSyncId: '' }
  showEditor.value = true
}

function openEdit(n: SyncNode) {
  editing.value = n
  form.value = { type: n.type, title: n.title, url: n.url || '', parentSyncId: n.parentSyncId || '' }
  showEditor.value = true
}

const folderOptions = computed(() =>
  activeNodes.value
    .filter((n) => n.type === 'folder' && n.syncId !== editing.value.syncId)
    .map((n) => ({ label: n.title, value: n.syncId }))
)

async function save() {
  if (!form.value.title.trim()) {
    message.warning(t('bookmarks.nameRequired'))
    return
  }
  if (form.value.type === 'bookmark' && !form.value.url.trim()) {
    message.warning(t('bookmarks.linkRequired'))
    return
  }
  saving.value = true
  try {
    const op = editing.value.syncId ? 'update' : 'create'
    const res = await pushChanges([{
      op,
      syncId: editing.value.syncId || undefined,
      type: form.value.type,
      title: form.value.title.trim(),
      url: form.value.type === 'bookmark' ? form.value.url.trim() : undefined,
      parentSyncId: form.value.parentSyncId || undefined,
    }])
    nodes.value = res.nodes
    status.value.rev = res.rev
    status.value.hasSynced = true
    showEditor.value = false
    message.success(op === 'create' ? t('bookmarks.created') : t('bookmarks.saved'))
  } catch (e: any) {
    message.error(e.response?.data?.error || t('bookmarks.saveFailed'))
  } finally {
    saving.value = false
  }
}

// ── 删除（墓碑语义） ──────────────────────
function confirmDelete(n: SyncNode) {
  const d = window.confirm(t('bookmarks.deleteConfirm', { name: n.title }))
  if (!d) return
  pushChanges([{ op: 'delete', syncId: n.syncId }])
    .then((res) => {
      nodes.value = res.nodes
      status.value.rev = res.rev
      message.success(t('bookmarks.deleted'))
    })
    .catch((e: any) => message.error(e.response?.data?.error || t('bookmarks.deleteFailed')))
}

// ── 拉取 ────────────────────────────────
async function refresh() {
  loading.value = true
  try {
    const res = await pullTree()
    nodes.value = res.nodes
    status.value.rev = res.rev
    status.value.hasSynced = true
    message.success(t('bookmarks.syncComplete'))
  } catch (e: any) {
    message.error(e.response?.data?.error || t('bookmarks.syncFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    status.value = await getStatus()
    if (status.value.configured && status.value.hasSynced) {
      const res = await getTree()
      nodes.value = res.nodes
    }
  } catch {
    // 未配置等错误由 empty 状态引导
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.bookmarks-page {
  min-height: 100vh;
  padding: 28px 32px;
  box-sizing: border-box;
  background: var(--sd-bg, #f5f6f8);
}
[data-theme='dark'] .bookmarks-page {
  background: #141416;
  color: #e8e8ea;
}
.bm-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}
.bm-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.bm-back {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 10px;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: inherit;
}
.bm-back:hover {
  background: rgba(128, 128, 128, 0.14);
}
.bm-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}
.bm-header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.bm-search {
  width: 240px;
}
.bm-tree {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.bm-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
