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
	"fmt"
)

func (c *Client) Login(ctx context.Context, clientID, clientSecret string) error {
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("ota: client_id and client_secret are required")
	}
	var resp loginResponse
	err := c.request(
		ctx,
		"POST",
		"/api/token/issue",
		nil,
		loginRequest{ClientID: clientID, ClientSecret: clientSecret},
		&resp,
	)
	if err != nil {
		return err
	}
	if resp.AccessToken == "" {
		return fmt.Errorf("ota: login returned empty access_token")
	}
	c.token = resp.AccessToken
	c.logger("logged in to %s", c.base)
	return nil
}
