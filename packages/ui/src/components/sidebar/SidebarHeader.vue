<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '#/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <header
    data-slot="sidebar-header"
    data-sidebar="header"
    role="region"
    aria-label="Sidebar header"
    :class="cn('flex flex-col gap-2 p-2', props.class)"
  >
    <slot />
  </header>
</template>
