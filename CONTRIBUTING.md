# Contributing to zeus-ai

First of all, thank you for considering contributing to zeus-ai! Your help is essential for making this tool better for everyone.

This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How Can I Contribute?

### Reporting Bugs

- **Ensure the bug was not already reported** by searching on GitHub under [Issues](https://github.com/amosehiguese/zeus-ai/issues).
- If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/amosehiguese/zeus-ai/issues/new). Be sure to include a **title and clear description**, as much relevant information as possible, and a **code sample** or an **executable test case** demonstrating the expected behavior that is not occurring.

### Suggesting Enhancements

- **Ensure the enhancement was not already suggested** by searching on GitHub under [Issues](https://github.com/amosehiguese/zeus-ai/issues).
- If you're unable to find an open issue for your enhancement, [open a new one](https://github.com/amosehiguese/zeus-ai/issues/new). Be sure to include a **title and clear description**, as much relevant information as possible, and a **mockup** if applicable.

### Pull Requests

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run the tests (`make test`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.20 or higher
- Git
- Optional: Ollama for local LLM testing

### Setting Up Development Environment

1. Clone your forked repository
   ```bash
   git clone https://github.com/amosehiguese/zeus-ai.git
   cd zeus-ai
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Build the project
   ```bash
   make build
   ```

4. Run tests
   ```bash
   make test
   ```

## Project Structure

```
zeus-ai/
├── cmd/                    # Application entry points
│   └── zeus-ai/            # Main CLI application
│       └── main.go
├── internal/               # Private application code
│   ├── config/             # Configuration handling
│   └── utils/              # Internal utilities
├── pkg/                    # Public libraries that can be used by external applications
│   ├── git/                # Git operations
│   ├── llm/                # LLM provider integrations
│   ├── prompt/             # Prompt templates
│   └── terminal/           # Terminal UI utilities
├── assets/                 # Static assets (images, etc.)
├── .github/                # GitHub workflows and templates
├── Makefile                # Build automation
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksum
├── LICENSE                 # License file
├── README.md               # Project documentation
└── CONTRIBUTING.md         # Contribution guidelines
```

## Coding Guidelines

### Go Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Format your code with `gofmt -s`
- Use meaningful variable names and add comments where necessary
- Write unit tests for new features

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- always sign your commits
- Reference issues and pull requests liberally after the first line
- Consider starting the commit message with an applicable type prefix:
  - `feat:` for new features
  - `fix:` for bug fixes
  - `docs:` for documentation changes
  - `style:` for formatting changes
  - `refactor:` for code refactoring
  - `test:` for adding tests
  - `chore:` for maintenance tasks

## Testing

- Write unit tests for new features
- Ensure all tests pass before submitting a pull request
- Consider adding integration tests for complex features

To run tests:

```bash
make test
```

## Documentation

- Update the README.md with details of changes to the interface
- Update the examples when adding new features
- Add proper godoc comments to exported functions and types

## License

By contributing to zeus-ai, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).