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
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/styles"
)

// RenderStepsHeader renders metadata about the current run
func RenderStepsHeader(run *github.RunDetails) string {
	if run == nil {
		return ""
	}

	// User
	userStr := fmt.Sprintf("👤 %s", run.Actor.Login)
	userStyled := styles.MetadataDim.Render(userStr)

	// Commit (short hash)
	commitHash := run.HeadSHA
	if len(commitHash) > 7 {
		commitHash = commitHash[:7]
	}
	commitStr := fmt.Sprintf("@ %s", commitHash)
	commitStyled := styles.MetadataDim.Render(commitStr)

	// Branch
	branchStr := fmt.Sprintf("\\ %s", run.HeadBranch)
	branchStyled := styles.MetadataDim.Render(branchStr)

	// Status
	statusSym := styles.GetStatusSymbol(run.Status, run.Conclusion)
	statusTxt := styles.GetStatusText(run.Status, run.Conclusion)
	statusStyled := styles.GetStatusStyle(run.Status, run.Conclusion).Render(fmt.Sprintf("%s %s", statusSym, statusTxt))

	// Duration
	durationStr := formatDuration(run.RunStartedAt, run.UpdatedAt)
	durationStyled := styles.MetadataDim.Render(fmt.Sprintf("⏱ %s", durationStr))

	// Join all metadata with spacing
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		" ",
		userStyled,
		"  ",
		commitStyled,
		"  ",
		branchStyled,
		"  ",
		statusStyled,
		"  ",
		durationStyled,
	)

	return header
}

// RenderStepsList renders the steps timeline for a job
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
		statusSym := styles.GetStepStatusSymbol(string(step.Status), string(step.Conclusion))
		statusStyle := styles.GetStatusStyle(string(step.Status), string(step.Conclusion))

		line := fmt.Sprintf(" %s  %s",
			statusStyle.Render(statusSym),
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

// Local helper to avoid import cycle with helpers package
func formatDuration(start, end time.Time) string {
	if start.IsZero() {
		return "-"
	}
	var d time.Duration
	if end.IsZero() {
		d = time.Since(start)
	} else {
		d = end.Sub(start)
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours >= 1 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes >= 1 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
