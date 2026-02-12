import {
  getModels, postModels, putModelsModelByModelId,
  deleteModelsModelByModelId, getProvidersByIdModels,
} from '@memoh/sdk'
import type { ModelsAddRequest } from '@memoh/sdk'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import { type ModelInfo } from '@memoh/shared'
import type { Ref } from 'vue'

// ---- Types ----

export type CreateModelRequest = ModelsAddRequest

// ---- Query: models by provider ----

export function useModelList(providerId: Ref<string | undefined>) {
  const queryCache = useQueryCache()

  const query = useQuery({
    key: ['model'],
    query: async () => {
      const { data } = await getProvidersByIdModels({
        path: { id: providerId.value! },
        throwOnError: true,
      })
      return data as ModelInfo[]
    },
  })

  return {
    ...query,
    invalidate: () => queryCache.invalidateQueries({ key: ['model'] }),
  }
}

// ---- Query: all models (cross-provider) ----

export function useAllModels() {
  return useQuery({
    key: ['all-models'],
    query: async () => {
      const { data } = await getModels({ throwOnError: true })
      return data as ModelInfo[]
    },
  })
}

// ---- Mutations ----

export function useCreateModel() {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async (body: CreateModelRequest) => {
      const { data } = await postModels({ body, throwOnError: true })
      return data
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['model'], exact: true }),
  })
}

export function useUpdateModel() {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async ({ modelId, data }: { modelId: string; data: Partial<CreateModelRequest> }) => {
      const { data: res } = await putModelsModelByModelId({
        path: { modelId },
        body: data,
        throwOnError: true,
      })
      return res
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['model'] }),
  })
}

export function useDeleteModel() {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async (modelName: string) => {
      await deleteModelsModelByModelId({
        path: { modelId: modelName },
        throwOnError: true,
      })
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['model'] }),
  })
}
