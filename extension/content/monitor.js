// SunDash Monitor - Content Script
// Injected into monitored pages to detect content changes

(function() {
  'use strict'

  let lastContentHash = null
  let observer = null
  let debounceTimer = null

  // Calculate content hash
  function getContentHash() {
    const content = document.body.innerText || document.body.textContent || ''
    return simpleHash(content)
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

  // Notify background of content change
  function notifyChange() {
    const hash = getContentHash()
    if (hash !== lastContentHash) {
      lastContentHash = hash
      chrome.runtime.sendMessage({
        type: 'PAGE_CONTENT_HASH',
        data: { hash, url: window.location.href }
      }).catch(() => {})
    }
  }

  // Debounced notification
  function debouncedNotify() {
    if (debounceTimer) clearTimeout(debounceTimer)
    debounceTimer = setTimeout(notifyChange, 1000)
  }

  // Start observing DOM changes
  function startObserver() {
    if (observer) return

    observer = new MutationObserver((mutations) => {
      // Filter out insignificant mutations
      const significant = mutations.some(m => {
        // Ignore mutations in our own UI elements
        if (m.target.closest && m.target.closest('#sundash-monitor-badge')) return false
        // Ignore attribute-only changes on non-content elements
        if (m.type === 'attributes' && !m.target.classList?.contains('content')) return false
        return true
      })

      if (significant) {
        debouncedNotify()
      }
    })

    observer.observe(document.body, {
      childList: true,
      subtree: true,
      characterData: true,
    })
  }

  // Stop observing
  function stopObserver() {
    if (observer) {
      observer.disconnect()
      observer = null
    }
  }

  // Listen for messages from background
  chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
    switch (msg.type) {
      case 'GET_CONTENT_HASH':
        sendResponse({ hash: getContentHash() })
        break

      case 'START_OBSERVING':
        startObserver()
        sendResponse({ success: true })
        break

      case 'STOP_OBSERVING':
        stopObserver()
        sendResponse({ success: true })
        break
    }
  })

  // Initial content hash
  lastContentHash = getContentHash()

  // Start observing if monitor is enabled
  chrome.storage.local.get('sundash_monitor_enabled', (data) => {
    if (data.sundash_monitor_enabled) {
      startObserver()
    }
  })

  // Listen for storage changes
  chrome.storage.onChanged.addListener((changes) => {
    if (changes.sundash_monitor_enabled) {
      if (changes.sundash_monitor_enabled.newValue) {
        startObserver()
      } else {
        stopObserver()
      }
    }
  })

  console.log('[SunDash Monitor] Content script loaded for:', window.location.href)
})()
