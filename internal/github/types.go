package github

import "time"

type WorkflowState string

const (
	Active             WorkflowState = "active"
	Deleted            WorkflowState = "deleted"
	DisabledFork       WorkflowState = "disabled_fork"
	DisabledInactivity WorkflowState = "disabled_inactivity"
	DisabledManually   WorkflowState = "disabled_manually"
)

type JobStatus string

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusInProgress JobStatus = "in_progress"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusWaiting    JobStatus = "waiting"
	JobStatusRequested  JobStatus = "requested"
	JobStatusPending    JobStatus = "pending"
)

type JobConclusion string

const (
	JobConclusionSuccess        JobConclusion = "success"
	JobConclusionFailure        JobConclusion = "failure"
	JobConclusionNeutral        JobConclusion = "neutral"
	JobConclusionCancelled      JobConclusion = "cancelled"
	JobConclusionSkipped        JobConclusion = "skipped"
	JobConclusionTimedOut       JobConclusion = "timed_out"
	JobConclusionActionRequired JobConclusion = "action_required"
)

type StepStatus string

const (
	StepStatusQueued     StepStatus = "queued"
	StepStatusInProgress StepStatus = "in_progress"
	StepStatusCompleted  StepStatus = "completed"
)

type Workflow struct {
	ID        int64         `json:"id"`
	NodeID    string        `json:"node_id"`
	Name      string        `json:"name"`
	Path      string        `json:"path"`
	State     WorkflowState `json:"state"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	HTMLURL   string        `json:"html_url"`
}

type WorkflowsResponse struct {
	TotalCount int        `json:"total_count"`
	Workflows  []Workflow `json:"workflows"`
}

type RunDetails struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	NodeID       string    `json:"node_id"`
	HeadBranch   string    `json:"head_branch"`
	HeadSHA      string    `json:"head_sha"`
	Path         string    `json:"path"`
	DisplayTitle string    `json:"display_title"`
	RunNumber    int       `json:"run_number"`
	Event        string    `json:"event"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	WorkflowID   int64     `json:"workflow_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RunStartedAt time.Time `json:"run_started_at"`
	Actor        Actor     `json:"actor"`
}

type Actor struct {
	Login     string `json:"login"`
	Url       string `json:"url"`
	Type      string `json:"type"`
	AvatarURL string `json:"avatar_url"`
}

type RunsResponse struct {
	TotalCount   int          `json:"total_count"`
	WorkflowRuns []RunDetails `json:"workflow_runs"`
}

type JobDetails struct {
	ID          int64         `json:"id"`
	RunID       int64         `json:"run_id"`
	Name        string        `json:"name"`
	Status      JobStatus     `json:"status"`
	Conclusion  JobConclusion `json:"conclusion"`
	HTMLURL     string        `json:"html_url"`
	StartedAt   time.Time     `json:"started_at"`
	CompletedAt time.Time     `json:"completed_at"`
	Steps       []Step        `json:"steps"`
}

type Step struct {
	Name        string     `json:"name"`
	Status      StepStatus `json:"status"`
	Conclusion  string     `json:"conclusion"`
	Number      int        `json:"number"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt time.Time  `json:"completed_at"`
}

type JobsResponse struct {
	TotalCount int          `json:"total_count"`
	Jobs       []JobDetails `json:"jobs"`
}
