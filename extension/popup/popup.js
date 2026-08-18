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
    return `<img src="https://www.google.com/s2/favicons?domain=${hostname}&sz=32" width="20" height="20" onerror="this.style.display='none'">`
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

// Listen for updates from background
chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'MONITOR_UPDATE') {
    const item = monitorList.find(m => m.url === msg.data.url)
    if (item) {
      item.changed = msg.data.changed
      item.lastCheck = Date.now()
      chrome.storage.local.set({ [MONITOR_LIST_KEY]: monitorList })
      updateUI()
    }
  }
})

// Initialize
init()
