# Contributing to gcorelabscloud-go

We love your input! We want to make contributing to gcorelabscloud-go as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code follows the existing style
6. Issue that pull request!

## Development Setup

1. Install Go 1.21 or later
2. Clone the repository:
   ```bash
   git clone https://github.com/G-Core/gcorelabscloud-go.git
   cd gcorelabscloud-go
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run tests:
   ```bash
   make test
   ```

## Code Style

- Follow standard Go conventions and [Effective Go](https://golang.org/doc/effective_go)
- Run `make lint` to check your code style
- Write descriptive commit messages in the [conventional commits](https://www.conventionalcommits.org/) style
- Include comments on exported functions and types

## Testing

- Add tests for any new code you write
- Tests should be clear and maintainable
- Run the full test suite with `make test`
- Run integration tests with `make integration` (requires Gcore Cloud credentials)

## Pull Request Process

1. Update the README.md with details of changes to the interface, if applicable
2. Update any relevant documentation
3. The PR will be merged once you have the sign-off of at least one maintainer

## Any contributions you make will be under the Mozilla Public License Version 2.0

In short, when you submit code changes, your submissions are understood to be under the same [Mozilla Public License Version 2.0](LICENSE) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using GitHub's [issue tracker](https://github.com/G-Core/gcorelabscloud-go/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/G-Core/gcorelabscloud-go/issues/new/choose).

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## License

By contributing, you agree that your contributions will be licensed under its Mozilla Public License Version 2.0. 