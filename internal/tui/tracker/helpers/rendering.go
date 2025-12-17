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

package helpers

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/PraveenGongada/catalyst/internal/styles"
)

// RenderPaneWithTitle renders a bordered pane with a title in the header
func RenderPaneWithTitle(title, content string, width, height int, isFocused bool) string {
	borderColor := styles.GetBorderColor(isFocused)

	// Create box style with borders but no top border (we'll render custom top with title)
	boxStyle := lipgloss.NewStyle().
		Width(width - 2).   // Account for left and right borders
		Height(height - 2). // Account for top and bottom borders
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	renderedBox := boxStyle.Render(content)

	// Format: ╭─ Title ──────────╮ (left-aligned, lazygit style)
	borderWidth := width - 2
	titleWithSpaces := " " + title + " "
	titleWidth := len(titleWithSpaces)

	// Calculate remaining space for dashes on the right
	remainingSpace := max(borderWidth-titleWidth, 0)

	// Define styles
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)
	titleStyle := lipgloss.NewStyle().Foreground(styles.ColorTextMain).Bold(true)

	header := borderStyle.Render("╭") +
		borderStyle.Render("─") +
		titleStyle.Render(title) +
		borderStyle.Render(" ") +
		borderStyle.Render(strings.Repeat("─", remainingSpace)) +
		borderStyle.Render("╮")

	return lipgloss.JoinVertical(lipgloss.Top, header, renderedBox)
}

// RenderPaneLoader renders a centered loading spinner with message
func RenderPaneLoader(spinnerView, message string, width, height int) string {
	spinnerStyle := lipgloss.NewStyle().Foreground(styles.ColorTextMain)
	messageStyle := styles.MetadataDim.MarginTop(1)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		spinnerStyle.Render(spinnerView),
		messageStyle.Render(message),
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}