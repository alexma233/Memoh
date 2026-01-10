import { tool } from 'ai'
import { z } from 'zod'

export interface GetMemoryToolParams {
  searchMemory: (query: string) => Promise<object[]>
}

export const getMemoryTools = ({ searchMemory }: GetMemoryToolParams) => {
  const searchMemoryTool = tool({
    description: 'Search chat history in the memory',
    inputSchema: z.object({
      query: z.string().describe('The query to search the memory'),
    }),
    execute: async ({ query }) => {
      const memory = await searchMemory(query)
      console.log(memory)
      return {
        success: true,
        memories: memory,
      }
    },
  })

  return {
    'search-memory': searchMemoryTool,
  }
}