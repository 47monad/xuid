# Releasing

Releases are cut from `main` as annotated, semver-prefixed Git tags (`vX.Y.Z`).
The Go module proxy and pkg.go.dev discover a release from the tag alone; a
GitHub Release is created automatically by
[`.github/workflows/release.yml`](.github/workflows/release.yml).

## Version policy

- Pre-1.0 (`v0.x.y`): breaking changes may land in minor bumps (`v0.2.0`).
- `v1.0.0` and later: breaking changes require a new major version and a module
  path ending in `/v2`, `/v3`, … (`module github.com/47monad/xuid/v2`).
- Patch releases (`v0.1.1`) are for backwards-compatible fixes only.

## Steps

1. Ensure `main` is green and up to date:

   ```bash
   git checkout main
   git pull --ff-only
   git status --porcelain        # must be empty
   go mod tidy -diff             # must be empty
   go build ./...
   go vet ./...
   go test -race ./...
   ```

2. Update `CHANGELOG.md`: move items out of `Unreleased` into a new version
   section dated today, and update the compare links at the bottom.

3. Commit the changelog (and any release housekeeping) and push to `main`.

4. Create and push the annotated tag:

   ```bash
   git tag -a v0.1.0 -m "xuid v0.1.0"
   git push origin v0.1.0
   ```

   Use `git tag -s` instead of `-a` for a GPG-signed tag.

5. Verify the release:

   ```bash
   curl -s https://proxy.golang.org/github.com/47monad/xuid/@v/list
   curl -s https://proxy.golang.org/github.com/47monad/xuid/@v/v0.1.0.info
   ```

   Then open <https://pkg.go.dev/github.com/47monad/xuid@v0.1.0>. If the proxy
   version list still looks empty it is cached briefly; requesting the explicit
   `.info` URL or running `go get` bypasses the stale cache.

## Notes

- Tags must be valid semver with a leading `v` (`v0.1.0`, not `0.1.0`).
- The module lives at the repository root, so the tag needs no subdirectory
  prefix. For a module in a subdirectory the tag would be `<subdir>/vX.Y.Z`.
- The minimum supported Go version is declared in `go.mod`; mention it in the
  release notes when it changes.
