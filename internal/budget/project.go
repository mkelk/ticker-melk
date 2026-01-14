package budget

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ProjectBudget tracks cumulative metrics for a project across multiple runs.
type ProjectBudget struct {
	// Project is the project identifier (e.g., "2026-01-14-6453-auth").
	Project string `json:"project"`

	// Iterations is the total number of iterations across all runs.
	Iterations int `json:"iterations"`

	// Tokens is the total token count (input + output).
	Tokens int `json:"tokens"`

	// Cost is the cumulative cost in USD.
	Cost float64 `json:"cost"`

	// LastUpdated is when this record was last updated.
	LastUpdated time.Time `json:"last_updated"`
}

// ProjectStore manages persistent storage of project budget data.
type ProjectStore struct {
	dir     string
	file    string
	budgets map[string]*ProjectBudget
	mu      sync.RWMutex
}

// NewProjectStore creates a new project store with the default directory.
func NewProjectStore() *ProjectStore {
	return &ProjectStore{
		dir:     ".ticker",
		file:    "projects.json",
		budgets: make(map[string]*ProjectBudget),
	}
}

// NewProjectStoreWithDir creates a project store with a custom directory.
func NewProjectStoreWithDir(dir string) *ProjectStore {
	return &ProjectStore{
		dir:     dir,
		file:    "projects.json",
		budgets: make(map[string]*ProjectBudget),
	}
}

// filePath returns the full path to the projects.json file.
func (s *ProjectStore) filePath() string {
	return filepath.Join(s.dir, s.file)
}

// Load reads project budgets from disk.
func (s *ProjectStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath())
	if err != nil {
		if os.IsNotExist(err) {
			// No file yet, start fresh
			s.budgets = make(map[string]*ProjectBudget)
			return nil
		}
		return fmt.Errorf("reading project budgets: %w", err)
	}

	var budgets []*ProjectBudget
	if err := json.Unmarshal(data, &budgets); err != nil {
		return fmt.Errorf("unmarshaling project budgets: %w", err)
	}

	s.budgets = make(map[string]*ProjectBudget)
	for _, b := range budgets {
		s.budgets[b.Project] = b
	}

	return nil
}

// Save writes project budgets to disk.
func (s *ProjectStore) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Ensure directory exists
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("creating project store directory: %w", err)
	}

	// Convert map to slice for JSON
	budgets := make([]*ProjectBudget, 0, len(s.budgets))
	for _, b := range s.budgets {
		budgets = append(budgets, b)
	}

	data, err := json.MarshalIndent(budgets, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling project budgets: %w", err)
	}

	if err := os.WriteFile(s.filePath(), data, 0644); err != nil {
		return fmt.Errorf("writing project budgets: %w", err)
	}

	return nil
}

// Add accumulates usage for a project.
// Creates a new entry if the project doesn't exist.
func (s *ProjectStore) Add(project string, iterations, tokens int, cost float64) {
	if project == "" {
		return // Ignore entries without a project
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	budget, ok := s.budgets[project]
	if !ok {
		budget = &ProjectBudget{Project: project}
		s.budgets[project] = budget
	}

	budget.Iterations += iterations
	budget.Tokens += tokens
	budget.Cost += cost
	budget.LastUpdated = time.Now()
}

// Get returns the budget for a project, or nil if not found.
func (s *ProjectStore) Get(project string) *ProjectBudget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	budget, ok := s.budgets[project]
	if !ok {
		return nil
	}

	// Return a copy
	copy := *budget
	return &copy
}

// All returns all project budgets.
func (s *ProjectStore) All() []*ProjectBudget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*ProjectBudget, 0, len(s.budgets))
	for _, b := range s.budgets {
		copy := *b
		result = append(result, &copy)
	}

	return result
}

// FormatSummary returns a human-readable summary for a project.
func (b *ProjectBudget) FormatSummary() string {
	return fmt.Sprintf("Project: %s\n  Iterations: %d\n  Tokens: %s\n  Cost: $%.2f",
		b.Project,
		b.Iterations,
		formatTokens(b.Tokens),
		b.Cost,
	)
}

// formatTokens formats token count with K/M suffixes.
func formatTokens(tokens int) string {
	if tokens >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(tokens)/1_000_000)
	}
	if tokens >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(tokens)/1_000)
	}
	return fmt.Sprintf("%d", tokens)
}
