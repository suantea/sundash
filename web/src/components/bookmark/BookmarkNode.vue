<template>
  <div class="bm-node">
    <div class="bm-row" :class="{ hovered: hover }" @mouseenter="hover = true" @mouseleave="hover = false">
      <template v-if="isFolder">
        <button class="bm-caret" @click="toggle">
          <Icon :icon="collapsed ? 'mdi:chevron-right' : 'mdi:chevron-down'" :width="16" :height="16" />
        </button>
        <Icon icon="mdi:folder" :width="17" :height="17" class="bm-folder-icon" />
        <span class="bm-label" @dblclick="emit('edit', node)">{{ node.title }}</span>
      </template>
      <template v-else>
        <span class="bm-caret-spacer"></span>
        <Icon icon="mdi:bookmark" :width="15" :height="15" class="bm-bookmark-icon" />
        <a
          class="bm-label bm-link"
          :href="node.url || '#'"
          target="_blank"
          rel="noopener noreferrer"
          :title="node.url"
        >{{ node.title }}</a>
      </template>

      <span v-if="isFolder && matchesKeyword" class="bm-match-count">{{ matchedChildrenCount }} 个匹配</span>

      <div v-if="hover" class="bm-actions">
        <button class="bm-act" :title="$t('common.edit')" @click="emit('edit', node)">
          <Icon icon="mdi:pencil-outline" :width="14" :height="14" />
        </button>
        <button class="bm-act danger" :title="$t('common.delete')" @click="emit('delete', node)">
          <Icon icon="mdi:delete-outline" :width="14" :height="14" />
        </button>
      </div>
    </div>

    <div v-if="isFolder && !collapsed && hasVisibleChildren" class="bm-children">
      <BookmarkNode
        v-for="child in visibleChildren"
        :key="child.syncId"
        :node="child"
        :nodes="nodes"
        :keyword="keyword"
        @edit="(n: SyncNode) => emit('edit', n)"
        @delete="(n: SyncNode) => emit('delete', n)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Icon } from '@iconify/vue'
import type { SyncNode } from '../../api/bmsync'

const props = defineProps<{
  node: SyncNode
  nodes: Map<string, SyncNode>
  keyword: string
}>()

const emit = defineEmits<{
  (e: 'edit', node: SyncNode): void
  (e: 'delete', node: SyncNode): void
}>()

const hover = ref(false)
const collapsed = ref(false)
const isFolder = computed(() => props.node.type === 'folder')

const children = computed(() => {
  const out: SyncNode[] = []
  props.nodes.forEach((n) => {
    if (n.parentSyncId === props.node.syncId && !n.deletedAt) out.push(n)
  })
  out.sort((a, b) => a.index - b.index || a.title.localeCompare(b.title, 'zh'))
  return out
})

const kw = computed(() => props.keyword.trim().toLowerCase())
const matchesKeyword = computed(() => {
  if (!kw.value) return false
  const n = props.node
  return n.title.toLowerCase().includes(kw.value) || (n.url || '').toLowerCase().includes(kw.value)
})

function subtreeMatchCount(n: SyncNode): number {
  let count = 0
  if (n.title.toLowerCase().includes(kw.value) || (n.url || '').toLowerCase().includes(kw.value)) count++
  props.nodes.forEach((c) => {
    if (c.parentSyncId === n.syncId && !c.deletedAt) count += subtreeMatchCount(c)
  })
  return count
}

const matchedChildrenCount = computed(() => (kw.value ? subtreeMatchCount(props.node) : 0))

// 搜索时：命中节点或其子命中节点都要展开
const searchExpanded = computed(() => kw.value && (matchesKeyword.value || children.value.some((c) => subtreeMatchCount(c) > 0)))

const visibleChildren = computed(() => {
  if (!kw.value) return children.value
  return children.value.filter((c) => subtreeMatchCount(c) > 0 || c.title.toLowerCase().includes(kw.value) || (c.url || '').toLowerCase().includes(kw.value))
})

const hasVisibleChildren = computed(() => visibleChildren.value.length > 0)

function toggle() {
  collapsed.value = !collapsed.value
}
</script>

<style scoped>
.bm-node {
  display: flex;
  flex-direction: column;
}
.bm-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 10px;
  min-height: 34px;
}
.bm-row.hovered {
  background: rgba(128, 128, 128, 0.12);
}
.bm-caret {
  width: 20px;
  height: 20px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: inherit;
  border-radius: 6px;
}
.bm-caret:hover {
  background: rgba(128, 128, 128, 0.18);
}
.bm-caret-spacer {
  width: 20px;
  flex: none;
}
.bm-folder-icon {
  color: #f5b50a;
  flex: none;
}
.bm-bookmark-icon {
  color: var(--sd-primary, #007aff);
  flex: none;
}
.bm-label {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
}
.bm-link {
  color: inherit;
  text-decoration: none;
}
.bm-link:hover {
  color: var(--sd-primary, #007aff);
}
.bm-match-count {
  font-size: 12px;
  color: rgba(128, 128, 128, 0.9);
  flex: none;
}
.bm-actions {
  display: flex;
  gap: 4px;
  flex: none;
}
.bm-act {
  width: 26px;
  height: 26px;
  border: none;
  border-radius: 7px;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: inherit;
}
.bm-act:hover {
  background: rgba(128, 128, 128, 0.2);
}
.bm-act.danger:hover {
  background: rgba(229, 72, 77, 0.15);
  color: #e5484d;
}
.bm-children {
  margin-left: 16px;
  padding-left: 8px;
  border-left: 1px solid rgba(128, 128, 128, 0.18);
}
</style>
