import {
  getBots, getBotsById, postBots, deleteBotsById, putBotsById,
  getBotsByIdChecks, getBotsByBotIdContainer, postBotsByBotIdContainer,
  deleteBotsByBotIdContainer, postBotsByBotIdContainerStart,
  postBotsByBotIdContainerStop, getBotsByBotIdContainerSnapshots,
  postBotsByBotIdContainerSnapshots,
} from '@memoh/sdk'
import type {
  BotsBot, BotsListBotsResponse, BotsCreateBotRequest, BotsUpdateBotRequest,
  BotsBotCheck, BotsListChecksResponse,
  HandlersGetContainerResponse, HandlersCreateContainerRequest,
  HandlersCreateContainerResponse, HandlersListSnapshotsResponse,
  HandlersCreateSnapshotRequest, HandlersCreateSnapshotResponse,
} from '@memoh/sdk'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import type { Ref } from 'vue'

// ---- Types (re-export with aliases for backward compat) ----

export type BotStatus = 'creating' | 'ready' | 'deleting'
export type BotCheckState = 'ok' | 'issue' | 'unknown'

export type BotInfo = BotsBot
export type ListBotsResponse = BotsListBotsResponse
export type CreateBotRequest = BotsCreateBotRequest
export type UpdateBotRequest = BotsUpdateBotRequest
export type BotCheck = BotsBotCheck
export type BotCheckListResponse = BotsListChecksResponse
export type BotContainerInfo = HandlersGetContainerResponse
export type CreateBotContainerRequest = HandlersCreateContainerRequest
export type CreateBotContainerResponse = HandlersCreateContainerResponse
export type BotContainerSnapshotListResponse = HandlersListSnapshotsResponse
export type CreateBotSnapshotRequest = HandlersCreateSnapshotRequest
export type CreateBotSnapshotResponse = HandlersCreateSnapshotResponse

export const BOT_PENDING_STATUSES: readonly BotStatus[] = ['creating', 'deleting']

export function isBotPendingStatus(status: string | undefined | null): boolean {
  return BOT_PENDING_STATUSES.includes((status ?? 'ready') as BotStatus)
}

export interface ContainerActionResponse {
  started?: boolean
  stopped?: boolean
}

// ---- Query: list bots ----

export function useBotList() {
  const queryCache = useQueryCache()

  const query = useQuery({
    key: ['bots'],
    query: async (): Promise<BotInfo[]> => {
      const { data } = await getBots({ throwOnError: true })
      return data.items ?? []
    },
  })

  return {
    ...query,
    invalidate: () => queryCache.invalidateQueries({ key: ['bots'] }),
  }
}

// ---- Query: bot detail ----

export function useBotDetail(botId: Ref<string>) {
  return useQuery({
    key: () => ['bot', botId.value],
    query: async () => {
      const { data } = await getBotsById({
        path: { id: botId.value },
        throwOnError: true,
      })
      return data
    },
    enabled: () => !!botId.value,
  })
}

// ---- Mutations ----

export function useCreateBot() {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async (body: CreateBotRequest) => {
      const { data } = await postBots({ body, throwOnError: true })
      return data
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['bots'] }),
  })
}

export function useDeleteBot() {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async (botId: string) => {
      await deleteBotsById({ path: { id: botId }, throwOnError: true })
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['bots'] }),
  })
}

export function useUpdateBot() {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async ({ id, ...body }: UpdateBotRequest & { id: string }) => {
      const { data } = await putBotsById({
        path: { id },
        body,
        throwOnError: true,
      })
      return data
    },
    onSettled: () => {
      queryCache.invalidateQueries({ key: ['bots'] })
      queryCache.invalidateQueries({ key: ['bot'] })
    },
  })
}

export async function fetchBotChecks(botId: string): Promise<BotCheck[]> {
  const { data } = await getBotsByIdChecks({
    path: { id: botId },
    throwOnError: true,
  })
  return data.items ?? []
}

export async function fetchBotContainer(botId: string): Promise<BotContainerInfo> {
  const { data } = await getBotsByBotIdContainer({
    path: { bot_id: botId },
    throwOnError: true,
  })
  return data
}

export async function createBotContainer(
  botId: string,
  payload: CreateBotContainerRequest = {},
): Promise<CreateBotContainerResponse> {
  const { data } = await postBotsByBotIdContainer({
    path: { bot_id: botId },
    body: payload,
    throwOnError: true,
  })
  return data
}

export async function deleteBotContainer(botId: string): Promise<void> {
  await deleteBotsByBotIdContainer({
    path: { bot_id: botId },
    throwOnError: true,
  })
}

export async function startBotContainer(botId: string): Promise<ContainerActionResponse> {
  const { data } = await postBotsByBotIdContainerStart({
    path: { bot_id: botId },
    throwOnError: true,
  })
  return data as ContainerActionResponse
}

export async function stopBotContainer(botId: string): Promise<ContainerActionResponse> {
  const { data } = await postBotsByBotIdContainerStop({
    path: { bot_id: botId },
    throwOnError: true,
  })
  return data as ContainerActionResponse
}

export async function listBotContainerSnapshots(
  botId: string,
  snapshotter?: string,
): Promise<BotContainerSnapshotListResponse> {
  const { data } = await getBotsByBotIdContainerSnapshots({
    path: { bot_id: botId },
    query: snapshotter?.trim() ? { snapshotter: snapshotter.trim() } : undefined,
    throwOnError: true,
  })
  return data
}

export async function createBotContainerSnapshot(
  botId: string,
  payload: CreateBotSnapshotRequest,
): Promise<CreateBotSnapshotResponse> {
  const { data } = await postBotsByBotIdContainerSnapshots({
    path: { bot_id: botId },
    body: payload,
    throwOnError: true,
  })
  return data
}
