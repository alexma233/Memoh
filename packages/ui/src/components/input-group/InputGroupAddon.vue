<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import type { InputGroupVariants } from '.'
import { cn } from '#/lib/utils'
import { inputGroupAddonVariants } from '.'

const props = withDefaults(defineProps<{
  align?: InputGroupVariants['align']
  ariaLabel?: string
  class?: HTMLAttributes['class']
}>(), {
  align: 'inline-start',
  ariaLabel: '',
})

function handleInputGroupAddonClick(e: MouseEvent) {
  const currentTarget = e.currentTarget as HTMLElement | null
  const target = e.target as HTMLElement | null
  if (target && currentTarget && target !== currentTarget && target.closest('button')) {
    return
  }
  if (currentTarget && currentTarget?.parentElement) {
    currentTarget.parentElement?.querySelector('input')?.focus()
  }
}
</script>

<template>
  <button
    type="button"
    data-slot="input-group-addon"
    :data-align="props.align"
    :class="cn('border-0 bg-transparent p-0 text-inherit', inputGroupAddonVariants({ align: props.align }), props.class)"
    :aria-label="props.ariaLabel || undefined"
    @click="handleInputGroupAddonClick"
  >
    <slot />
  </button>
</template>
