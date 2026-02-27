package linear

import (
	"context"
	"fmt"
	"net/http"

	graphql "github.com/hasura/go-graphql-client"
)

const linearAPIURL = "https://api.linear.app/graphql"

// Client is a Linear GraphQL API client.
type Client struct {
	gql *graphql.Client
}

// NewClient creates a new Linear API client authenticated with the given API key.
func NewClient(apiKey string) *Client {
	httpClient := &http.Client{
		Transport: &authTransport{token: apiKey, base: http.DefaultTransport},
	}
	return &Client{
		gql: graphql.NewClient(linearAPIURL, httpClient),
	}
}

// authTransport injects the Authorization header on every request.
type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", t.token)
	return t.base.RoundTrip(req)
}

// ListIssues retrieves issues, optionally filtered.
func (c *Client) ListIssues(ctx context.Context, filter IssueFilter) ([]Issue, error) {
	var query struct {
		Issues struct {
			Nodes []struct {
				ID         graphql.String
				Identifier graphql.String
				Title      graphql.String
				Priority   graphql.Int
				URL        graphql.String
				State      struct {
					ID   graphql.String
					Name graphql.String
					Type graphql.String
				}
				Assignee *struct {
					ID          graphql.String
					Name        graphql.String
					DisplayName graphql.String
					Email       graphql.String
				}
				Team struct {
					ID   graphql.String
					Name graphql.String
					Key  graphql.String
				}
				CreatedAt graphql.String
				UpdatedAt graphql.String
			}
		} `graphql:"issues(filter: $filter, first: 50)"`
	}

	filterInput, err := buildIssueFilter(filter)
	if err != nil {
		return nil, err
	}

	variables := map[string]interface{}{
		"filter": filterInput,
	}

	if err := c.gql.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("listing issues: %w", normalizeError(err))
	}

	var issues []Issue
	for _, n := range query.Issues.Nodes {
		issue := Issue{
			ID:         string(n.ID),
			Identifier: string(n.Identifier),
			Title:      string(n.Title),
			Priority:   int(n.Priority),
			URL:        string(n.URL),
			State: IssueState{
				ID:   string(n.State.ID),
				Name: string(n.State.Name),
				Type: string(n.State.Type),
			},
			Team: Team{
				ID:   string(n.Team.ID),
				Name: string(n.Team.Name),
				Key:  string(n.Team.Key),
			},
			CreatedAt: string(n.CreatedAt),
			UpdatedAt: string(n.UpdatedAt),
		}
		if n.Assignee != nil {
			issue.Assignee = &User{
				ID:          string(n.Assignee.ID),
				Name:        string(n.Assignee.Name),
				DisplayName: string(n.Assignee.DisplayName),
				Email:       string(n.Assignee.Email),
			}
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

// GetIssue retrieves a single issue by its identifier (e.g. "ENG-123").
func (c *Client) GetIssue(ctx context.Context, issueID string) (*Issue, error) {
	var query struct {
		Issue struct {
			ID          graphql.String
			Identifier  graphql.String
			Title       graphql.String
			Description graphql.String
			Priority    graphql.Int
			URL         graphql.String
			State       struct {
				ID   graphql.String
				Name graphql.String
				Type graphql.String
			}
			Assignee *struct {
				ID          graphql.String
				Name        graphql.String
				DisplayName graphql.String
				Email       graphql.String
			}
			Team struct {
				ID   graphql.String
				Name graphql.String
				Key  graphql.String
			}
			CreatedAt graphql.String
			UpdatedAt graphql.String
		} `graphql:"issue(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(issueID),
	}

	if err := c.gql.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("getting issue: %w", normalizeError(err))
	}

	issue := &Issue{
		ID:          string(query.Issue.ID),
		Identifier:  string(query.Issue.Identifier),
		Title:       string(query.Issue.Title),
		Description: string(query.Issue.Description),
		Priority:    int(query.Issue.Priority),
		URL:         string(query.Issue.URL),
		State: IssueState{
			ID:   string(query.Issue.State.ID),
			Name: string(query.Issue.State.Name),
			Type: string(query.Issue.State.Type),
		},
		Team: Team{
			ID:   string(query.Issue.Team.ID),
			Name: string(query.Issue.Team.Name),
			Key:  string(query.Issue.Team.Key),
		},
		CreatedAt: string(query.Issue.CreatedAt),
		UpdatedAt: string(query.Issue.UpdatedAt),
	}
	if query.Issue.Assignee != nil {
		issue.Assignee = &User{
			ID:          string(query.Issue.Assignee.ID),
			Name:        string(query.Issue.Assignee.Name),
			DisplayName: string(query.Issue.Assignee.DisplayName),
			Email:       string(query.Issue.Assignee.Email),
		}
	}
	return issue, nil
}

// CreateIssue creates a new Linear issue.
func (c *Client) CreateIssue(ctx context.Context, input IssueCreateInput) (*Issue, error) {
	var mutation struct {
		IssueCreate struct {
			Success graphql.Boolean
			Issue   struct {
				ID         graphql.String
				Identifier graphql.String
				Title      graphql.String
				URL        graphql.String
				State      struct {
					Name graphql.String
				}
				Team struct {
					Key graphql.String
				}
				CreatedAt graphql.String
			}
		} `graphql:"issueCreate(input: $input)"`
	}

	mutInput := buildCreateInput(input)
	variables := map[string]interface{}{
		"input": mutInput,
	}

	if err := c.gql.Mutate(ctx, &mutation, variables); err != nil {
		return nil, fmt.Errorf("creating issue: %w", normalizeError(err))
	}

	if !bool(mutation.IssueCreate.Success) {
		return nil, fmt.Errorf("creating issue: operation returned success=false")
	}

	n := mutation.IssueCreate.Issue
	return &Issue{
		ID:         string(n.ID),
		Identifier: string(n.Identifier),
		Title:      string(n.Title),
		URL:        string(n.URL),
		State:      IssueState{Name: string(n.State.Name)},
		Team:       Team{Key: string(n.Team.Key)},
		CreatedAt:  string(n.CreatedAt),
	}, nil
}

// UpdateIssue updates an existing Linear issue.
func (c *Client) UpdateIssue(ctx context.Context, issueID string, input IssueUpdateInput) (*Issue, error) {
	var mutation struct {
		IssueUpdate struct {
			Success graphql.Boolean
			Issue   struct {
				ID         graphql.String
				Identifier graphql.String
				Title      graphql.String
				URL        graphql.String
				State      struct {
					Name graphql.String
				}
				Priority  graphql.Int
				UpdatedAt graphql.String
			}
		} `graphql:"issueUpdate(id: $id, input: $input)"`
	}

	mutInput := buildUpdateInput(input)
	variables := map[string]interface{}{
		"id":    graphql.String(issueID),
		"input": mutInput,
	}

	if err := c.gql.Mutate(ctx, &mutation, variables); err != nil {
		return nil, fmt.Errorf("updating issue: %w", normalizeError(err))
	}

	if !bool(mutation.IssueUpdate.Success) {
		return nil, fmt.Errorf("updating issue: operation returned success=false")
	}

	n := mutation.IssueUpdate.Issue
	return &Issue{
		ID:         string(n.ID),
		Identifier: string(n.Identifier),
		Title:      string(n.Title),
		URL:        string(n.URL),
		State:      IssueState{Name: string(n.State.Name)},
		Priority:   int(n.Priority),
		UpdatedAt:  string(n.UpdatedAt),
	}, nil
}

// DeleteIssue deletes a Linear issue.
func (c *Client) DeleteIssue(ctx context.Context, issueID string) error {
	var mutation struct {
		IssueDelete struct {
			Success graphql.Boolean
		} `graphql:"issueDelete(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(issueID),
	}

	if err := c.gql.Mutate(ctx, &mutation, variables); err != nil {
		return fmt.Errorf("deleting issue: %w", normalizeError(err))
	}

	if !bool(mutation.IssueDelete.Success) {
		return fmt.Errorf("deleting issue: operation returned success=false")
	}

	return nil
}

// ListTeams returns all teams the authenticated user has access to.
func (c *Client) ListTeams(ctx context.Context) ([]Team, error) {
	var query struct {
		Teams struct {
			Nodes []struct {
				ID   graphql.String
				Name graphql.String
				Key  graphql.String
			}
		} `graphql:"teams"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("listing teams: %w", normalizeError(err))
	}

	var teams []Team
	for _, n := range query.Teams.Nodes {
		teams = append(teams, Team{
			ID:   string(n.ID),
			Name: string(n.Name),
			Key:  string(n.Key),
		})
	}
	return teams, nil
}

// GetViewer returns the currently authenticated user.
func (c *Client) GetViewer(ctx context.Context) (*User, error) {
	var query struct {
		Viewer struct {
			ID          graphql.String
			Name        graphql.String
			DisplayName graphql.String
			Email       graphql.String
		}
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("getting viewer: %w", normalizeError(err))
	}

	return &User{
		ID:          string(query.Viewer.ID),
		Name:        string(query.Viewer.Name),
		DisplayName: string(query.Viewer.DisplayName),
		Email:       string(query.Viewer.Email),
	}, nil
}

// IssueFilterInput is the GraphQL input type for filtering issues.
// Only exported for use in GraphQL variable marshaling.
type IssueFilterInput struct {
	Team     *IDComparator     `json:"team,omitempty"`
	Assignee *NullableIDComp   `json:"assignee,omitempty"`
	State    *StateFilterInput `json:"state,omitempty"`
}

// IDComparator filters by an ID equality.
type IDComparator struct {
	ID *IDComp `json:"id,omitempty"`
}

// IDComp holds an EQ comparator.
type IDComp struct {
	Eq string `json:"eq"`
}

// NullableIDComp filters by nullable ID.
type NullableIDComp struct {
	ID *IDComp `json:"id,omitempty"`
}

// StateFilterInput filters by state name.
type StateFilterInput struct {
	Name *StringComparator `json:"name,omitempty"`
}

// StringComparator holds a string equality comparator.
type StringComparator struct {
	Eq string `json:"eq"`
}

func buildIssueFilter(f IssueFilter) (IssueFilterInput, error) {
	var input IssueFilterInput
	if f.TeamID != "" {
		input.Team = &IDComparator{ID: &IDComp{Eq: f.TeamID}}
	}
	if f.State != "" {
		input.State = &StateFilterInput{Name: &StringComparator{Eq: f.State}}
	}
	return input, nil
}

// IssueCreateInputGQL is the GraphQL input for issueCreate.
type IssueCreateInputGQL struct {
	Title       graphql.String  `json:"title"`
	TeamID      graphql.String  `json:"teamId"`
	Description *graphql.String `json:"description,omitempty"`
	Priority    *graphql.Int    `json:"priority,omitempty"`
	AssigneeID  *graphql.String `json:"assigneeId,omitempty"`
	ParentID    *graphql.String `json:"parentId,omitempty"`
}

func buildCreateInput(input IssueCreateInput) IssueCreateInputGQL {
	gqlInput := IssueCreateInputGQL{
		Title:  graphql.String(input.Title),
		TeamID: graphql.String(input.TeamID),
	}
	if input.Description != "" {
		d := graphql.String(input.Description)
		gqlInput.Description = &d
	}
	if input.Priority != nil {
		p := graphql.Int(*input.Priority)
		gqlInput.Priority = &p
	}
	if input.AssigneeID != "" {
		a := graphql.String(input.AssigneeID)
		gqlInput.AssigneeID = &a
	}
	if input.ParentID != "" {
		pid := graphql.String(input.ParentID)
		gqlInput.ParentID = &pid
	}
	return gqlInput
}

// IssueUpdateInputGQL is the GraphQL input for issueUpdate.
type IssueUpdateInputGQL struct {
	Title       *graphql.String `json:"title,omitempty"`
	Description *graphql.String `json:"description,omitempty"`
	StateID     *graphql.String `json:"stateId,omitempty"`
	Priority    *graphql.Int    `json:"priority,omitempty"`
	AssigneeID  *graphql.String `json:"assigneeId,omitempty"`
}

func buildUpdateInput(input IssueUpdateInput) IssueUpdateInputGQL {
	gqlInput := IssueUpdateInputGQL{}
	if input.Title != nil {
		t := graphql.String(*input.Title)
		gqlInput.Title = &t
	}
	if input.Description != nil {
		d := graphql.String(*input.Description)
		gqlInput.Description = &d
	}
	if input.StateID != nil {
		s := graphql.String(*input.StateID)
		gqlInput.StateID = &s
	}
	if input.Priority != nil {
		p := graphql.Int(*input.Priority)
		gqlInput.Priority = &p
	}
	if input.AssigneeID != nil {
		a := graphql.String(*input.AssigneeID)
		gqlInput.AssigneeID = &a
	}
	return gqlInput
}

// normalizeError extracts a human-readable message from a GraphQL error.
func normalizeError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s", err.Error())
}
