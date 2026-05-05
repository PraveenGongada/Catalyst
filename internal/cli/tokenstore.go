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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PraveenGongada/catalyst/internal/ota"
)

type storedToken struct {
	AccessToken  string `json:"access_token"`
	BaseURL      string `json:"base_url,omitempty"`
	Organisation string `json:"organisation,omitempty"`
	SavedAt      string `json:"saved_at"`
}

var errTokenNotFound = errors.New("ota token not found; run `catalyst ota login` first")

func tokenFilePath() (string, error) {
	if isCI() {
		return "/tmp/catalyst_ota_tokens.json", nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get cwd: %w", err)
	}
	return filepath.Join(cwd, ".catalyst", "credentials.json"), nil
}

func isCI() bool {
	v := strings.TrimSpace(os.Getenv("CI"))
	if v == "" {
		return false
	}
	switch strings.ToLower(v) {
	case "0", "false", "no":
		return false
	}
	return true
}

func saveToken(t storedToken) (string, error) {
	path, err := tokenFilePath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	t.SavedAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode token: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	if !isCI() {
		ensureGitignoreEntry(".catalyst")
	}
	return path, nil
}

func newAuthedClient(opts ...ota.Option) (*ota.Client, error) {
	tok, err := loadStoredToken()
	if err != nil {
		if !errors.Is(err, errTokenNotFound) {
			return nil, err
		}
		tok, err = autoLoginAndCache()
		if err != nil {
			return nil, err
		}
	}
	if tok.BaseURL == "" || tok.Organisation == "" {
		return nil, fmt.Errorf("cached token is missing base_url/organisation; re-run `catalyst ota login`")
	}
	reauth := ota.WithReAuth(func(_ context.Context) (string, error) {
		if path, pathErr := tokenFilePath(); pathErr == nil {
			_ = os.Remove(path)
		}
		fresh, err := autoLoginAndCache()
		if err != nil {
			return "", err
		}
		return fresh.AccessToken, nil
	})
	clientOpts := append([]ota.Option{ota.WithToken(tok.AccessToken), reauth}, opts...)
	return ota.New(tok.BaseURL, tok.Organisation, clientOpts...), nil
}

func autoLoginAndCache() (*storedToken, error) {
	baseURL := os.Getenv(envAirborneBaseURL)
	organisation := os.Getenv(envAirborneOrganisation)
	clientID := os.Getenv(envAirborneClientID)
	clientSecret := os.Getenv(envAirborneClientSecret)
	if baseURL == "" || organisation == "" || clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("no cached token; run `catalyst ota login` or set %s, %s, %s, %s",
			envAirborneBaseURL, envAirborneOrganisation, envAirborneClientID, envAirborneClientSecret)
	}
	client := ota.New(baseURL, organisation)
	if err := client.Login(context.Background(), clientID, clientSecret); err != nil {
		return nil, fmt.Errorf("auto-login: %w", err)
	}
	tok := storedToken{
		AccessToken:  client.Token(),
		BaseURL:      baseURL,
		Organisation: organisation,
	}
	if _, err := saveToken(tok); err != nil {
		return nil, fmt.Errorf("cache token: %w", err)
	}
	return &tok, nil
}

func loadStoredToken() (*storedToken, error) {
	path, err := tokenFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errTokenNotFound
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var t storedToken
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if t.AccessToken == "" {
		return nil, fmt.Errorf("token file %s contains empty access_token", path)
	}
	return &t, nil
}

// Best-effort; callers ignore the result so a missing or unwritable .gitignore never aborts login.
func ensureGitignoreEntry(entry string) {
	if _, err := os.Stat(".git"); err != nil {
		return
	}
	const name = ".gitignore"
	existing, err := os.ReadFile(name)
	content := string(existing)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return
	}
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == entry {
			return
		}
	}
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += entry + "\n"
	_ = os.WriteFile(name, []byte(content), 0o644)
}
