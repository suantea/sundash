import { ref, computed } from 'vue'
import { usePanelStore } from '../stores/panel'
import type { Card, PanelGroup } from '../types'

/**
 * Collapse/hide state for groups and cards, persisted to localStorage.
 * Extracted from Home.vue.
 */
export function useGroupVisibility() {
  const panelStore = usePanelStore()

  const collapsedGroupIds = ref<Set<string>>(new Set(JSON.parse(localStorage.getItem('sundash-collapsed-groups') || '[]')))
  const hiddenGroupIds = ref<Set<string>>(new Set(JSON.parse(localStorage.getItem('sundash-hidden-groups') || '[]')))
  const hiddenCardIds = ref<Set<string>>(new Set(JSON.parse(localStorage.getItem('sundash-hidden-cards') || '[]')))
  const showHiddenGroups = ref(false)
  const showHiddenCardsGroup = ref<Set<string>>(new Set())

  function persist(key: string, set: Set<string>) {
    localStorage.setItem(key, JSON.stringify([...set]))
  }

  function isCollapsed(id: string) { return collapsedGroupIds.value.has(id) }
  function toggleCollapse(id: string) {
    const s = new Set(collapsedGroupIds.value)
    s.has(id) ? s.delete(id) : s.add(id)
    collapsedGroupIds.value = s
    persist('sundash-collapsed-groups', s)
  }

  function isGroupHidden(id: string) { return hiddenGroupIds.value.has(id) }
  function toggleHideGroup(id: string) {
    const s = new Set(hiddenGroupIds.value)
    s.has(id) ? s.delete(id) : s.add(id)
    hiddenGroupIds.value = s
    persist('sundash-hidden-groups', s)
  }

  function isCardHidden(id: string) { return hiddenCardIds.value.has(id) }
  function toggleHideCard(id: string) {
    const s = new Set(hiddenCardIds.value)
    s.has(id) ? s.delete(id) : s.add(id)
    hiddenCardIds.value = s
    persist('sundash-hidden-cards', s)
  }

  const visibleGroups = computed(() => panelStore.groups.filter(g => !hiddenGroupIds.value.has(g.id)))

  function getGroupName(id: string) {
    return panelStore.groups.find(g => g.id === id)?.name || id.slice(0, 8)
  }

  function hiddenCardsCount(group: PanelGroup) {
    return (group.cards || []).filter((c: Card) => hiddenCardIds.value.has(c.id)).length
  }

  function toggleShowHiddenCards(groupId: string) {
    const s = new Set(showHiddenCardsGroup.value)
    s.has(groupId) ? s.delete(groupId) : s.add(groupId)
    showHiddenCardsGroup.value = s
  }

  function groupCards(group: PanelGroup): Card[] {
    // Copy before sorting so the store state is never mutated in place.
    const cards = [...(group.cards || [])].sort((a, b) => a.sort_order - b.sort_order)
    if (showHiddenCardsGroup.value.has(group.id)) return cards
    return cards.filter((c: Card) => !hiddenCardIds.value.has(c.id))
  }

  return {
    collapsedGroupIds, hiddenGroupIds, hiddenCardIds, showHiddenGroups, showHiddenCardsGroup,
    isCollapsed, toggleCollapse, isGroupHidden, toggleHideGroup, isCardHidden, toggleHideCard,
    visibleGroups, getGroupName, hiddenCardsCount, toggleShowHiddenCards, groupCards,
  }
}
