/*
 * Copyright 2025 Praveen Kumar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tracker

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/PraveenGongada/catalyst/internal/cache"
	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/styles"
)

type ViewMode int

const (
	WorkflowsView ViewMode = iota
	RunDetailsView
)

type FocusPane int

const (
	LeftPane FocusPane = iota
	RightPane
)

type WorkflowRuns struct {
	Pages       map[int][]github.RunDetails
	TotalCount  int
	CurrentPage int
}

type WorkflowTrackerModel struct {
	apiClient      *github.APIClient
	workflowsCache *cache.DiskCache[string, []github.Workflow]
	pinnedCache    *cache.DiskCache[string, []cache.WorkflowPin]
	runsCache      *cache.InMemoryCache[int64, WorkflowRuns]
	jobsCache      *cache.InMemoryCache[int64, []github.JobDetails]
	spinner        spinner.Model
	viewMode       ViewMode

	// View 1: Workflows → Runs
	workflows         []github.Workflow
	filteredWorkflows []github.Workflow
	selectedWorkflow  int
	runs              []github.RunDetails
	selectedRun       int
	focusPane         FocusPane

	// Pagination
	runsPage       int
	runsTotalCount int
	jumpQuery      string
	jumpFocused    bool

	filterQuery   string
	filterFocused bool

	// View 2: Jobs → Steps
	currentRun        *github.RunDetails
	jobs              []github.JobDetails
	selectedJob       int
	workflowsViewport viewport.Model
	runsViewport      viewport.Model
	jobsViewport      viewport.Model
	stepsViewport     viewport.Model

	// UI state
	width            int
	height           int
	loading          bool // Initial workflows loading
	loadingWorkflows bool // Refreshing workflows
	loadingRuns      bool // Loading runs for selected workflow
	loadingJobs      bool // Loading jobs for selected run
	error            error
	leftPaneWidth    int
	rightPaneWidth   int
	lastPollTime     time.Time
	pollInterval     time.Duration
}

func NewWorkflowTrackerModel(repository string) *WorkflowTrackerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.ColorTextMain)

	vp := viewport.New(0, 0)
	vp.HighPerformanceRendering = false

	return &WorkflowTrackerModel{
		apiClient:         github.NewAPIClient(repository),
		workflowsCache:    cache.NewDiskCache[string, []github.Workflow]("workflows"),
		pinnedCache:       cache.NewDiskCache[string, []cache.WorkflowPin]("pinned_workflows"),
		runsCache:         cache.NewInMemoryCache[int64, WorkflowRuns](),
		jobsCache:         cache.NewInMemoryCache[int64, []github.JobDetails](),
		spinner:           s,
		viewMode:          WorkflowsView,
		selectedWorkflow:  0,
		selectedRun:       0,
		selectedJob:       0,
		workflowsViewport: vp,
		runsViewport:      vp,
		jobsViewport:      vp,
		stepsViewport:     vp,
		width:             styles.TrackerDefaultWidth,
		height:            styles.TrackerDefaultHeight,
		loading:           true,
		pollInterval:      10 * time.Second,
		focusPane:         LeftPane,
		filterQuery:       "",
		filterFocused:     false,
		filteredWorkflows: []github.Workflow{},
		runsPage:          1,
		runsTotalCount:    0,
		jumpQuery:         "",
		jumpFocused:       false,
	}
}

func (m *WorkflowTrackerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadWorkflows(),
		m.tick(),
	)
}

// Data loading methods

func (m *WorkflowTrackerModel) loadWorkflows() tea.Cmd {
	return m.loadWorkflowsInternal(false)
}

func (m *WorkflowTrackerModel) loadWorkflowsInternal(forceRefresh bool) tea.Cmd {
	return func() tea.Msg {
		var workflows []github.Workflow

		// Try cache first if not forcing refresh
		if !forceRefresh && m.workflowsCache != nil {
			workflows, _ = m.workflowsCache.Get(m.apiClient.Repository)
		}

		// If no cache or force refresh, fetch from API
		if workflows == nil || forceRefresh {
			res, apiErr := m.apiClient.ListWorkflows()
			if apiErr != nil {
				return WorkflowsLoadedMsg{err: apiErr}
			}
			workflows = res.Workflows

			// Update cache
			if m.workflowsCache != nil {
				_ = m.workflowsCache.Set(m.apiClient.Repository, workflows)
			}
		}

		// Sort workflows with pinned items first
		workflows = m.sortWorkflows(workflows)

		return WorkflowsLoadedMsg{workflows: workflows}
	}
}

func (m *WorkflowTrackerModel) loadRunsForWorkflow() tea.Cmd {
	return m.loadRunsForPage(1)
}

func (m *WorkflowTrackerModel) loadRunsForPage(page int) tea.Cmd {
	return func() tea.Msg {
		workflows := m.getDisplayWorkflows()
		if m.selectedWorkflow >= len(workflows) {
			return RunsLoadedMsg{}
		}

		wf := workflows[m.selectedWorkflow]

		// Check cache first
		if cached, found := m.runsCache.Get(wf.ID); found {
			if runs, ok := cached.Pages[page]; ok {
				// Update current page in cache since we are navigating to it
				cached.CurrentPage = page
				m.runsCache.Set(wf.ID, cached)
				
				return RunsLoadedMsg{
					runs:       runs,
					totalCount: cached.TotalCount,
					page:       page,
				}
			}
		}

		// Fetch from API
		res, err := m.apiClient.ListWorkflowRuns(wf.ID, 20, page)
		if err != nil {
			return RunsLoadedMsg{err: err}
		}

		// Update cache
		m.updateRunsCache(wf.ID, page, res.WorkflowRuns, res.TotalCount)

		return RunsLoadedMsg{
			runs:       res.WorkflowRuns,
			totalCount: res.TotalCount,
			page:       page,
		}
	}
}

func (m *WorkflowTrackerModel) loadRunsForWorkflowInternal(forceRefresh bool) tea.Cmd {
	return func() tea.Msg {
		workflows := m.getDisplayWorkflows()
		if m.selectedWorkflow >= len(workflows) {
			return RunsLoadedMsg{}
		}

		wf := workflows[m.selectedWorkflow]
		page := m.runsPage
		if page < 1 {
			page = 1
		}

		// Check cache first if not forcing refresh
		if !forceRefresh {
			if cached, found := m.runsCache.Get(wf.ID); found {
				if runs, ok := cached.Pages[page]; ok {
					return RunsLoadedMsg{
						runs:       runs,
						totalCount: cached.TotalCount,
						page:       page,
					}
				}
			}
		}

		// Fetch from API
		res, err := m.apiClient.ListWorkflowRuns(wf.ID, 20, page)
		if err != nil {
			return RunsLoadedMsg{err: err}
		}

		// Update cache
		m.updateRunsCache(wf.ID, page, res.WorkflowRuns, res.TotalCount)

		return RunsLoadedMsg{
			runs:       res.WorkflowRuns,
			totalCount: res.TotalCount,
			page:       page,
		}
	}
}

func (m *WorkflowTrackerModel) updateRunsCache(workflowID int64, page int, runs []github.RunDetails, totalCount int) {
	if cached, found := m.runsCache.Get(workflowID); found {
		cached.Pages[page] = runs
		cached.TotalCount = totalCount
		cached.CurrentPage = page
		m.runsCache.Set(workflowID, cached)
	} else {
		pages := make(map[int][]github.RunDetails)
		pages[page] = runs
		m.runsCache.Set(workflowID, WorkflowRuns{
			Pages:       pages,
			TotalCount:  totalCount,
			CurrentPage: page,
		})
	}
}

func (m *WorkflowTrackerModel) loadJobsForRun() tea.Cmd {
	return m.loadJobsForRunInternal(false)
}

func (m *WorkflowTrackerModel) loadJobsForRunInternal(forceRefresh bool) tea.Cmd {
	return func() tea.Msg {
		if m.currentRun == nil {
			return JobsLoadedMsg{}
		}

		// Check cache first if not forcing refresh
		if !forceRefresh {
			if cachedJobs, found := m.jobsCache.Get(m.currentRun.ID); found {
				return JobsLoadedMsg{jobs: cachedJobs}
			}
		}

		// Fetch from API
		res, err := m.apiClient.ListJobs(m.currentRun.ID)
		if err != nil {
			return JobsLoadedMsg{err: err}
		}

		// Update cache
		m.jobsCache.Set(m.currentRun.ID, res.Jobs)

		return JobsLoadedMsg{jobs: res.Jobs}
	}
}

func (m *WorkflowTrackerModel) tick() tea.Cmd {
	return tea.Tick(m.pollInterval, func(t time.Time) tea.Msg { return PollTickMsg(t) })
}

func (m *WorkflowTrackerModel) shouldPoll() bool {
	if time.Since(m.lastPollTime) < m.pollInterval {
		return false
	}
	switch m.viewMode {
	case WorkflowsView:
		for _, r := range m.runs {
			if r.Status == "in_progress" || r.Status == "queued" || r.Conclusion == "running" {
				return true
			}
		}
	case RunDetailsView:
		if m.currentRun != nil && (m.currentRun.Status == "in_progress" || m.currentRun.Status == "queued" || m.currentRun.Conclusion == "running") {
			return true
		}
	}
	return false
}

// Workflow management methods

func (m *WorkflowTrackerModel) togglePin(workflowID int64) {
	if m.pinnedCache == nil {
		return
	}

	_ = m.pinnedCache.Update(m.apiClient.Repository, func(pinned []cache.WorkflowPin) []cache.WorkflowPin {
		for i, pin := range pinned {
			if pin.WorkflowID == workflowID {
				return append(pinned[:i], pinned[i+1:]...)
			}
		}

		newPin := cache.WorkflowPin{
			WorkflowID: workflowID,
			PinnedAt:   time.Now(),
		}
		return append([]cache.WorkflowPin{newPin}, pinned...)
	})
}

func (m *WorkflowTrackerModel) sortWorkflows(workflows []github.Workflow) []github.Workflow {
	if m.pinnedCache == nil || len(workflows) == 0 {
		return workflows
	}

	pinnedWorkflows, found := m.pinnedCache.Get(m.apiClient.Repository)
	if !found || len(pinnedWorkflows) == 0 {
		return workflows
	}

	pinnedMap := make(map[int64]int)
	for i, pin := range pinnedWorkflows {
		pinnedMap[pin.WorkflowID] = i
	}

	var pinned, unpinned []github.Workflow
	for _, wf := range workflows {
		if _, ok := pinnedMap[wf.ID]; ok {
			pinned = append(pinned, wf)
		} else {
			unpinned = append(unpinned, wf)
		}
	}

	// Sort pinned by pin order
	for i := 0; i < len(pinned)-1; i++ {
		for j := i + 1; j < len(pinned); j++ {
			if pinnedMap[pinned[i].ID] > pinnedMap[pinned[j].ID] {
				pinned[i], pinned[j] = pinned[j], pinned[i]
			}
		}
	}

	// Sort unpinned alphabetically
	for i := 0; i < len(unpinned)-1; i++ {
		for j := i + 1; j < len(unpinned); j++ {
			if unpinned[i].Name > unpinned[j].Name {
				unpinned[i], unpinned[j] = unpinned[j], unpinned[i]
			}
		}
	}

	return append(pinned, unpinned...)
}

func (m *WorkflowTrackerModel) getDisplayWorkflows() []github.Workflow {
	if len(m.filteredWorkflows) > 0 {
		return m.filteredWorkflows
	}
	return m.workflows
}

func (m *WorkflowTrackerModel) applyFilter() {
	if strings.TrimSpace(m.filterQuery) == "" {
		m.filteredWorkflows = m.workflows
		return
	}

	lowerQuery := strings.ToLower(m.filterQuery)
	filtered := make([]github.Workflow, 0)
	for _, wf := range m.workflows {
		if strings.Contains(strings.ToLower(wf.Name), lowerQuery) {
			filtered = append(filtered, wf)
		}
	}
	m.filteredWorkflows = filtered
}
