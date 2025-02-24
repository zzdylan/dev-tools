<template>
  <div class="home">
    <div ref="menuListRef" class="tools-grid">
      <router-link
        v-for="item in menuItems"
        :key="item.id"
        :to="item.path"
        class="tool-card"
      >
        <div class="drag-handle">⋮⋮</div>
        <div class="tool-icon">{{ item.icon }}</div>
        <h3>{{ item.title }}</h3>
        <p>{{ item.description }}</p>
        <el-switch
          class="visibility-switch"
          v-model="item.visible"
          @change="(val: boolean) => handleVisibilityChange(item.id, val)"
        />
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useToolsStore } from '../stores/tools'
import { storeToRefs } from 'pinia'
import Sortable from 'sortablejs'
import { ref, onMounted } from 'vue'

const store = useToolsStore()
const { menuConfig } = storeToRefs(store)
const menuListRef = ref<HTMLElement | null>(null)

const menuItems = computed({
  get: () => [...menuConfig.value.items].sort((a, b) => a.order - b.order),
  set: (val) => store.updateMenuOrder(val),
})

const handleVisibilityChange = (id: string, visible: boolean) => {
  store.updateMenuItemVisibility(id, visible)
}

const handleDragEnd = (evt: Sortable.SortableEvent) => {
  const items = [...menuItems.value]
  const [movedItem] = items.splice(evt.oldIndex!, 1)
  items.splice(evt.newIndex!, 0, movedItem)
  store.updateMenuOrder(items)
}

onMounted(() => {
  if (menuListRef.value) {
    new Sortable(menuListRef.value, {
      handle: '.drag-handle',
      onEnd: handleDragEnd,
      animation: 150,
      ghostClass: 'sortable-ghost',
    })
  }
})
</script>

<style scoped>
.home {
  padding: 20px;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.tool-card {
  position: relative;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s;
}

.tool-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.tool-icon {
  font-size: 24px;
  margin-bottom: 12px;
}

.tool-card h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #111827;
}

.tool-card p {
  margin: 0;
  font-size: 14px;
  color: #6b7280;
}

.visibility-switch {
  position: absolute;
  top: 10px;
  right: 10px;
}

.drag-handle {
  position: absolute;
  top: 10px;
  left: 10px;
  cursor: move;
  color: #9ca3af;
  font-size: 16px;
}

.sortable-ghost {
  opacity: 0.5;
  background: #e5e7eb;
}
</style>
