// ABOUTME: API client for communicating with the Ritual server
// ABOUTME: Handles HTTP requests and SSE streaming for real-time updates

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Task represents a scheduled ritual task
type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Prompt    string    `json:"prompt"`
	Schedule  string    `json:"schedule"`
	Model     string    `json:"model"`
	Output    string    `json:"output"`
	Status    string    `json:"status"`
	NextRun   time.Time `json:"nextRun"`
	LastRun   time.Time `json:"lastRun"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LogEntry represents an execution log entry
type LogEntry struct {
	ID         string    `json:"id"`
	TaskID     string    `json:"taskId"`
	TaskName   string    `json:"taskName"`
	Output     string    `json:"output"`
	Status     string    `json:"status"`
	Error      string    `json:"error,omitempty"`
	ExecutedAt time.Time `json:"executedAt"`
}

// GetTasks retrieves all tasks
func (c *Client) GetTasks() ([]Task, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/tasks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// CreateTask creates a new task
func (c *Client) CreateTask(task Task) (*Task, error) {
	body, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(c.baseURL+"/api/tasks", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var created Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateTask updates an existing task
func (c *Client) UpdateTask(id string, task Task) (*Task, error) {
	body, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/api/tasks/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var updated Task
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteTask deletes a task
func (c *Client) DeleteTask(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/api/tasks/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetLogs retrieves execution logs
func (c *Client) GetLogs() ([]LogEntry, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/logs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var logs []LogEntry
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, err
	}

	return logs, nil
}
