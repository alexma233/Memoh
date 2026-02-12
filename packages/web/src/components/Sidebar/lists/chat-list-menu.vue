<template>
  <section>
    <SidebarGroup>
      <SidebarGroupLabel :class="collapsedHiddenClass">
        <span>{{ $t('sidebar.bots') }}</span>
      </SidebarGroupLabel>
      <SidebarGroupContent>
        <SidebarMenu v-if="bots.length > 0">
          <SidebarMenuItem
            v-for="bot in bots"
            :key="bot.id"
          >
            <SidebarMenuButton
              :is-active="botId === bot.id"
              :tooltip="botLabel(bot)"
              :disabled="loadingChats || initializing || isBotPending(bot)"
              @click="onSelectBot(bot.id)"
            >
              <Avatar class="size-5 shrink-0">
                <AvatarImage
                  v-if="bot.avatar_url"
                  :src="bot.avatar_url"
                  :alt="botLabel(bot)"
                />
                <AvatarFallback class="text-[10px]">
                  {{ botLabel(bot).slice(0, 2).toUpperCase() }}
                </AvatarFallback>
              </Avatar>
              <span :class="['truncate', collapsedHiddenClass]">{{ botLabel(bot) }}</span>
              <FontAwesomeIcon
                v-if="isBotPending(bot)"
                :icon="['fas', 'spinner']"
                class="ml-auto size-3 animate-spin text-muted-foreground"
              />
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
        <div
          v-else
          :class="['px-2 py-2 text-xs text-muted-foreground', collapsedHiddenClass]"
        >
          {{ $t('bots.emptyTitle') }}
        </div>
      </SidebarGroupContent>
    </SidebarGroup>
  </section>
</template>

<script setup lang="ts">
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@memoh/ui'
import { computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useChatList } from '@/store/chat-list'
import { storeToRefs } from 'pinia'
import type { Bot } from '@/composables/api/useChat'
import type { SidebarListProps } from './types'

const props = withDefaults(defineProps<SidebarListProps>(), {
  collapsible: false,
})

const router = useRouter()
const route = useRoute()
const chatStore = useChatList()
const {
  botId,
  bots,
  loadingChats,
  initializing,
} = storeToRefs(chatStore)

const collapsedHiddenClass = computed(() => (
  props.collapsible ? 'group-data-[state=collapsed]:hidden' : ''
))

onMounted(() => {
  void chatStore.initialize().catch(() => undefined)
})

function botLabel(bot: Bot): string {
  return bot.display_name?.trim() || bot.id
}

function isBotPending(bot: Bot): boolean {
  return bot.status === 'creating' || bot.status === 'deleting'
}

async function onSelectBot(targetBotID: string) {
  if (!targetBotID) {
    return
  }
  try {
    if (botId.value !== targetBotID) {
      await chatStore.selectBot(targetBotID)
    }
    if (route.name !== 'chat') {
      await router.push({ name: 'chat' })
    }
  } catch {
    return
  }
}
</script>

