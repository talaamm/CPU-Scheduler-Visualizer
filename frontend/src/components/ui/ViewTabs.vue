<template>
  <div class="view-tabs">
    <button
      v-for="tab in tabs"
      :key="tab.id"
      class="view-tab"
      :class="{ active: modelValue === tab.id }"
      @click="$emit('update:modelValue', tab.id)"
    >
      {{ tab.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
export interface TabDef {
  id: string
  label: string
}

defineProps<{
  modelValue: string
  tabs: TabDef[]
}>()

defineEmits<{
  'update:modelValue': [id: string]
}>()
</script>

<style scoped>
.view-tabs {
  display: flex;
  border-bottom: 1px solid #2E3348;
  padding: 0 20px;
  background: #1A1D27;
  flex-shrink: 0;
}

.view-tab {
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  color: #94A3B8;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  padding: 12px 16px;
  transition: all 0.15s;
  margin-bottom: -1px;
}
.view-tab:hover { color: #F1F5F9; }
.view-tab.active { color: #4F8EF7; border-bottom-color: #4F8EF7; }
</style>
