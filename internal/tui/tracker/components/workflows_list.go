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
	"strings"

	"github.com/PraveenGongada/catalyst/internal/cache"
	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker/helpers"
)

func RenderWorkflowsHeader(repository string) string {
	repo := strings.Split(repository, "/")
	if len(repo) != 2 {
		return repository
	}
	return styles.TableHeader.PaddingLeft(1).Render("Repo: " + repo[1])
}

func RenderWorkflowsList(
	workflows []github.Workflow,
	selectedWorkflow int,
	isFocused bool,
	paneWidth int,
	repository string,
	pinnedCache *cache.DiskCache[string, []cache.WorkflowPin],
) string {
	var b strings.Builder

	if len(workflows) == 0 {
		return styles.MetadataDim.Render("No workflows found")
	}

	var pinnedWorkflows []cache.WorkflowPin
	if pinnedCache != nil {
		pinnedWorkflows, _ = pinnedCache.Get(repository)
	}

	pinnedMap := make(map[int64]bool)
	for _, pin := range pinnedWorkflows {
		pinnedMap[pin.WorkflowID] = true
	}

	pinIndicator := "📌"

	for i, wf := range workflows {
		isSelected := i == selectedWorkflow
		isPinned := pinnedMap[wf.ID]

		var displayName string
		maxNameLen := max(paneWidth-4, 0)

		if isPinned {
			displayName = helpers.Truncate(wf.Name, maxNameLen)
			padding := max(paneWidth-len(displayName)-5, 0)
			displayName = " " + displayName + strings.Repeat(
				" ",
				padding,
			) + " " + pinIndicator + " "
		} else {
			displayName = helpers.Truncate(wf.Name, maxNameLen)
			padding := max(paneWidth-len(displayName)-1, 0)
			displayName = " " + displayName + strings.Repeat(" ", padding)
		}

		var line string
		if isSelected && isFocused {
			line = styles.ListItemSelected.Render(displayName)
		} else if isSelected {
			line = styles.ListItemSelectedText.Render(displayName)
		} else {
			line = styles.ListItemUnselected.Render(displayName)
		}

		b.WriteString(line + "\n")
	}
	return b.String()
}
