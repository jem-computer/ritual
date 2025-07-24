// ABOUTME: Sets up BullMQ queue for scheduled task execution
// ABOUTME: Handles job processing and scheduling with Redis backend

import { Queue, Worker } from 'bullmq';
import Redis from 'ioredis';

// Redis connection with retry logic
const connection = new Redis({
  host: process.env.REDIS_HOST || 'localhost',
  port: parseInt(process.env.REDIS_PORT || '6379'),
  maxRetriesPerRequest: null,
  enableOfflineQueue: false,
  retryStrategy: (times) => {
    if (times > 3) {
      console.error('Redis connection failed after 3 attempts');
      return null;
    }
    return Math.min(times * 100, 3000);
  },
});

// Track Redis connection status
let redisConnected = false;

connection.on('connect', () => {
  console.log('Connected to Redis');
  redisConnected = true;
});

connection.on('error', (err) => {
  console.error('Redis connection error:', err.message);
  redisConnected = false;
});

export function isRedisConnected() {
  return redisConnected;
}

// Create queue for scheduled tasks
export const taskQueue = new Queue('ritual-tasks', {
  connection,
  defaultJobOptions: {
    removeOnComplete: {
      age: 3600, // keep completed jobs for 1 hour
      count: 100, // keep max 100 completed jobs
    },
    removeOnFail: {
      age: 24 * 3600, // keep failed jobs for 24 hours
    },
  },
});

// Note: BullMQ v5+ doesn't require QueueScheduler anymore
// It's handled automatically by the Worker

// Job processor
export const taskWorker = new Worker(
  'ritual-tasks',
  async (job) => {
    const { taskId, prompt, outputChannels } = job.data;
    
    console.log(`[${new Date().toISOString()}] Executing task ${taskId}:`);
    console.log(`  Prompt: ${prompt}`);
    console.log(`  Output channels: ${outputChannels.join(', ')}`);
    
    // TODO: In the future, this will:
    // 1. Execute the AI prompt
    // 2. Send results to specified output channels
    // For now, just log the execution
    
    return {
      taskId,
      executedAt: new Date().toISOString(),
      status: 'completed',
      message: `Task ${taskId} executed successfully`,
    };
  },
  {
    connection,
    concurrency: 5, // Process up to 5 jobs concurrently
  }
);

// Graceful shutdown
export async function closeQueue() {
  await taskWorker.close();
  await taskQueue.close();
  await connection.quit();
}

// Error handling
taskWorker.on('failed', (job, err) => {
  console.error(`Job ${job?.id} failed:`, err);
});

taskWorker.on('completed', (job) => {
  console.log(`Job ${job.id} completed successfully`);
});