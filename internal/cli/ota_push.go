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
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/PraveenGongada/catalyst/internal/ota"
)

type otaPushFlags struct {
	namespace string
	platform  string
	tag       string
	bundle    string
	assetsDir string
	zip       bool
}

func newOtaPushCmd() *cobra.Command {
	f := &otaPushFlags{}
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Upload a prebuilt RN bundle (and optionally assets) to the OTA server as a new package",
		Long: `Upload a prebuilt React Native bundle to the OTA server and register
it as a new package under the given namespace.

The CLI does not build the bundle — produce it yourself (for example with
` + "`npx react-native bundle` followed by optional `hermesc` / zip)" + ` and
pass the resulting file via --bundle. Any assets emitted alongside it can
be uploaded together by pointing --assets-dir at their directory; every
regular file under it will be uploaded with its relative path.

Authentication: run ` + "`catalyst ota login`" + ` earlier in the same
session. The cached token — along with the base URL and organisation — is
read automatically. This command does not touch catalyst.yaml; pass the
namespace directly (e.g. from ` + "`${{ matrix.ota_namespace }}`" + `
in a GitHub Actions matrix fan-out).`,
		RunE: func(cmd *cobra.Command, _ []string) error { return runOtaPush(cmd, f) },
	}
	fl := cmd.Flags()
	fl.StringVar(&f.namespace, "namespace", "", "OTA application namespace to push to (required)")
	fl.StringVar(&f.platform, "platform", "", "android | ios (required; picks the remote bundle filename)")
	fl.StringVar(&f.tag, "tag", "", "Semver tag for this package (required)")
	fl.StringVar(&f.bundle, "bundle", "", "Path to the prebuilt bundle file to upload as the package index (required)")
	fl.StringVar(&f.assetsDir, "assets-dir", "", "Optional directory of asset files to upload alongside the bundle")
	fl.BoolVar(&f.zip, "zip", false, "Zip --bundle before uploading (inner filename = remote bundle name)")
	_ = cmd.MarkFlagRequired("namespace")
	_ = cmd.MarkFlagRequired("platform")
	_ = cmd.MarkFlagRequired("tag")
	_ = cmd.MarkFlagRequired("bundle")
	return cmd
}

func runOtaPush(cmd *cobra.Command, f *otaPushFlags) error {
	ctx := cmd.Context()
	platform := strings.ToLower(f.platform)
	if platform != "android" && platform != "ios" {
		return fmt.Errorf("--platform must be android or ios (got %q)", f.platform)
	}
	if _, err := os.Stat(f.bundle); err != nil {
		return fmt.Errorf("--bundle: %w", err)
	}
	assets, err := collectAssets(f.assetsDir)
	if err != nil {
		return err
	}

	client, err := newAuthedClient(ota.WithLogger(func(format string, args ...any) {
		fmt.Fprintf(cmd.ErrOrStderr(), "[ota] "+format+"\n", args...)
	}))
	if err != nil {
		return err
	}
	if _, err := client.EnsureApplication(ctx, f.namespace); err != nil {
		return fmt.Errorf("ensure application %s: %w", f.namespace, err)
	}

	bundleRemote := bundleRemoteName(platform)
	uploadPath := f.bundle
	if f.zip {
		zipPath, err := zipBundleFile(f.bundle, bundleRemote)
		if err != nil {
			return fmt.Errorf("zip bundle: %w", err)
		}
		defer os.Remove(zipPath)
		uploadPath = zipPath
	}
	bundleUpload, err := client.Upload(ctx, uploadPath, bundleRemote, f.tag)
	if err != nil {
		return fmt.Errorf("upload bundle: %w", err)
	}

	assetIDs, err := uploadAssets(ctx, client, assets, f.tag)
	if err != nil {
		return err
	}

	pkg, err := client.CreatePackage(ctx, bundleUpload.ID, assetIDs, f.tag)
	if err != nil {
		return fmt.Errorf("create package: %w", err)
	}
	fmt.Fprintf(
		cmd.OutOrStdout(),
		"pushed namespace=%s platform=%s tag=%s version=%d bundle=%s assets=%d\n",
		f.namespace,
		platform,
		pkg.Tag,
		pkg.Version,
		bundleRemote,
		len(assets),
	)
	return nil
}

type asset struct {
	localPath  string
	remotePath string
}

func collectAssets(dir string) ([]asset, error) {
	if dir == "" {
		return nil, nil
	}
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("--assets-dir: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("--assets-dir must be a directory (got %s)", dir)
	}
	var out []asset
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(dir, path)
		if relErr != nil {
			return relErr
		}
		out = append(out, asset{localPath: path, remotePath: filepath.ToSlash(rel)})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk --assets-dir: %w", err)
	}
	return out, nil
}

func uploadAssets(
	ctx context.Context,
	client *ota.Client,
	assets []asset,
	tag string,
) ([]string, error) {
	ids := make([]string, 0, len(assets))
	for _, a := range assets {
		r, err := client.Upload(ctx, a.localPath, a.remotePath, tag)
		if err != nil {
			return nil, fmt.Errorf("upload asset %s: %w", a.remotePath, err)
		}
		ids = append(ids, r.ID)
	}
	return ids, nil
}

// bundleRemoteName is the server-side file_path for the bundle. Keeping this
// stable regardless of file format (JS, Hermes bytecode, zipped) preserves
// version continuity on the server and matches the name the APK-bundled
// release_config.json points at.
func bundleRemoteName(platform string) string {
	if platform == "ios" {
		return "main.jsbundle"
	}
	return "index.android.bundle"
}

func zipBundleFile(src, innerName string) (string, error) {
	tmp, err := os.CreateTemp("", "catalyst-ota-*.zip")
	if err != nil {
		return "", err
	}
	cleanup := func() { _ = tmp.Close(); _ = os.Remove(tmp.Name()) }

	w := zip.NewWriter(tmp)
	entry, err := w.Create(innerName)
	if err != nil {
		cleanup()
		return "", err
	}
	srcF, err := os.Open(src)
	if err != nil {
		cleanup()
		return "", err
	}
	defer srcF.Close()
	if _, err := io.Copy(entry, srcF); err != nil {
		cleanup()
		return "", err
	}
	if err := w.Close(); err != nil {
		cleanup()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}
