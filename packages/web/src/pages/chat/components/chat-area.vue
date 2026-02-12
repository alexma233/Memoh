<template>
  <div class="flex-1 flex flex-col h-full min-w-0">
    <!-- No bot selected -->
    <div
      v-if="!currentBotId"
      class="flex-1 flex items-center justify-center text-muted-foreground"
    >
      <div class="text-center">
        <p class="text-lg">{{ $t('chat.selectBot') }}</p>
        <p class="text-sm mt-1">{{ $t('chat.selectBotHint') }}</p>
      </div>
    </div>

    <template v-else>
      <!-- Messages -->
      <div
        ref="scrollContainer"
        class="flex-1 overflow-y-auto"
      >
        <div class="max-w-3xl mx-auto px-4 py-6 space-y-6">
          <!-- Empty state -->
          <div
            v-if="messages.length === 0"
            class="flex items-center justify-center min-h-[300px]"
          >
            <p class="text-muted-foreground text-lg">
              {{ $t('chat.greeting') }}
            </p>
          </div>

          <!-- Message list -->
          <MessageItem
            v-for="msg in messages"
            :key="msg.id"
            :message="msg"
          />

        </div>
      </div>

      <!-- Input -->
      <div class="border-t p-4">
        <div class="max-w-3xl mx-auto relative">
          <Textarea
            v-model="inputText"
            class="pr-16 min-h-[60px] max-h-[200px] resize-none"
            :placeholder="$t('chat.inputPlaceholder')"
            :disabled="!currentBotId"
            @keydown.enter.exact="handleKeydown"
          />
          <div class="absolute right-2 bottom-2">
            <Button
              v-if="!streaming"
              size="sm"
              :disabled="!inputText.trim() || !currentBotId"
              @click="handleSend"
            >
              <FontAwesomeIcon :icon="['fas', 'paper-plane']" class="size-3.5" />
            </Button>
            <Button
              v-else
              size="sm"
              variant="destructive"
              @click="chatStore.abort()"
            >
              <FontAwesomeIcon :icon="['fas', 'spinner']" class="size-3.5 animate-spin" />
            </Button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { Textarea, Button } from '@memoh/ui'
import { useChatList } from '@/store/chat-list'
import { storeToRefs } from 'pinia'
import MessageItem from './message-item.vue'

const chatStore = useChatList()
const { messages, streaming, currentBotId } = storeToRefs(chatStore)

const inputText = ref('')
const scrollContainer = ref<HTMLElement>()

let userScrolledUp = false

function scrollToBottom(smooth = true) {
  nextTick(() => {
    const el = scrollContainer.value
    if (!el) return
    el.scrollTo({
      top: el.scrollHeight,
      behavior: smooth ? 'smooth' : 'instant',
    })
  })
}

function handleScroll() {
  const el = scrollContainer.value
  if (!el) return
  const distanceFromBottom = el.scrollHeight - el.clientHeight - el.scrollTop
  userScrolledUp = distanceFromBottom > 50
}

watch(
  () => {
    const last = messages.value[messages.value.length - 1]
    return last?.blocks.reduce((acc, b) => {
      if (b.type === 'text') return acc + b.content.length
      if (b.type === 'thinking') return acc + b.content.length
      return acc + 1
    }, 0) ?? 0
  },
  () => {
    if (!userScrolledUp) scrollToBottom()
  },
)

watch(
  () => messages.value.length,
  () => {
    userScrolledUp = false
    scrollToBottom()
  },
)

watch(scrollContainer, (el) => {
  if (el) el.addEventListener('scroll', handleScroll, { passive: true })
}, { immediate: true })

function handleKeydown(e: KeyboardEvent) {
  if (e.isComposing) return
  e.preventDefault()
  handleSend()
}

function handleSend() {
  const text = inputText.value.trim()
  if (!text || streaming.value) return
  inputText.value = ''
  chatStore.sendMessage(text)
}
</script>
