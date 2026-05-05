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
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/PraveenGongada/catalyst/internal/config"
)

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

func Execute(version, commit, date string) {
	buildVersion, buildCommit, buildDate = version, commit, date

	rewriteLegacyArgs()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := NewRootCmd().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	var cfgPath string

	root := &cobra.Command{
		Use:   "catalyst",
		Short: "Matrix + OTA orchestration for React Native monorepos",
		Long: "Catalyst coordinates build matrix dispatches and OTA pushes across many\n" +
			"React Native app variants from a single catalyst.yaml.",
		SilenceUsage:  true,
		SilenceErrors: false,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if skipConfigLoad(cmd) {
				return nil
			}
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("error loading configuration: %w", err)
			}
			cmd.SetContext(WithConfig(cmd.Context(), cfg))
			return nil
		},
		RunE: runTUI,
	}
	root.PersistentFlags().
		StringVar(&cfgPath, "config", "", "Path to catalyst.yaml (default: $CATALYST_CONFIG or ./catalyst.yaml)")

	root.AddCommand(newTUICmd(), newExtractCmd(), newOtaCmd(), newVersionCmd())
	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf(
				"Catalyst %s (commit: %s, built: %s)\n",
				buildVersion,
				buildCommit,
				buildDate,
			)
		},
	}
}

func skipConfigLoad(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "version", "help", "completion", "__complete":
		return true
	case "login", "push", "latest-tag":
		return true
	}
	for c := cmd; c != nil; c = c.Parent() {
		if c.Name() == "completion" {
			return true
		}
	}
	return false
}

func rewriteLegacyArgs() {
	if len(os.Args) < 2 {
		return
	}

	if os.Args[1] == "-help" {
		os.Args[1] = "--help"
		return
	}
	norm := func(s string) string {
		if strings.HasPrefix(s, "--") {
			return "-" + strings.TrimPrefix(s, "--")
		}
		return s
	}
	first := norm(os.Args[1])

	legacyPrefixes := []string{"-extract", "-format", "-apps", "-version", "-config"}
	matched := false
	for _, p := range legacyPrefixes {
		if first == p || strings.HasPrefix(first, p+"=") {
			matched = true
			break
		}
	}
	if !matched {
		return
	}

	var (
		extract     string
		format      string
		apps        string
		configPath  string
		wantVersion bool
	)

	i := 1
	for i < len(os.Args) {
		arg := norm(os.Args[i])
		switch {
		case arg == "-version":
			wantVersion = true
			i++
		case arg == "-extract" && i+1 < len(os.Args):
			extract = os.Args[i+1]
			i += 2
		case strings.HasPrefix(arg, "-extract="):
			extract = strings.TrimPrefix(arg, "-extract=")
			i++
		case arg == "-format" && i+1 < len(os.Args):
			format = os.Args[i+1]
			i += 2
		case strings.HasPrefix(arg, "-format="):
			format = strings.TrimPrefix(arg, "-format=")
			i++
		case arg == "-apps" && i+1 < len(os.Args):
			apps = os.Args[i+1]
			i += 2
		case strings.HasPrefix(arg, "-apps="):
			apps = strings.TrimPrefix(arg, "-apps=")
			i++
		case arg == "-config" && i+1 < len(os.Args):
			configPath = os.Args[i+1]
			i += 2
		case strings.HasPrefix(arg, "-config="):
			configPath = strings.TrimPrefix(arg, "-config=")
			i++
		default:
			return
		}
	}

	var rewritten []string
	if wantVersion {
		rewritten = []string{"version"}
	} else if extract != "" {
		rewritten = []string{"extract", "--workflow", extract}
		if format != "" {
			rewritten = append(rewritten, "--format", format)
		}
		if apps != "" {
			rewritten = append(rewritten, "--apps", apps)
		}
	} else if configPath != "" {
		rewritten = []string{}
	} else {
		return
	}
	if configPath != "" {
		rewritten = append([]string{"--config", configPath}, rewritten...)
	}

	fmt.Fprintf(
		os.Stderr,
		"warning: legacy flag form detected; rewriting `%s` as `%s` (will be removed in a future release)\n",
		strings.Join(os.Args[1:], " "),
		strings.Join(rewritten, " "),
	)
	os.Args = append([]string{os.Args[0]}, rewritten...)
}
