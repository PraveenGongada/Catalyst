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

package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/helpers"
)

func RenderStepsHeader(run *github.RunDetails) string {
	if run == nil {
		return ""
	}

	userStr := fmt.Sprintf("%s", run.Actor.Login)
	userStyled := styles.MetadataDim.Render(userStr)

	commitHash := run.HeadSHA
	if len(commitHash) > 7 {
		commitHash = commitHash[:7]
	}
	commitStr := fmt.Sprintf("%s", commitHash)
	commitStyled := styles.MetadataDim.Render(commitStr)

	branchStr := fmt.Sprintf(" %s", run.HeadBranch)
	branchStyled := styles.MetadataDim.Render(branchStr)

	statusText := styles.GetStatusText(string(run.Status), string(run.Conclusion))
	durationStr := helpers.FormatDuration(run.RunStartedAt, run.UpdatedAt)
	durationStyled := styles.MetadataDim.Render(fmt.Sprintf("⏱ %s", durationStr))

	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		" ",
		userStyled,
		"  ",
		commitStyled,
		"  ",
		branchStyled,
		"  ",
		statusText,
		"  ",
		durationStyled,
	)

	return header
}

func RenderStepsList(jobs []github.JobDetails, selectedJob int) string {
	if selectedJob >= len(jobs) || len(jobs) == 0 {
		return styles.MetadataDim.Render("No build steps found.")
	}
	job := jobs[selectedJob]

	if len(job.Steps) == 0 {
		return styles.MetadataDim.Render("No build steps found.")
	}

	var stepLines []string

	for i, step := range job.Steps {
		statusSymbol := styles.GetStepSymbol(string(step.Status), string(step.Conclusion))
		line := fmt.Sprintf(" %s  %s",
			statusSymbol,
			lipgloss.NewStyle().Foreground(styles.ColorTextMain).Render(step.Name),
		)
		stepLines = append(stepLines, line)

		if i < len(job.Steps)-1 {
			connector := " " + styles.TimelineConnector.Render("│")
			stepLines = append(stepLines, connector)
		}
	}

	return strings.Join(stepLines, "\n")
}
