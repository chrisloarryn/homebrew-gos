# Release Process Documentation

This document describes how to create and manage releases for the GOS CLI.

## 🚀 Creating a New Release

### Automated Release (Recommended)

Use the provided release script:

```bash
# Create and push a new release
./release.sh 1.0.0
```

Or using Make:

```bash
# Create and push a new release
make release VERSION='1.0.0'
```

### Manual Release

1. **Ensure clean working directory:**
   ```bash
   git status
   # Should show no uncommitted changes
   ```

2. **Run tests:**
   ```bash
   make test
   ```

3. **Create and push tag:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

## 📦 What Happens During Release

When a tag is pushed, GitHub Actions automatically:

1. **Builds binaries** for multiple platforms:
   - macOS (Intel and Apple Silicon)
   - Linux (x86_64 and ARM64)
   - Windows (x86_64)

2. **Creates packages:**
   - Tar.gz archives for Unix systems
   - Zip archives for Windows
   - Debian packages (.deb)
   - RPM packages (.rpm)

3. **Docker images:**
   - Multi-platform Docker images
   - Tagged with version and latest

4. **Homebrew formula:**
   - Updates the Homebrew tap automatically

5. **Release notes:**
   - Auto-generated changelog from commits
   - Organized by type (features, bug fixes, etc.)

## 🏷️ Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR.MINOR.PATCH** (e.g., 1.2.3)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Examples:
- `1.0.0` - Initial stable release
- `1.1.0` - New features added
- `1.1.1` - Bug fixes
- `2.0.0` - Breaking changes

## 📋 Pre-Release Checklist

Before creating a release:

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] CHANGELOG is updated (if manual)
- [ ] Version follows semantic versioning
- [ ] No breaking changes in minor/patch releases
- [ ] Git working directory is clean

## 🔧 Release Configuration

### GoReleaser Configuration

The `.goreleaser.yml` file configures:
- Build targets and flags
- Archive formats
- Package generation
- Docker image creation
- Release notes generation
- Homebrew tap updates

### GitHub Actions

The `.github/workflows/release.yml` file:
- Triggers on tag pushes
- Runs tests before release
- Uses GoReleaser for building
- Publishes to multiple platforms

## 📊 Monitoring Releases

### GitHub Actions Dashboard
Monitor release progress at:
```
https://github.com/cristobalcontreras/homebrew-gos/actions
```

### Release Page
View releases at:
```
https://github.com/cristobalcontreras/homebrew-gos/releases
```

### Docker Hub
Docker images are available at:
```
https://hub.docker.com/r/cristobalcontreras/gos
```

## 🐛 Troubleshooting Releases

### Failed Release Build

1. Check GitHub Actions logs
2. Verify GoReleaser configuration
3. Ensure all required secrets are set

### Missing Assets

1. Check `.goreleaser.yml` configuration
2. Verify build targets
3. Check for build errors in logs

### Docker Issues

1. Verify Docker Hub credentials
2. Check Dockerfile syntax
3. Ensure base image availability

## 🔐 Required Secrets

For full release functionality, ensure these secrets are set in GitHub:

- `GITHUB_TOKEN` - Automatically provided by GitHub
- `DOCKERHUB_USERNAME` - Docker Hub username (optional)
- `DOCKERHUB_TOKEN` - Docker Hub access token (optional)

## 📈 Release Analytics

Track release adoption through:
- GitHub release download counts
- Docker Hub pull statistics
- Homebrew install analytics (if available)

## 🔄 Hotfix Releases

For urgent bug fixes:

1. Create hotfix branch from main
2. Apply minimal fix
3. Test thoroughly
4. Create patch release (e.g., 1.1.1 → 1.1.2)
5. Follow normal release process

## 📝 Release Notes Template

Release notes are auto-generated but can be customized in `.goreleaser.yml`:

```markdown
## 🚀 GOS CLI v1.0.0

### 🎉 What's New
- New feature descriptions

### 🐛 Bug Fixes
- Bug fix descriptions

### 🔧 Improvements
- Enhancement descriptions

### 📦 Installation
[Installation instructions]
```

## 🏁 Post-Release Tasks

After a successful release:

1. Verify all assets are available
2. Test installation on different platforms
3. Update documentation if needed
4. Announce release (if applicable)
5. Monitor for issues

---

For questions about the release process, check the GitHub Issues or contact the maintainers.
