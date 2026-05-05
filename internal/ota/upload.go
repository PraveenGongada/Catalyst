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
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (c *Client) Upload(
	ctx context.Context,
	localPath, remotePath, tag string,
) (*UploadResult, error) {
	if remotePath == "" {
		return nil, fmt.Errorf("ota: Upload requires remotePath")
	}
	f, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("ota: open %s: %w", localPath, err)
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("ota: stat %s: %w", localPath, err)
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, fmt.Errorf("ota: hash %s: %w", localPath, err)
	}
	checksum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("ota: rewind %s: %w", localPath, err)
	}

	q := url.Values{}
	q.Set("file_path", remotePath)
	if tag != "" {
		q.Set("tag", tag)
	}

	u := c.base + "/api/file/upload?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, "POST", u, io.NopCloser(f))
	if err != nil {
		return nil, fmt.Errorf("ota: new upload request: %w", err)
	}
	req.ContentLength = stat.Size()
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-checksum", checksum)
	req.GetBody = func() (io.ReadCloser, error) {
		if _, seekErr := f.Seek(0, io.SeekStart); seekErr != nil {
			return nil, fmt.Errorf("ota: rewind %s: %w", localPath, seekErr)
		}
		return io.NopCloser(f), nil
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out UploadResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("ota: decode upload response: %w", err)
	}
	if out.ID == "" || out.FilePath == "" {
		return nil, fmt.Errorf("ota: upload returned empty id/file_path")
	}
	c.logger("uploaded %s (id=%s, size=%d)", out.FilePath, out.ID, stat.Size())
	return &out, nil
}
