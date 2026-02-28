package issue_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jasonwashburn/linear-cli/cmd/issue"
)

func executeCommand(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	cmd := issue.NewIssueCmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestIssueList_NoFlags(t *testing.T) {
	out, err := executeCommand("list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `Would list issues (team="", state="", assignee="")`
	if !contains(out, want) {
		t.Errorf("got %q, want output to contain %q", out, want)
	}
}

func TestIssueList_WithFlags(t *testing.T) {
	out, err := executeCommand("list", "--team", "eng", "--state", "In Progress", "--assignee", "bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `Would list issues (team="eng", state="In Progress", assignee="bob")`
	if !contains(out, want) {
		t.Errorf("got %q, want output to contain %q", out, want)
	}
}

func TestIssueGet(t *testing.T) {
	out, err := executeCommand("get", "ENG-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `Would get issue "ENG-123"`
	if !contains(out, want) {
		t.Errorf("got %q, want output to contain %q", out, want)
	}
}

func TestIssueGet_MissingArg(t *testing.T) {
	_, err := executeCommand("get")
	if err == nil {
		t.Error("expected error when issue-id argument is missing")
	}
}

func TestIssueCreate_RequiredTitle(t *testing.T) {
	out, err := executeCommand("create", "--title", "My Issue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `Would create issue (title="My Issue", team="", description="", priority=0, assignee="")`
	if !contains(out, want) {
		t.Errorf("got %q, want output to contain %q", out, want)
	}
}

func TestIssueCreate_AllFlags(t *testing.T) {
	out, err := executeCommand("create", "--title", "Bug", "--team", "backend", "--description", "desc", "--priority", "2", "--assignee", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `Would create issue (title="Bug", team="backend", description="desc", priority=2, assignee="alice")`
	if !contains(out, want) {
		t.Errorf("got %q, want output to contain %q", out, want)
	}
}

func TestIssueCreate_MissingTitle(t *testing.T) {
	_, err := executeCommand("create")
	if err == nil {
		t.Error("expected error when --title flag is missing")
	}
}

func TestIssueUpdate(t *testing.T) {
	out, err := executeCommand("update", "ENG-456", "--title", "New Title", "--state", "Done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `Would update issue "ENG-456" (title="New Title", description="", state="Done", priority=0, assignee="")`
	if !contains(out, want) {
		t.Errorf("got %q, want output to contain %q", out, want)
	}
}

func TestIssueUpdate_MissingArg(t *testing.T) {
	_, err := executeCommand("update")
	if err == nil {
		t.Error("expected error when issue-id argument is missing")
	}
}

func TestIssueDelete(t *testing.T) {
	out, err := executeCommand("delete", "ENG-789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `Would delete issue "ENG-789"`
	if !contains(out, want) {
		t.Errorf("got %q, want output to contain %q", out, want)
	}
}

func TestIssueDelete_MissingArg(t *testing.T) {
	_, err := executeCommand("delete")
	if err == nil {
		t.Error("expected error when issue-id argument is missing")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
