import { getChannels, getBotsByIdChannelByPlatform, putBotsByIdChannelByPlatform } from '@memoh/sdk'
import type {
  HandlersChannelMeta, ChannelChannelCapabilities, ChannelConfigSchema,
  ChannelFieldSchema, ChannelFieldType, ChannelChannelConfig,
  ChannelUpsertConfigRequest,
} from '@memoh/sdk'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import type { Ref } from 'vue'

// ---- Types (re-export SDK types) ----

export type FieldSchema = ChannelFieldSchema
export type FieldType = ChannelFieldType
export type ConfigSchema = ChannelConfigSchema
export type ChannelCapabilities = ChannelChannelCapabilities
export type ChannelMeta = HandlersChannelMeta
export type ChannelConfig = ChannelChannelConfig
export type UpsertConfigRequest = ChannelUpsertConfigRequest

export interface BotChannelItem {
  meta: ChannelMeta
  config: ChannelConfig | null
  configured: boolean
}

// ---- Query: channel type metadata ----

export function useChannelMetas() {
  return useQuery({
    key: ['channel-metas'],
    query: async () => {
      const { data } = await getChannels({ throwOnError: true })
      return data
    },
  })
}

// ---- Query: bot channel configs (combined) ----

export function useBotChannels(botId: Ref<string>) {
  const queryCache = useQueryCache()

  const query = useQuery({
    key: () => ['bot-channels', botId.value],
    query: async (): Promise<BotChannelItem[]> => {
      const { data: metas } = await getChannels({ throwOnError: true })
      const configurableTypes = (metas as ChannelMeta[]).filter((m) => !m.configless)
      const results = await Promise.all(
        configurableTypes.map(async (meta) => {
          try {
            const { data: config } = await getBotsByIdChannelByPlatform({
              path: { id: botId.value, platform: meta.type! },
              throwOnError: true,
            })
            return { meta, config, configured: true } as BotChannelItem
          } catch {
            return { meta, config: null, configured: false } as BotChannelItem
          }
        }),
      )
      return results
    },
    enabled: () => !!botId.value,
  })

  return {
    ...query,
    invalidate: () => queryCache.invalidateQueries({ key: ['bot-channels', botId.value] }),
  }
}

// ---- Mutation: upsert bot channel config ----

export function useUpsertBotChannel(botId: Ref<string>) {
  const queryCache = useQueryCache()

  return useMutation({
    mutation: async ({ platform, data }: { platform: string; data: UpsertConfigRequest }) => {
      const { data: res } = await putBotsByIdChannelByPlatform({
        path: { id: botId.value, platform },
        body: data,
        throwOnError: true,
      })
      return res
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['bot-channels', botId.value] }),
  })
}
