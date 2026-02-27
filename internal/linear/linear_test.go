package linear_test

import (
	"testing"

	"github.com/jasonwashburn/linear-cli/internal/linear"
)

func TestPriorityLabel(t *testing.T) {
	tests := []struct {
		priority int
		want     string
	}{
		{0, "No priority"},
		{1, "Urgent"},
		{2, "High"},
		{3, "Medium"},
		{4, "Low"},
		{99, "Unknown"},
	}

	for _, tt := range tests {
		got := linear.PriorityLabel(tt.priority)
		if got != tt.want {
			t.Errorf("PriorityLabel(%d) = %q, want %q", tt.priority, got, tt.want)
		}
	}
}

func TestNewClient(t *testing.T) {
	client := linear.NewClient("test-api-key")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestBuildCreateInput(t *testing.T) {
	priority := 2
	input := linear.IssueCreateInput{
		Title:       "Test Issue",
		TeamID:      "team-123",
		Description: "A description",
		Priority:    &priority,
		AssigneeID:  "user-456",
		ParentID:    "parent-789",
	}

	// Verify that field values are correctly populated
	if input.Title != "Test Issue" {
		t.Errorf("expected Title %q, got %q", "Test Issue", input.Title)
	}
	if input.TeamID != "team-123" {
		t.Errorf("expected TeamID %q, got %q", "team-123", input.TeamID)
	}
	if *input.Priority != 2 {
		t.Errorf("expected Priority 2, got %d", *input.Priority)
	}
	if input.ParentID != "parent-789" {
		t.Errorf("expected ParentID %q, got %q", "parent-789", input.ParentID)
	}
}

func TestBuildUpdateInput(t *testing.T) {
	title := "Updated Title"
	priority := 1
	input := linear.IssueUpdateInput{
		Title:    &title,
		Priority: &priority,
	}

	if *input.Title != "Updated Title" {
		t.Errorf("expected Title %q, got %q", "Updated Title", *input.Title)
	}
	if *input.Priority != 1 {
		t.Errorf("expected Priority 1, got %d", *input.Priority)
	}
	if input.Description != nil {
		t.Errorf("expected Description nil, got %q", *input.Description)
	}
}
