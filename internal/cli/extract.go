/*
 * Copyright 2026 Praveen Kumar
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

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PraveenGongada/catalyst/internal/extractor"
)

func newExtractCmd() *cobra.Command {
	var (
		workflowKey string
		format      string
		appsFilter  string
	)

	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract matrices for a workflow and print them to stdout",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg := ConfigFromContext(cmd.Context())
			if cfg == nil {
				return fmt.Errorf("configuration was not loaded")
			}
			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("invalid configuration: %w", err)
			}
			if format != "json" && format != "yaml" {
				return fmt.Errorf("invalid format %q. Supported: json, yaml", format)
			}
			if _, ok := cfg.GitHub.Workflows[workflowKey]; !ok {
				available := make([]string, 0, len(cfg.GitHub.Workflows))
				for k := range cfg.GitHub.Workflows {
					available = append(available, k)
				}
				return fmt.Errorf(
					"workflow %q not found. Available: %v",
					workflowKey,
					available,
				)
			}

			matrices, err := extractor.ExtractWorkflowMatrices(cfg, workflowKey, appsFilter)
			if err != nil {
				return fmt.Errorf("error extracting matrices: %w", err)
			}
			if len(matrices) == 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), "Warning: no matrices found for workflow %q\n", workflowKey)
			}
			out, err := extractor.FormatOutput(matrices, format)
			if err != nil {
				return fmt.Errorf("error formatting output: %w", err)
			}
			fmt.Print(out)
			return nil
		},
	}

	cmd.Flags().StringVarP(&workflowKey, "workflow", "w", "", "Workflow key defined under github.workflows in catalyst.yaml")
	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format (json|yaml)")
	cmd.Flags().StringVarP(&appsFilter, "apps", "a", "", "Comma-separated list of app names to include (default: all)")
	_ = cmd.MarkFlagRequired("workflow")
	return cmd
}
