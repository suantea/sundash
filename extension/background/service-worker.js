// SunDash Monitor - Background Service Worker

const MONITOR_ENABLED_KEY = 'sundash_monitor_enabled'
const MONITOR_LIST_KEY = 'sundash_monitor_list'
const CHECK_INTERVAL = 30000 // 30 seconds

let checkTimer = null
let monitorEnabled = false
let monitorList = []

// Initialize
chrome.runtime.onInstalled.addListener(async () => {
  console.log('SunDash Monitor installed')
  await loadSettings()
  if (monitorEnabled) {
    startMonitoring()
  }
})

chrome.runtime.onStartup.addListener(async () => {
  await loadSettings()
  if (monitorEnabled) {
    startMonitoring()
  }
})

// Load settings from storage
async function loadSettings() {
  const data = await chrome.storage.local.get([MONITOR_ENABLED_KEY, MONITOR_LIST_KEY])
  monitorEnabled = data[MONITOR_ENABLED_KEY] || false
  monitorList = data[MONITOR_LIST_KEY] || []
}

// Message handling
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  switch (msg.type) {
    case 'START_MONITOR':
      monitorEnabled = true
      startMonitoring()
      break

    case 'STOP_MONITOR':
      monitorEnabled = false
      stopMonitoring()
      break

    case 'ADD_MONITOR':
      addMonitorItem(msg.data)
      break

    case 'REMOVE_MONITOR':
      removeMonitorItem(msg.data)
      break

    case 'CLEAR_MONITORS':
      monitorList = []
      saveMonitorList()
      break

    case 'PAGE_CONTENT_HASH':
      // Received from content script
      handleContentHash(sender.tab?.url, msg.data)
      break
  }
})

// Start monitoring
function startMonitoring() {
  if (checkTimer) return
  checkTimer = setInterval(checkAllPages, CHECK_INTERVAL)
  console.log('Monitoring started')
}

// Stop monitoring
function stopMonitoring() {
  if (checkTimer) {
    clearInterval(checkTimer)
    checkTimer = null
  }
  console.log('Monitoring stopped')
}

// Add monitor item
async function addMonitorItem(item) {
  const exists = monitorList.some(m => m.url === item.url)
  if (exists) return

  monitorList.push({
    ...item,
    lastCheck: null,
    lastHash: null,
    changed: false,
  })

  await saveMonitorList()

  // Try to get initial content hash
  requestContentHash(item.url)
}

// Remove monitor item
async function removeMonitorItem(item) {
  monitorList = monitorList.filter(m => m.url !== item.url)
  await saveMonitorList()
}

// Save monitor list
async function saveMonitorList() {
  await chrome.storage.local.set({ [MONITOR_LIST_KEY]: monitorList })
}

// Check all monitored pages
async function checkAllPages() {
  for (const item of monitorList) {
    try {
      await checkPage(item)
    } catch (err) {
      console.error(`Error checking ${item.url}:`, err)
    }
  }
}

// Check a single page
async function checkPage(item) {
  // Try to get content hash from active tabs
  const tabs = await chrome.tabs.query({ url: item.url })

  if (tabs.length > 0) {
    // Page is open in a tab, ask content script for hash
    try {
      const response = await chrome.tabs.sendMessage(tabs[0].id, { type: 'GET_CONTENT_HASH' })
      if (response && response.hash) {
        handleContentHash(item.url, { hash: response.hash })
      }
    } catch (err) {
      // Content script might not be loaded
      console.log('Content script not available for:', item.url)
    }
  } else {
    // Page is not open, fetch and hash (with timeout + size cap)
    try {
      const text = await fetchLimitedText(item.url, 2 * 1024 * 1024, 8000)
      const hash = simpleHash(text)
      handleContentHash(item.url, { hash })
    } catch (err) {
      console.log('Could not fetch:', item.url)
    }
  }
}

// Fetches a URL with a timeout and returns at most maxBytes of text.
async function fetchLimitedText(url, maxBytes, timeoutMs) {
  const controller = new AbortController()
  const timer = setTimeout(() => controller.abort(), timeoutMs)
  try {
    const response = await fetch(url, { signal: controller.signal })
    if (!response.ok || !response.body) {
      throw new Error(`HTTP ${response.status}`)
    }
    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let received = 0
    let text = ''
    while (received < maxBytes) {
      const { done, value } = await reader.read()
      if (done) break
      text += decoder.decode(value, { stream: true })
      received += value.length
    }
    return text
  } finally {
    clearTimeout(timer)
  }
}

// Handle content hash
function handleContentHash(url, data) {
  const item = monitorList.find(m => m.url === url)
  if (!item) return

  const { hash } = data
  const changed = item.lastHash !== null && item.lastHash !== hash

  item.lastHash = hash
  item.lastCheck = Date.now()
  item.changed = changed

  saveMonitorList()

  // Notify popup if open
  chrome.runtime.sendMessage({
    type: 'MONITOR_UPDATE',
    data: { url, changed }
  }).catch(() => {})

  // Show notification if changed
  if (changed) {
    chrome.notifications.create(`change-${Date.now()}`, {
      type: 'basic',
      iconUrl: 'icons/icon128.png',
      title: 'SunDash Monitor',
      message: `页面内容已更新: ${item.title || url}`,
      priority: 1,
    })
  }
}

// Request content hash from content script
function requestContentHash(url) {
  chrome.tabs.query({ url }, (tabs) => {
    if (tabs.length > 0) {
      chrome.tabs.sendMessage(tabs[0].id, { type: 'GET_CONTENT_HASH' })
        .catch(() => {})
    }
  })
}

// Simple hash function
function simpleHash(str) {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash
  }
  return hash.toString(36)
}

// Periodic check — only when monitoring was previously enabled (e.g. after a
// service worker restart). The onInstalled/onStartup listeners already handle
// explicit start, so this must NOT unconditionally start the timer.
loadSettings().then(() => {
  if (monitorEnabled) {
    startMonitoring()
  }
})
