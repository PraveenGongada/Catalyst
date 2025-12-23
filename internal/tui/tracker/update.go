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
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/PraveenGongada/catalyst/internal/github"
)

func (m *WorkflowTrackerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// --- 1. Calculate Layout Widths ---
		m.leftPaneWidth = int(float64(m.width) * 0.30)
		if m.leftPaneWidth < 35 {
			m.leftPaneWidth = 35
		}

		m.rightPaneWidth = m.width - m.leftPaneWidth

		// --- 2. Calculate Vertical Heights ---
		// We need to fit: Header + Viewport + Border(2) + Footer(2)
		// Total Available for (Header + Viewport) = m.Height - 4
		availableContentHeight := m.height - 4

		// Constants for header heights based on render logic:
		// Workflows Header: "WORKFLOWS\n\n" = 3 lines
		// Runs Header: "RUNS\n\n" + HeaderRow + "\n" = 5 lines
		// Jobs Header: "JOBS\n\n" = 3 lines
		// Steps Header: "STEPS...\n\n" = 3 lines

		// View 1: Workflows
		m.workflowsViewport.Width = m.leftPaneWidth - 2
		m.workflowsViewport.Height = availableContentHeight - 1 // Header is 3 lines

		// View 1: Runs
		m.runsViewport.Width = m.rightPaneWidth - 2
		m.runsViewport.Height = availableContentHeight - 2 // Header is 5 lines (includes table columns)

		// View 2: Jobs
		m.jobsViewport.Width = m.leftPaneWidth - 2

		m.jobsViewport.Height = m.height - 3

		// View 2: Steps
		m.stepsViewport.Width = m.rightPaneWidth - 4
		m.stepsViewport.Height = m.height - 5

	case tea.KeyMsg:
		if m.jumpFocused {
			switch msg.String() {
			case "esc":
				m.jumpFocused = false
				m.jumpQuery = ""
				return m, nil
			case "enter":
				m.jumpFocused = false
				// Parse page number
				// Simple implementation: scan int
				var page int
				_, err := fmt.Sscanf(m.jumpQuery, "%d", &page)
				if err == nil && page > 0 {
					totalPages := (m.runsTotalCount + 19) / 20
					if page <= totalPages {
						m.loadingRuns = true
						m.runsPage = page
						return m, m.loadRunsForPage(page)
					}
				}
				m.jumpQuery = ""
				return m, nil
			default:
				switch msg.Type {
				case tea.KeyBackspace:
					if len(m.jumpQuery) > 0 {
						m.jumpQuery = m.jumpQuery[:len(m.jumpQuery)-1]
					}
				case tea.KeyRunes:
					// Only allow digits
					if msg.Runes[0] >= '0' && msg.Runes[0] <= '9' {
						m.jumpQuery += string(msg.Runes)
					}
				}
				return m, nil
			}
		}

		if m.filterFocused {
			switch msg.String() {
			case "esc":
				// Clear filter and reset
				m.filterQuery = ""
				m.filterFocused = false
				m.filteredWorkflows = m.workflows
				m.selectedWorkflow = 0
				if len(m.workflows) > 0 {
					m.loadingRuns = true
					return m, m.loadRunsForWorkflow()
				}
				return m, nil
			case "enter":
				// Apply filter and load runs
				m.filterFocused = false
				m.applyFilter()
				m.selectedWorkflow = 0
				if len(m.getDisplayWorkflows()) > 0 {
					m.loadingRuns = true
					return m, m.loadRunsForWorkflow()
				}
				return m, nil
			default:
				switch msg.Type {
				case tea.KeyBackspace:
					if len(m.filterQuery) > 0 {
						m.filterQuery = m.filterQuery[:len(m.filterQuery)-1]
					}
				case tea.KeyRunes:
					m.filterQuery += string(msg.Runes)
				}
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "/":
			if m.viewMode == WorkflowsView && m.focusPane == LeftPane {
				m.filterFocused = true
			}

		case "r":
			switch m.viewMode {
			case WorkflowsView:
				switch m.focusPane {
				case LeftPane:
					// Refresh workflows from API
					m.loadingWorkflows = true
					return m, m.loadWorkflowsInternal(true)
				case RightPane:
					// Refresh runs from API
					m.loadingRuns = true
					return m, m.loadRunsForWorkflowInternal(true)
				}
			case RunDetailsView:
				// Refresh jobs from API in run details view
				m.loadingJobs = true
				return m, m.loadJobsForRunInternal(true)
			}

		case "P":
			// Pin/unpin selected workflow
			if m.viewMode == WorkflowsView && m.focusPane == LeftPane {
				workflows := m.getDisplayWorkflows()
				if len(workflows) > 0 && m.selectedWorkflow < len(workflows) {
					wf := workflows[m.selectedWorkflow]
					if m.pinnedCache != nil {
						m.togglePin(wf.ID)
						return m, m.loadWorkflowsInternal(false)
					}
				}
			}

		case "up", "k":
			return m, m.handleUpKey()

		case "down", "j":
			return m, m.handleDownKey()

		case "tab":
			m.focusPane = (m.focusPane + 1) % 2

		case "enter":
			return m, m.handleEnterKey()

		case "esc":
			if m.viewMode == RunDetailsView {
				m.currentRun = nil
				m.jobs = []github.JobDetails{}
				m.selectedJob = 0
				m.viewMode = WorkflowsView
				m.pollInterval = 10 * time.Second
				m.focusPane = LeftPane
				return m, m.loadRunsForWorkflow()
			}

		case "left", "h":
			if m.focusPane == RightPane {
				m.focusPane = LeftPane
			}

		case "right", "l":
			if m.focusPane == LeftPane {
				m.focusPane = RightPane
			}

		case "n":
			if m.viewMode == WorkflowsView && m.focusPane == RightPane {
				totalPages := (m.runsTotalCount + 19) / 20
				if m.runsPage < totalPages {
					m.loadingRuns = true
					m.runsPage++
					return m, m.loadRunsForPage(m.runsPage)
				}
			}

		case "p": 
			if m.viewMode == WorkflowsView && m.focusPane == RightPane {
				if m.runsPage > 1 {
					m.loadingRuns = true
					m.runsPage--
					return m, m.loadRunsForPage(m.runsPage)
				}
			}

		case "z":
			if m.viewMode == WorkflowsView && m.focusPane == RightPane {
				m.jumpFocused = true
				m.jumpQuery = ""
			}
		}

	case WorkflowsLoadedMsg:
		m.loading = false
		m.loadingWorkflows = false
		if msg.err != nil {
			m.error = msg.err
		} else {
			m.workflows = msg.workflows
			m.applyFilter()
			workflows := m.getDisplayWorkflows()
			if len(workflows) > 0 && m.selectedWorkflow < len(workflows) {
				wf := workflows[m.selectedWorkflow]
				// Reset pagination on workflow switch/load
				m.runsPage = 1
				m.runsTotalCount = 0

				if _, found := m.runsCache.Get(wf.ID); !found {
					m.loadingRuns = true
					return m, m.loadRunsForWorkflow()
				}
			}
		}

	case RunsLoadedMsg:
		m.loading = false
		m.loadingRuns = false
		if msg.err != nil {
			m.error = msg.err
		} else {
			m.runs = msg.runs
			m.runsPage = msg.page
			m.runsTotalCount = msg.totalCount
			m.runsViewport.SetYOffset(0) // Reset scroll on page change
			m.selectedRun = 0

			// Ensure current page is saved to cache
			workflows := m.getDisplayWorkflows()
			if len(workflows) > 0 && m.selectedWorkflow < len(workflows) {
				wfID := workflows[m.selectedWorkflow].ID
				if cached, found := m.runsCache.Get(wfID); found {
					cached.CurrentPage = msg.page
					m.runsCache.Set(wfID, cached)
				}
			}
		}

	case JobsLoadedMsg:
		m.loading = false
		m.loadingJobs = false
		if msg.err != nil {
			m.error = msg.err
		} else {
			m.jobs = msg.jobs
		}

	case PollTickMsg:
		cmds = append(cmds, m.tick())
		if m.shouldPoll() {
			switch m.viewMode {
			case WorkflowsView:
				// Poll with force refresh to get latest status of in-progress runs
				cmds = append(cmds, m.loadRunsForWorkflowInternal(true))
			case RunDetailsView:
				// Poll with force refresh to get latest status of in-progress jobs
				cmds = append(cmds, m.loadJobsForRunInternal(true))
			default:
				return m, nil
			}
			m.lastPollTime = time.Now()
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	var vpCmd tea.Cmd
	m.workflowsViewport, vpCmd = m.workflowsViewport.Update(msg)
	cmds = append(cmds, vpCmd)
	m.runsViewport, vpCmd = m.runsViewport.Update(msg)
	cmds = append(cmds, vpCmd)
	m.jobsViewport, vpCmd = m.jobsViewport.Update(msg)
	cmds = append(cmds, vpCmd)
	m.stepsViewport, vpCmd = m.stepsViewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}
