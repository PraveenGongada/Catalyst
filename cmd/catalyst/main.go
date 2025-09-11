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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PraveenGongada/catalyst/internal/config"
	"github.com/PraveenGongada/catalyst/internal/extractor"
	"github.com/PraveenGongada/catalyst/internal/tui"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func main() {
	configPath := flag.String(
		"config",
		"",
		"Path to the configuration file (default: $CATALYST_CONFIG or ./catalyst.yaml)",
	)

	versionFlag := flag.Bool(
		"version",
		false,
		"Print version information and exit",
	)

	extractWorkflow := flag.String(
		"extract",
		"",
		"Extract matrices for the specified workflow key",
	)

	outputFormat := flag.String(
		"format",
		"json",
		"Output format for extracted matrices (json|yaml)",
	)

	appsFilter := flag.String(
		"apps",
		"",
		"Comma-separated list of app names to filter (e.g., App1,App2). If empty, all apps are included",
	)

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Catalyst %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
		os.Exit(0)
	}

	if *extractWorkflow != "" {
		if err := handleExtractCommand(*configPath, *extractWorkflow, *outputFormat, *appsFilter); err != nil {
			fmt.Fprintf(os.Stderr, "Error extracting matrices: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := tui.Start(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting Catalyst: %v\n", err)
		os.Exit(1)
	}
}

func handleExtractCommand(configPath, workflowKey, format, appsFilter string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if format != "json" && format != "yaml" {
		return fmt.Errorf("invalid format '%s'. Supported formats: json, yaml", format)
	}

	if _, exists := cfg.GitHub.Workflows[workflowKey]; !exists {
		return fmt.Errorf("workflow '%s' not found in configuration. Available workflows: %v",
			workflowKey, getAvailableWorkflows(cfg))
	}

	extractedMatrices, err := extractor.ExtractWorkflowMatrices(cfg, workflowKey, appsFilter)
	if err != nil {
		return fmt.Errorf("error extracting matrices: %w", err)
	}

	if len(extractedMatrices) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: No matrices found for workflow '%s'\n", workflowKey)
	}

	output, err := extractor.FormatOutput(extractedMatrices, format)
	if err != nil {
		return fmt.Errorf("error formatting output: %w", err)
	}

	fmt.Print(output)
	return nil
}

func getAvailableWorkflows(cfg *config.Config) []string {
	var workflows []string
	for key := range cfg.GitHub.Workflows {
		workflows = append(workflows, key)
	}
	return workflows
}
