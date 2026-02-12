import { defineStore } from 'pinia'
import { computed, reactive, ref, watch } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import type { user, robot } from '@memoh/shared'
import {
  createChat,
  deleteChat as requestDeleteChat,
  type Bot as ChatBot,
  type ChatSummary,
  type Message as PersistedMessage,
  fetchBots,
  fetchMessages,
  fetchChats,
  extractAssistantTexts,
  extractMessageText,
  streamMessage,
  streamMessageEvents,
} from '@/composables/api/useChat'

export const useChatList = defineStore('chatList', () => {
  const defaultProcessingMessage = 'Message received, processing...'
  const chatList = reactive<(user | robot)[]>([])
  const chats = ref<ChatSummary[]>([])
  const loading = ref(false)
  const loadingChats = ref(false)
  const loadingOlder = ref(false)
  const hasMoreOlder = ref(true)
  const initializing = ref(false)
  const botId = useLocalStorage<string | null>('chat-bot-id', null)
  const chatId = useLocalStorage<string | null>('chat-id', null)
  const bots = ref<ChatBot[]>([])
  const participantChats = computed(() =>
    chats.value.filter((item) => (item.access_mode ?? 'participant') === 'participant'),
  )
  const observedChats = computed(() =>
    chats.value.filter((item) => item.access_mode === 'channel_identity_observed'),
  )
  const activeChat = computed(() =>
    chats.value.find((item) => item.id === chatId.value) ?? null,
  )
  const activeChatReadOnly = computed(() => activeChat.value?.access_mode === 'channel_identity_observed')
  let messageEventsController: AbortController | null = null
  let messageEventsLoopVersion = 0
  let messageEventsSince = ''

  // Watch for botId changes to re-initialize
  watch(botId, (newBotId) => {
    if (newBotId) {
      void initialize()
    } else {
      stopMessageEvents()
      messageEventsSince = ''
      chats.value = []
      chatId.value = null
      replaceMessages([])
    }
  })

  const nextId = () => `${Date.now()}-${Math.floor(Math.random() * 1000)}`
  const isPendingBot = (bot: ChatBot | null | undefined) => (
    bot?.status === 'creating'
    || bot?.status === 'deleting'
  )

  const resolveBotIdentityLabel = (targetBotID?: string | null) => {
    const activeBotID = targetBotID ?? botId.value
    if (!activeBotID) {
      return 'Assistant'
    }
    const currentBot = bots.value.find((item) => item.id === activeBotID)
    return currentBot?.display_name?.trim() || currentBot?.id || 'Assistant'
  }

  const sleep = (ms: number) => new Promise<void>((resolve) => {
    window.setTimeout(resolve, ms)
  })

  const stopMessageEvents = () => {
    messageEventsLoopVersion += 1
    if (messageEventsController) {
      messageEventsController.abort()
      messageEventsController = null
    }
  }

  const updateMessageEventsSince = (createdAt?: string) => {
    const value = (createdAt ?? '').trim()
    if (!value) {
      return
    }
    if (!messageEventsSince) {
      messageEventsSince = value
      return
    }
    const currentMs = Date.parse(messageEventsSince)
    const nextMs = Date.parse(value)
    if (Number.isNaN(nextMs)) {
      return
    }
    if (Number.isNaN(currentMs) || nextMs > currentMs) {
      messageEventsSince = value
    }
  }

  const updateMessageEventsSinceFromRows = (rows: PersistedMessage[]) => {
    messageEventsSince = ''
    for (const row of rows) {
      updateMessageEventsSince(row.created_at)
    }
  }

  const hasMessageWithID = (messageID: string) => {
    const targetID = messageID.trim()
    if (!targetID) {
      return false
    }
    return chatList.some((item) => String(item.id ?? '').trim() === targetID)
  }

  const appendRealtimeMessage = (raw: PersistedMessage) => {
    updateMessageEventsSince(raw.created_at)

    const platform = (raw.platform ?? '').trim().toLowerCase()
    if (platform === 'web') {
      // Web messages are already rendered by the request-scoped stream path.
      return
    }

    const messageID = String(raw.id ?? '').trim()
    if (messageID && hasMessageWithID(messageID)) {
      return
    }

    const item = toChatItem(raw)
    if (!item) {
      return
    }
    chatList.push(item)
    chatList.sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime())
    if (chatId.value) {
      touchChat(chatId.value)
    }
  }

  const handleMessageStreamEvent = (targetBotID: string, event: Record<string, unknown>) => {
    const eventType = String(event.type ?? '').toLowerCase()
    if (eventType !== 'message_created') {
      return
    }

    const eventBotID = String(event.bot_id ?? '').trim()
    if (eventBotID && eventBotID !== targetBotID) {
      return
    }

    const payload = event.message
    if (!payload || typeof payload !== 'object') {
      return
    }

    const raw = payload as PersistedMessage
    const payloadBotID = String(raw.bot_id ?? '').trim()
    if (payloadBotID && payloadBotID !== targetBotID) {
      return
    }
    appendRealtimeMessage(raw)
  }

  const startMessageEvents = (targetBotID: string) => {
    const normalizedBotID = targetBotID.trim()
    stopMessageEvents()
    if (!normalizedBotID) {
      return
    }

    const controller = new AbortController()
    messageEventsController = controller
    const loopVersion = messageEventsLoopVersion

    const run = async () => {
      let retryDelayMs = 1000
      while (!controller.signal.aborted && messageEventsLoopVersion === loopVersion) {
        try {
          await streamMessageEvents(
            normalizedBotID,
            controller.signal,
            (event) => {
              handleMessageStreamEvent(normalizedBotID, event as Record<string, unknown>)
            },
            messageEventsSince || undefined,
          )
          retryDelayMs = 1000
          if (!controller.signal.aborted && messageEventsLoopVersion === loopVersion) {
            await sleep(300)
          }
        } catch (error) {
          if (controller.signal.aborted || messageEventsLoopVersion !== loopVersion) {
            return
          }
          console.error('Message events stream failed:', error)
          await sleep(retryDelayMs)
          retryDelayMs = Math.min(retryDelayMs * 2, 5000)
        }
      }
    }

    void run()
  }

  const addUserMessage = (text: string) => {
    chatList.push({
      description: text,
      time: new Date(),
      action: 'user',
      id: nextId(),
    })
  }

  const addRobotMessage = (text: string, state: robot['state'] = 'complete') => {
    const id = nextId()
    chatList.push({
      description: text,
      time: new Date(),
      action: 'robot',
      id,
      type: resolveBotIdentityLabel(),
      state,
    })
    return id
  }

  const updateRobotMessage = (id: string, patch: Partial<robot>) => {
    const target = chatList.find(
      (item): item is robot => item.action === 'robot' && String(item.id) === id,
    )
    if (target) {
      Object.assign(target, patch)
    }
  }

  const ensureBot = async () => {
    try {
      const botsList = await fetchBots()
      bots.value = botsList
      if (!botsList.length) {
        botId.value = null
        return null
      }
      // If we have a persisted botId and it's still valid, use it
      if (botId.value) {
        const persisted = botsList.find((b) => b.id === botId.value)
        if (persisted && !isPendingBot(persisted)) {
          return botId.value
        }
      }
      const firstReadyBot = botsList.find((item) => !isPendingBot(item))
      if (firstReadyBot) {
        botId.value = firstReadyBot.id
        return botId.value
      }
      // Fallback to the first bot when all bots are in a pending lifecycle state.
      botId.value = botsList[0]!.id
      return botId.value
    } catch (error) {
      console.error('Failed to fetch bots:', error)
      return botId.value // Fallback to whatever we have
    }
  }

  const replaceMessages = (items: (user | robot)[]) => {
    chatList.splice(0, chatList.length, ...items)
  }

  const toChatItem = (raw: Awaited<ReturnType<typeof fetchMessages>>[number]): user | robot | null => {
    if (raw.role !== 'user' && raw.role !== 'assistant') {
      return null
    }

    const text = extractMessageText(raw)
    if (!text) {
      return null
    }

    const createdAt = raw.created_at ? new Date(raw.created_at) : new Date()
    const time = Number.isNaN(createdAt.getTime()) ? new Date() : createdAt
    const itemID = raw.id || nextId()

    const platform = (raw.platform ?? '').trim().toLowerCase()
    const channelTag = platform && platform !== 'web' ? platform : undefined

    if (raw.role === 'user') {
      return {
        description: text,
        time,
        action: 'user',
        id: itemID,
        ...(channelTag && { platform: channelTag }),
      }
    }

    return {
      description: text,
      time,
      action: 'robot',
      id: itemID,
      type: resolveBotIdentityLabel(raw.bot_id || botId.value),
      state: 'complete',
      ...(channelTag && { platform: channelTag }),
    }
  }

  const PAGE_SIZE = 30

  const loadMessages = async (targetBotID: string, targetChatID: string) => {
    const rows = await fetchMessages(targetBotID, targetChatID, { limit: PAGE_SIZE })
    const items = rows
      .map(toChatItem)
      .filter((item): item is user | robot => item !== null)
    replaceMessages(items)
    hasMoreOlder.value = true
    updateMessageEventsSinceFromRows(rows)
  }

  const loadOlderMessages = async (): Promise<number> => {
    const currentBotID = botId.value ?? ''
    const currentChatID = chatId.value ?? ''
    if (!currentBotID || !currentChatID || loadingOlder.value || !hasMoreOlder.value) {
      return 0
    }
    const first = chatList[0]
    if (!first?.time) {
      return 0
    }
    const before = typeof first.time === 'object' && first.time instanceof Date
      ? first.time.toISOString()
      : new Date(first.time).toISOString()
    loadingOlder.value = true
    try {
      const rows = await fetchMessages(currentBotID, currentChatID, { limit: PAGE_SIZE, before })
      const items = rows
        .map(toChatItem)
        .filter((item): item is user | robot => item !== null)
      if (items.length < PAGE_SIZE) {
        hasMoreOlder.value = false
      }
      chatList.unshift(...items)
      return items.length
    } finally {
      loadingOlder.value = false
    }
  }

  const initialize = async () => {
    if (initializing.value) {
      return
    }

    initializing.value = true
    loadingChats.value = true
    stopMessageEvents()
    try {
      const currentBotID = await ensureBot()
      if (!currentBotID) {
        messageEventsSince = ''
        chats.value = []
        chatId.value = null
        replaceMessages([])
        return
      }
      const visibleChats = await fetchChats(currentBotID)
      chats.value = visibleChats

      if (visibleChats.length === 0) {
        messageEventsSince = ''
        chatId.value = null
        replaceMessages([])
        return
      }

      const activeChatID = chatId.value && visibleChats.some((item) => item.id === chatId.value)
        ? chatId.value
        : visibleChats[0]!.id
      chatId.value = activeChatID
      await loadMessages(currentBotID, activeChatID)
      startMessageEvents(currentBotID)
    } finally {
      loadingChats.value = false
      initializing.value = false
    }
  }

  const selectBot = async (targetBotID: string) => {
    if (botId.value === targetBotID) {
      return
    }
    botId.value = targetBotID
    chatId.value = null
    await initialize()
  }

  const createNewChat = async () => {
    loadingChats.value = true
    try {
      const currentBotID = await ensureBot()
      if (!currentBotID) return
      const created = await createChat(currentBotID)
      chats.value = [created, ...chats.value.filter((item) => item.id !== created.id)]
      chatId.value = created.id
      replaceMessages([])
    } finally {
      loadingChats.value = false
    }
  }

  const removeChat = async (targetChatID: string) => {
    const deletingChatID = targetChatID.trim()
    if (!deletingChatID) {
      return
    }

    loadingChats.value = true
    try {
      const currentBotID = botId.value ?? ''
      if (!currentBotID) {
        throw new Error('Bot not selected')
      }
      await requestDeleteChat(currentBotID, deletingChatID)
      const remainingChats = chats.value.filter((item) => item.id !== deletingChatID)
      chats.value = remainingChats

      if (chatId.value !== deletingChatID) {
        return
      }

      if (remainingChats.length === 0) {
        chatId.value = null
        replaceMessages([])
        return
      }

      const nextChatID = remainingChats[0]!.id
      chatId.value = nextChatID
      await loadMessages(currentBotID, nextChatID)
    } finally {
      loadingChats.value = false
    }
  }

  const selectChat = async (targetChatID: string) => {
    const nextChatID = targetChatID.trim()
    if (!nextChatID || nextChatID === chatId.value) {
      return
    }

    chatId.value = nextChatID
    loadingChats.value = true
    try {
      const currentBotID = botId.value ?? ''
      if (!currentBotID) {
        throw new Error('Bot not selected')
      }
      await loadMessages(currentBotID, nextChatID)
    } finally {
      loadingChats.value = false
    }
  }

  const ensureActiveChat = async () => {
    if (chatId.value) {
      return
    }
    const currentBotID = botId.value ?? await ensureBot()
    if (!currentBotID) {
      throw new Error('Bot not ready')
    }
    const created = await createChat(currentBotID)
    chats.value = [created, ...chats.value.filter((item) => item.id !== created.id)]
    chatId.value = created.id
    replaceMessages([])
  }

  const touchChat = (targetChatID: string) => {
    const index = chats.value.findIndex((item) => item.id === targetChatID)
    if (index < 0) {
      return
    }
    const [target] = chats.value.splice(index, 1)
    if (!target) {
      return
    }
    target.updated_at = new Date().toISOString()
    chats.value.unshift(target)
  }

  const sendMessage = async (text: string) => {
    const trimmed = text.trim()
    if (!trimmed) return

    loading.value = true
    let thinkingId: string | null = null
    try {
      await ensureActiveChat()
      const activeChatID = chatId.value!
      if (activeChatReadOnly.value) {
        throw new Error('Chat is read-only')
      }
      addUserMessage(trimmed)

      thinkingId = addRobotMessage(defaultProcessingMessage, 'thinking')
      const currentThinkingID = thinkingId
      let streamedText = ''
      const activeBotID = botId.value!
      const finalResponse = await streamMessage(
        activeBotID,
        activeChatID,
        trimmed,
        (delta) => {
          if (!delta) {
            return
          }
          streamedText += delta
          updateRobotMessage(currentThinkingID, {
            description: streamedText,
            state: 'generate',
          })
        },
        (status) => {
          if (status !== 'started') {
            return
          }
          updateRobotMessage(currentThinkingID, {
            description: defaultProcessingMessage,
            state: 'thinking',
          })
        },
      )

      if (streamedText.trim()) {
        updateRobotMessage(currentThinkingID, {
          description: streamedText.trim(),
          state: 'complete',
        })
        touchChat(activeChatID)
        return
      }

      const assistantTexts = extractAssistantTexts(finalResponse?.messages ?? [])
      if (assistantTexts.length === 0) {
        updateRobotMessage(currentThinkingID, {
          description: 'No textual response.',
          state: 'complete',
        })
        touchChat(activeChatID)
        return
      }

      updateRobotMessage(currentThinkingID, {
        description: assistantTexts[0]!,
        state: 'complete',
      })
      for (const textItem of assistantTexts.slice(1)) {
        addRobotMessage(textItem)
      }
      touchChat(activeChatID)
    } catch (error) {
      const reason = error instanceof Error ? error.message : 'Unknown error'
      if (thinkingId) {
        updateRobotMessage(thinkingId, {
          description: `Failed to send message: ${reason}`,
          state: 'complete',
        })
      } else {
        addRobotMessage(`Failed to send message: ${reason}`)
      }
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    chatList,
    chats,
    participantChats,
    observedChats,
    chatId,
    botId,
    bots,
    activeChat,
    activeChatReadOnly,
    loading,
    loadingChats,
    loadingOlder,
    hasMoreOlder,
    initializing,
    initialize,
    selectBot,
    selectChat,
    createNewChat,
    removeChat,
    deleteChat: removeChat,
    sendMessage,
    loadOlderMessages,
  }
})
