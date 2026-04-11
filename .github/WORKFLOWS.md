# GitHub Actions Workflows

This directory contains the GitHub Actions workflows for the better-fg project.

## Workflows

### 🔄 CI (`ci.yml`)
**Triggers**: Push to main, Pull requests
- **Prerequisites Check**: Detects if go.mod was modified and skips expensive jobs on non-bot commits
- **Go Testing**: Runs tests on multiple Go versions (1.21, 1.22, 1.23)
- **Linting**: golangci-lint, gofmt, go vet
- **Nix Build**: Tests Nix package builds
- **Cross-platform Build**: Tests compilation for multiple platforms

**Optimization**: When Dependabot updates go.mod, the first CI run skips expensive jobs (test, lint, build). The `gomod2nix` workflow commits the updated gomod2nix.toml, which triggers a second CI run that executes all jobs with the correct state.

### 📦 gomod2nix (`gomod2nix.yml`)
**Triggers**: Pull requests (when go.mod is modified)
- **Check**: Detects if go.mod was modified in the PR
- **Update**: Regenerates gomod2nix.toml using Nix
- **Commit**: Signs and commits the updated file as the bot user
- **Skip**: Exits cleanly if go.mod was not modified

This workflow runs before the full CI suite and ensures gomod2nix.toml is always in sync with go.mod changes.

### 📝 Commit Linting (`lint-commits.yml`)
**Triggers**: Push to main, Pull requests
- Validates commit message format: `type: description`
- Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Enforces project-specific rules (no unauthorized code attributions)

### 🚀 Release (`release.yml`)
**Triggers**: Tag push (`v*.*.*`)
- **Build Binaries**: Cross-compiles for Linux, macOS, Windows (amd64/arm64)
- **Create Release**: Auto-generates changelog and uploads binaries
- **Update Flake**: Creates PR to update flake.nix with new version and vendorHash

### 📦 Dependencies (`deps.yml`)
**Triggers**: Weekly schedule (Sundays), Manual dispatch
- **Go Dependencies**: Updates Go modules and creates PR
- **Nix Dependencies**: Updates flake inputs and creates PR
- **Security Audit**: Runs vulnerability scans

## Configuration Files

### `.golangci.yml`
Comprehensive Go linting configuration with:
- Code quality checks
- Security scanning
- Style enforcement
- Performance optimizations

### `commitlint.config.js`
Commit message validation rules:
- Enforces conventional commit format
- Project-specific attribution rules
- Message length limits

### `dependabot.yml`
Automated dependency updates for:
- Go modules (weekly)
- GitHub Actions (weekly)

## Release Process

1. **Create and push a tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Automated workflow**:
   - Builds binaries for all platforms
   - Creates GitHub release with changelog
   - Uploads release assets
   - Creates PR to update flake.nix

3. **Manual steps**:
   - Review and merge the flake.nix update PR
   - Announce the release

## Security

- All workflows use pinned action versions
- Minimal permissions (GITHUB_TOKEN only)
- Vulnerability scanning included
- Dependency review on PRs

## Caching

- Go modules cache for faster builds
- Nix store cache using GitHub Actions cache
- Cross-platform build artifacts

## Required Secrets

- `GITHUB_TOKEN`: Automatically provided by GitHub
- `GH_PAT` (optional): Personal Access Token for the bot account to trigger workflows (falls back to GITHUB_TOKEN)
- `GPG_PRIVATE_KEY` (optional): GPG key for signing automated commits
- `GPG_PASSPHRASE` (optional): Passphrase for the GPG key

## Configuration Variables

The following variables can be configured at the organization or repository level:

- `BOT_NAME` (optional): Name of the bot user for automated commits (default: `super-smooth-bot`)
- `BOT_EMAIL` (optional): Email of the bot user for automated commits (default: `super-smooth-bot@anomaly.co`)

## Customization

To customize for your repository:

1. Update `dependabot.yml` with your GitHub username
2. Add your username to workflow reviewers/assignees
3. Adjust Go version matrix in CI if needed
4. Configure additional caching strategies if desired
