import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api'
import type { PanelGroup, Card, PanelData } from '../types'

export const usePanelStore = defineStore('panel', () => {
  const groups = ref<PanelGroup[]>([])
  const loading = ref(false)

  async function fetchPanel() {
    loading.value = true
    try {
      const res = await api.get<PanelData>('panels')
      groups.value = res.data.groups || []
    } catch (e) {
      console.error('Failed to load panel:', e)
    } finally {
      loading.value = false
    }
  }

  // Direct state fill used by the home page bootstrap request.
  function setGroups(groupsData: PanelGroup[]) {
    groups.value = groupsData || []
    loading.value = false
  }

  async function createGroup(name: string) {
    const res = await api.post<PanelGroup>('panels/groups', { name })
    groups.value.push({ ...res.data, cards: [] })
    return res.data
  }

  async function updateGroup(id: string, data: { name?: string; sort_order?: number }) {
    await api.put(`/api/panels/groups/${id}`, data)
    const group = groups.value.find(g => g.id === id)
    if (group) {
      if (data.name !== undefined) group.name = data.name
      if (data.sort_order !== undefined) group.sort_order = data.sort_order
    }
  }

  async function deleteGroup(id: string) {
    await api.delete(`/api/panels/groups/${id}`)
    groups.value = groups.value.filter(g => g.id !== id)
  }

  async function createCard(data: {
    group_id: string; title: string; url: string;
    url_internal?: string; icon?: string; icon_color?: string;
    bg_color?: string; description?: string; open_type?: string;
  }) {
    const res = await api.post<Card>('panels/cards', data)
    const group = groups.value.find(g => g.id === data.group_id)
    if (group) {
      if (!group.cards) group.cards = []
      group.cards.push(res.data)
    }
    return res.data
  }

  async function updateCard(id: string, data: Partial<Card>) {
    await api.put(`/api/panels/cards/${id}`, data)
    for (const group of groups.value) {
      if (group.cards) {
        const card = group.cards.find(c => c.id === id)
        if (card) {
          Object.assign(card, data)
          break
        }
      }
    }
  }

  async function deleteCard(id: string) {
    await api.delete(`/api/panels/cards/${id}`)
    for (const group of groups.value) {
      if (group.cards) {
        group.cards = group.cards.filter(c => c.id !== id)
      }
    }
  }

  async function reorder(groupOrders: { id: string; sort_order: number }[], cardOrders: { id: string; group_id: string; sort_order: number }[]) {
    await api.put('panels/reorder', { group_orders: groupOrders, card_orders: cardOrders })
    // Update local state
    for (const go of groupOrders) {
      const group = groups.value.find(g => g.id === go.id)
      if (group) group.sort_order = go.sort_order
    }
    groups.value.sort((a, b) => a.sort_order - b.sort_order)
    for (const co of cardOrders) {
      for (const group of groups.value) {
        if (group.id === co.group_id && group.cards) {
          const card = group.cards.find(c => c.id === co.id)
          if (card) {
            card.sort_order = co.sort_order
            card.group_id = co.group_id
          }
        }
      }
    }
  }

  function getCardUrl(card: Card, networkMode: 'internal' | 'external'): string {
    if (networkMode === 'internal' && card.url_internal) {
      return card.url_internal
    }
    return card.url
  }

  // Batch update card colors
  async function batchUpdateCardColors(data: { cardIds: string[]; icon_color?: string; bg_color?: string }) {
    const { cardIds, icon_color, bg_color } = data
    // Update each card via API
    for (const cardId of cardIds) {
      const updateData: Partial<Card> = {}
      if (icon_color !== undefined) updateData.icon_color = icon_color
      if (bg_color !== undefined) updateData.bg_color = bg_color
      await api.put(`/api/panels/cards/${cardId}`, updateData)
      
      // Update local state
      for (const group of groups.value) {
        if (group.cards) {
          const card = group.cards.find(c => c.id === cardId)
          if (card) {
            Object.assign(card, updateData)
            break
          }
        }
      }
    }
  }

  // Settings
  const settings = ref<Record<string, string>>({})

  async function fetchSettings() {
    try {
      const res = await api.get<Record<string, string>>('settings')
      settings.value = res.data || {}
    } catch (e) {
      console.error('Failed to load settings:', e)
    }
  }

  async function updateSetting(key: string, value: string) {
    await api.put('settings', { key, value })
    settings.value[key] = value
  }

  async function batchUpdateSettings(data: Record<string, string>) {
    await api.put('settings/batch', { settings: data })
    Object.assign(settings.value, data)
  }

  function getSetting(key: string): string {
    return settings.value[key] || ''
  }

  return {
    groups, loading,
    fetchPanel, setGroups, createGroup, updateGroup, deleteGroup,
    createCard, updateCard, deleteCard,
    reorder, getCardUrl, batchUpdateCardColors,
    settings, fetchSettings, updateSetting, batchUpdateSettings, getSetting,
  }
})
