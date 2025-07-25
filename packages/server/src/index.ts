// ABOUTME: Main entry point for the Ritual server
// ABOUTME: Handles HTTP API and task scheduling

import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { serve } from "@hono/node-server";
import { taskQueue, closeQueue, isRedisConnected } from "./queue.js";
import { parseSchedule } from "./schedule-parser.js";
import {
	initDatabase,
	closeDatabase,
	getAllTasks,
	getTaskById,
	createTask,
	updateTask,
	deleteTask,
	getAllExecutionLogs,
	type Task,
	type ExecutionLog,
} from "./db.js";
import { initAIService } from "./ai-service.js";

const app = new Hono();

// Middleware
app.use("/*", cors());
app.use(logger());


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
	const tasks = await getAllTasks();
	for (const task of tasks) {
		if (task.status === "ACTIVE") {
			await scheduleTask(task);
		}
	}
}

// Health check
app.get("/health", (c) => {
	return c.json({ 
		status: "ok",
		redis: isRedisConnected() ? "connected" : "disconnected",
		database: "connected"
	});
});

// Task routes
app.get("/api/tasks", async (c) => {
	const tasks = await getAllTasks();
	return c.json(tasks);
});

app.post("/api/tasks", async (c) => {
	const body = await c.req.json();
	
	const newTask = await createTask({
		...body,
		nextRun: null,
		lastRun: null,
	});
	
	// Schedule the task if it's active
	if (newTask.status === "ACTIVE") {
		await scheduleTask(newTask);
	}
	
	return c.json(newTask, 201);
});

app.put("/api/tasks/:id", async (c) => {
	const id = c.req.param("id");
	const body = await c.req.json();

	const oldTask = await getTaskById(id);
	if (!oldTask) {
		return c.json({ error: "Task not found" }, 404);
	}

	// Handle scheduling changes
	if (body.status !== undefined || body.schedule !== undefined) {
		const willBeActive = body.status !== undefined ? body.status === "ACTIVE" : oldTask.status === "ACTIVE";
		const scheduleChanged = body.schedule !== undefined && body.schedule !== oldTask.schedule;
		
		// Unschedule old job if needed
		if (oldTask.status === "ACTIVE" && (!willBeActive || scheduleChanged)) {
			await unscheduleTask(oldTask);
		}
		
		// Schedule new job if needed
		if (willBeActive && (oldTask.status !== "ACTIVE" || scheduleChanged)) {
			const taskToSchedule = { ...oldTask, ...body };
			await scheduleTask(taskToSchedule);
		}
	}

	const updatedTask = await updateTask(id, body);
	if (!updatedTask) {
		return c.json({ error: "Failed to update task" }, 500);
	}
	
	return c.json(updatedTask);
});

app.delete("/api/tasks/:id", async (c) => {
	const id = c.req.param("id");

	const task = await getTaskById(id);
	if (!task) {
		return c.json({ error: "Task not found" }, 404);
	}
	
	// Unschedule the task
	await unscheduleTask(task);
	
	const deleted = await deleteTask(id);
	if (!deleted) {
		return c.json({ error: "Failed to delete task" }, 500);
	}
	
	return c.body(null, 204);
});

// Log routes
app.get("/api/logs", async (c) => {
	const logs = await getAllExecutionLogs();
	return c.json(logs);
});

// Start server
const port = process.env.PORT || 8080;

async function startServer() {
	try {
		// Initialize database
		await initDatabase();
		
		// Initialize AI service
		initAIService();
		
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
	await closeDatabase();
	process.exit(0);
});

process.on("SIGINT", async () => {
	console.log("SIGINT received, shutting down gracefully...");
	await closeQueue();
	await closeDatabase();
	process.exit(0);
});

startServer();

