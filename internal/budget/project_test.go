package budget

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectStore_AddAndGet(t *testing.T) {
	store := NewProjectStore()

	// Add some usage
	store.Add("project-a", 5, 10000, 1.50)
	store.Add("project-a", 3, 5000, 0.75)
	store.Add("project-b", 2, 3000, 0.50)

	// Check project A
	budgetA := store.Get("project-a")
	if budgetA == nil {
		t.Fatal("expected budget for project-a")
	}
	if budgetA.Iterations != 8 {
		t.Errorf("expected 8 iterations, got %d", budgetA.Iterations)
	}
	if budgetA.Tokens != 15000 {
		t.Errorf("expected 15000 tokens, got %d", budgetA.Tokens)
	}
	if budgetA.Cost != 2.25 {
		t.Errorf("expected cost 2.25, got %.2f", budgetA.Cost)
	}

	// Check project B
	budgetB := store.Get("project-b")
	if budgetB == nil {
		t.Fatal("expected budget for project-b")
	}
	if budgetB.Iterations != 2 {
		t.Errorf("expected 2 iterations, got %d", budgetB.Iterations)
	}

	// Check non-existent project
	budgetC := store.Get("project-c")
	if budgetC != nil {
		t.Error("expected nil for non-existent project")
	}
}

func TestProjectStore_EmptyProject(t *testing.T) {
	store := NewProjectStore()

	// Adding with empty project should be ignored
	store.Add("", 5, 10000, 1.50)

	all := store.All()
	if len(all) != 0 {
		t.Errorf("expected 0 projects, got %d", len(all))
	}
}

func TestProjectStore_SaveLoad(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create store and add data
	store := NewProjectStoreWithDir(tmpDir)
	store.Add("test-project", 10, 50000, 5.00)
	store.Add("another-project", 3, 15000, 1.50)

	// Save
	if err := store.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Verify file exists
	filePath := filepath.Join(tmpDir, "projects.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("projects.json not created")
	}

	// Create new store and load
	store2 := NewProjectStoreWithDir(tmpDir)
	if err := store2.Load(); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// Verify data loaded correctly
	budget := store2.Get("test-project")
	if budget == nil {
		t.Fatal("expected budget for test-project after load")
	}
	if budget.Iterations != 10 {
		t.Errorf("expected 10 iterations after load, got %d", budget.Iterations)
	}
	if budget.Tokens != 50000 {
		t.Errorf("expected 50000 tokens after load, got %d", budget.Tokens)
	}
	if budget.Cost != 5.00 {
		t.Errorf("expected cost 5.00 after load, got %.2f", budget.Cost)
	}
}

func TestProjectStore_LoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewProjectStoreWithDir(tmpDir)

	// Load from non-existent file should succeed (start fresh)
	if err := store.Load(); err != nil {
		t.Fatalf("load from non-existent file failed: %v", err)
	}

	all := store.All()
	if len(all) != 0 {
		t.Errorf("expected empty store, got %d projects", len(all))
	}
}

func TestProjectStore_All(t *testing.T) {
	store := NewProjectStore()
	store.Add("project-1", 1, 1000, 0.10)
	store.Add("project-2", 2, 2000, 0.20)
	store.Add("project-3", 3, 3000, 0.30)

	all := store.All()
	if len(all) != 3 {
		t.Errorf("expected 3 projects, got %d", len(all))
	}

	// Verify all projects are present
	found := make(map[string]bool)
	for _, b := range all {
		found[b.Project] = true
	}
	for _, p := range []string{"project-1", "project-2", "project-3"} {
		if !found[p] {
			t.Errorf("project %s not found in All()", p)
		}
	}
}

func TestProjectBudget_FormatSummary(t *testing.T) {
	tests := []struct {
		name     string
		budget   ProjectBudget
		contains []string
	}{
		{
			name: "small tokens",
			budget: ProjectBudget{
				Project:    "test-project",
				Iterations: 5,
				Tokens:     500,
				Cost:       0.50,
			},
			contains: []string{"test-project", "5", "500", "$0.50"},
		},
		{
			name: "thousands of tokens",
			budget: ProjectBudget{
				Project:    "big-project",
				Iterations: 47,
				Tokens:     1_200_000,
				Cost:       12.34,
			},
			contains: []string{"big-project", "47", "1.2M", "$12.34"},
		},
		{
			name: "K tokens",
			budget: ProjectBudget{
				Project:    "medium-project",
				Iterations: 10,
				Tokens:     50_000,
				Cost:       5.00,
			},
			contains: []string{"medium-project", "10", "50.0K", "$5.00"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := tt.budget.FormatSummary()
			for _, s := range tt.contains {
				if !contains(summary, s) {
					t.Errorf("expected summary to contain %q, got: %s", s, summary)
				}
			}
		})
	}
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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
