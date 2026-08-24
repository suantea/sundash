<template>
  <div v-if="appStore.showRSSWidget" class="rss-widget" :class="{ 'expanded': isExpanded }">
    <!-- Compact bar -->
    <div class="rss-bar" @click="isExpanded = !isExpanded">
      <div class="rss-item count" :title="'订阅数量: ' + feeds.length">
        <Icon icon="mdi:rss" :size="13" />
        <span>{{ feeds.length }}</span>
      </div>
      <div class="rss-item add" @click="showAddFeed = true" title="新增订阅">
        <Icon icon="mdi:plus" :size="13" />
      </div>
      <div class="rss-expand">
        <Icon :icon="isExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'" :size="14" />
      </div>
    </div>

    <!-- Expanded detail panel -->
    <Transition name="slide">
      <div v-if="isExpanded" class="rss-detail">
        <div class="rss-header">
          <div class="rss-title">我的订阅</div>
          <div class="rss-actions">
            <button class="rss-action-btn" @click="showAddFeed = true">
              <Icon icon="mdi:plus" /> 新建
            </button>
            <button class="rss-action-btn" @click="refreshAll" :disabled="feeds.length === 0">
              <Icon icon="mdi:refresh" /> 刷新全部
            </button>
          </div>
        </div>
        <div class="rss-list">
          <div v-for="feed in feeds" :key="feed.id" class="rss-item">
            <div class="rss-header">
              <div class="rss-title">{{ feed.title }}</div>
              <div class="rss-actions">
                <button class="rss-action-btn" @click="toggleFeed(feed)" :title="feed.isExpanded ? '收起' : '展开'">
                  <Icon :icon="feed.isExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'" />
                </button>
                <button class="rss-action-btn" @click="editFeed(feed.id)">
                  <Icon icon="mdi:lead-pencil" />
                </button>
                <button class="rss-action-btn" @confirm="deleteFeed(feed.id)" @cancel="clearDeleteConfirm" :loading="deletingFeedId === feed.id" @show="showDeleteConfirm = true">
                  <Icon icon="mdi:delete" />
                </button>
              </div>
            </div>
            <!-- Feed content (expanded) -->
            <div v-if="feed.isExpanded" class="rss-content">
              <div v-if="feed.items.length === 0" class="empty-state">
                暂无文章，点击刷新获取最新内容。
              </div>
              <div v-else class="rss-feed-items">
                <div v-for="item in feed.items" :key="item.id" class="rss-item">
                  <div class="rss-item-title">
                    <a :href="item.link" target="_blank" rel="noopener">{{ item.title }}</a>
                  </div>
                  <div class="rss-item-meta">
                    <span class="rss-item-date">{{ formatDate(item.pub_date) }}</span>
                    <span v-if="item.author" class="rss-item-author">By {{ item.author }}</span>
                  </div>
                  <div class="rss-item-description" v-html="item.description"></div>
                </div>
              </div>
            </div>
          </div>
          <div v-if="feeds.length === 0" class="empty-state">
            您还没有添加任何 RSS 订阅，点击「新建」添加第一个订阅吧！
          </div>
        </div>
      </div>
    </Transition>

    <!-- Add feed modal -->
    <Teleport to="body">
      <div v-if="showAddFeed" class="modal-backdrop" @click.self="showAddFeed = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>添加 RSS 订阅</h3>
            <button class="modal-close" @click="showAddFeed = false">
              <Icon icon="mdi:close" />
            </button>
          </div>
          <div class="modal-body">
            <input v-model="newFeedUrl" placeholder="输入 RSS Feed URL..." type="text" class="rss-input" @keyup.enter="addFeed" />
          </div>
          <div class="modal-footer">
            <button class="modal-btn" @click="showAddFeed = false">取消</button>
            <button class="modal-btn primary" @click="addFeed">保存</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Edit feed modal -->
    <Teleport to="body">
      <div v-if="showEditFeed" class="modal-backdrop" @click.self="showEditFeed = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>编辑 RSS 订阅</h3>
            <button class="modal-close" @click="showEditFeed = false">
              <Icon icon="mdi:close" />
            </button>
          </div>
          <div class="modal-body">
            <input v-model="editFeedUrl" placeholder="输入 RSS Feed URL..." type="text" class="rss-input" @keyup.enter="updateFeed" />
          </div>
          <div class="modal-footer">
            <button class="modal-btn" @click="showEditFeed = false">取消</button>
            <button class="modal-btn primary" @click="updateFeed">保存</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete confirm modal -->
    <Teleport to="body">
      <div v-if="showDeleteConfirm" class="modal-backdrop" @click.self="showDeleteConfirm = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>确认删除</h3>
          </div>
          <div class="modal-body">
            确定要删除此 RSS 订阅吗？此操作将同时删除所有已获取的文章。
          </div>
          <div class="modal-footer">
            <button class="modal-btn" @click="showDeleteConfirm = false">取消</button>
            <button class="modal-btn danger" @click="confirmDeleteFeed">删除</button>
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
import axios from '@/api'

const appStore = useAppStore()
const feeds = ref<any[]>([])
const isExpanded = ref(false)
const showAddFeed = ref(false)
const showEditFeed = ref(false)
const showDeleteConfirm = ref(false)
const deletingFeedId = ref<string>('')
const newFeedUrl = ref('')
const editFeedUrl = ref('')
let timer: ReturnType<typeof setInterval> | null = null

async function fetchFeeds() {
  try {
    const token = localStorage.getItem('sundash-token')
    const res = await axios.get('/api/rss', {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (res.status === 200) {
      feeds.value = res.data.map((feed: any) => ({
        ...feed,
        isExpanded: false, // track expansion state locally
        items: feed.items || [] // ensure items exist
      }))
    }
  } catch (err) {
    console.error('Failed to fetch RSS feeds:', err)
  }
}

async function addFeed() {
  const url = newFeedUrl.value.trim()
  if (!url) return
  try {
    const token = localStorage.getItem('sundash-token')
    await axios.post('/api/rss', { url }, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    newFeedUrl.value = ''
    showAddFeed.value = false
    await fetchFeeds()
  } catch (err) {
    console.error('Failed to add RSS feed:', err)
  }
}

async function updateFeed() {
  if (!editFeedId || !editFeedUrl.value.trim()) return
  try {
    const token = localStorage.getItem('sundash-token')
    await axios.put(`/api/rss/${editFeedId}`, { url: editFeedUrl.value.trim() }, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    editFeedId = ''
    editFeedUrl.value = ''
    showEditFeed.value = false
    await fetchFeeds()
  } catch (err) {
    console.error('Failed to update RSS feed:', err)
  }
}

async function deleteFeed(feedId: string) {
  editFeedId = feedId
  showDeleteConfirm.value = true
}

async function confirmDeleteFeed() {
  if (!editFeedId) return
  try {
    const token = localStorage.getItem('sundash-token')
    await axios.delete(`/api/rss/${editFeedId}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    editFeedId = ''
    showDeleteConfirm.value = false
    await fetchFeeds()
  } catch (err) {
    console.error('Failed to delete RSS feed:', err)
  }
}

async function refreshFeed(feed: any) {
  try {
    const token = localStorage.getItem('sundash-token')
    await axios.put(`/api/rss/${feed.id}`, { url: feed.url }, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    await fetchFeeds()
  } catch (err) {
    console.error('Failed to refresh RSS feed:', err)
  }
}

async function refreshAll() {
  try {
    const token = localStorage.getItem('sundash-token')
    // We'll refresh each feed by calling update (which triggers a fetch)
    for (const feed of feeds.value) {
      await axios.put(`/api/rss/${feed.id}`, { url: feed.url }, {
        headers: token ? { Authorization: `Bearer ${token}` } : {}
      })
    }
    await fetchFeeds()
  } catch (err) {
    console.error('Failed to refresh all RSS feeds:', err)
  }
}

function toggleFeed(feed: any) {
  feed.isExpanded = !feed.isExpanded
}

function editFeed(feedId: string) {
  const feed = feeds.value.find(f => f.id === feedId)
  if (feed) {
    editFeedId = feedId
    editFeedUrl.value = feed.url
    showEditFeed.value = true
  }
}

function formatDate(dateString: string | undefined): string {
  if (!dateString) return '--:--'
  try {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
  } catch {
    return dateString
  }
}

onMounted(() => {
  fetchFeeds()
  timer = setInterval(fetchFeeds, 5 * 60 * 1000) // 5 minutes
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.rss-widget {
  background: var(--sd-bg-elevated);
  border: 1px solid var(--sd-border);
  border-radius: 12px;
  overflow: hidden;
  margin: 0 16px 12px;
  backdrop-filter: blur(10px);
}

.rss-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 14px;
  cursor: pointer;
  user-select: none;
}

.rss-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--sd-text-secondary);
}

.rss-item span {
  font-variant-numeric: tabular-nums;
}

.rss-item.add {
  margin-left: auto;
  cursor: pointer;
  color: var(--sd-color-success);
}

.rss-item.add:hover {
  color: var(--sd-color-success-dark);
}

.rss-expand {
  margin-left: auto;
  color: var(--sd-text-tertiary);
}

/* Expanded detail */
.rss-detail {
  padding: 12px 14px 14px;
  border-top: 1px solid var(--sd-border);
}

.rss-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.rss-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--sd-text-primary);
}

.rss-actions {
  display: flex;
  gap: 8px;
}

.rss-action-btn {
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

.rss-action-btn:hover {
  background: var(--sd-bg-base);
  color: var(--sd-text-primary);
}

.rss-action-btn Icon {
  font-size: 14px;
}

.rss-list {
  max-height: 400px;
  overflow-y: auto;
}

.rss-item {
  background: var(--sd-bg-subtle);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 8px;
  border: 1px solid var(--sd-border);
}

.rss-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.rss-item-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--sd-text-primary);
}

.rss-item-title a {
  color: var(--sd-text-primary);
  text-decoration: none;
}

.rss-item-title a:hover {
  text-decoration: underline;
}

.rss-item-meta {
  font-size: 12px;
  color: var(--sd-text-tertiary);
  display: flex;
  gap: 8px;
  align-items: center;
}

.rss-item-date {
  font-variant-numeric: tabular-nums;
}

.rss-item-author {
  font-style: italic;
}

.rss-item-description {
  font-size: 13px;
  color: var(--sd-text-primary);
  line-height: 1.5;
  margin-top: 8px;
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

.rss-input {
  width: 100%;
  height: 36px;
  border: 1px solid var(--sd-border);
  border-radius: 4px;
  background: transparent;
  font-size: 14px;
  color: var(--sd-text-primary);
  padding: 0 8px;
}

.rss-input:focus {
  outline: none;
  border-color: var(--sd-primary);
  box-shadow: 0 0 0 2px rgba(0,122,255,0.2);
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

.modal-btn.danger {
  background: var(--sd-color-danger);
  color: white;
}

.modal-btn.danger:hover {
  background: var(--sd-color-danger-dark);
}

/* Scrollbar styling */
.rss-list::-webkit-scrollbar {
  width: 6px;
}

.rss-list::-webkit-scrollbar-track {
  background: var(--sd-bg-base);
}

.rss-list::-webkit-scrollbar-thumb {
  background: var(--sd-border);
  border-radius: 3px;
}

.rss-list::-webkit-scrollbar-thumb:hover {
  background: var(--sd-text-tertiary);
}
</style>