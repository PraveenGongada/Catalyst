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
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/PraveenGongada/catalyst/internal/github"
)

func (m *WorkflowTrackerModel) handleUpKey() tea.Cmd {
	switch m.viewMode {
	case WorkflowsView:
		if m.focusPane == LeftPane {
			if m.selectedWorkflow > 0 {
				m.selectedWorkflow--
				m.updateWorkFlowsViewOffset()
				return m.handleWorkflowSelectionChange()
			}
		} else {
			if m.selectedRun > 0 {
				m.selectedRun--
				if m.selectedRun < m.runsViewport.YOffset {
					m.runsViewport.SetYOffset(m.selectedRun)
				}
			}
		}
	case RunDetailsView:
		if m.focusPane == LeftPane {
			if m.selectedJob > 0 {
				m.selectedJob--
				if m.selectedJob < m.jobsViewport.YOffset {
					m.jobsViewport.SetYOffset(m.selectedJob)
				}
			}
		} else {
			m.stepsViewport.LineUp(1)
		}
	}
	return nil
}

func (m *WorkflowTrackerModel) handleDownKey() tea.Cmd {
	switch m.viewMode {
	case WorkflowsView:
		if m.focusPane == LeftPane {
			workflows := m.getDisplayWorkflows()
			if m.selectedWorkflow < len(workflows)-1 {
				m.selectedWorkflow++
				m.updateWorkFlowsViewOffset()
				return m.handleWorkflowSelectionChange()
			}
		} else {
			if m.selectedRun < len(m.runs)-1 {
				m.selectedRun++
				if m.selectedRun >= m.runsViewport.YOffset+m.runsViewport.Height {
					m.runsViewport.SetYOffset(m.selectedRun - m.runsViewport.Height + 1)
				}
			}
		}
	case RunDetailsView:
		if m.focusPane == LeftPane {
			if m.selectedJob < len(m.jobs)-1 {
				m.selectedJob++
				if m.selectedJob >= m.jobsViewport.YOffset+m.jobsViewport.Height {
					m.jobsViewport.SetYOffset(m.selectedJob - m.jobsViewport.Height + 1)
				}
			}
		} else {
			m.stepsViewport.LineDown(1)
		}
	}
	return nil
}

func (m *WorkflowTrackerModel) handleEnterKey() tea.Cmd {
	if m.viewMode == WorkflowsView && len(m.runs) > 0 {
		m.currentRun = &m.runs[m.selectedRun]
		m.viewMode = RunDetailsView
		m.selectedJob = 0
		m.focusPane = LeftPane
		m.pollInterval = 5 * time.Second
		m.loadingJobs = true
		return m.loadJobsForRun()
	}
	return nil
}

func (m *WorkflowTrackerModel) updateWorkFlowsViewOffset() {
	if m.selectedWorkflow < m.workflowsViewport.YOffset {
		m.workflowsViewport.SetYOffset(m.selectedWorkflow)
	} else if m.selectedWorkflow >= m.workflowsViewport.YOffset+m.workflowsViewport.Height {
		m.workflowsViewport.SetYOffset(m.selectedWorkflow - m.workflowsViewport.Height + 1)
	}
}

func (m *WorkflowTrackerModel) handleWorkflowSelectionChange() tea.Cmd {
	m.selectedRun = 0
	m.runsViewport.SetYOffset(0)

	workflows := m.getDisplayWorkflows()
	if m.selectedWorkflow >= len(workflows) {
		return nil
	}

	wf := workflows[m.selectedWorkflow]

	if cached, found := m.runsCache.Get(wf.ID); found {
		targetPage := max(cached.CurrentPage, 1)

		if runs, ok := cached.Pages[targetPage]; ok {
			m.runs = runs
			m.runsPage = targetPage
			m.runsTotalCount = cached.TotalCount
			return nil
		} else {
			m.runs = []github.RunDetails{}
			m.loadingRuns = true
			m.runsPage = targetPage
			return m.loadRunsForPage(targetPage)
		}
	}

	m.runs = []github.RunDetails{}
	m.runsPage = 1
	m.runsTotalCount = 0
	m.loadingRuns = true
	return m.loadRunsForPage(1)
}
