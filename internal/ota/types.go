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

type loginRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type loginResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
}

type createApplicationRequest struct {
	Application string `json:"application"`
}

type ListPackagesResponse struct {
	Data       []PackageSummary `json:"data"`
	TotalPages int              `json:"total_pages"`
	TotalItems int              `json:"total_items"`
}

type PackageSummary struct {
	Tag     string   `json:"tag,omitempty"`
	Version int      `json:"version"`
	Index   string   `json:"index"`
	Files   []string `json:"files"`
}

type UploadResult struct {
	ID       string `json:"id"`
	FilePath string `json:"file_path"`
	URL      string `json:"url"`
	Version  int    `json:"version,omitempty"`
	Tag      string `json:"tag,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Status   string `json:"status,omitempty"`
}

type createPackageRequest struct {
	Index string   `json:"index"`
	Files []string `json:"files"`
	Tag   string   `json:"tag,omitempty"`
}

type PackageResult struct {
	Tag     string   `json:"tag,omitempty"`
	Version int      `json:"version"`
	Index   string   `json:"index"`
	Files   []string `json:"files"`
}

type apiError struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
