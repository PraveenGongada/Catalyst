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

func (c *Client) CreatePackage(
	ctx context.Context,
	indexID string,
	fileIDs []string,
	tag string,
) (*PackageResult, error) {
	if indexID == "" {
		return nil, fmt.Errorf("ota: CreatePackage requires indexID")
	}
	if fileIDs == nil {
		fileIDs = []string{}
	}
	var out PackageResult
	err := c.request(
		ctx,
		"POST",
		"/api/packages",
		nil,
		createPackageRequest{Index: indexID, Files: fileIDs, Tag: tag},
		&out,
	)
	if err != nil {
		return nil, err
	}
	c.logger("created package tag=%s version=%d", out.Tag, out.Version)
	return &out, nil
}
