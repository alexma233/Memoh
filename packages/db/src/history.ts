import { pgTable, timestamp, uuid, jsonb, text } from 'drizzle-orm/pg-core'

export const history = pgTable(
  'history', 
  {
    id: uuid('id').primaryKey().defaultRandom(),
    messages: jsonb('messages').notNull(),
    timestamp: timestamp('timestamp').notNull(),
    user: text('user').notNull(),
  }
)