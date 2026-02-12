import { getBotsByBotIdSettings, putBotsByBotIdSettings } from '@memoh/sdk'
import type { SettingsSettings, SettingsUpsertRequest } from '@memoh/sdk'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import type { Ref } from 'vue'

// ---- Types ----

export type BotSettings = SettingsSettings
export type UpsertBotSettingsRequest = SettingsUpsertRequest

// ---- Query ----

export function useBotSettings(botId: Ref<string>) {
  return useQuery({
    key: () => ['bot-settings', botId.value],
    query: async () => {
      const { data } = await getBotsByBotIdSettings({
        path: { bot_id: botId.value },
        throwOnError: true,
      })
      return data
    },
    enabled: () => !!botId.value,
  })
}

// ---- Mutation ----

export function useUpdateBotSettings(botId: Ref<string>) {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async (body: UpsertBotSettingsRequest) => {
      const { data } = await putBotsByBotIdSettings({
        path: { bot_id: botId.value },
        body,
        throwOnError: true,
      })
      return data
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['bot-settings', botId.value] }),
  })
}
