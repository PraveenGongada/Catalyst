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
)

type latestTagFlags struct {
	namespace string
	fallback  string
}

func newOtaLatestTagCmd() *cobra.Command {
	f := &latestTagFlags{}
	cmd := &cobra.Command{
		Use:   "latest-tag",
		Short: "Print the most recent package tag for an OTA namespace (or a fallback if none)",
		Long: `Reads the latest package tag from the OTA server for the given
namespace. Useful in CI version-bump scripts.

Requires a prior ` + "`catalyst ota login`" + ` in the same session (the
cached token — along with the base URL and organisation — is read
automatically). This command does not touch catalyst.yaml; pass the
namespace directly (e.g. from ` + "`${{ matrix.ota_namespace }}`" + `).

If no package exists yet, the --fallback value (default "0.0.0") is
printed instead.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			client, err := newAuthedClient()
			if err != nil {
				return err
			}
			probe, err := client.EnsureApplication(ctx, f.namespace)
			if err != nil {
				return fmt.Errorf("ensure application %s: %w", f.namespace, err)
			}
			tag := f.fallback
			if probe != nil && len(probe.Data) > 0 && probe.Data[0].Tag != "" {
				tag = probe.Data[0].Tag
			}
			fmt.Fprintln(cmd.OutOrStdout(), tag)
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringVar(&f.namespace, "namespace", "", "OTA application namespace (required)")
	fl.StringVar(&f.fallback, "fallback", "0.0.0", "Value to print when no package exists yet")
	_ = cmd.MarkFlagRequired("namespace")
	return cmd
}
