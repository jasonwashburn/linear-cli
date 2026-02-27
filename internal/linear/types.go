package linear

// Issue represents a Linear issue.
type Issue struct {
	ID          string
	Identifier  string
	Title       string
	Description string
	State       IssueState
	Priority    int
	Assignee    *User
	Team        Team
	CreatedAt   string
	UpdatedAt   string
	URL         string
}

// IssueState represents the state of an issue.
type IssueState struct {
	ID   string
	Name string
	Type string
}

// User represents a Linear user.
type User struct {
	ID          string
	Name        string
	DisplayName string
	Email       string
}

// Team represents a Linear team.
type Team struct {
	ID   string
	Name string
	Key  string
}

// IssueCreateInput holds parameters for creating an issue.
type IssueCreateInput struct {
	Title       string
	TeamID      string
	Description string
	Priority    *int
	AssigneeID  string
	ParentID    string
}

// IssueUpdateInput holds parameters for updating an issue.
type IssueUpdateInput struct {
	Title       *string
	Description *string
	StateID     *string
	Priority    *int
	AssigneeID  *string
}

// IssueFilter holds filter parameters for listing issues.
type IssueFilter struct {
	TeamID     string
	State      string
	AssigneeID string
}

// PriorityLabel returns a human-readable label for a priority value.
func PriorityLabel(p int) string {
	switch p {
	case 0:
		return "No priority"
	case 1:
		return "Urgent"
	case 2:
		return "High"
	case 3:
		return "Medium"
	case 4:
		return "Low"
	default:
		return "Unknown"
	}
}
