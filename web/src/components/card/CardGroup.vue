<template>
  <div class="card-group">
    <div class="group-header">
      <div class="group-title">
        <h3>{{ group.name }}</h3>
        <span class="group-count">{{ group.cards?.length || 0 }}</span>
      </div>
      <div class="group-actions">
        <button class="group-action-btn" @click="$emit('addCard', group.id)" :title="t('home.addBookmark')">
          <Icon icon="mdi:plus" :size="16" />
        </button>
        <n-dropdown :options="menuOptions" @select="handleMenu" trigger="click" placement="bottom-end">
          <button class="group-action-btn">
            <Icon icon="mdi:dots-horizontal" :size="16" />
          </button>
        </n-dropdown>
      </div>
    </div>

    <draggable
      :list="localCards"
      :group="{ name: 'cards', pull: false, put: true }"
      item-key="id"
      class="cards-grid"
      :style="gridStyle"
      ghost-class="card-ghost"
      drag-class="card-drag"
      @end="onDragEnd"
    >
      <template #item="{ element }">
        <CardItem
          :card="element"
          @click="openCard(element)"
          @edit="$emit('editCard', element)"
          @delete="$emit('deleteCard', element.id)"
        />
      </template>
    </draggable>

    <div class="cards-grid cards-grid-footer" :style="gridStyle">
      <button class="card-add" @click="$emit('addCard', group.id)">
        <Icon icon="mdi:plus" :size="20" color="var(--sd-text-tertiary)" />
        <span>{{ t('home.addBookmark') }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import draggable from 'vuedraggable'
import type { PanelGroup, Card } from '../../types'
import { useAppStore } from '../../stores/app'
import { usePanelStore } from '../../stores/panel'
import CardItem from './CardItem.vue'

const { t } = useI18n()
const props = defineProps<{ group: PanelGroup }>()
const emit = defineEmits<{
  editGroup: [group: { id: string; name: string }]
  deleteGroup: [id: string]
  addCard: [groupId: string]
  editCard: [card: Card]
  deleteCard: [cardId: string]
  reorder: [groupId: string, cardIds: string[]]
}>()

const appStore = useAppStore()
const panelStore = usePanelStore()

const localCards = computed({
  get: () => [...(props.group.cards || [])].sort((a, b) => a.sort_order - b.sort_order),
  set: () => {},
})

// Grid style based on cardsPerRow setting
const gridStyle = computed(() => {
  const count = parseInt(appStore.cardsPerRow) || 5
  return { gridTemplateColumns: `repeat(${count}, 1fr)` }
})

const menuOptions = [
  { label: t('home.editGroup'), key: 'edit', icon: () => h(Icon, { icon: 'mdi:pencil' }) },
  { label: t('home.deleteGroup'), key: 'delete', icon: () => h(Icon, { icon: 'mdi:delete' }) },
]

function handleMenu(key: string) {
  if (key === 'edit') {
    emit('editGroup', { id: props.group.id, name: props.group.name })
  } else if (key === 'delete') {
    emit('deleteGroup', props.group.id)
  }
}

function openCard(card: Card) {
  const url = panelStore.getCardUrl(card, appStore.networkMode)
  if (card.open_type === 'iframe') {
    window.open(url, '_blank', 'width=1024,height=768')
  } else {
    window.open(url, '_blank')
  }
}

function onDragEnd() {
  const cardIds = localCards.value.map(c => c.id)
  emit('reorder', props.group.id, cardIds)
}
</script>

<style scoped>
.card-group {
  background: var(--sd-bg-surface);
  border: 1px solid var(--sd-border);
  border-radius: var(--sd-radius-lg);
  overflow: hidden;
  transition: box-shadow var(--sd-duration-normal) var(--sd-ease);
}

.card-group:hover {
  box-shadow: var(--sd-shadow-sm);
}

.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sd-space-4) var(--sd-space-5);
}

.group-title {
  display: flex;
  align-items: center;
  gap: var(--sd-space-3);
}

.group-title h3 {
  font-size: var(--sd-text-lg);
  font-weight: var(--sd-weight-semibold);
  color: var(--sd-text-primary);
  margin: 0;
}

.group-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 22px;
  height: 22px;
  padding: 0 6px;
  background: var(--sd-primary-light);
  color: var(--sd-primary);
  border-radius: var(--sd-radius-full);
  font-size: var(--sd-text-xs);
  font-weight: var(--sd-weight-semibold);
}

.group-actions {
  display: flex;
  align-items: center;
  gap: var(--sd-space-1);
}

.group-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: none;
  background: transparent;
  color: var(--sd-text-secondary);
  border-radius: var(--sd-radius-sm);
  cursor: pointer;
  transition: all var(--sd-duration-fast) var(--sd-ease);
}

.group-action-btn:hover {
  background: var(--sd-primary-light);
  color: var(--sd-primary);
}

.cards-grid {
  display: grid;
  gap: var(--sd-space-3);
  padding: 0 var(--sd-space-5) var(--sd-space-5);
  min-height: 8px;
}

.cards-grid-footer {
  padding-top: 0;
}

.card-add {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--sd-space-2);
  min-height: 130px;
  border: 1.5px dashed var(--sd-border);
  border-radius: var(--sd-radius-md);
  background: transparent;
  cursor: pointer;
  transition: all var(--sd-duration-normal) var(--sd-ease);
  color: var(--sd-text-tertiary);
  font-size: var(--sd-text-sm);
  font-family: var(--sd-font);
  font-weight: var(--sd-weight-medium);
}

.card-add:hover {
  border-color: var(--sd-primary);
  color: var(--sd-primary);
  background: var(--sd-primary-light);
}

:deep(.card-ghost) {
  opacity: 0.3;
}

:deep(.card-drag) {
  opacity: 0.9;
  transform: rotate(1deg) scale(1.02);
  box-shadow: var(--sd-shadow-xl);
}
</style>
