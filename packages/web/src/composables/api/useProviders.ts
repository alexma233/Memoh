import {
  getProviders, postProviders, putProvidersById, deleteProvidersById,
} from '@memoh/sdk'
import type { ProvidersCreateRequest, ProvidersUpdateRequest } from '@memoh/sdk'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import { type ProviderInfo } from '@memoh/shared'
import type { Ref } from 'vue'

// ---- Types ----

export type ProviderWithId = ProviderInfo & { id: string }
export type CreateProviderRequest = ProvidersCreateRequest
export type UpdateProviderRequest = ProvidersUpdateRequest

// ---- Query: provider list ----

export function useProviderList(clientType: Ref<string>) {
  return useQuery({
    key: ['provider'],
    query: async () => {
      const { data } = await getProviders({
        query: { client_type: clientType.value },
        throwOnError: true,
      })
      return data as ProviderWithId[]
    },
  })
}

/** Fetch all providers (no filter). */
export function useAllProviders() {
  return useQuery({
    key: ['all-providers'],
    query: async () => {
      const { data } = await getProviders({ throwOnError: true })
      return data as ProviderWithId[]
    },
  })
}

// ---- Mutations ----

export function useCreateProvider() {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async (body: CreateProviderRequest) => {
      const { data } = await postProviders({ body, throwOnError: true })
      return data
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['provider'] }),
  })
}

export function useUpdateProvider(providerId: Ref<string | undefined>) {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async (body: UpdateProviderRequest) => {
      const { data } = await putProvidersById({
        path: { id: providerId.value! },
        body,
        throwOnError: true,
      })
      return data
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['provider'] }),
  })
}

export function useDeleteProvider(providerId: Ref<string | undefined>) {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: async () => {
      await deleteProvidersById({
        path: { id: providerId.value! },
        throwOnError: true,
      })
    },
    onSettled: () => queryCache.invalidateQueries({ key: ['provider'] }),
  })
}
