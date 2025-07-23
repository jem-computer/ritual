// ABOUTME: Main entry point for the Ritual server
// ABOUTME: Handles HTTP API and task scheduling

import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { serve } from "@hono/node-server";

const app = new Hono();

// Middleware
app.use("/*", cors());
app.use(logger());

// Mock data store (in-memory for now)
const mockTasks = [
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

// Health check
app.get("/health", (c) => {
	return c.json({ status: "ok" });
});

// Task routes
app.get("/api/tasks", (c) => {
	return c.json(mockTasks);
});

app.post("/api/tasks", async (c) => {
	const body = await c.req.json();
	const newTask = {
		id: String(mockTasks.length + 1),
		...body,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
	};
	mockTasks.push(newTask);
	return c.json(newTask, 201);
});

app.put("/api/tasks/:id", async (c) => {
	const id = c.req.param("id");
	const body = await c.req.json();

	const taskIndex = mockTasks.findIndex((t) => t.id === id);
	if (taskIndex === -1) {
		return c.json({ error: "Task not found" }, 404);
	}

	mockTasks[taskIndex] = {
		...mockTasks[taskIndex],
		...body,
		id, // Preserve the ID
		updatedAt: new Date().toISOString(),
	};

	return c.json(mockTasks[taskIndex]);
});

app.delete("/api/tasks/:id", (c) => {
	const id = c.req.param("id");

	const taskIndex = mockTasks.findIndex((t) => t.id === id);
	if (taskIndex === -1) {
		return c.json({ error: "Task not found" }, 404);
	}

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
console.log(`Server running on port ${port}`);

serve({
	fetch: app.fetch,
	port: Number(port),
});

