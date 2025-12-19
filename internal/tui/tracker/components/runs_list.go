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

	"github.com/PraveenGongada/catalyst/internal/constants"
	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/helpers"
)

func RenderRunsHeader(availWidth int) string {
	remainingWidth := max(availWidth-constants.RunItemsWidth, 0)

	tableHeaderStyle := styles.TableHeader.PaddingLeft(
		remainingWidth / constants.RunItemsWidthFactor,
	)
	// No need to add the extra padding to the first item
	hID := styles.TableHeader.PaddingLeft(1).
		Render(fmt.Sprintf(" %-*s", constants.RunIdWidth-1, "ID"))
	hBranch := tableHeaderStyle.Render(fmt.Sprintf(" %-*s", constants.RunBranchWidth-1, "Branch"))
	hUser := tableHeaderStyle.Render(fmt.Sprintf(" %-*s", constants.RunUserWidth-1, "User"))
	hStatus := tableHeaderStyle.Render(fmt.Sprintf(" %-*s", constants.RunStatusWidth-1, "Status"))
	hTime := tableHeaderStyle.Render(fmt.Sprintf(" %-*s", constants.RunTimeWidth-1, "Time"))
	hDuration := tableHeaderStyle.Render(
		fmt.Sprintf(" %-*s", constants.RunDurationWidth-1, "Duration"),
	)

	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		hID,
		hUser,
		hStatus,
		hBranch,
		hTime,
		hDuration,
	)

	return header
}

func RenderRunsList(
	runs []github.RunDetails,
	selectedRun int,
	isFocused bool,
	availWidth int,
	loading bool,
	spinnerView string,
) string {
	var b strings.Builder

	if loading {
		return styles.MetadataDim.Render(spinnerView + " Loading runs...")
	}
	if len(runs) == 0 {
		return styles.MetadataDim.PaddingLeft(2).Render("No runs found.")
	}

	remainingWidth := max(availWidth-constants.RunItemsWidth, 0)

	for i, run := range runs {
		isSelected := i == selectedRun

		idStr := fmt.Sprintf("#%d", run.RunNumber)

		statusSymbolStr := styles.GetStatusSymbol(
			string(run.Status),
			string(run.Conclusion),
		)
		statusTextStr := styles.GetStatusText(string(run.Status),
			string(run.Conclusion),
		)
		branchStr := helpers.Truncate(run.HeadBranch, constants.RunBranchWidth-2)
		userStr := helpers.Truncate(run.Actor.Login, constants.RunUserWidth-2)
		timeStr := helpers.FormatRelativeTime(run.RunStartedAt)
		durationStr := helpers.FormatDuration(run.RunStartedAt, run.UpdatedAt)

		var textStyle lipgloss.Style
		if isSelected && isFocused {
			textStyle = styles.ListItemSelected
		} else if isSelected {
			textStyle = styles.ListItemSelectedText
		} else {
			textStyle = styles.ListItemUnselected
		}

		statusSymbol := textStyle.Render(statusSymbolStr)
		coloredId := textStyle.Render(fmt.Sprintf(" %-*s", constants.RunIdWidth-1, idStr))
		// We no need to add the Extra Width to the first item
		textStyle = textStyle.PaddingLeft(remainingWidth / constants.RunItemsWidthFactor)
		coloredBranch := textStyle.Render(
			fmt.Sprintf(" %-*s", constants.RunBranchWidth-1, branchStr),
		)
		coloredUser := textStyle.Render(fmt.Sprintf(" %-*s", constants.RunUserWidth-1, userStr))
		coloredStatus := textStyle.Render(
			fmt.Sprintf(" %-*s", constants.RunStatusWidth-1, statusTextStr),
		)
		coloredTime := textStyle.Render(fmt.Sprintf(" %-*s", constants.RunTimeWidth-1, timeStr))
		coloredDuration := textStyle.Render(
			fmt.Sprintf(" %-*s", constants.RunDurationWidth-1, durationStr),
		)

		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			statusSymbol,
			coloredId,
			coloredUser,
			coloredStatus,
			coloredBranch,
			coloredTime,
			coloredDuration,
		)

		b.WriteString(line + "\n")
	}

	return b.String()
}
