#!/usr/bin/env node
// ABOUTME: Development script that starts server in background then launches TUI
// ABOUTME: Ensures TUI has full terminal control while server runs separately

import { spawn } from 'child_process';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __dirname = dirname(fileURLToPath(import.meta.url));
const rootDir = join(__dirname, '..');

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
    serverProcess.kill();
    tuiProcess.kill();
    process.exit(0);
  });

  process.on('SIGTERM', () => {
    serverProcess.kill();
    tuiProcess.kill();
    process.exit(0);
  });
}, 2000); // Give server 2 seconds to start