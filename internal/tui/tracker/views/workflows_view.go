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

package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/PraveenGongada/catalyst/internal/cache"
	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/components"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/helpers"
)

// RenderWorkflowsView renders the workflows ↔ runs two-pane view
func RenderWorkflowsView(
	displayWorkflows []github.Workflow,
	allWorkflows []github.Workflow,
	selectedWorkflow int,
	runs []github.RunDetails,
	selectedRun int,
	focusPane int, // 0=LeftPane, 1=RightPane
	leftPaneWidth, rightPaneWidth, height int,
	loadingWorkflows, loadingRuns bool,
	spinnerView string,
	repository string,
	pinnedCache *cache.DiskCache[string, []cache.WorkflowPin],
	workflowsViewportContent, runsViewportContent string,
	filterQuery string,
	filterFocused bool,
	runsPage, runsTotalCount int,
	jumpQuery string,
	jumpFocused bool,
) string {
	paneHeight := height - 3

	// Calculate content widths (pane width - 2 for borders)
	leftContentWidth := max(leftPaneWidth-2, 0)
	rightContentWidth := max(rightPaneWidth-2, 0)

	// Left pane content (workflows)
	var leftPaneContent string
	leftBorderColor := styles.GetBorderColor(focusPane == 0)
	if loadingWorkflows {
		// Calculate content height (pane height - header - border)
		contentHeight := paneHeight - 2
		loaderContent := helpers.RenderPaneLoader(spinnerView, "Refreshing workflows...", leftContentWidth-2, contentHeight)
		leftPaneContent = lipgloss.JoinVertical(
			lipgloss.Top,
			components.RenderWorkflowsHeader(repository),
			lipgloss.NewStyle().Foreground(leftBorderColor).Render(strings.Repeat("─", leftContentWidth)),
			loaderContent,
		)
	} else {
		leftPaneContent = lipgloss.JoinVertical(
			lipgloss.Top,
			components.RenderWorkflowsHeader(repository),
			lipgloss.NewStyle().Foreground(leftBorderColor).Render(strings.Repeat("─", leftContentWidth)),
			workflowsViewportContent,
		)
	}

	// Right pane content (runs)
	var rightPaneContent string
	rightBorderColor := styles.GetBorderColor(focusPane == 1)

	// Calculate pagination info
	totalPages := 0
	if runsTotalCount > 0 {
		totalPages = (runsTotalCount + 19) / 20
	}
	paginationInfo := fmt.Sprintf("Page %d/%d", runsPage, totalPages)
	if totalPages == 0 {
		paginationInfo = ""
	}

	paginationStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMeta).
		Align(lipgloss.Center).
		Width(rightContentWidth)

	paginationBar := paginationStyle.Render(paginationInfo)

	if loadingRuns {
		// Calculate content height (pane height - header - border - pagination)
		contentHeight := paneHeight - 3
		loaderContent := helpers.RenderPaneLoader(spinnerView, "Loading runs...", rightContentWidth-2, contentHeight)
		rightPaneContent = lipgloss.JoinVertical(
			lipgloss.Top,
			components.RenderRunsHeader(rightContentWidth),
			lipgloss.NewStyle().Foreground(rightBorderColor).Render(strings.Repeat("─", rightContentWidth)),
			loaderContent,
			paginationBar,
		)
	} else {
		rightPaneContent = lipgloss.JoinVertical(
			lipgloss.Top,
			components.RenderRunsHeader(rightContentWidth),
			lipgloss.NewStyle().Foreground(rightBorderColor).Render(strings.Repeat("─", rightContentWidth)),
			runsViewportContent,
			paginationBar,
		)
	}

	// Render panes with titles
	leftPane := helpers.RenderPaneWithTitle("Workflows", leftPaneContent, leftPaneWidth, paneHeight, focusPane == 0)
	rightPane := helpers.RenderPaneWithTitle("Runs", rightPaneContent, rightPaneWidth, paneHeight, focusPane == 1)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	var footer string
	if filterFocused {
		searchText := "🔍 Search: " + filterQuery
		footer = lipgloss.NewStyle().
			Foreground(styles.ColorTextMain).
			PaddingTop(1).
			Render(searchText)
	} else if jumpFocused {
		jumpText := "Go to page: " + jumpQuery
		footer = lipgloss.NewStyle().
			Foreground(styles.ColorTextMain).
			PaddingTop(1).
			Render(jumpText)
	} else {
		footer = renderHelpBar(false, focusPane == 0, rightPaneWidth+leftPaneWidth)
	}

	return content + footer
}

func renderHelpBar(inRunDetailsView bool, leftPaneFocused bool, width int) string {
	var help string
	if inRunDetailsView {
		help = "↑/↓ navigate • r refresh • tab/←/→ switch pane • esc back to runs • q quit"
	} else {
		if leftPaneFocused {
			help = "↑/↓ navigate • / filter • p pin • r refresh • tab/←/→ switch pane • enter view jobs • q quit"
		} else {
			help = "↑/↓ navigate • n/P next/prev page • j jump to page • r refresh • tab/←/→ switch pane • enter view jobs • q quit"
		}
	}
	return styles.HelpText.Width(width).Render(help)
}
