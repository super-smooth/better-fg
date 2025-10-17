# better-fg

A more userfriendly and interactive version of `fg`.

## Features

- Interactive TUI for selecting from multiple background jobs
- Fuzzy searching for jobs
- Support for `fg %<number>` and `fg <number>`
- Falls back to normal `fg` behavior when there is only one background job

## Installation

### Using Homebrew (Recommended)

```bash
brew install Lewenhaupt/tap/better-fg
```

### Using Install Script

```bash
curl -sSL https://raw.githubusercontent.com/Lewenhaupt/better-fg/main/install.sh | bash
```

### Using Nix

```bash
# Install directly from the repository
nix profile install github:Lewenhaupt/better-fg

# Or run without installing
nix run github:Lewenhaupt/better-fg
```

### Manual Installation

Download the latest binary from the [releases page](https://github.com/Lewenhaupt/better-fg/releases) and place it in your `PATH`.

## Usage

To use `better-fg`, you need to initialize it in your shell. This is required because the tool needs to interact with your shell's job control.

1.  Make sure the `better-fg` binary is in your system's `PATH`.
2.  Add the following line to your shell's configuration file (e.g., `~/.zshrc`, `~/.bashrc`):

    ```bash
    eval "$(better-fg init)"
    ```

3.  Restart your shell or source your configuration file (e.g., `source ~/.zshrc`).

Now you can use the `bfg` command to interactively select background jobs, or continue using `fg` for the traditional behavior.

## Development

### Prerequisites

- Go 1.21+
- Nix (optional, for development environment)

### Setup

```bash
# Clone repository
git clone <repository-url>
cd better-fg

# Setup development environment
direnv allow  # if using Nix
# or
nix develop  # if using Nix without direnv
# or
go mod download  # if using Go directly

# Setup git hooks for commit message validation
./scripts/setup-git-hooks.sh

# Run tests
go test ./...

# Build using Go
go build -o better-fg ./cmd/better-fg

# Or build using Nix
nix build .#default
```
