<script setup lang="ts">
import { ref } from 'vue'
import { onClickOutside } from '@vueuse/core'

interface Props {
  align?: 'left' | 'right'
  width?: string
}

const props = withDefaults(defineProps<Props>(), {
  align: 'left',
  width: '12rem',
})

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement>()

onClickOutside(dropdownRef, () => {
  isOpen.value = false
})

function toggle() {
  isOpen.value = !isOpen.value
}

function close() {
  isOpen.value = false
}

defineExpose({ close, toggle })
</script>

<template>
  <div ref="dropdownRef" class="relative inline-block">
    <div @click="toggle" class="cursor-pointer">
      <slot name="trigger" />
    </div>
    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="opacity-0 scale-95 translate-y-1"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute z-50 mt-1 rounded-lg border border-custom-border-200 bg-custom-background-100 py-1 shadow-custom-lg"
        :class="props.align === 'right' ? 'right-0' : 'left-0'"
        :style="{ minWidth: props.width }"
      >
        <slot :close="close" />
      </div>
    </Transition>
  </div>
</template>
