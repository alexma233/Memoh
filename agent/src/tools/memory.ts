import { tool } from 'ai'
import { AuthFetcher } from '..'
import type { IdentityContext } from '../types'
import { z } from 'zod'

export type MemoryToolParams = {
  fetch: AuthFetcher
  identity: IdentityContext
}

type MemorySearchItem = {
  id?: string
  memory?: string
  score?: number
  createdAt?: string
  metadata?: {
    source?: string
  }
}

export const getMemoryTools = ({ fetch, identity }: MemoryToolParams) => {
  const searchMemory = tool({
    description: 'Search for memories',
    inputSchema: z.object({
      query: z.string().describe('The query to search for memories'),
      limit: z.number().int().positive().max(50).optional(),
    }),
    execute: async ({ query, limit }) => {
      const botId = identity.botId.trim()
      if (!botId) {
        throw new Error('botId is required to search memory')
      }
      const response = await fetch(`/bots/${botId}/memory/search`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          query,
          limit,
        }),
      })
      const data = await response.json()
      const results = Array.isArray(data?.results)
        ? (data.results as MemorySearchItem[])
        : []
      const simplified = results.map((item) => ({
        id: item?.id,
        memory: item?.memory,
        score: item?.score,
      }))
      return {
        query,
        total: simplified.length,
        results: simplified,
      }
    },
  })

  return {
    'search_memory': searchMemory,
  }
}