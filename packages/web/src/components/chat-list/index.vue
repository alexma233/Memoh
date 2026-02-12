<template>
  <div
    ref="displayContainer"
    class="flex flex-col gap-4 min-w-0 max-w-full"
  >
    <template
      v-for="chatItem in chatList"
      :key="chatItem.id"
    >
      <UserChat
        v-if="chatItem.action === 'user'"
        :user-say="chatItem"
      />
      <RobotChat
        v-if="chatItem.action === 'robot'"
        :robot-say="chatItem"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import UserChat from './user-chat/index.vue'
import RobotChat from './robot-chat/index.vue'
import { inject, nextTick, ref, watch, onMounted, onUnmounted } from 'vue'
import { useChatList } from '@/store/chat-list'
import { storeToRefs } from 'pinia'
import { useAutoScroll } from '@/composables/useAutoScroll'

const { chatList, sendMessage, loadOlderMessages } = useChatList()
const { loading, loadingChats, loadingOlder, hasMoreOlder } = storeToRefs(useChatList())

const chatSay = inject('chatSay', ref(''))

watch(chatSay, async () => {
  if (chatSay.value) {
    const text = chatSay.value
    chatSay.value = ''
    try {
      await sendMessage(text)
      scheduleScrollToBottom()
    } catch {
      // ignore
    }
  }
}, { immediate: true })

const displayContainer = ref<HTMLElement>()
const { scrollToBottom, getScrollParent } = useAutoScroll(displayContainer, loading)

const LOAD_OLDER_THRESHOLD = 200

function scrollToBottomNow() {
  const el = getScrollParent()
  if (!el) return
  el.scrollTop = el.scrollHeight - el.clientHeight
}

function scheduleScrollToBottom() {
  nextTick(() => {
    scrollToBottomNow()
    requestAnimationFrame(() => {
      scrollToBottomNow()
      requestAnimationFrame(() => scrollToBottomNow())
    })
    setTimeout(() => scrollToBottomNow(), 100)
    setTimeout(() => scrollToBottomNow(), 350)
  })
}

watch(loadingChats, (isLoading, wasLoading) => {
  if (wasLoading && !isLoading && chatList.length > 0) {
    scheduleScrollToBottom()
  }
})

watch(() => chatList.length, (len, prevLen) => {
  if (prevLen === 0 && len > 0 && !loadingChats.value) {
    scheduleScrollToBottom()
  }
})

let loadOlderScheduled = false
let scrollContainer: HTMLElement | null = null

function onScroll() {
  if (!scrollContainer || loadingOlder.value || !hasMoreOlder.value || chatList.length === 0) return
  if (scrollContainer.scrollTop > LOAD_OLDER_THRESHOLD) return
  if (loadOlderScheduled) return
  loadOlderScheduled = true
  const oldScrollHeight = scrollContainer.scrollHeight
  const oldScrollTop = scrollContainer.scrollTop
  loadOlderMessages().then((n) => {
    if (n > 0) {
      nextTick(() => {
        if (scrollContainer) {
          scrollContainer.scrollTop = oldScrollTop + (scrollContainer.scrollHeight - oldScrollHeight)
        }
      })
    }
  }).finally(() => {
    loadOlderScheduled = false
  })
}

onMounted(() => {
  nextTick(() => {
    scrollContainer = getScrollParent() ?? null
    if (scrollContainer) {
      scrollContainer.addEventListener('scroll', onScroll, { passive: true })
    }
    if (chatList.length > 0 && !loadingChats.value) {
      scheduleScrollToBottom()
    }
  })
})

onUnmounted(() => {
  if (scrollContainer) {
    scrollContainer.removeEventListener('scroll', onScroll)
    scrollContainer = null
  }
})
</script>
