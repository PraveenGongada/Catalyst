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
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/PraveenGongada/catalyst/internal/util"
)

type otaMatricesFlags struct {
	platform    string
	apps        string
	environment string
}

type MatrixEntry struct {
	Variant   string `json:"variant"`
	Target    string `json:"target"` // alias of Variant for cross-workflow compatibility
	Platform  string `json:"platform"`
	Env       string `json:"env"`
	Namespace string `json:"ota_namespace"`
}

func newOtaMatricesCmd() *cobra.Command {
	f := &otaMatricesFlags{}
	cmd := &cobra.Command{
		Use:   "matrices",
		Short: "Print OTA push matrix entries as a JSON array (for GitHub Actions fan-out)",
		Long: `Emit the set of (variant, platform, ota_namespace) entries that
` + "`catalyst ota push`" + ` should be invoked for, based on catalyst.yaml.

Output is a JSON array suitable for feeding into a GitHub Actions matrix via
` + "`matrix.include: ${{ fromJson(...) }}`" + `. Each entry has these fields:

  variant        — app name (e.g. "Cumta")
  target         — alias for variant (kept for parity with other workflows)
  platform       — "android" or "ios"
  env            — catalyst.yaml block this entry came from (--env)
  ota_namespace  — from matrix.<Variant>.<Platform>.<Env>.matrix

Variants without an ota_namespace in the requested environment block
are skipped silently.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			cfg := ConfigFromContext(ctx)
			if cfg == nil {
				return fmt.Errorf("configuration was not loaded")
			}
			var platforms []string
			switch strings.ToLower(f.platform) {
			case "android":
				platforms = []string{"android"}
			case "ios":
				platforms = []string{"ios"}
			case "both":
				platforms = []string{"android", "ios"}
			default:
				return fmt.Errorf("invalid --platform %q (want android|ios|both)", f.platform)
			}
			allowed := util.ParseAllowedApps(f.apps)

			variants := make([]string, 0, len(cfg.Matrix))
			for v := range cfg.Matrix {
				variants = append(variants, v)
			}
			sort.Strings(variants)

			entries := make([]MatrixEntry, 0, len(variants)*len(platforms))
			for _, variant := range variants {
				if !util.IsAppAllowed(variant, allowed) {
					continue
				}
				for _, p := range platforms {
					ns := cfg.OTANamespace(variant, p, f.environment)
					if ns == "" {
						continue
					}
					entries = append(entries, MatrixEntry{
						Variant:   variant,
						Target:    variant,
						Platform:  strings.ToLower(p),
						Env:       f.environment,
						Namespace: ns,
					})
				}
			}
			if len(entries) == 0 {
				return fmt.Errorf("no matrix entries matched --apps %q / --platform %q / --env %q", f.apps, f.platform, f.environment)
			}
			data, err := json.Marshal(entries)
			if err != nil {
				return fmt.Errorf("encode matrices: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringVar(&f.platform, "platform", "both", "android | ios | both")
	fl.StringVar(&f.apps, "apps", "", "Comma-separated variant filter (empty or 'all' selects every variant)")
	fl.StringVar(&f.environment, "env", defaultOTAEnvironment, "catalyst.yaml matrix environment to read (Production | Debug)")
	return cmd
}
