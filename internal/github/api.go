/*
 * Copyright 2025 Praveen Kumar
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

package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type APIClient struct {
	Repository string
}

func NewAPIClient(repository string) *APIClient {
	return &APIClient{
		Repository: repository,
	}
}

func (c *APIClient) ListWorkflows() (*WorkflowsResponse, error) {
	endpoint := fmt.Sprintf("/repos/%s/actions/workflows", c.Repository)

	cmd := exec.Command("gh", "api", endpoint)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	var response WorkflowsResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse workflows response: %w", err)
	}

	return &response, nil
}

func (c *APIClient) ListWorkflowRuns(workflowID int64, perPage, page int) (*RunsResponse, error) {
	endpoint := fmt.Sprintf("/repos/%s/actions/workflows/%d/runs?per_page=%d&page=%d",
		c.Repository, workflowID, perPage, page)

	cmd := exec.Command("gh", "api", endpoint)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list workflow runs: %w", err)
	}

	var response RunsResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse runs response: %w", err)
	}

	return &response, nil
}

func (c *APIClient) GetRunDetails(runID int64) (*RunDetails, error) {
	endpoint := fmt.Sprintf("/repos/%s/actions/runs/%d", c.Repository, runID)

	cmd := exec.Command("gh", "api", endpoint)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get run details: %w", err)
	}

	var details RunDetails
	if err := json.Unmarshal(output, &details); err != nil {
		return nil, fmt.Errorf("failed to parse run details: %w", err)
	}

	return &details, nil
}

func (c *APIClient) ListJobs(runID int64) (*JobsResponse, error) {
	endpoint := fmt.Sprintf("/repos/%s/actions/runs/%d/jobs", c.Repository, runID)

	cmd := exec.Command("gh", "api", endpoint)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	var response JobsResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse jobs response: %w", err)
	}

	return &response, nil
}
