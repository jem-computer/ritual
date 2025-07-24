// ABOUTME: Main entry point for the Ritual server
// ABOUTME: Handles HTTP API and task scheduling

import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { serve } from "@hono/node-server";
import { taskQueue, closeQueue, isRedisConnected } from "./queue.js";
import { parseSchedule } from "./schedule-parser.js";

const app = new Hono();

// Middleware
app.use("/*", cors());
app.use(logger());

// Types
interface Task {
	id: string;
	name: string;
	status: "ACTIVE" | "PAUSED";
	prompt: string;
	schedule: string;
	output: string;
	model: string;
	nextRun: string | null;
	lastRun: string | null;
	createdAt: string;
	updatedAt: string;
	jobId?: string; // BullMQ job ID for scheduled tasks
}

// Mock data store (in-memory for now)
const mockTasks: Task[] = [
	{
		id: "1",
		name: "Daily Task Summary",
		status: "ACTIVE",
		prompt: "Summarize today's Things tasks in iambic pentameter",
		schedule: "daily at 8:00 AM",
		output: "SMS to +1234567890",
		model: "gpt-4",
		nextRun: new Date(Date.now() + 12 * 60 * 60 * 1000).toISOString(),
		lastRun: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
		createdAt: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
		updatedAt: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
	},
	{
		id: "2",
		name: "Team Commit Summary",
		status: "ACTIVE",
		prompt: "Send my team a summary of today's commits",
		schedule: "daily at 6:00 PM",
		output: "Slack #dev-team",
		model: "gpt-3.5-turbo",
		nextRun: new Date(Date.now() + 6 * 60 * 60 * 1000).toISOString(),
		lastRun: new Date(Date.now() - 18 * 60 * 60 * 1000).toISOString(),
		createdAt: new Date(Date.now() - 14 * 24 * 60 * 60 * 1000).toISOString(),
		updatedAt: new Date(Date.now() - 18 * 60 * 60 * 1000).toISOString(),
	},
	{
		id: "3",
		name: "Weekly Report",
		status: "PAUSED",
		prompt:
			"Generate a weekly productivity report based on my calendar and tasks",
		schedule: "weekly on Friday at 5:00 PM",
		output: "Email to me@example.com",
		model: "gpt-4",
		nextRun: new Date(Date.now() + 72 * 60 * 60 * 1000).toISOString(),
		lastRun: null,
		createdAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
		updatedAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
	},
];

// Helper functions for job management
async function scheduleTask(task: Task) {
	if (task.status !== "ACTIVE") {
		return;
	}

	if (!isRedisConnected()) {
		console.warn(`Cannot schedule task ${task.id}: Redis not connected`);
		return;
	}

	try {
		const { cron } = parseSchedule(task.schedule);
		
		// Remove existing job if it exists
		if (task.jobId) {
			await taskQueue.removeRepeatableByKey(task.jobId);
		}

		// Schedule new job
		const job = await taskQueue.add(
			`task-${task.id}`,
			{
				taskId: task.id,
				prompt: task.prompt,
				outputChannels: [task.output],
				model: task.model,
			},
			{
				repeat: {
					pattern: cron,
				},
				repeatJobKey: task.id, // Use task ID as the repeat key
			}
		);

		// Store the repeat key for later removal
		task.jobId = task.id;
		console.log(`Scheduled task ${task.id} with cron: ${cron}`);
	} catch (error) {
		console.error(`Failed to schedule task ${task.id}:`, error);
	}
}

async function unscheduleTask(task: Task) {
	if (task.jobId) {
		try {
			// Remove all repeatable jobs for this task
			const repeatableJobs = await taskQueue.getRepeatableJobs();
			for (const job of repeatableJobs) {
				if (job.key.includes(`task-${task.id}`)) {
					await taskQueue.removeRepeatableByKey(job.key);
				}
			}
			console.log(`Unscheduled task ${task.id}`);
		} catch (error) {
			console.error(`Failed to unschedule task ${task.id}:`, error);
		}
	}
}

// Initialize scheduled tasks on startup
async function initializeScheduledTasks() {
	for (const task of mockTasks) {
		if (task.status === "ACTIVE") {
			await scheduleTask(task);
		}
	}
}

// Health check
app.get("/health", (c) => {
	return c.json({ 
		status: "ok",
		redis: isRedisConnected() ? "connected" : "disconnected"
	});
});

// Task routes
app.get("/api/tasks", (c) => {
	return c.json(mockTasks);
});

app.post("/api/tasks", async (c) => {
	const body = await c.req.json();
	const newTask: Task = {
		id: String(mockTasks.length + 1),
		...body,
		nextRun: null,
		lastRun: null,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
	};
	
	mockTasks.push(newTask);
	
	// Schedule the task if it's active
	if (newTask.status === "ACTIVE") {
		await scheduleTask(newTask);
	}
	
	return c.json(newTask, 201);
});

app.put("/api/tasks/:id", async (c) => {
	const id = c.req.param("id");
	const body = await c.req.json();

	const taskIndex = mockTasks.findIndex((t) => t.id === id);
	if (taskIndex === -1) {
		return c.json({ error: "Task not found" }, 404);
	}

	const oldTask = mockTasks[taskIndex];
	const updatedTask = {
		...oldTask,
		...body,
		id, // Preserve the ID
		updatedAt: new Date().toISOString(),
	};

	// Handle scheduling changes
	if (oldTask.status !== updatedTask.status || oldTask.schedule !== updatedTask.schedule) {
		// Unschedule old job
		await unscheduleTask(oldTask);
		
		// Schedule new job if active
		if (updatedTask.status === "ACTIVE") {
			await scheduleTask(updatedTask);
		}
	}

	mockTasks[taskIndex] = updatedTask;
	return c.json(updatedTask);
});

app.delete("/api/tasks/:id", async (c) => {
	const id = c.req.param("id");

	const taskIndex = mockTasks.findIndex((t) => t.id === id);
	if (taskIndex === -1) {
		return c.json({ error: "Task not found" }, 404);
	}

	const task = mockTasks[taskIndex];
	
	// Unschedule the task
	await unscheduleTask(task);
	
	mockTasks.splice(taskIndex, 1);
	return c.body(null, 204);
});

// Mock logs
const mockLogs = [
	{
		id: "1",
		taskId: "1",
		taskName: "Daily Task Summary",
		output:
			"Today's tasks in verse, a summary terse:\nTwo bugs were fixed with careful thought,\nThree features planned, though time was short.\nThe standup call ran fifteen long,\nBut progress made was rather strong.",
		status: "SUCCESS",
		executedAt: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
	},
	{
		id: "2",
		taskId: "2",
		taskName: "Team Commit Summary",
		output:
			"Team commits today:\n- Fixed authentication bug in login flow\n- Added dark mode support to dashboard\n- Updated dependencies to latest versions\n- Improved error handling in API client",
		status: "SUCCESS",
		executedAt: new Date(Date.now() - 18 * 60 * 60 * 1000).toISOString(),
	},
];

// Log routes
app.get("/api/logs", (c) => {
	return c.json(mockLogs);
});

// Start server
const port = process.env.PORT || 8080;

async function startServer() {
	try {
		// Initialize scheduled tasks
		await initializeScheduledTasks();
		
		console.log(`Server running on port ${port}`);
		
		serve({
			fetch: app.fetch,
			port: Number(port),
		});
	} catch (error) {
		console.error("Failed to start server:", error);
		process.exit(1);
	}
}

// Graceful shutdown
process.on("SIGTERM", async () => {
	console.log("SIGTERM received, shutting down gracefully...");
	await closeQueue();
	process.exit(0);
});

process.on("SIGINT", async () => {
	console.log("SIGINT received, shutting down gracefully...");
	await closeQueue();
	process.exit(0);
});

startServer();

