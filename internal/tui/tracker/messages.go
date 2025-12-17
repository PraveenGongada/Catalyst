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

package tracker

import (
	"time"

	"github.com/PraveenGongada/catalyst/internal/github"
)

// PollTickMsg is sent periodically to trigger polling
type PollTickMsg time.Time

// WorkflowsLoadedMsg is sent when workflows are loaded from API
type WorkflowsLoadedMsg struct {
	workflows []github.Workflow
	err       error
}

// RunsLoadedMsg is sent when runs are loaded from API
type RunsLoadedMsg struct {
	runs       []github.RunDetails
	totalCount int
	page       int
	err        error
}

// JobsLoadedMsg is sent when jobs are loaded from API
type JobsLoadedMsg struct {
	jobs []github.JobDetails
	err  error
}
