<template>
  <div v-if="appStore.showMemo" class="memo-widget" :class="{ 'expanded': isExpanded }">
    <!-- Compact bar -->
    <div class="memo-bar" @click="isExpanded = !isExpanded">
      <div class="memo-item count" :title="t('memo.count', { count: memos.length })">
        <Icon icon="mdi:note-text" :size="13" />
        <span>{{ memos.length }}</span>
      </div>
      <div class="memo-item add" @click="showAddMemo = true" :title="t('memo.add')">
        <Icon icon="mdi:plus" :size="13" />
      </div>
      <div class="memo-expand">
        <Icon :icon="isExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'" :size="14" />
      </div>
    </div>

    <!-- Expanded detail panel -->
    <Transition name="slide">
      <div v-if="isExpanded" class="memo-detail">
        <div class="memo-header">
          <div class="memo-title">{{ t('memo.title') }}</div>
          <div class="memo-actions">
            <button class="memo-action-btn" @click="showAddMemo = true">
              <Icon icon="mdi:plus" /> {{ t('memo.add') }}
            </button>
            <button class="memo-action-btn" @click="archiveAll" :disabled="memos.length === 0">
              <Icon icon="mdi:archive" /> {{ t('memo.archiveAll') }}
            </button>
          </div>
        </div>
        <div class="memo-list">
          <div v-for="memo in memos" :key="memo.id" class="memo-item">
            <div class="memo-content">{{ memo.content }}</div>
            <div class="memo-footer">
              <span class="memo-time">{{ formatTime(memo.updated_at) }}</span>
              <div class="memo-actions">
                <button class="memo-action-btn" @click="toggleArchive(memo)" :title="memo.is_archived ? t('memo.unarchive') : t('memo.archive')">
                  <Icon :icon="memo.is_archived ? 'mdi:folder-outline' : 'mdi:folder'" />
                </button>
                <button class="memo-action-btn" @click="deleteMemo(memo.id)">
                  <Icon icon="mdi:delete" />
                </button>
              </div>
            </div>
          </div>
          <div v-if="memos.length === 0" class="empty-state">
            {{ t('memo.empty') }}
          </div>
        </div>
      </div>
    </Transition>

    <!-- Add memo modal -->
    <Teleport to="body">
      <div v-if="showAddMemo" class="modal-backdrop" @click.self="showAddMemo = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>{{ t('memo.add') }}</h3>
            <button class="modal-close" @click="showAddMemo = false">
              <Icon icon="mdi:close" />
            </button>
          </div>
          <div class="modal-body">
            <textarea v-model="newMemoContent" :placeholder="t('memo.placeholder')" rows="4" class="memo-input" @keyup.enter="saveMemo" />
          </div>
          <div class="modal-footer">
            <button class="modal-btn" @click="showAddMemo = false">{{ t('common.cancel') }}</button>
            <button class="modal-btn primary" @click="saveMemo">{{ t('common.save') }}</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAppStore } from '../../stores/app'
import { useI18n } from 'vue-i18n'
import { api as axios } from '@/api'

const appStore = useAppStore()
const { t } = useI18n()
const memos = ref<any[]>([])
const isExpanded = ref(false)
const showAddMemo = ref(false)
const newMemoContent = ref('')
let timer: ReturnType<typeof setInterval> | null = null

async function fetchMemos() {
  try {
    const token = localStorage.getItem('sundash-token')
    const res = await axios.get('memo', {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (res.status === 200) {
      memos.value = res.data
    }
  } catch (err) {
    console.error('Failed to fetch memos:', err)
  }
}

async function saveMemo() {
  const content = newMemoContent.value.trim()
  if (!content) return
  try {
    const token = localStorage.getItem('sundash-token')
    await axios.post('memo', { content }, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    newMemoContent.value = ''
    showAddMemo.value = false
    await fetchMemos()
  } catch (err) {
    console.error('Failed to save memo:', err)
  }
}

async function toggleArchive(memo: any) {
  try {
    const token = localStorage.getItem('sundash-token')
    await axios.put(`memo/${memo.id}/archive`, { archived: !memo.is_archived }, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    await fetchMemos()
  } catch (err) {
    console.error('Failed to toggle archive:', err)
  }
}

async function deleteMemo(id: string) {
  if (!confirm(t('memo.confirmDelete'))) return
  try {
    const token = localStorage.getItem('sundash-token')
    await axios.delete(`memo/${id}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    await fetchMemos()
  } catch (err) {
    console.error('Failed to delete memo:', err)
  }
}

async function archiveAll() {
  if (!confirm(t('memo.confirmArchiveAll'))) return
  try {
    const token = localStorage.getItem('sundash-token')
    // We'll archive each memo individually for simplicity
    for (const memo of memos.value) {
      if (!memo.is_archived) {
        await axios.put(`memo/${memo.id}/archive`, { archived: true }, {
          headers: token ? { Authorization: `Bearer ${token}` } : {}
        })
      }
    }
    await fetchMemos()
  } catch (err) {
    console.error('Failed to archive all:', err)
  }
}

function formatTime(timeString: string | undefined): string {
  if (!timeString) return '--:--'
  try {
    const date = new Date(timeString)
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } catch {
    return timeString
  }
}

onMounted(() => {
  fetchMemos()
  timer = setInterval(fetchMemos, 30 * 1000) // 30 seconds
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.memo-widget {
  background: var(--sd-bg-elevated);
  border: 1px solid var(--sd-border);
  border-radius: 12px;
  overflow: hidden;
  margin: 0 16px 12px;
  backdrop-filter: blur(10px);
}

.memo-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 14px;
  cursor: pointer;
  user-select: none;
}

.memo-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--sd-text-secondary);
}

.memo-item span {
  font-variant-numeric: tabular-nums;
}

.memo-item.add {
  margin-left: auto;
  cursor: pointer;
  color: var(--sd-color-success);
}

.memo-item.add:hover {
  color: var(--sd-color-success-dark);
}

.memo-expand {
  margin-left: auto;
  color: var(--sd-text-tertiary);
}

/* Expanded detail */
.memo-detail {
  padding: 12px 14px 14px;
  border-top: 1px solid var(--sd-border);
}

.memo-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.memo-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--sd-text-primary);
}

.memo-actions {
  display: flex;
  gap: 8px;
}

.memo-action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--sd-text-tertiary);
  background: transparent;
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.memo-action-btn:hover {
  background: var(--sd-bg-base);
  color: var(--sd-text-primary);
}

.memo-action-btn Icon {
  font-size: 14px;
}

.memo-list {
  max-height: 300px;
  overflow-y: auto;
}

.memo-item {
  background: var(--sd-bg-subtle);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 8px;
  border: 1px solid var(--sd-border);
}

.memo-content {
  font-size: 14px;
  color: var(--sd-text-primary);
  line-height: 1.5;
  margin-bottom: 8px;
  word-break: break-word;
}

.memo-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--sd-text-tertiary);
}

.memo-time {
  font-variant-numeric: tabular-nums;
}

.memo-actions .memo-action-btn {
  padding: 2px 6px;
}

.empty-state {
  text-align: center;
  color: var(--sd-text-tertiary);
  font-size: 14px;
  padding: 20px;
}

/* Modal */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: var(--sd-bg-elevated);
  border-radius: 12px;
  width: 90%;
  max-width: 400px;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--sd-border);
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--sd-text-primary);
}

.modal-close {
  background: transparent;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: var(--sd-text-tertiary);
}

.modal-close:hover {
  color: var(--sd-text-primary);
}

.modal-body {
  flex: 1;
  padding: 16px;
}

.memo-input {
  width: 100%;
  height: 100%;
  border: none;
  resize: none;
  background: transparent;
  font-size: 14px;
  color: var(--sd-text-primary);
  padding: 8px;
}

.memo-input:focus {
  outline: none;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--sd-border);
}

.modal-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.modal-btn:hover {
  opacity: 0.9;
}

.modal-btn.primary {
  background: var(--sd-primary);
  color: white;
}

.modal-btn.primary:hover {
  background: var(--sd-primary-dark);
}

/* Scrollbar styling */
.memo-list::-webkit-scrollbar {
  width: 6px;
}

.memo-list::-webkit-scrollbar-track {
  background: var(--sd-bg-base);
}

.memo-list::-webkit-scrollbar-thumb {
  background: var(--sd-border);
  border-radius: 3px;
}

.memo-list::-webkit-scrollbar-thumb:hover {
  background: var(--sd-text-tertiary);
}
</style>
