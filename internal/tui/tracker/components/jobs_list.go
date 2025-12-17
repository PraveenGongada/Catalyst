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

func RenderJobsList(
	jobs []github.JobDetails,
	selectedJob int,
	focusPane bool,
	paneWidth int,
) string {
	var b strings.Builder

	if len(jobs) == 0 {
		return styles.MetadataDim.Render("No jobs found.")
	}

	for i, job := range jobs {
		isSelected := i == selectedJob
		isFocused := focusPane

		var textStyle lipgloss.Style
		if isSelected && isFocused {
			textStyle = styles.ListItemSelected
		} else if isSelected {
			textStyle = styles.ListItemSelectedText
		} else {
			textStyle = styles.ListItemUnselected
		}

		statusSymbol := textStyle.Render(styles.GetColoredStatusSymbol(string(job.Status), string(job.Conclusion)))
		jobName := helpers.Truncate(job.Name, paneWidth-2)
		name := textStyle.Render(fmt.Sprintf(" %-*s", paneWidth-2, jobName))

		line := lipgloss.JoinHorizontal(lipgloss.Left, statusSymbol, name)

		b.WriteString(line + "\n")
	}
	return b.String()
}
