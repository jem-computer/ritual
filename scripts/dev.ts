#!/usr/bin/env node
// ABOUTME: Development script that starts Redis, server in background then launches TUI
// ABOUTME: Ensures TUI has full terminal control while server runs separately

import { spawn, execSync } from 'child_process';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __dirname = dirname(fileURLToPath(import.meta.url));
const rootDir = join(__dirname, '..');

// Check if Redis is already running
function isRedisRunning() {
  try {
    execSync('redis-cli ping', { stdio: 'ignore' });
    return true;
  } catch {
    return false;
  }
}

// Start Redis if not already running
let redisProcess = null;
if (!isRedisRunning()) {
  console.log('Starting Redis...');
  redisProcess = spawn('redis-server', [], {
    stdio: ['ignore', 'pipe', 'pipe'],
    detached: false
  });
  
  redisProcess.stdout.on('data', (data) => {
    if (process.env.DEBUG) {
      console.log('[REDIS]', data.toString().trim());
    }
  });
  
  redisProcess.stderr.on('data', (data) => {
    console.error('[REDIS ERROR]', data.toString());
  });
  
  // Wait for Redis to be ready
  let redisReady = false;
  let attempts = 0;
  while (!redisReady && attempts < 10) {
    try {
      execSync('redis-cli ping', { stdio: 'ignore' });
      redisReady = true;
      console.log('Redis is ready');
    } catch {
      attempts++;
      execSync('sleep 0.5');
    }
  }
  
  if (!redisReady) {
    console.error('Failed to start Redis');
    process.exit(1);
  }
} else {
  console.log('Redis is already running');
}

// Start server in background
console.log('Starting server...');
const serverProcess = spawn('bun', ['run', 'dev'], {
  cwd: join(rootDir, 'packages/server'),
  stdio: ['ignore', 'pipe', 'pipe'],
  detached: false
});

// Buffer server output
let serverReady = false;
serverProcess.stdout.on('data', (data) => {
  const output = data.toString();
  if (output.includes('Server running') || output.includes('localhost:8080')) {
    serverReady = true;
  }
  // Optionally log server output to a file or debug mode
  if (process.env.DEBUG) {
    console.log('[SERVER]', output.trim());
  }
});

serverProcess.stderr.on('data', (data) => {
  console.error('[SERVER ERROR]', data.toString());
});

// Wait a bit for server to start, then launch TUI
setTimeout(() => {
  console.log('Starting TUI...');
  const tuiProcess = spawn('go', ['run', 'cmd/ritual/main.go', '--server', 'http://localhost:8080'], {
    cwd: join(rootDir, 'packages/tui'),
    stdio: 'inherit' // This gives TUI full control of the terminal
  });

  tuiProcess.on('exit', (code) => {
    console.log('TUI exited with code', code);
    serverProcess.kill();
    process.exit(code);
  });

  // Handle cleanup
  process.on('SIGINT', () => {
    console.log('\nShutting down...');
    tuiProcess.kill();
    serverProcess.kill();
    if (redisProcess) {
      redisProcess.kill();
    }
    process.exit(0);
  });

  process.on('SIGTERM', () => {
    tuiProcess.kill();
    serverProcess.kill();
    if (redisProcess) {
      redisProcess.kill();
    }
    process.exit(0);
  });
}, 2000); // Give server 2 seconds to start