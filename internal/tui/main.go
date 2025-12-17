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

package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/PraveenGongada/catalyst/internal/config"
	"github.com/PraveenGongada/catalyst/internal/github"
	"github.com/PraveenGongada/catalyst/internal/styles"
	"github.com/PraveenGongada/catalyst/internal/tui/tracker"
	"github.com/PraveenGongada/catalyst/internal/types"
)

type MainModel struct {
	subModels    map[types.Stage]tea.Model
	currentStage types.Stage
	width        int
	height       int
	config       *config.Config

	selectedApps         []string
	selectedPlatforms    []string
	selectedEnvironments []string
}

func NewMainModel(cfg *config.Config) *MainModel {
	model := &MainModel{
		currentStage:         types.AppSelectStage,
		height:               styles.DefaultHeight,
		width:                styles.DefaultWidth,
		config:               cfg,
		selectedApps:         []string{},
		selectedPlatforms:    []string{},
		selectedEnvironments: []string{},
	}
	model.subModels = getAllModels(model)
	return model
}

func (m *MainModel) Init() tea.Cmd {
	return m.subModels[m.currentStage].Init()
}

func (m *MainModel) View() string {
	return m.subModels[m.currentStage].View()
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.AppStyle.GetFrameSize()
		m.width = msg.Width - h
		m.height = msg.Height - v
		m.subModels = getAllModels(m)
		return m, nil
	}
	model := m.subModels[m.currentStage]
	var cmd tea.Cmd
	model, cmd = model.Update(msg)
	return model, cmd
}

func (m *MainModel) moveToNextStage() {
	if int(m.currentStage) < len(m.subModels)-1 {
		m.currentStage++
		model := m.subModels[m.currentStage]
		model.Init()
	}
}

func (m *MainModel) moveToPreviousStage() {
	if int(m.currentStage) > 0 {
		m.currentStage--
		model := m.subModels[m.currentStage]
		model.Init()
	}
}

func (m *MainModel) SetSelectedApps(apps []string) {
	m.selectedApps = apps
}

func (m *MainModel) SetSelectedPlatforms(platforms []string) {
	m.selectedPlatforms = platforms
}

func (m *MainModel) SetSelectedEnvironments(envs []string) {
	m.selectedEnvironments = envs
}

func (m *MainModel) GetSelectedApps() []string {
	return m.selectedApps
}

func (m *MainModel) GetSelectedPlatforms() []string {
	return m.selectedPlatforms
}

func (m *MainModel) GetSelectedEnvironments() []string {
	return m.selectedEnvironments
}

func (m *MainModel) UpdateSelectionsFromModel(
	stage types.Stage,
	selections []types.SelectableItem,
) {
	var selected []string
	for _, item := range selections {
		if item.Selected {
			selected = append(selected, item.Title)
		}
	}

	switch stage {
	case types.AppSelectStage:
		m.selectedApps = selected
	case types.OSSelectStage:
		m.selectedPlatforms = selected
	case types.EnvSelectStage:
		m.selectedEnvironments = selected
	}
}

func getAllModels(m *MainModel) map[types.Stage]tea.Model {
	return map[types.Stage]tea.Model{
		types.AppSelectStage: NewAppSelectModel(m),
		types.OSSelectStage:  NewOSSelectModel(m),
		types.EnvSelectStage: NewEnvSelectModel(m),
		types.InputStage:     NewInputsModel(m),
		types.ConfirmStage:   NewConfirmModel(m),
	}
}

func Start(configPath string, trackFlag bool) error {
	if err := github.IsGHInstalled(); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if trackFlag {
		trackerModel := tracker.NewWorkflowTrackerModel(cfg.GitHub.Repository)
		_, err = tea.NewProgram(trackerModel, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	} else {
		mainModel := NewMainModel(cfg)
		_, err = tea.NewProgram(mainModel).Run()
	}
	return err
}
