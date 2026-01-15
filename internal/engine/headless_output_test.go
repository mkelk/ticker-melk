package engine

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mkelk/ticker-melk/internal/budget"
	"github.com/mkelk/ticker-melk/internal/ticks"
	"github.com/mkelk/ticker-melk/internal/verify"
)

func TestHeadlessOutput_Start(t *testing.T) {
	epic := &ticks.Epic{ID: "abc123", Title: "Test Epic"}

	t.Run("human readable format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.Start(epic, 50, 10.0)

		output := buf.String()
		if !strings.Contains(output, "[START]") {
			t.Error("expected [START] prefix")
		}
		if !strings.Contains(output, "abc123") {
			t.Error("expected epic ID in output")
		}
		if !strings.Contains(output, "Test Epic") {
			t.Error("expected epic title in output")
		}
	})

	t.Run("jsonl format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		out.Start(epic, 50, 10.0)

		// Parse JSON
		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if data["type"] != "start" {
			t.Errorf("expected type=start, got %v", data["type"])
		}
		if data["epic_id"] != "abc123" {
			t.Errorf("expected epic_id=abc123, got %v", data["epic_id"])
		}
	})

	t.Run("with epic prefix", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "xyz")
		out.SetWriter(&buf)

		out.Start(epic, 50, 10.0)

		output := buf.String()
		if !strings.Contains(output, "[xyz]") {
			t.Error("expected [xyz] prefix for multi-epic mode")
		}
	})
}

func TestHeadlessOutput_Task(t *testing.T) {
	task := &ticks.Task{ID: "task1", Title: "Do Something"}

	t.Run("human readable format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.Task(task, 3)

		output := buf.String()
		if !strings.Contains(output, "[TASK]") {
			t.Error("expected [TASK] prefix")
		}
		if !strings.Contains(output, "task1") {
			t.Error("expected task ID")
		}
		if !strings.Contains(output, "Do Something") {
			t.Error("expected task title")
		}
		if !strings.Contains(output, "iteration 3") {
			t.Error("expected iteration number")
		}
	})

	t.Run("jsonl format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		out.Task(task, 3)

		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if data["type"] != "task" {
			t.Errorf("expected type=task, got %v", data["type"])
		}
		if data["task_id"] != "task1" {
			t.Errorf("expected task_id=task1, got %v", data["task_id"])
		}
		if data["iteration"].(float64) != 3 {
			t.Errorf("expected iteration=3, got %v", data["iteration"])
		}
	})
}

func TestHeadlessOutput_Output(t *testing.T) {
	t.Run("human readable streams directly", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.Output("Hello world")

		// Output should be streamed without prefix
		if buf.String() != "Hello world" {
			t.Errorf("expected 'Hello world', got %q", buf.String())
		}
	})

	t.Run("jsonl wraps in output type", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		out.Output("Hello world")

		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if data["type"] != "output" {
			t.Errorf("expected type=output, got %v", data["type"])
		}
		if data["text"] != "Hello world" {
			t.Errorf("expected text='Hello world', got %v", data["text"])
		}
	})
}

func TestHeadlessOutput_Error(t *testing.T) {
	t.Run("human readable format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.Error(testError{"something went wrong"})

		output := buf.String()
		if !strings.Contains(output, "[ERROR]") {
			t.Error("expected [ERROR] prefix")
		}
		if !strings.Contains(output, "something went wrong") {
			t.Error("expected error message")
		}
	})

	t.Run("jsonl format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		out.Error(testError{"something went wrong"})

		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if data["type"] != "error" {
			t.Errorf("expected type=error, got %v", data["type"])
		}
		if data["error"] != "something went wrong" {
			t.Errorf("expected error message, got %v", data["error"])
		}
	})
}

func TestHeadlessOutput_Signal(t *testing.T) {
	tests := []struct {
		signal   Signal
		expected string
	}{
		{SignalComplete, "COMPLETE"},
		{SignalBlocked, "BLOCKED"},
		{SignalEject, "EJECT"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			var buf bytes.Buffer
			out := NewHeadlessOutput(false, "")
			out.SetWriter(&buf)

			out.Signal(tt.signal, "test reason")

			output := buf.String()
			if !strings.Contains(output, "["+tt.expected+"]") {
				t.Errorf("expected [%s] prefix in %q", tt.expected, output)
			}
			if !strings.Contains(output, "test reason") {
				t.Error("expected reason in output")
			}
		})
	}
}

func TestHeadlessOutput_Complete(t *testing.T) {
	result := &RunResult{
		EpicID:      "abc123",
		Iterations:  10,
		Duration:    5 * time.Second,
		TotalCost:   1.23,
		TotalTokens: 7000,
		ExitReason:  "all tasks completed",
		Signal:      SignalComplete,
	}

	t.Run("human readable format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.Complete(result)

		output := buf.String()
		if !strings.Contains(output, "[COMPLETE]") {
			t.Error("expected [COMPLETE] prefix")
		}
		if !strings.Contains(output, "abc123") {
			t.Error("expected epic ID")
		}
		if !strings.Contains(output, "10 iterations") {
			t.Error("expected iteration count")
		}
		if !strings.Contains(output, "$1.23") {
			t.Error("expected cost")
		}
	})

	t.Run("jsonl format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		out.Complete(result)

		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if data["type"] != "complete" {
			t.Errorf("expected type=complete, got %v", data["type"])
		}
		if data["epic_id"] != "abc123" {
			t.Errorf("expected epic_id=abc123, got %v", data["epic_id"])
		}
		if data["iterations"].(float64) != 10 {
			t.Errorf("expected iterations=10, got %v", data["iterations"])
		}
	})
}

func TestHeadlessOutput_TaskComplete(t *testing.T) {
	t.Run("passed verification", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.TaskComplete("task1", true)

		output := buf.String()
		if !strings.Contains(output, "[TASK_COMPLETE]") {
			t.Error("expected [TASK_COMPLETE] prefix")
		}
		if !strings.Contains(output, "task1") {
			t.Error("expected task ID")
		}
		if !strings.Contains(output, "closed") {
			t.Error("expected 'closed' status")
		}
	})

	t.Run("failed verification", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.TaskComplete("task1", false)

		output := buf.String()
		if !strings.Contains(output, "reopened") {
			t.Error("expected 'reopened' for failed verification")
		}
	})
}

func TestHeadlessOutput_Verify(t *testing.T) {
	t.Run("verify start", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.VerifyStart("task1")

		output := buf.String()
		if !strings.Contains(output, "[VERIFY]") {
			t.Error("expected [VERIFY] prefix")
		}
	})

	t.Run("verify end passed", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		results := &verify.Results{AllPassed: true}
		out.VerifyEnd("task1", results)

		output := buf.String()
		if !strings.Contains(output, "passed") {
			t.Error("expected 'passed' in output")
		}
	})

	t.Run("verify end failed", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		results := &verify.Results{
			AllPassed: false,
			Results: []*verify.Result{
				{Verifier: "git", Passed: false, Output: "uncommitted changes"},
			},
		}
		out.VerifyEnd("task1", results)

		output := buf.String()
		if !strings.Contains(output, "failed") {
			t.Error("expected 'failed' in output")
		}
	})
}

func TestHeadlessOutput_Interrupted(t *testing.T) {
	t.Run("human readable format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.Interrupted()

		output := buf.String()
		if !strings.Contains(output, "[INTERRUPTED]") {
			t.Error("expected [INTERRUPTED] prefix")
		}
	})

	t.Run("jsonl format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		out.Interrupted()

		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if data["type"] != "interrupted" {
			t.Errorf("expected type=interrupted, got %v", data["type"])
		}
	})
}

func TestHeadlessOutput_ProjectSummary(t *testing.T) {
	pb := &budget.ProjectBudget{
		Project:    "2026-01-14-6453-auth",
		Iterations: 47,
		Tokens:     1200000,
		Cost:       12.34,
	}

	t.Run("human readable format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.ProjectSummary(pb)

		output := buf.String()
		if !strings.Contains(output, "[PROJECT]") {
			t.Error("expected [PROJECT] prefix")
		}
		if !strings.Contains(output, "2026-01-14-6453-auth") {
			t.Error("expected project name")
		}
		if !strings.Contains(output, "47") {
			t.Error("expected iterations")
		}
		if !strings.Contains(output, "1.2M") {
			t.Error("expected formatted tokens")
		}
		if !strings.Contains(output, "$12.34") {
			t.Error("expected cost")
		}
	})

	t.Run("jsonl format", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		out.ProjectSummary(pb)

		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if data["type"] != "project_summary" {
			t.Errorf("expected type=project_summary, got %v", data["type"])
		}
		if data["project"] != "2026-01-14-6453-auth" {
			t.Errorf("expected project=2026-01-14-6453-auth, got %v", data["project"])
		}
		if data["iterations"].(float64) != 47 {
			t.Errorf("expected iterations=47, got %v", data["iterations"])
		}
		if data["tokens"].(float64) != 1200000 {
			t.Errorf("expected tokens=1200000, got %v", data["tokens"])
		}
		if data["cost"].(float64) != 12.34 {
			t.Errorf("expected cost=12.34, got %v", data["cost"])
		}
	})

	t.Run("nil budget does not output", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.ProjectSummary(nil)

		if buf.Len() > 0 {
			t.Error("expected no output for nil budget")
		}
	})

	t.Run("empty project does not output", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(false, "")
		out.SetWriter(&buf)

		out.ProjectSummary(&budget.ProjectBudget{Project: ""})

		if buf.Len() > 0 {
			t.Error("expected no output for empty project")
		}
	})
}

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		tokens   int
		expected string
	}{
		{500, "500"},
		{1000, "1.0K"},
		{1500, "1.5K"},
		{50000, "50.0K"},
		{1000000, "1.0M"},
		{1200000, "1.2M"},
		{2500000, "2.5M"},
	}

	for _, tt := range tests {
		result := formatTokens(tt.tokens)
		if result != tt.expected {
			t.Errorf("formatTokens(%d) = %s, expected %s", tt.tokens, result, tt.expected)
		}
	}
}

func TestHeadlessOutput_Complete_WithProject(t *testing.T) {
	result := &RunResult{
		EpicID:      "abc123",
		Project:     "test-project",
		Iterations:  10,
		Duration:    5 * time.Second,
		TotalCost:   1.23,
		TotalTokens: 7000,
		ExitReason:  "all tasks completed",
		Signal:      SignalComplete,
	}

	t.Run("jsonl includes project field", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		out.Complete(result)

		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if data["project"] != "test-project" {
			t.Errorf("expected project=test-project, got %v", data["project"])
		}
	})

	t.Run("jsonl omits empty project", func(t *testing.T) {
		var buf bytes.Buffer
		out := NewHeadlessOutput(true, "")
		out.SetWriter(&buf)

		resultNoProject := &RunResult{
			EpicID:      "abc123",
			Iterations:  10,
			Duration:    5 * time.Second,
			TotalCost:   1.23,
			TotalTokens: 7000,
			ExitReason:  "all tasks completed",
			Signal:      SignalComplete,
		}
		out.Complete(resultNoProject)

		var data map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if _, ok := data["project"]; ok {
			t.Error("expected no project field for empty project")
		}
	})
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e testError) Error() string {
	return e.msg
}
