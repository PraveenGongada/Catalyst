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
	"os"

	"github.com/spf13/cobra"

	"github.com/PraveenGongada/catalyst/internal/ota"
)

const (
	envAirborneClientID     = "AIRBORNE_CLIENT_ID"
	envAirborneClientSecret = "AIRBORNE_CLIENT_SECRET"
	envAirborneBaseURL      = "AIRBORNE_BASE_URL"
	envAirborneOrganisation = "AIRBORNE_ORGANISATION"
)

type otaLoginFlags struct {
	clientID     string
	clientSecret string
	baseURL      string
	organisation string
}

func newOtaLoginCmd() *cobra.Command {
	f := &otaLoginFlags{}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to the OTA server and cache the token for subsequent commands",
		Long: `Logs in using AIRBORNE_CLIENT_ID / AIRBORNE_CLIENT_SECRET (or
--client-id / --client-secret) and caches the access token so subsequent ` +
			"`catalyst ota push` and `catalyst ota latest-tag`" + ` invocations
can run without re-authenticating.

Cache location:
  CI=true  → /tmp/catalyst_ota_tokens.json
  else     → <cwd>/.catalyst/credentials.json

Run this once per CI job; every other ota command reads from it transparently.

This command does not read catalyst.yaml; pass server coordinates via
--base-url / --organisation or the matching AIRBORNE_* environment variables.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			baseURL, err := resolveFlag("--base-url", f.baseURL, envAirborneBaseURL)
			if err != nil {
				return err
			}
			organisation, err := resolveFlag("--organisation", f.organisation, envAirborneOrganisation)
			if err != nil {
				return err
			}
			clientID, err := resolveFlag("--client-id", f.clientID, envAirborneClientID)
			if err != nil {
				return err
			}
			clientSecret, err := resolveFlag("--client-secret", f.clientSecret, envAirborneClientSecret)
			if err != nil {
				return err
			}
			client := ota.New(baseURL, organisation)
			if err := client.Login(ctx, clientID, clientSecret); err != nil {
				return fmt.Errorf("login: %w", err)
			}
			path, err := saveToken(storedToken{
				AccessToken:  client.Token(),
				BaseURL:      baseURL,
				Organisation: organisation,
			})
			if err != nil {
				return fmt.Errorf("cache token: %w", err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "✅ logged in; token cached at %s\n", path)
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringVar(&f.clientID, "client-id", "", "Airborne client id (falls back to $AIRBORNE_CLIENT_ID)")
	fl.StringVar(&f.clientSecret, "client-secret", "", "Airborne client secret (falls back to $AIRBORNE_CLIENT_SECRET)")
	fl.StringVar(&f.baseURL, "base-url", "", "Airborne API base URL (falls back to $AIRBORNE_BASE_URL)")
	fl.StringVar(&f.organisation, "organisation", "", "Airborne organisation slug (falls back to $AIRBORNE_ORGANISATION)")
	return cmd
}

func resolveFlag(flagName, explicit, envVar string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if v := os.Getenv(envVar); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("%s is required (pass it directly or set $%s)", flagName, envVar)
}
