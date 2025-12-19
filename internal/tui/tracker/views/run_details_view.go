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
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/components"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/helpers"
)

func RenderRunDetailsView(
	currentRun *github.RunDetails,
	jobs []github.JobDetails,
	selectedJob int,
	focusPane int,
	leftPaneWidth, rightPaneWidth, height int,
	loadingJobs bool,
	spinnerView string,
	jobsViewportContent, stepsViewportContent string,
) string {
	paneHeight := height - 3

	// Left pane content: show loader if loading jobs, otherwise show jobs list
	leftContentWidth := leftPaneWidth - 2
	var leftPaneContent string
	var rightPaneContent string
	if loadingJobs {
		leftPaneContent = helpers.RenderPaneLoader(
			spinnerView,
			"Loading jobs...",
			leftContentWidth-2,
			paneHeight,
		)
		rightPaneContent = helpers.RenderPaneLoader(
			spinnerView,
			"Loading steps...",
			rightPaneWidth-2,
			paneHeight,
		)
	} else {
		leftPaneContent = jobsViewportContent
		rightContentWidth := rightPaneWidth - 2
		rightBorderColor := styles.GetBorderColor(focusPane == 1)
		rightPaneContent = lipgloss.JoinVertical(
			lipgloss.Top,
			components.RenderStepsHeader(currentRun),
			lipgloss.NewStyle().Foreground(rightBorderColor).Render(strings.Repeat("─", rightContentWidth)),
			stepsViewportContent,
		)
	}

	leftPane := helpers.RenderPaneWithTitle(
		"Jobs",
		leftPaneContent,
		leftPaneWidth,
		paneHeight,
		focusPane == 0,
	)
	rightPane := helpers.RenderPaneWithTitle(
		"Steps",
		rightPaneContent,
		rightPaneWidth,
		paneHeight,
		focusPane == 1,
	)

	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	return panes + renderHelpBar(true, focusPane == 0, rightPaneWidth+leftPaneWidth)
}
