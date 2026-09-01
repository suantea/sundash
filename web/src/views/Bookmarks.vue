<template>
  <div class="bookmarks-page">
    <header class="bm-header">
      <div class="bm-header-left">
        <button class="bm-back" @click="router.push('/')" title="返回主页">
          <Icon icon="mdi:arrow-left" :size="18" />
        </button>
        <h1 class="bm-title">书签同步</h1>
        <n-tag v-if="!status.configured" type="warning" size="small" round>
          未配置同步服务器
        </n-tag>
        <n-tag v-else-if="status.hasSynced" type="success" size="small" round>
          已同步 · rev {{ status.rev }}
        </n-tag>
        <n-tag v-else type="info" size="small" round>
          尚未同步
        </n-tag>
      </div>
      <div class="bm-header-actions">
        <n-input v-model:value="keyword" placeholder="搜索书签…" clearable class="bm-search">
          <template #prefix><Icon icon="mdi:magnify" :size="15" /></template>
        </n-input>
        <n-button size="medium" @click="refresh">
          <template #icon><Icon icon="mdi:refresh" :size="16" /></template>
          同步
        </n-button>
        <n-button size="medium" type="primary" @click="openCreate('folder')">
          <template #icon><Icon icon="mdi:folder-plus" :size="16" /></template>
          新建文件夹
        </n-button>
        <n-button size="medium" type="primary" secondary @click="openCreate('bookmark')">
          <template #icon><Icon icon="mdi:bookmark-plus" :size="16" /></template>
          新建书签
        </n-button>
      </div>
    </header>

    <!-- 未配置引导 -->
    <n-empty v-if="!status.configured" description="请先在「设置 → 管理面板」中填写 bookmark-sync 服务器地址与 Token">
      <template #extra>
        <n-button size="small" @click="router.push('/admin')">前往设置</n-button>
      </template>
    </n-empty>

    <!-- 已配置但无数据 -->
    <n-empty v-else-if="loading" description="正在同步书签…">
      <template #extra><n-spin size="small" /></template>
    </n-empty>
    <n-empty v-else-if="!hasContent" description="暂无书签，点击右上角新建，或先在 Chrome / Safari 扩展中导入">
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
      :title="editing.syncId ? '编辑' : (form.type === 'folder' ? '新建文件夹' : '新建书签')"
    >
      <n-form label-placement="top">
        <n-form-item label="名称">
          <n-input v-model:value="form.title" placeholder="名称" />
        </n-form-item>
        <n-form-item v-if="form.type === 'bookmark'" label="链接">
          <n-input v-model:value="form.url" placeholder="https://…" />
        </n-form-item>
        <n-form-item label="所在文件夹">
          <n-select
            v-model:value="form.parentSyncId"
            :options="folderOptions"
            placeholder="根目录"
            clearable
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="bm-modal-footer">
          <n-button @click="showEditor = false">取消</n-button>
          <n-button type="primary" :loading="saving" @click="save">保存</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import { useMessage } from 'naive-ui'
import { Icon } from '@iconify/vue'
import router from '../router'
import BookmarkNode from '../components/bookmark/BookmarkNode.vue'
import { getStatus, getTree, pullTree, pushChanges, type SyncNode } from '../api/bmsync'

const message = useMessage()

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
    message.warning('请填写名称')
    return
  }
  if (form.value.type === 'bookmark' && !form.value.url.trim()) {
    message.warning('请填写链接')
    return
  }
  saving.value = true
  try {
    const op = editing.value.syncId ? 'update' : 'create'
    const { data } = await pushChanges([{
      op,
      syncId: editing.value.syncId || undefined,
      type: form.value.type,
      title: form.value.title.trim(),
      url: form.value.type === 'bookmark' ? form.value.url.trim() : undefined,
      parentSyncId: form.value.parentSyncId || undefined,
    }])
    nodes.value = data.nodes
    status.value.rev = data.rev
    status.value.hasSynced = true
    showEditor.value = false
    message.success(op === 'create' ? '已创建' : '已保存')
  } catch (e: any) {
    message.error(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

// ── 删除（墓碑语义） ──────────────────────
function confirmDelete(n: SyncNode) {
  const d = window.confirm(`确定删除「${n.title}」吗？\n删除会同步到其他电脑，Chrome / Safari 中的书签也会被移除。`)
  if (!d) return
  pushChanges([{ op: 'delete', syncId: n.syncId }])
    .then(({ data }) => {
      nodes.value = data.nodes
      status.value.rev = data.rev
      message.success('已删除（已同步到其他设备）')
    })
    .catch((e: any) => message.error(e.response?.data?.error || '删除失败'))
}

// ── 拉取 ────────────────────────────────
async function refresh() {
  loading.value = true
  try {
    const { data } = await pullTree()
    nodes.value = data.nodes
    status.value.rev = data.rev
    status.value.hasSynced = true
    message.success('同步完成')
  } catch (e: any) {
    message.error(e.response?.data?.error || '同步失败')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    status.value = await getStatus()
    if (status.value.configured && status.value.hasSynced) {
      const { data } = await getTree()
      nodes.value = data.nodes
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
