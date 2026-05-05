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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	base   string
	org    string
	app    string
	token  string
	http   *http.Client
	logger func(format string, args ...any)
	reauth func(ctx context.Context) (string, error)
}

type Option func(*Client)

func WithLogger(f func(format string, args ...any)) Option {
	return func(c *Client) { c.logger = f }
}

func WithToken(token string) Option {
	return func(c *Client) { c.token = token }
}

func WithReAuth(fn func(ctx context.Context) (string, error)) Option {
	return func(c *Client) { c.reauth = fn }
}

func (c *Client) Token() string { return c.token }

func New(baseURL, org string, opts ...Option) *Client {
	c := &Client{
		base:   strings.TrimRight(baseURL, "/"),
		org:    org,
		http:   &http.Client{Timeout: 10 * time.Minute},
		logger: func(string, ...any) {},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.doAttempt(req, true)
}

func (c *Client) doAttempt(req *http.Request, allowReauth bool) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.org != "" {
		req.Header.Set("x-organisation", c.org)
	}
	if c.app != "" {
		req.Header.Set("x-application", c.app)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ota: %s %s: %w", req.Method, req.URL.Path, err)
	}
	if resp.StatusCode == http.StatusUnauthorized && allowReauth && c.reauth != nil && replayable(req) {
		resp.Body.Close()
		newTok, reauthErr := c.reauth(req.Context())
		if reauthErr != nil {
			return nil, reauthErr
		}
		c.token = newTok
		if req.GetBody != nil {
			body, bodyErr := req.GetBody()
			if bodyErr != nil {
				return nil, fmt.Errorf("ota: replay body: %w", bodyErr)
			}
			req.Body = body
		}
		return c.doAttempt(req, false)
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, decodeError(resp, req.Method, req.URL.Path)
	}
	return resp, nil
}

func replayable(req *http.Request) bool {
	return req.Body == nil || req.Body == http.NoBody || req.GetBody != nil
}

func (c *Client) request(
	ctx context.Context,
	method, path string,
	query url.Values,
	body any,
	out any,
) error {
	u := c.base + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("ota: encode body: %w", err)
		}
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reader)
	if err != nil {
		return fmt.Errorf("ota: new request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("ota: decode %s %s: %w", method, path, err)
	}
	return nil
}

type HTTPError struct {
	Method     string
	Path       string
	StatusCode int
	Message    string
	Body       string
}

func (e *HTTPError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = e.Body
	}
	if msg == "" {
		msg = http.StatusText(e.StatusCode)
	}
	return fmt.Sprintf("ota: %s %s: %d %s", e.Method, e.Path, e.StatusCode, msg)
}

func decodeError(resp *http.Response, method, path string) error {
	const limit = 8 * 1024
	body, _ := io.ReadAll(io.LimitReader(resp.Body, limit))
	e := &HTTPError{
		Method:     method,
		Path:       path,
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}
	var env apiError
	if json.Unmarshal(body, &env) == nil {
		switch {
		case env.Message != "":
			e.Message = env.Message
		case env.Error != "":
			e.Message = env.Error
		}
	}
	return e
}
