// SunDash Monitor - Popup Script

const MONITOR_ENABLED_KEY = 'sundash_monitor_enabled'
const MONITOR_LIST_KEY = 'sundash_monitor_list'

// DOM elements
const statusDot = document.getElementById('monitor-status')
const statusText = document.getElementById('monitor-text')
const monitorCount = document.getElementById('monitor-count')
const pageTitle = document.getElementById('page-title')
const pageUrl = document.getElementById('page-url')
const monitorItems = document.getElementById('monitor-items')
const btnToggleMonitor = document.getElementById('btn-toggle-monitor')
const btnAddMonitor = document.getElementById('btn-add-monitor')
const btnClearAll = document.getElementById('btn-clear-all')
const btnOpenPanel = document.getElementById('btn-open-panel')

let monitorEnabled = false
let monitorList = []
let currentTab = null

// Initialize
async function init() {
  // Get current tab
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true })
  currentTab = tab

  if (tab) {
    pageTitle.textContent = tab.title || '未知页面'
    pageUrl.textContent = tab.url || ''
    // Pre-fill bookmark form
    document.getElementById('bookmark-title').value = tab.title || ''
    document.getElementById('bookmark-url').value = tab.url || ''
  }

  // Load settings
  const data = await chrome.storage.local.get([MONITOR_ENABLED_KEY, MONITOR_LIST_KEY])
  monitorEnabled = data[MONITOR_ENABLED_KEY] || false
  monitorList = data[MONITOR_LIST_KEY] || []

  updateUI()
}

function updateUI() {
  // Update status
  statusDot.classList.toggle('active', monitorEnabled)
  statusText.textContent = monitorEnabled ? '监控已开启' : '监控已关闭'
  btnToggleMonitor.classList.toggle('active', monitorEnabled)
  monitorCount.textContent = monitorList.length

  // Update monitor list
  if (monitorList.length === 0) {
    monitorItems.innerHTML = '<div class="empty-state">暂无监控项目</div>'
    return
  }

  monitorItems.innerHTML = monitorList.map((item, index) => `
    <div class="monitor-item" data-index="${index}">
      <div class="item-icon">${getFavicon(item.url)}</div>
      <div class="item-info">
        <div class="item-title">${escapeHtml(item.title)}</div>
        <div class="item-url">${escapeHtml(truncateUrl(item.url))}</div>
      </div>
      <span class="item-status ${item.changed ? 'changed' : 'unchanged'}">
        ${item.changed ? '已更新' : '正常'}
      </span>
      <button class="item-remove" data-index="${index}" title="移除">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
        </svg>
      </button>
    </div>
  `).join('')

  // Add event listeners
  monitorItems.querySelectorAll('.item-remove').forEach(btn => {
    btn.addEventListener('click', (e) => {
      e.stopPropagation()
      const index = parseInt(btn.dataset.index)
      removeFromMonitor(index)
    })
  })

  monitorItems.querySelectorAll('.monitor-item').forEach(item => {
    item.addEventListener('click', () => {
      const index = parseInt(item.dataset.index)
      const monitorItem = monitorList[index]
      if (monitorItem) {
        chrome.tabs.create({ url: monitorItem.url })
      }
    })
  })
}

// Toggle monitor
btnToggleMonitor.addEventListener('click', async () => {
  monitorEnabled = !monitorEnabled
  await chrome.storage.local.set({ [MONITOR_ENABLED_KEY]: monitorEnabled })

  // Notify background
  chrome.runtime.sendMessage({
    type: monitorEnabled ? 'START_MONITOR' : 'STOP_MONITOR'
  })

  updateUI()
})

// Add current page to monitor
btnAddMonitor.addEventListener('click', async () => {
  if (!currentTab || !currentTab.url) return

  // Check if already monitored
  const exists = monitorList.some(item => item.url === currentTab.url)
  if (exists) {
    return
  }

  const newItem = {
    url: currentTab.url,
    title: currentTab.title || currentTab.url,
    addedAt: Date.now(),
    lastCheck: null,
    lastHash: null,
    changed: false,
  }

  monitorList.push(newItem)
  await chrome.storage.local.set({ [MONITOR_LIST_KEY]: monitorList })

  // Notify background to start monitoring this URL
  chrome.runtime.sendMessage({
    type: 'ADD_MONITOR',
    data: newItem
  })

  updateUI()
})

// Clear all
btnClearAll.addEventListener('click', async () => {
  if (monitorList.length === 0) return
  monitorList = []
  await chrome.storage.local.set({ [MONITOR_LIST_KEY]: monitorList })
  chrome.runtime.sendMessage({ type: 'CLEAR_MONITORS' })
  updateUI()
})

// Remove from monitor
async function removeFromMonitor(index) {
  const removed = monitorList.splice(index, 1)
  await chrome.storage.local.set({ [MONITOR_LIST_KEY]: monitorList })
  if (removed[0]) {
    chrome.runtime.sendMessage({
      type: 'REMOVE_MONITOR',
      data: removed[0]
    })
  }
  updateUI()
}

// Open panel
btnOpenPanel.addEventListener('click', () => {
  chrome.tabs.create({ url: 'http://localhost:3000' })
})

// Helpers
function getFavicon(url) {
  try {
    const hostname = new URL(url).hostname
    return `<img src="https://favicon.im/${hostname}" width="20" height="20" onerror="this.style.display='none'">`
  } catch {
    return '?'
  }
}

function truncateUrl(url) {
  try {
    const u = new URL(url)
    return u.hostname + u.pathname.substring(0, 30)
  } catch {
    return url.substring(0, 40)
  }
}

function escapeHtml(str) {
  const div = document.createElement('div')
  div.textContent = str
  return div.innerHTML
}

// ==================== Bookmark Feature ====================
const SERVER_URL_KEY = 'sundash_server_url'
const SERVER_TOKEN_KEY = 'sundash_token'

const btnBookmark = document.getElementById('btn-bookmark')
const bookmarkSection = document.getElementById('bookmark-section')
const bookmarkGroup = document.getElementById('bookmark-group')
const bookmarkTitle = document.getElementById('bookmark-title')
const bookmarkUrl = document.getElementById('bookmark-url')
const btnSaveBookmark = document.getElementById('btn-save-bookmark')
const btnCancelBookmark = document.getElementById('btn-cancel-bookmark')
const bookmarkStatus = document.getElementById('bookmark-status')

let groupsLoaded = false

async function getServerConfig() {
  const data = await chrome.storage.local.get([SERVER_URL_KEY, SERVER_TOKEN_KEY])
  return {
    url: data[SERVER_URL_KEY] || 'http://localhost:3000',
    token: data[SERVER_TOKEN_KEY] || '',
  }
}

async function loadGroups() {
  if (groupsLoaded) return
  const config = await getServerConfig()
  if (!config.token) {
    bookmarkGroup.innerHTML = '<option value="">请先在设置中配置 Token</option>'
    return
  }
  try {
    const res = await fetch(`${config.url}/api/panels`, {
      headers: { Authorization: `Bearer ${config.token}` },
    })
    if (!res.ok) throw new Error('Failed')
    const data = await res.json()
    const groups = data.groups || []
    bookmarkGroup.innerHTML = groups.map(g =>
      `<option value="${g.id}">${escapeHtml(g.name)}</option>`
    ).join('')
    groupsLoaded = true
  } catch {
    bookmarkGroup.innerHTML = '<option value="">加载分组失败，请检查连接</option>'
  }
}

// Toggle bookmark section
btnBookmark.addEventListener('click', async () => {
  bookmarkSection.classList.toggle('hidden')
  if (!bookmarkSection.classList.contains('hidden')) {
    await loadGroups()
  }
})

btnCancelBookmark.addEventListener('click', () => {
  bookmarkSection.classList.add('hidden')
  bookmarkStatus.classList.add('hidden')
})

// Save bookmark
btnSaveBookmark.addEventListener('click', async () => {
  const groupId = bookmarkGroup.value
  const title = bookmarkTitle.value.trim()
  const url = bookmarkUrl.value.trim()

  if (!groupId || !title || !url) {
    showBookmarkStatus('请填写完整信息', 'error')
    return
  }

  const config = await getServerConfig()
  if (!config.token) {
    showBookmarkStatus('请先在设置中配置 Token', 'error')
    return
  }

  btnSaveBookmark.disabled = true
  btnSaveBookmark.textContent = '保存中...'

  try {
    // Auto-fetch icon
    let icon = ''
    try {
      const iconRes = await fetch(`${config.url}/api/favicon?url=${encodeURIComponent(url)}`, {
        headers: { Authorization: `Bearer ${config.token}` },
      })
      if (iconRes.ok) {
        const iconData = await iconRes.json()
        icon = iconData.icon_name || iconData.favicon_url || ''
        // Also auto-fill title if empty and server returned one
        if (iconData.title && !bookmarkTitle.value.trim()) {
          bookmarkTitle.value = iconData.title
        }
      }
    } catch {}

    const res = await fetch(`${config.url}/api/panels/cards`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${config.token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        group_id: groupId,
        title,
        url,
        icon,
        open_type: 'new_tab',
      }),
    })

    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.error || '保存失败')
    }

    showBookmarkStatus('✅ 收藏成功！', 'success')
    setTimeout(() => bookmarkSection.classList.add('hidden'), 1500)
  } catch (e) {
    showBookmarkStatus(`❌ ${e.message}`, 'error')
  } finally {
    btnSaveBookmark.disabled = false
    btnSaveBookmark.textContent = '保存书签'
  }
})

function showBookmarkStatus(msg, type) {
  bookmarkStatus.textContent = msg
  bookmarkStatus.className = `bookmark-status ${type}`
  bookmarkStatus.classList.remove('hidden')
}

// Initialize
init()
