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

package constants

import util "github.com/PraveenGongada/catalyst/internal/util"

const (
	RegexInputPlaceholder = `{{inputs\.([^}]+)}}`
	RunIdWidth            = 9
	RunStatusWidth        = 12
	RunTimeWidth          = 16
	RunDurationWidth      = 12
	RunUserWidth          = 22
	RunBranchWidth        = 34
	JobNameWidth          = 24
	JobDurationWidth      = 12
	JobStatusWidth        = 12
)

var (
	RunItemsList        = []int{RunIdWidth, RunStatusWidth, RunTimeWidth, RunDurationWidth, RunUserWidth, RunBranchWidth}
	RunItemsWidth       = util.Sum(RunItemsList)
	RunItemsWidthFactor = len(RunItemsList) - 1 // -1 because we don't want to add padding to the first item
)
