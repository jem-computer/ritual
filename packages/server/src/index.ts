// ABOUTME: Main entry point for the Ritual server
// ABOUTME: Handles HTTP API and task scheduling

import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'
import { serve } from '@hono/node-server'

const app = new Hono()

// Middleware
app.use('/*', cors())
app.use(logger())

// Health check
app.get('/health', (c) => {
  return c.json({ status: 'ok' })
})

// Task routes
app.get('/api/tasks', (c) => {
  // TODO: Return all tasks
  return c.json([])
})

app.post('/api/tasks', async (c) => {
  // TODO: Create new task
  const body = await c.req.json()
  return c.json({ id: '1', ...body }, 201)
})

app.put('/api/tasks/:id', async (c) => {
  // TODO: Update task
  const id = c.req.param('id')
  const body = await c.req.json()
  return c.json({ id, ...body })
})

app.delete('/api/tasks/:id', (c) => {
  // TODO: Delete task
  return c.body(null, 204)
})

// Log routes
app.get('/api/logs', (c) => {
  // TODO: Return execution logs
  return c.json([])
})

// Start server
const port = process.env.PORT || 8080
console.log(`Server running on port ${port}`)

serve({
  fetch: app.fetch,
  port: Number(port),
})