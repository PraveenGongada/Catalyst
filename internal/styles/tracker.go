/*
 * Copyright 2025 Praveen Kumar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package styles

import "github.com/charmbracelet/lipgloss"

const (
	TrackerDefaultHeight = 30
	TrackerDefaultWidth  = 100
)

var (
	// Theme colors
	ColorBorder         = lipgloss.Color("#bd93f9") // Bright Purple
	ColorBorderInactive = lipgloss.Color("#b2b2b2") // Light grey
	ColorTextMain       = lipgloss.Color("#ffffff") // White text
	ColorTextMeta       = lipgloss.Color("#888888") // Gray metadata text
	ColorHighlight      = lipgloss.Color("#bd93f9") // Purple accent for cursors/text
	ColorSelectionBg    = lipgloss.Color("#3d3d3d") // 50% opacity white background

	// Status colors
	ColorSuccess        = lipgloss.Color("#3FB950") // Green
	ColorFailure        = lipgloss.Color("#F85149") // Red
	ColorInProgress     = lipgloss.Color("#F1E05A") // Yellow
	ColorQueued         = lipgloss.Color("#58A6FF") // Blue
	ColorCancelled      = lipgloss.Color("#8B949E") // Gray
	ColorSkipped        = lipgloss.Color("#6E7681") // Dark Gray
	ColorTimedOut       = lipgloss.Color("#D29922") // Orange
	ColorActionRequired = lipgloss.Color("#D29922") // Orange
	ColorNeutral        = lipgloss.Color("#8B949E") // Gray
)

// Status symbols
const (
	SymbolSuccess        = "" // Checkmark
	SymbolFailure        = "" // X mark
	SymbolInProgress     = "" // Spinner
	SymbolQueued         = "" // Clock
	SymbolWaiting        = "" // Pause
	SymbolCancelled      = "" // Stop sign
	SymbolSkipped        = "" // Skip forward
	SymbolTimedOut       = "" // Stopwatch
	SymbolActionRequired = "⚠" // Warning
	SymbolNeutral        = "" // Question mark

	// UniCode symbols for steps symbols
	StepsSymbolSuccess   = "✔"
	StepSymbolFailure    = "✖"
	StepSymbolCancelled  = "✖"
	StepSymbolSkipped    = "⊘"
	StepSymbolInProgress = "●"
	StepSymbolQueued     = "○"
	StepSymbolWaiting    = "⏸"
	StepSymbolNeutral    = ""
)

var (
	TrackerHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorTextMain).
			MarginBottom(1)

	TableHeader = lipgloss.NewStyle().
			Foreground(ColorTextMain).
			Bold(true)

	TrackerSubheader = lipgloss.NewStyle().
				Foreground(ColorTextMeta).
				Italic(true).
				MarginLeft(1)

	PaneStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	ListItemHighlighted = lipgloss.NewStyle().
				Foreground(ColorBorder).
				Bold(true)

	ListItemSelectedText = lipgloss.NewStyle().
				Foreground(ColorBorder).
				Bold(true)

	ListItemSelected = lipgloss.NewStyle().
				Foreground(ColorTextMain).
				Background(ColorSelectionBg).
				Bold(true)

	ListItemUnselected = lipgloss.NewStyle().
				Foreground(ColorTextMeta)

	MetadataDim = lipgloss.NewStyle().
			Foreground(ColorTextMeta)

	RunHeaderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			BorderBottom(true).
			MarginBottom(0).
			Padding(0, 1)

	BadgeSuccess        = lipgloss.NewStyle().Foreground(ColorSuccess)
	BadgeFailure        = lipgloss.NewStyle().Foreground(ColorFailure)
	BadgeCancelled      = lipgloss.NewStyle().Foreground(ColorCancelled)
	BadgeSkipped        = lipgloss.NewStyle().Foreground(ColorSkipped)
	BadgeTimedOut       = lipgloss.NewStyle().Foreground(ColorTimedOut)
	BadgeActionRequired = lipgloss.NewStyle().Foreground(ColorActionRequired)
	BadgeNeutral        = lipgloss.NewStyle().Foreground(ColorNeutral)
	BadgeInProgress     = lipgloss.NewStyle().Foreground(ColorInProgress)
	BadgeQueued         = lipgloss.NewStyle().Foreground(ColorQueued)
	BadgeMetadata       = lipgloss.NewStyle().Foreground(ColorTextMeta)

	HelpText = lipgloss.NewStyle().
			Foreground(ColorTextMeta).
			Align(lipgloss.Center).
			PaddingTop(1)

	TimelineConnector = lipgloss.NewStyle().
				Foreground(ColorSkipped)
)

func GetStatusStyle(status, conclusion string) lipgloss.Style {
	if status == "completed" {
		switch conclusion {
		case "success":
			return BadgeSuccess
		case "failure":
			return BadgeFailure
		case "cancelled":
			return BadgeCancelled
		case "skipped":
			return BadgeSkipped
		case "timed_out":
			return BadgeTimedOut
		case "action_required":
			return BadgeActionRequired
		case "neutral":
			return BadgeNeutral
		default:
			return BadgeSkipped
		}
	}

	switch status {
	case "in_progress":
		return BadgeInProgress
	case "queued", "pending", "requested":
		return BadgeQueued
	case "waiting":
		return BadgeQueued
	default:
		return BadgeMetadata
	}
}

func GetStatusSymbol(status, conclusion string) string {
	if status == "completed" {
		switch conclusion {
		case "success":
			return SymbolSuccess
		case "failure":
			return SymbolFailure
		case "cancelled":
			return SymbolCancelled
		case "skipped":
			return SymbolSkipped
		case "timed_out":
			return SymbolTimedOut
		case "action_required":
			return SymbolActionRequired
		case "neutral":
			return SymbolNeutral
		default:
			return SymbolNeutral
		}
	}

	switch status {
	case "in_progress":
		return SymbolInProgress
	case "queued", "pending", "requested":
		return SymbolQueued
	case "waiting":
		return SymbolWaiting
	default:
		return SymbolNeutral
	}
}

func GetStepStatusSymbol(status, conclusion string) string {
	if status == "completed" {
		switch conclusion {
		case "success":
			return StepsSymbolSuccess
		case "failure":
			return StepSymbolFailure
		case "cancelled":
			return StepSymbolCancelled
		case "skipped":
			return StepSymbolSkipped
		case "timed_out":
			return StepSymbolSkipped
		case "neutral":
			return StepSymbolNeutral
		default:
			return StepSymbolNeutral
		}
	}

	switch status {
	case "in_progress":
		return StepSymbolInProgress
	case "queued", "pending", "requested":
		return StepSymbolQueued
	case "waiting":
		return StepSymbolWaiting
	default:
		return StepSymbolNeutral
	}
}

func GetStatusText(status, conclusion string) string {
	if status == "completed" {
		switch conclusion {
		case "success":
			return "SUCCESS"
		case "failure":
			return "FAILURE"
		case "cancelled":
			return "CANCELLED"
		case "skipped":
			return "SKIPPED"
		case "timed_out":
			return "TIMED OUT"
		case "action_required":
			return "ACTION REQUIRED"
		case "neutral":
			return "NEUTRAL"
		default:
			return conclusion
		}
	}

	switch status {
	case "in_progress":
		return "RUNNING"
	case "queued", "pending", "requested":
		return "QUEUED"
	case "waiting":
		return "WAITING"
	default:
		return status
	}
}

func GetColoredStatusSymbol(status, conclusion string) string {
	symbol := GetStatusSymbol(status, conclusion)
	style := GetStatusStyle(status, conclusion)
	return style.Render(symbol)
}

func GetStepColoredSymbol(status, conclusion string) string {
	if status == "in_progress" {
		return lipgloss.NewStyle().Foreground(ColorInProgress).Render(SymbolInProgress)
	}
	if status == "queued" {
		return lipgloss.NewStyle().Foreground(ColorQueued).Render(SymbolQueued)
	}

	switch conclusion {
	case "failure":
		return lipgloss.NewStyle().Foreground(ColorFailure).Render(SymbolFailure)
	case "skipped":
		return lipgloss.NewStyle().Foreground(ColorSkipped).Render(SymbolSkipped)
	case "cancelled":
		return lipgloss.NewStyle().Foreground(ColorCancelled).Render(SymbolCancelled)
	case "success":
		return lipgloss.NewStyle().Foreground(ColorSuccess).Render(SymbolSuccess)
	default:
		return lipgloss.NewStyle().Foreground(ColorNeutral).Render(SymbolNeutral)
	}
}

// GetBorderColor returns the appropriate border color based on focus state
func GetBorderColor(isFocused bool) lipgloss.Color {
	if isFocused {
		return ColorBorder
	}
	return ColorBorderInactive
}
