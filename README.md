# 🚀 Catalyst

<div align="center">
  
[![Release](https://img.shields.io/github/v/release/PraveenGongada/catalyst?style=flat-square)](https://github.com/PraveenGongada/catalyst/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/PraveenGongada/catalyst)](https://goreportcard.com/report/github.com/PraveenGongada/catalyst)
[![Go Version](https://img.shields.io/github/go-mod/go-version/PraveenGongada/catalyst?style=flat-square)](go.mod)
[![License](https://img.shields.io/github/license/PraveenGongada/catalyst?style=flat-square)](LICENSE)

</div>

Catalyst is an elegant terminal UI tool that simplifies triggering GitHub Actions workflows with matrix configurations. It's designed to streamline the deployment process for mobile applications across multiple platforms and environments.

![Catalyst Demo](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/catalyst-demo.gif)

## ✨ Features

- 📱 Select multiple apps for deployment
- 📲 Choose target platforms (iOS, Android)
- 🌿 Select environments (Debug, Production)
- 🔧 Configure dynamic input values
- 📝 Add changelog information for deployments
- 🔀 Specify target branch for workflow execution
- 🔍 Preview matrix configurations before triggering
- 🔄 Trigger GitHub Actions workflows with complex matrix configurations
- 🔧 Extract matrix configurations for external use (GitHub UI, CI/CD automation)
- 🌐 Maintain consistent configurations across different trigger methods
- 📤 Push React Native bundles to an Airborne OTA server straight from your matrix config

## 📦 Installation

### Prerequisites

- [GitHub CLI](https://cli.github.com/) (gh) - Required for triggering workflows
  - Install from: https://cli.github.com/manual/installation
  - Authenticate with: `gh auth login`

### GitHub Actions

Use the [setup-catalyst](https://github.com/PraveenGongada/setup-catalyst) action in your workflows to install Catalyst:

```yaml
- name: Setup Catalyst
  uses: PraveenGongada/setup-catalyst@v1
  with:
    version: "latest" # or specify a version like 'v1.0.0'
```

### Homebrew

```bash
brew tap PraveenGongada/tap
brew install catalyst
```

### Go

```bash
go install github.com/PraveenGongada/catalyst/cmd/catalyst@latest
```

### Linux

```bash
# Download the latest release (replace X.Y.Z with actual version)
curl -L https://github.com/PraveenGongada/catalyst/releases/download/vX.Y.Z/catalyst_Linux_x86_64.tar.gz -o catalyst.tar.gz

# Extract the binary
tar -xzf catalyst.tar.gz

# Move to a directory in your PATH
sudo mv catalyst /usr/local/bin/
```

### Manual Installation

Download the appropriate binary for your platform from the [releases page](https://github.com/PraveenGongada/catalyst/releases).

## 🔧 Configuration

Catalyst uses a YAML configuration file to define your matrix configurations. By default, it looks for `catalyst.yaml` in the current directory.

### Configuration File Location

Catalyst looks for the configuration file in the following order:

1. Path specified with the `-config` flag
2. Path set in the `CATALYST_CONFIG` environment variable
3. `catalyst.yaml` in the current or git root directory

### Adding to Your Shell Profile

You can optionally set the configuration path using an environment variable. Add the following to your `.zshrc`, `.bashrc`, or equivalent:

```bash
# Set the path to your catalyst config
export CATALYST_CONFIG="/path/to/your/catalyst.yaml"
```

### Sample Configuration

Here's a simplified example of a configuration file:

```yaml
# GitHub workflow metadata
github:
  repository: "your-org/mobile-apps"
  workflows:
    ios_debug:
      name: "iOS Debug Workflow"
      file: "ios_debug.yml"
    ios_prod:
      name: "iOS Production Workflow"
      file: "ios_prod.yml"

# Dynamic inputs that can be referenced in matrices
inputs:
  version:
    description: "App version"
    required: true
    default: "1.0.0"

# Matrix configurations
matrix:
  MyApp:
    ios:
      Debug:
        workflow: "ios_debug"
        matrix:
          bundle_id: "com.example.myapp.debug"
          version: "{{inputs.version}}"
      Prod:
        workflow: "ios_prod"
        matrix:
          bundle_id: "com.example.myapp"
          version: "{{inputs.version}}"
          ota_namespace: "myapp_ios" # optional, used by `catalyst ota`
```

## 🚀 Usage

Run Catalyst from your terminal:

```bash
catalyst
```

Or specify a custom configuration path:

```bash
catalyst -config /path/to/your/catalyst.yaml
```

Check the version:

```bash
catalyst -version
```

### Matrix Extraction

Catalyst provides a powerful matrix extraction feature that allows you to generate the exact matrix configurations without going through the interactive interface. This is particularly useful when you want to:

- **Use GitHub UI**: Extract matrices and manually trigger workflows through GitHub's web interface
- **CI/CD Integration**: Incorporate Catalyst configurations into automated pipelines
- **Debugging & Validation**: Inspect matrix configurations before deployment
- **Flexibility**: Maintain consistent configurations while choosing your preferred trigger method

Extract matrices for a specific workflow:

```bash
# Extract as JSON (default) - perfect for GitHub workflow inputs
catalyst -extract ios_dev

# Extract as YAML for better readability
catalyst -extract android_prod -format yaml

# Use with a custom config file
catalyst -config /path/to/config.yaml -extract ios_prod
```

This approach ensures that whether you use Catalyst's TUI, the extraction feature, or any other method, your deployments remain consistent and follow the same configuration patterns defined in your `catalyst.yaml` file.

### OTA Push (Airborne)

If you ship React Native apps and use [Airborne](https://github.com/juspay/airborne) for OTA updates, Catalyst can also fan out a bundle push across the same matrix you use for native builds. Add an `ota_namespace` to whichever matrix entries you want to push, then in CI:

```bash
# Emit only the entries that have an ota_namespace, as JSON for matrix.include
catalyst ota matrices --platform both --env Production

# Per matrix entry: log in once, then push the prebuilt bundle
catalyst ota login --base-url "$AIRBORNE_BASE_URL" --organisation "$AIRBORNE_ORG"
catalyst ota push --namespace "$NS" --platform ios --tag "$TAG" --bundle main.jsbundle
```

See `catalyst ota --help` for the full surface (`login`, `push`, `latest-tag`, `matrices`).

## 📋 GitHub Actions Workflow Setup

To use Catalyst with your GitHub Actions, you need to set up a repository dispatch event in your workflow:

```yaml
# .github/workflows/example.yml
name: Example Deployment Workflow

on:
  workflow_dispatch:
    inputs:
      payload:
        description: "JSON payload containing matrix configurations"
        required: true
      change_log:
        description: "Changelog for this deployment"
        required: true

jobs:
  deploy:
    runs-on: ubuntu-latest
    # Parse the matrices JSON from the payload
    strategy:
      matrix:
        include: ${{ fromJson(fromJson(inputs.payload).matrices) }}

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      # Add your deployment steps here
```

## 📸 Interface Screenshots

Here's a visual walkthrough of the Catalyst interface:

### App Selection

![App Selection Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/app-selection.png)

### Platform Selection

![Platform Selection Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/platform-selection.png)

### Environment Selection

![Environment Selection Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/env-selection.png)

### Input Configuration

![Input Configuration Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/input-configuration.png)

### Deployment Summary

![Deployment Summary Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/deployment-summary.png)

### Matrix Preview

![Matrix Preview Screen](https://raw.githubusercontent.com/PraveenGongada/catalyst/refs/heads/main/docs/images/matrix-preview.png)

## 👨‍💻 Development

### Prerequisites

- Go 1.22 or higher
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/praveengongada/catalyst.git
cd catalyst

# Build the binary (main is in cmd/catalyst)
go build -o catalyst ./cmd/catalyst

# Run the application
./catalyst
```

### Making a Release

Catalyst uses GoReleaser to automate the release process:

1. Update the version in your code
2. Create and push a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. The GitHub Actions workflow will automatically build and publish the release

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Third-party library attributions can be found in the [NOTICE](NOTICE) file.

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/PraveenGongada">Praveen Kumar</a></p>
  <p>
    <a href="https://linkedin.com/in/praveengongada">LinkedIn</a> •
    <a href="https://praveengongada.com">Website</a>
  </p>
</div>
