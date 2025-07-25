// ABOUTME: Seeds the database with sample tasks for development
// ABOUTME: Run with: bun run src/seed.ts

import { initDatabase, createTask, closeDatabase } from './db.js';

async function seedDatabase() {
  console.log('Seeding database...');
  
  try {
    await initDatabase();
    
    // Sample tasks
    const sampleTasks = [
      {
        name: "Daily Task Summary",
        status: "ACTIVE" as const,
        prompt: "Summarize today's Things tasks in iambic pentameter",
        schedule: "daily at 8:00 AM",
        output: "SMS to +1234567890",
        model: "claude-3-5-sonnet-20241022",
        nextRun: new Date(Date.now() + 12 * 60 * 60 * 1000).toISOString(),
        lastRun: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
      },
      {
        name: "Team Commit Summary",
        status: "ACTIVE" as const,
        prompt: "Send my team a summary of today's commits",
        schedule: "daily at 6:00 PM",
        output: "Slack #dev-team",
        model: "claude-3-5-haiku-20241022",
        nextRun: new Date(Date.now() + 6 * 60 * 60 * 1000).toISOString(),
        lastRun: new Date(Date.now() - 18 * 60 * 60 * 1000).toISOString(),
      },
      {
        name: "Weekly Report",
        status: "PAUSED" as const,
        prompt: "Generate a weekly productivity report based on my calendar and tasks",
        schedule: "weekly on Friday at 5:00 PM",
        output: "Email to me@example.com",
        model: "claude-3-5-sonnet-20241022",
        nextRun: new Date(Date.now() + 72 * 60 * 60 * 1000).toISOString(),
        lastRun: null,
      },
    ];
    
    for (const task of sampleTasks) {
      const created = await createTask(task);
      console.log(`Created task: ${created.name}`);
    }
    
    console.log('Database seeded successfully!');
  } catch (error) {
    console.error('Failed to seed database:', error);
    process.exit(1);
  } finally {
    await closeDatabase();
  }
}

// Run if called directly
if (import.meta.main) {
  seedDatabase();
}