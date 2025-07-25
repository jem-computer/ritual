// ABOUTME: LibSQL database setup and schema management
// ABOUTME: Handles tasks, execution logs, and migrations

import { createClient, type Client } from '@libsql/client';
import { z } from 'zod';

// Database schema types
export const TaskSchema = z.object({
  id: z.string(),
  name: z.string(),
  status: z.enum(['ACTIVE', 'PAUSED']),
  prompt: z.string(),
  schedule: z.string(),
  output: z.string(),
  model: z.string(),
  nextRun: z.string().nullable(),
  lastRun: z.string().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
  jobId: z.string().nullable().optional(),
});

export type Task = z.infer<typeof TaskSchema>;

export const ExecutionLogSchema = z.object({
  id: z.string(),
  taskId: z.string(),
  taskName: z.string(),
  prompt: z.string(),
  output: z.string(),
  status: z.enum(['SUCCESS', 'FAILURE']),
  error: z.string().nullable(),
  executedAt: z.string(),
  duration: z.number(), // milliseconds
});

export type ExecutionLog = z.infer<typeof ExecutionLogSchema>;

// Database client
let db: Client;

export async function initDatabase() {
  // Create database connection
  const dbPath = process.env.DATABASE_URL || 'file:./ritual.db';
  db = createClient({
    url: dbPath,
  });

  // Run migrations
  await runMigrations();
  
  console.log('Database initialized');
  return db;
}

async function runMigrations() {
  // Create tasks table
  await db.execute(`
    CREATE TABLE IF NOT EXISTS tasks (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      status TEXT CHECK(status IN ('ACTIVE', 'PAUSED')) NOT NULL,
      prompt TEXT NOT NULL,
      schedule TEXT NOT NULL,
      output TEXT NOT NULL,
      model TEXT NOT NULL,
      next_run TEXT,
      last_run TEXT,
      created_at TEXT NOT NULL,
      updated_at TEXT NOT NULL,
      job_id TEXT
    )
  `);

  // Create execution_logs table
  await db.execute(`
    CREATE TABLE IF NOT EXISTS execution_logs (
      id TEXT PRIMARY KEY,
      task_id TEXT NOT NULL,
      task_name TEXT NOT NULL,
      prompt TEXT NOT NULL,
      output TEXT NOT NULL,
      status TEXT CHECK(status IN ('SUCCESS', 'FAILURE')) NOT NULL,
      error TEXT,
      executed_at TEXT NOT NULL,
      duration INTEGER NOT NULL,
      FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
    )
  `);

  // Create indexes
  await db.execute(`
    CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)
  `);
  
  await db.execute(`
    CREATE INDEX IF NOT EXISTS idx_execution_logs_task_id ON execution_logs(task_id)
  `);
  
  await db.execute(`
    CREATE INDEX IF NOT EXISTS idx_execution_logs_executed_at ON execution_logs(executed_at)
  `);
}

// Task operations
export async function getAllTasks(): Promise<Task[]> {
  const result = await db.execute(`
    SELECT id, name, status, prompt, schedule, output, model,
           next_run as nextRun, last_run as lastRun,
           created_at as createdAt, updated_at as updatedAt, job_id as jobId
    FROM tasks
    ORDER BY created_at DESC
  `);
  
  return result.rows.map(row => TaskSchema.parse({
    ...row,
    jobId: row.jobId || undefined,
  }));
}

export async function getTaskById(id: string): Promise<Task | null> {
  const result = await db.execute({
    sql: `
      SELECT id, name, status, prompt, schedule, output, model,
             next_run as nextRun, last_run as lastRun,
             created_at as createdAt, updated_at as updatedAt, job_id as jobId
      FROM tasks
      WHERE id = ?
    `,
    args: [id],
  });
  
  if (result.rows.length === 0) {
    return null;
  }
  
  return TaskSchema.parse({
    ...result.rows[0],
    jobId: result.rows[0].jobId || undefined,
  });
}

export async function createTask(task: Omit<Task, 'id' | 'createdAt' | 'updatedAt'>): Promise<Task> {
  const id = crypto.randomUUID();
  const now = new Date().toISOString();
  
  await db.execute({
    sql: `
      INSERT INTO tasks (id, name, status, prompt, schedule, output, model,
                        next_run, last_run, created_at, updated_at, job_id)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `,
    args: [
      id,
      task.name,
      task.status,
      task.prompt,
      task.schedule,
      task.output,
      task.model,
      task.nextRun,
      task.lastRun,
      now,
      now,
      task.jobId || null,
    ],
  });
  
  const created = await getTaskById(id);
  if (!created) {
    throw new Error('Failed to create task');
  }
  
  return created;
}

export async function updateTask(id: string, updates: Partial<Omit<Task, 'id' | 'createdAt'>>): Promise<Task | null> {
  const existing = await getTaskById(id);
  if (!existing) {
    return null;
  }
  
  const updated = {
    ...existing,
    ...updates,
    updatedAt: new Date().toISOString(),
  };
  
  await db.execute({
    sql: `
      UPDATE tasks
      SET name = ?, status = ?, prompt = ?, schedule = ?, output = ?, model = ?,
          next_run = ?, last_run = ?, updated_at = ?, job_id = ?
      WHERE id = ?
    `,
    args: [
      updated.name,
      updated.status,
      updated.prompt,
      updated.schedule,
      updated.output,
      updated.model,
      updated.nextRun,
      updated.lastRun,
      updated.updatedAt,
      updated.jobId || null,
      id,
    ],
  });
  
  return updated;
}

export async function deleteTask(id: string): Promise<boolean> {
  const result = await db.execute({
    sql: 'DELETE FROM tasks WHERE id = ?',
    args: [id],
  });
  
  return result.rowsAffected > 0;
}

// Execution log operations
export async function getAllExecutionLogs(): Promise<ExecutionLog[]> {
  const result = await db.execute(`
    SELECT id, task_id as taskId, task_name as taskName, prompt, output,
           status, error, executed_at as executedAt, duration
    FROM execution_logs
    ORDER BY executed_at DESC
    LIMIT 100
  `);
  
  return result.rows.map(row => ExecutionLogSchema.parse(row));
}

export async function getTaskExecutionLogs(taskId: string): Promise<ExecutionLog[]> {
  const result = await db.execute({
    sql: `
      SELECT id, task_id as taskId, task_name as taskName, prompt, output,
             status, error, executed_at as executedAt, duration
      FROM execution_logs
      WHERE task_id = ?
      ORDER BY executed_at DESC
      LIMIT 50
    `,
    args: [taskId],
  });
  
  return result.rows.map(row => ExecutionLogSchema.parse(row));
}

export async function createExecutionLog(log: Omit<ExecutionLog, 'id'>): Promise<ExecutionLog> {
  const id = crypto.randomUUID();
  
  await db.execute({
    sql: `
      INSERT INTO execution_logs (id, task_id, task_name, prompt, output,
                                 status, error, executed_at, duration)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `,
    args: [
      id,
      log.taskId,
      log.taskName,
      log.prompt,
      log.output,
      log.status,
      log.error,
      log.executedAt,
      log.duration,
    ],
  });
  
  return { ...log, id };
}

export function getDatabase(): Client {
  if (!db) {
    throw new Error('Database not initialized. Call initDatabase() first.');
  }
  return db;
}

export async function closeDatabase() {
  if (db) {
    await db.close();
  }
}