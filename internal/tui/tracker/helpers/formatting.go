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
	"fmt"
	"time"
)

// Truncate truncates a string to maxLen, adding "..." if truncated
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func FormatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}

	ist := time.FixedZone("IST", 5*60*60+30*60)
	tIST := t.In(ist)
	now := time.Now().In(ist)

	timeFormat := "03:04 PM"

	// Check if it's today
	if tIST.Year() == now.Year() && tIST.YearDay() == now.YearDay() {
		return fmt.Sprintf("%-6s %s", "Today", tIST.Format(timeFormat))
	}

	// Check if it's in the current year
	if tIST.Year() == now.Year() {
		return fmt.Sprintf("%-6s %s", tIST.Format("Jan 2"), tIST.Format(timeFormat))
	}

	// Different year - include year
	return fmt.Sprintf("%-6s %s", tIST.Format("Jan 2, 2006"), tIST.Format(timeFormat))
}

func FormatDuration(start, end time.Time) string {
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