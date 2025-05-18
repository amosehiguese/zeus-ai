<br/>

<div align="center">
  Like this project ? Leave us a star ⭐
</div>

<br/>

<div align="center">
  <a href="https://github.com/amosehiguese/zeus-ai" target="_blank">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="assets/zeus-ai-logo.png">
    <img src="https://raw.githubusercontent.com/amosehiguese/zeus-ai/main/assets/zeus-ai-logo.png" alt="zeus-ai logo">
  </picture>
  </a>
</div>

<br/>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/amosehiguese/zeus-ai">
    <img src="https://goreportcard.com/badge/github.com/amosehiguese/zeus-ai" alt="Report Card">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT">
  </a>
  <a href="https://github.com/amosehiguese/zeus-ai/releases">
    <img src="https://img.shields.io/github/release/amosehiguese/zeus-ai.svg" alt="releases">
  </a>
</p>

zeus-ai is a Git-aware CLI tool that helps developers generate smart commit messages using LLM APIs. It analyzes your Git diff, sends it to an LLM (such as DeepSeek, Claude, or any local model via Ollama), receives multiple commit message suggestions, and allows you to select and edit the message before committing.

## Zeus Preview
![Zeus Demo](assets/zeusv.gif)

## 🚀 Features

- **Multiple LLM Provider Support**: Ollama, OpenRouter
- **Git Integration**: Works with staged/unstaged changes, supports signed commits
- **Smart Suggestions**: Generate multiple commit message options based on your changes
- **Flexible Configuration**: Config file, environment variables, command flags
- **Terminal UI**: Simple interface for selecting and editing commit messages

## 📋 Installation

### Using Go

```bash
# Clone the repository
git clone https://github.com/amosehiguese/zeus-ai.git

# Navigate to the directory
cd zeus-ai

# Build the binary
go build -o zeusctl ./zeusctl or make build

# Move to path (optional)
sudo mv zeusctl /usr/local/bin/
```

### Manual Download

Download the binary for your platform from the [releases page](https://github.com/amosehiguese/zeus-ai/releases).

## 🔧 Configuration

### Quick Start

Initialize a configuration file with guided setup:

```bash
zeus-ai init
```

This will create a `.zeusrc` file in the current directory with default settings.

### Configuration File

zeus-ai looks for a `.zeusrc` file in the following locations, in order:
1. Current directory
2. Any parent directory (up to the root)
3. Home directory (`~/.zeusrc`)

Example `.zeusrc` file:

```yaml
# LLM Provider configuration
provider: openrouter        
api_key: your-api-key-here
model: mistralai/mistral-small-3.1-24b-instruct:free

# Default commit style
default_style: conventional  # Options: conventional, simple

# Optional settings
editor: vim            # Overrides $EDITOR environment variable
sign_by_default: true  # Always sign commits
auto_stage: false      # Don't automatically stage all changes
```

### Environment Variables

All settings can be configured with environment variables, which take precedence over the config file:

```bash
# Required settings
export ZEUS_PROVIDER=openrouter
export ZEUS_API_KEY=your-api-key-here
export ZEUS_MODEL=mistralai/mistral-small-3.1-24b-instruct:free

# Optional settings
export ZEUS_DEFAULT_STYLE=conventional
export ZEUS_EDITOR=vim
export ZEUS_SIGN_BY_DEFAULT=true
export ZEUS_AUTO_STAGE=false
```

### Provider-Specific Configuration

#### Ollama (Local Models)
```yaml
provider: ollama
model: mistral  # or any model you have pulled in Ollama
# No API key needed for local Ollama
```

#### OpenRouter
```yaml
provider: openrouter
api_key: your-openrouter-api-key
model: mistralai/mistral-small-3.1-24b-instruct:free  # or any other model supported by OpenRouter
```

## 💻 Usage

### Basic Command

```bash
zeus-ai suggest
```

This will:
1. Check for staged changes or ask to use unstaged changes
2. Send the diff to the configured LLM
3. Generate commit message suggestions
4. Display them for you to choose or edit
5. Create the commit with your selected message

### Command Options

```bash
# Include detailed body text in suggestions
zeus-ai suggest --body

# Open the selected message in default editor
zeus-ai suggest --edit

# GPG-sign the commit message once you've approved the commit message
zeus-ai suggest --sign

# Display message suggestions but don't run commit
zeus-ai suggest --dry-run

# Automatically stage all changes
zeus-ai suggest --auto-stage

# Specify commit style (conventional or simple)
zeus-ai suggest --style conventional
```

### Conventional Commit Format

When using `--style conventional`, suggestions will follow the format:

```
type(scope): description
```

Where:
- `type` is one of: feat, fix, docs, style, refactor, test, chore
- `scope` is optional and specifies the section of the codebase
- `description` is a concise explanation of the change

### Usage Examples

#### Generate a simple commit message
```bash
zeus-ai suggest
```

#### Create a detailed conventional commit message and edit before committing
```bash
zeus-ai suggest --body --style conventional --edit
```

#### Preview suggestions without committing
```bash
zeus-ai suggest --dry-run
```

#### Automatically stage all changes and sign the commit
```bash
zeus-ai suggest --auto-stage --sign
```

## 🔄 Command Flow

1. **Check Git Status**: zeus-ai verifies you're in a Git repository and checks for staged changes
2. **Collect Diff**: Gets the Git diff of staged changes (or unstaged if specified)
3. **Generate Suggestions**: Sends the diff to the LLM to generate commit message suggestions
4. **Present Options**: Shows the suggestions with a simple selection interface
5. **Edit (Optional)**: Opens your selected message in an editor if --edit is used
6. **Commit**: Creates the Git commit with your chosen message

## 🧩 Integration with Git Aliases

Add zeus-ai to your Git workflow by setting up a Git alias:

```bash
git config --global alias.ai 'zeus-ai suggest'
```

Then use:

```bash
git ai
```

## 🔍 Troubleshooting

### Common Issues

#### API Key Authentication Error
```
Error: failed to generate suggestions: API returned error: {"error":{"type":"auth_error"}}
```
**Solution**: Check your API key is correct in the config file or environment variable.

#### Ollama Not Running
```
Error: failed to generate suggestions: Ollama server not running at http://localhost:11434
```
**Solution**: Start Ollama with `ollama serve` before using zeus-ai.

#### No Changes to Commit
```
Error: no changes to commit
```
**Solution**: Make sure you have changes in your working directory and stage them with `git add`.

## 🚧 Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for CLI functionality
- [Viper](https://github.com/spf13/viper) for configuration management
- [OpenRouter](https://openrouter.ai/), and [Ollama](https://ollama.ai/) for LLM APIs
