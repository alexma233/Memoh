import { pgTable, text, timestamp, integer } from 'drizzle-orm/pg-core'
import { sql } from 'drizzle-orm'

export const memory = pgTable('memory', {
  id: text('id').primaryKey(),
  memoryId: text('memory_id').notNull(),
  previousValue: text('previous_value'),
  newValue: text('new_value'),
  action: text('action').notNull(),
  createdAt: timestamp('created_at', { withTimezone: true }).default(sql`timezone('utc', now())`),
  updatedAt: timestamp('updated_at', { withTimezone: true }),
  isDeleted: integer('is_deleted').default(0),
})