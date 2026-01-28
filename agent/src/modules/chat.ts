import { Elysia, sse } from 'elysia'
import z from 'zod'
import { createAgent } from '../agent'
import { ClientType } from '../types'
import { ModelMessage } from 'ai'

const ChatBody = z.object({
  apiKey: z.string().min(1, 'API key is required'),
  baseUrl: z.string().min(1, 'Base URL is required'),
  model: z.string().min(1, 'Model is required'),
  clientType: z.enum([
    'openai',
    'anthropic',
    'google',
  ]),
  locale: z.string().optional(),
  language: z.string().optional(),
  maxSteps: z.number().optional(),
  maxContextLoadTime: z.number().min(1, 'Max context load time is required'),
  platforms: z.array(z.string()).optional(),
  currentPlatform: z.string().optional(),

  messages: z.array(z.any()),
  query: z.string().min(1, 'Query is required'),
})

const ScheduleBody = z.object({
  schedule: z.object({
    id: z.string().min(1, 'Schedule ID is required'),
    name: z.string().min(1, 'Schedule name is required'),
    description: z.string().min(1, 'Schedule description is required'),
    pattern: z.string().min(1, 'Schedule pattern is required'),
    maxCalls: z.number().optional(),
    command: z.string().min(1, 'Schedule command is required'),
  }),
}).and(ChatBody)

export const chatModule = new Elysia({ prefix: '/chat' })
  .post('/', async ({ body }) => {
    console.log('[Chat] request', {
      type: 'chat',
      clientType: body.clientType,
      model: body.model,
      baseUrl: body.baseUrl,
    })
    const { ask } = createAgent({
      apiKey: body.apiKey,
      baseUrl: body.baseUrl,
      model: body.model,
      clientType: body.clientType as ClientType,
      locale: body.locale,
      language: body.language,
      maxSteps: body.maxSteps,
      maxContextLoadTime: body.maxContextLoadTime,
      platforms: body.platforms,
      currentPlatform: body.currentPlatform,
    })
    try {
      const result = await ask({
        messages: body.messages as unknown as ModelMessage[],
        query: body.query,
      })
      console.log('[Chat] response', { type: 'chat', messages: result.messages?.length ?? 0 })
      return result
    } catch (error) {
      console.error('[Chat] error', {
        type: 'chat',
        clientType: body.clientType,
        model: body.model,
        baseUrl: body.baseUrl,
        error,
      })
      throw error
    }
  }, {
    body: ChatBody,
  })
  .post('/stream', async function* ({ body }) {
    console.log('[Chat] request', {
      type: 'stream',
      clientType: body.clientType,
      model: body.model,
      baseUrl: body.baseUrl,
    })
    const { stream } = createAgent({
      apiKey: body.apiKey,
      baseUrl: body.baseUrl,
      model: body.model,
      clientType: body.clientType as ClientType,
      locale: body.locale,
      language: body.language,
      maxSteps: body.maxSteps,
      maxContextLoadTime: body.maxContextLoadTime,
      platforms: body.platforms,
      currentPlatform: body.currentPlatform,
    })
    try {
      const streanGenerator = stream({
        messages: body.messages as unknown as ModelMessage[],
        query: body.query,
      })
      while (true) {
        const chunk = await streanGenerator.next()
        if (chunk.done) {
          console.log('[Chat] response', { type: 'stream', messages: chunk.value?.messages?.length ?? 0 })
          yield sse({
            type: 'done',
            data: chunk.value,
          })
          break
        }
        yield sse({
          type: 'delta',
          data: chunk.value
        })
      }
    } catch (error) {
      console.error('[Chat] error', {
        type: 'stream',
        clientType: body.clientType,
        model: body.model,
        baseUrl: body.baseUrl,
        error,
      })
      throw error
    }
  }, {
    body: ChatBody,
  })
  .post('/schedule', async ({ body }) => {
    const { triggerSchedule } = createAgent({
      apiKey: body.apiKey,
      baseUrl: body.baseUrl,
      model: body.model,
      clientType: body.clientType as ClientType,
      locale: body.locale,
      language: body.language,
      maxSteps: body.maxSteps,
      maxContextLoadTime: body.maxContextLoadTime,
      platforms: body.platforms,
      currentPlatform: body.currentPlatform,
    })
    return await triggerSchedule({
      messages: body.messages as unknown as ModelMessage[],
      query: body.query,
    }, body.schedule)
  }, {
    body: ScheduleBody,
  })