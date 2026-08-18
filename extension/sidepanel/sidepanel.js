// SunDash Monitor - Side Panel Script

const MONITOR_LIST_KEY = 'sundash_monitor_list'
const itemsContainer = document.getElementById('items')

async function loadItems() {
  const data = await chrome.storage.local.get(MONITOR_LIST_KEY)
  const list = data[MONITOR_LIST_KEY] || []

  if (list.length === 0) {
    itemsContainer.innerHTML = '<div class="empty">暂无监控项目<br><small>在弹出面板中添加</small></div>'
    return
  }

  itemsContainer.innerHTML = list.map(item => `
    <div class="item" data-url="${escapeAttr(item.url)}">
      <div class="item-icon">${getFavicon(item.url)}</div>
      <div class="item-info">
        <div class="item-title">${escapeHtml(item.title)}</div>
        <div class="item-url">${escapeHtml(truncateUrl(item.url))}</div>
      </div>
      <span class="item-badge ${item.changed ? 'changed' : 'ok'}">
        ${item.changed ? '已更新' : '正常'}
      </span>
    </div>
  `).join('')

  itemsContainer.querySelectorAll('.item').forEach(el => {
    el.addEventListener('click', () => {
      chrome.tabs.create({ url: el.dataset.url })
    })
  })
}

function getFavicon(url) {
  try {
    const hostname = new URL(url).hostname
    return `<img src="https://www.google.com/s2/favicons?domain=${hostname}&sz=32" width="20" height="20" onerror="this.outerHTML='?'">`
  } catch { return '?' }
}

function truncateUrl(url) {
  try {
    const u = new URL(url)
    return u.hostname + u.pathname.substring(0, 30)
  } catch { return url.substring(0, 40) }
}

function escapeHtml(str) {
  const div = document.createElement('div')
  div.textContent = str
  return div.innerHTML
}

function escapeAttr(str) {
  return str.replace(/"/g, '&quot;').replace(/'/g, '&#39;')
}

// Listen for updates
chrome.storage.onChanged.addListener((changes) => {
  if (changes[MONITOR_LIST_KEY]) {
    loadItems()
  }
})

chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === 'MONITOR_UPDATE') {
    loadItems()
  }
})

loadItems()
