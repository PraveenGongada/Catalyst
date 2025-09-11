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

package util

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"

	"github.com/PraveenGongada/catalyst/internal/types"
)

func GetStringFromSelection(items *[]types.SelectableItem) []string {
	itemStrings := []string{}
	for _, item := range *items {
		if item.Selected {
			itemStrings = append(itemStrings, item.Title)
		}
	}

	return itemStrings
}

func HasSelection(items []types.SelectableItem) bool {
	for _, item := range items {
		if item.Selected {
			return true
		}
	}
	return false
}

func NormalizeAppName(appName string) string {
	return strings.ToLower(strings.ReplaceAll(appName, " ", ""))
}

func ParseAllowedApps(appsFilter string) map[string]bool {
	if appsFilter == "" {
		return nil
	}

	allowedApps := make(map[string]bool)
	apps := strings.Split(appsFilter, ",")
	for _, app := range apps {
		trimmedApp := strings.TrimSpace(app)
		if trimmedApp != "" {
			normalizedApp := NormalizeAppName(trimmedApp)
			allowedApps[normalizedApp] = true
		}
	}
	return allowedApps
}

func IsAppAllowed(appName string, allowedApps map[string]bool) bool {
	if allowedApps == nil {
		return true
	}
	normalizedApp := NormalizeAppName(appName)
	return allowedApps[normalizedApp]
}

func AdditionalHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "select item"),
		),
		key.NewBinding(
			key.WithKeys("right", "enter", "ctrl+n"),
			key.WithHelp("→/ctrl+n/enter", "go next"),
		),
		key.NewBinding(
			key.WithKeys("left", "ctrl+p"),
			key.WithHelp("←/ctrl+p", "go back"),
		),
	}
}
