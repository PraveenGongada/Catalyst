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

package ota

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

// Returns the probe response when the app already existed, or nil when it
// was just created — so callers can skip a follow-up ListPackages.
func (c *Client) EnsureApplication(ctx context.Context, namespace string) (*ListPackagesResponse, error) {
	c.app = namespace

	probe, probeErr := c.ListPackages(ctx, 1)
	if probeErr == nil {
		c.logger("application %s exists", namespace)
		return probe, nil
	}

	var httpErr *HTTPError
	if errors.As(probeErr, &httpErr) && httpErr.StatusCode == 404 {
		c.logger("application %s not found; creating", namespace)
		if err := c.createApplication(ctx, namespace); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return nil, probeErr
}

func (c *Client) createApplication(ctx context.Context, namespace string) error {
	return c.request(
		ctx,
		"POST",
		"/api/organisations/applications/create",
		nil,
		createApplicationRequest{Application: namespace},
		nil,
	)
}

func (c *Client) ListPackages(ctx context.Context, count int) (*ListPackagesResponse, error) {
	q := url.Values{}
	if count > 0 {
		q.Set("count", strconv.Itoa(count))
	}
	var resp ListPackagesResponse
	if err := c.request(ctx, "GET", "/api/packages/list", q, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
