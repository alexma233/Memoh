import { tool } from 'ai'
import { z } from 'zod'
import { AuthFetcher } from '..'
import type { IdentityContext } from '../types'

export type ContactToolParams = {
  fetch: AuthFetcher
  identity: IdentityContext
}

export const getContactTools = ({ fetch, identity }: ContactToolParams) => {
  const botId = identity.botId.trim()

  const contactSearch = tool({
    description: 'Search contacts by name or alias',
    inputSchema: z.object({
      query: z.string().describe('The query to search for contacts'),
    }),
    execute: async ({ query }) => {
      const url = `/bots/${botId}/contacts?q=${encodeURIComponent(query)}`
      const response = await fetch(url)
      return response.json()
    },
  })

  const contactCreate = tool({
    description: 'Create a contact',
    inputSchema: z.object({
      name: z.string().describe('The display name of the contact'),
      alias: z.string().describe('The alias of the contact').optional(),
      tags: z.array(z.string()).describe('The tags of the contact').optional(),
    }),
    execute: async ({ name, alias, tags }) => {
      const response = await fetch(`/bots/${botId}/contacts`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          display_name: name,
          alias: alias,
          tags: tags ?? [],
        }),
      })
      return response.json()
    },
  })

  const contactUpdate = tool({
    description: 'Update a contact',
    inputSchema: z.object({
      contact_id: z.string().describe('The ID of the contact to update'),
      name: z.string().describe('The display name of the contact').optional(),
      alias: z.string().describe('The alias of the contact').optional(),
      tags: z.array(z.string()).describe('The tags of the contact').optional(),
    }),
    execute: async ({ contact_id, name, alias, tags }) => {
      const response = await fetch(`/bots/${botId}/contacts/${contact_id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          display_name: name,
          alias: alias,
          tags: tags ?? [],
        }),
      })
      return response.json()
    },
  })

  // const contactBindToken = tool({
  //   description: 'Issue a one-time bind token for a contact',
  //   inputSchema: z.object({
  //     contact_id: ContactID,
  //     target_platform: z.string().describe('The platform to bind the contact to'),
  //     target_external_id: z.string().describe('The external ID of the contact'),
  //     ttl_seconds: z.number().describe('The number of seconds the bind token is valid').optional(),
  //   }),
  //   execute: async ({ bot_id, contact_id, target_platform, target_external_id, ttl_seconds }) => {
  //     const response = await fetch(`/bots/${botId}/contacts/${contact_id}/bind_token`, {
  //       method: 'POST',
  //       headers: { 'Content-Type': 'application/json' },
  //       body: JSON.stringify({
  //         target_platform: target_platform,
  //         target_external_id: target_external_id,
  //         ttl_seconds: ttl_seconds,
  //       }),
  //     })
  //     return response.json()
  //   },
  // })

  const contactBind = tool({
    description: 'Bind a contact to a platform identity using a bind token',
    inputSchema: z.object({
      contact_id: z.string().describe('The ID of the contact to bind'),
      platform: z.string().describe('The platform to bind the contact to'),
      external_id: z.string().describe('The external ID of the contact'),
      bind_token: z.string().describe('The bind token to use'),
    }),
    execute: async ({ contact_id, platform, external_id, bind_token }) => {
      const response = await fetch(`/bots/${botId}/contacts/${contact_id}/bind`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          platform: platform,
          external_id: external_id,
          bind_token: bind_token,
        }),
      })
      return response.json()
    },
  })

  return {
    'contact_search': contactSearch,
    'contact_create': contactCreate,
    'contact_update': contactUpdate,
    // 'contact_bind_token': contactBindToken,
    'contact_bind': contactBind,
  }
}
