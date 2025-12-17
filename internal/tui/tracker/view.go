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

	"github.com/charmbracelet/lipgloss"

	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/components"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/views"
)

func (m *WorkflowTrackerModel) View() string {

	if m.loading && len(m.workflows) == 0 {
		return m.renderLoadingView()
	}

	if m.error != nil {
		return m.renderErrorView()
	}

	switch m.viewMode {
	case WorkflowsView:
		return m.renderWorkflowsView()
	case RunDetailsView:
		return m.renderRunDetailsView()
	default:
		return "Unknown view mode"
	}
}

func (m *WorkflowTrackerModel) renderLoadingView() string {
	title := styles.TrackerHeader.Render("WORKFLOW TRACKER")
	msg := styles.MetadataDim.Render(fmt.Sprintf("%s Loading...", m.spinner.View()))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, title+"\n"+msg)
}

func (m *WorkflowTrackerModel) renderErrorView() string {
	title := styles.BadgeFailure.Render("ERROR")
	msg := m.error.Error()
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, title+"\n"+msg)
}

func (m *WorkflowTrackerModel) renderWorkflowsView() string {
	// Render workflows list content
	m.renderWorkflowsList()

	// Render runs list content
	m.renderRunsList()

	// Get viewport contents
	workflowsViewportContent := m.workflowsViewport.View()
	runsViewportContent := m.runsViewport.View()

	return views.RenderWorkflowsView(
		m.getDisplayWorkflows(),
		m.workflows,
		m.selectedWorkflow,
		m.runs,
		m.selectedRun,
		int(m.focusPane),
		m.leftPaneWidth,
		m.rightPaneWidth,
		m.height,
		m.loadingWorkflows,
		m.loadingRuns,
		m.spinner.View(),
		m.apiClient.Repository,
		m.pinnedCache,
		workflowsViewportContent,
		runsViewportContent,
		m.filterQuery,
		m.filterFocused,
		m.runsPage,
		m.runsTotalCount,
		m.jumpQuery,
		m.jumpFocused,
	)
}

func (m *WorkflowTrackerModel) renderWorkflowsList() {
	content := components.RenderWorkflowsList(
		m.getDisplayWorkflows(),
		m.selectedWorkflow,
		m.focusPane == LeftPane,
		m.leftPaneWidth-2,
		m.apiClient.Repository,
		m.pinnedCache,
	)
	m.workflowsViewport.SetContent(content)
}

func (m *WorkflowTrackerModel) renderRunsList() {
	content := components.RenderRunsList(
		m.runs,
		m.selectedRun,
		m.focusPane == RightPane,
		m.rightPaneWidth-2,
		m.loading,
		m.spinner.View(),
	)
	m.runsViewport.SetContent(content)
}

func (m *WorkflowTrackerModel) renderRunDetailsView() string {
	// Render jobs list content
	m.renderJobsList()

	// Render steps list content
	m.renderStepsList()

	// Get viewport contents
	jobsViewportContent := m.jobsViewport.View()
	stepsViewportContent := m.stepsViewport.View()

	return views.RenderRunDetailsView(
		m.currentRun,
		m.jobs,
		m.selectedJob,
		int(m.focusPane),
		m.leftPaneWidth,
		m.rightPaneWidth,
		m.height,
		m.loadingJobs,
		m.spinner.View(),
		jobsViewportContent,
		stepsViewportContent,
	)
}

func (m *WorkflowTrackerModel) renderJobsList() {
	content := components.RenderJobsList(
		m.jobs,
		m.selectedJob,
		m.focusPane == LeftPane,
		m.leftPaneWidth-2,
	)
	m.jobsViewport.SetContent(content)
}

func (m *WorkflowTrackerModel) renderStepsList() {
	content := components.RenderStepsList(m.jobs, m.selectedJob)
	m.stepsViewport.SetContent(content)
}
