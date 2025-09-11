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

package extractor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/PraveenGongada/catalyst/internal/config"
	"github.com/PraveenGongada/catalyst/internal/constants"
	"github.com/PraveenGongada/catalyst/internal/util"
)

type OutputFormat struct {
	Matrices []map[string]interface{} `json:"matrices" yaml:"matrices"`
}

var inputPlaceholderPattern = regexp.MustCompile(constants.RegexInputPlaceholder)

func ExtractWorkflowMatrices(
	cfg *config.Config,
	workflowKey string,
	appsFilter string,
) ([]map[string]interface{}, error) {
	var matrices []map[string]interface{}

	allowedApps := util.ParseAllowedApps(appsFilter)

	for appName, appConfig := range cfg.Matrix {
		if !util.IsAppAllowed(appName, allowedApps) {
			continue
		}

		for _, platformConfig := range appConfig {
			for _, envConfig := range platformConfig {
				if envConfig.Workflow == workflowKey {
					filteredMatrix := filterInputPlaceholders(envConfig.Matrix)

					matrices = append(matrices, filteredMatrix)
				}
			}
		}
	}

	return matrices, nil
}

func filterInputPlaceholders(matrix map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})

	for key, value := range matrix {
		if strValue, ok := value.(string); ok {
			if inputPlaceholderPattern.MatchString(strValue) {
				continue
			}
		}

		filtered[key] = value
	}

	return filtered
}

func FormatOutput(matrices []map[string]interface{}, format string) (string, error) {
	output := OutputFormat{
		Matrices: matrices,
	}

	switch strings.ToLower(format) {
	case "json":
		jsonBytes, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return "", fmt.Errorf("error marshaling to JSON: %w", err)
		}
		return string(jsonBytes), nil

	case "yaml":
		yamlBytes, err := yaml.Marshal(output)
		if err != nil {
			return "", fmt.Errorf("error marshaling to YAML: %w", err)
		}
		return string(yamlBytes), nil

	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
