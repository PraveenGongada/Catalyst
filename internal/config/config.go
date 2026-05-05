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

package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/PraveenGongada/catalyst/internal/constants"
)

type Config struct {
	GitHub GitHubConfig                         `yaml:"github"`
	Inputs map[string]InputConfig               `yaml:"inputs"`
	Matrix map[string]map[string]PlatformConfig `yaml:"matrix"`
}

type GitHubConfig struct {
	Repository string                    `yaml:"repository"`
	Workflows  map[string]WorkflowConfig `yaml:"workflows"`
}

type WorkflowConfig struct {
	Name string `yaml:"name"`
	File string `yaml:"file"`
}

type InputConfig struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default"`
}

type PlatformConfig map[string]EnvironmentConfig

type EnvironmentConfig struct {
	Workflow string         `yaml:"workflow"`
	Matrix   map[string]any `yaml:"matrix"`
}

func Load(path string) (*Config, error) {
	if path == "" {
		if envPath := os.Getenv("CATALYST_CONFIG"); envPath != "" {
			path = envPath
		} else {
			path = "catalyst.yaml"
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if gitRoot := getGitRoot(); gitRoot != "" {
					gitRootPath := filepath.Join(gitRoot, "catalyst.yaml")
					if _, err := os.Stat(gitRootPath); err == nil {
						path = gitRootPath
					}
				}
			}
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func getGitRoot() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (c *Config) Validate() error {
	if c.GitHub.Repository == "" {
		return fmt.Errorf("GitHub repository is required")
	}
	if len(c.GitHub.Workflows) == 0 {
		return fmt.Errorf("at least one workflow is required")
	}

	if len(c.Matrix) == 0 {
		return fmt.Errorf("at least one app is required in matrix")
	}

	for app, platforms := range c.Matrix {
		if len(platforms) == 0 {
			return fmt.Errorf("app %s has no platforms", app)
		}

		for platform, environments := range platforms {
			if len(environments) == 0 {
				return fmt.Errorf("app %s platform %s has no environments", app, platform)
			}

			for env, config := range environments {
				if config.Workflow == "" {
					return fmt.Errorf(
						"app %s platform %s environment %s has no workflow",
						app,
						platform,
						env,
					)
				}

				if _, ok := c.GitHub.Workflows[config.Workflow]; !ok {
					return fmt.Errorf(
						"app %s platform %s environment %s references unknown workflow %s",
						app,
						platform,
						env,
						config.Workflow,
					)
				}
			}
		}
	}

	return nil
}

func (c *Config) GetApps() []string {
	apps := make([]string, 0, len(c.Matrix))
	for app := range c.Matrix {
		apps = append(apps, app)
	}
	return apps
}

func (c *Config) GetPlatforms(apps []string) []string {
	platformSet := make(map[string]bool)

	for _, app := range apps {
		if appConfig, ok := c.Matrix[app]; ok {
			for platform := range appConfig {
				platformSet[platform] = true
			}
		}
	}

	platforms := make([]string, 0, len(platformSet))
	for platform := range platformSet {
		platforms = append(platforms, platform)
	}

	return platforms
}

func (c *Config) GetEnvironments(apps []string, platforms []string) []string {
	envSet := make(map[string]bool)
	platformMap := make(map[string]string)

	for _, app := range apps {
		if appConfig, ok := c.Matrix[app]; ok {
			for platform := range appConfig {
				platformMap[strings.ToLower(platform)] = platform
			}
		}
	}

	for _, app := range apps {
		if appConfig, ok := c.Matrix[app]; ok {
			for _, selectedPlatform := range platforms {
				configPlatform := selectedPlatform
				if actualPlatform, ok := platformMap[strings.ToLower(selectedPlatform)]; ok {
					configPlatform = actualPlatform
				}

				if platformConfig, ok := appConfig[configPlatform]; ok {
					for env := range platformConfig {
						envSet[env] = true
					}
				}
			}
		}
	}

	environments := make([]string, 0, len(envSet))
	for env := range envSet {
		environments = append(environments, env)
	}

	return environments
}

var variablePattern = regexp.MustCompile(constants.RegexInputPlaceholder)

// Variant and platform lookups are case-insensitive; environment is case-sensitive.
func (c *Config) OTANamespace(variant, platform, environment string) string {
	appCfg, ok := lookupCI(c.Matrix, variant)
	if !ok {
		return ""
	}
	platCfg, ok := lookupCI(appCfg, platform)
	if !ok {
		return ""
	}
	envCfg, ok := platCfg[environment]
	if !ok {
		return ""
	}
	raw, ok := envCfg.Matrix["ota_namespace"]
	if !ok {
		return ""
	}
	if s, ok := raw.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", raw)
}

func lookupCI[V any](m map[string]V, key string) (V, bool) {
	if v, ok := m[key]; ok {
		return v, true
	}
	lk := strings.ToLower(key)
	for k, v := range m {
		if strings.ToLower(k) == lk {
			return v, true
		}
	}
	var zero V
	return zero, false
}

func (c *Config) SubstituteVariables(value string, inputValues map[string]string) string {
	return variablePattern.ReplaceAllStringFunc(value, func(match string) string {
		matches := variablePattern.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}

		varName := matches[1]

		if val, ok := inputValues[varName]; ok {
			return val
		}

		if input, ok := c.Inputs[varName]; ok {
			return input.Default
		}

		return match
	})
}
