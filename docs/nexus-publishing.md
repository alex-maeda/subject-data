# Publishing Go Modules to Nexus

The `data-model/` and `swagger/` directories are standalone Go modules that consumers can import. They are published to a Nexus `raw (hosted)` repository on `artifact-server`, which serves files following the Go module proxy protocol.

## Published Modules

| Module | Path | Version file |
|--------|------|-------------|
| data-model | `github.com/sovraai/subject-data/data-model` | `data-model/VERSION` |
| swagger | `github.com/sovraai/subject-data/swagger` | `swagger/VERSION` |

## Nexus Setup (One-Time)

Create a single repository in Nexus:

* **Name:** `go-internal`
* **Type:** raw (hosted)
* **Deployment policy:** Disable redeploy (prevents overwriting published versions)

## Versioning

Each module has a `VERSION` file containing a semver string (e.g. `v0.1.0`). To release a new version:

1. Update the `VERSION` file (e.g. `v0.1.0` → `v0.2.0`)
2. Commit and push to `main`
3. CI publishes the module to Nexus automatically

If the version already exists in Nexus, CI skips publishing — so pushes that don't bump the version are no-ops. Overwriting a published version is not possible (disabled at the Nexus level).

## CI Workflow

The `.github/workflows/publish-modules.yml` workflow runs on every push to `main`. For each module it:

1. Reads the `VERSION` file
2. Checks if that version already exists in Nexus (skips if so)
3. Builds and uploads the 3 files required by the Go module proxy protocol:
   * `@v/<version>.info` — version metadata JSON
   * `@v/<version>.mod` — the module's go.mod
   * `@v/<version>.zip` — module source in Go module zip format

### Required GitHub Configuration

Set up secrets and variables via the `gh` CLI:

The following org-level secrets are already configured:

* `NEXUS_CI_USERNAME` / `NEXUS_CI_PASSWORD` — Nexus credentials
* `TS_OAUTH_CLIENT_ID` / `TS_OAUTH_SECRET` — Tailscale OAuth for CI runners

Set the repository-level variable:

```bash
gh variable set NEXUS_GO_URL --body "http://artifact-server.tail5ab057.ts.net/nexus/repository/go-internal"
```

The workflow uses `tailscale/github-action@v3` to connect standard GitHub-hosted runners to the tailnet, matching the pattern used by bam_prototype's CI.

## Consumer Setup

### Environment

```bash
export GOPROXY="http://artifact-server.tail5ab057.ts.net/nexus/repository/go-internal/,https://proxy.golang.org,direct"
export GONOSUMDB="github.com/sovraai/*"
export GONOSUMCHECK="github.com/sovraai/*"
```

If Nexus requires authentication for reads, also set:

```bash
export GONOSUMDB="github.com/sovraai/*"
```

And add to `~/.netrc`:

```
machine artifact-server.tail5ab057.ts.net
  login <username>
  password <password>
```

### Usage in go.mod

```
require github.com/sovraai/subject-data/data-model v0.1.0
```

Then `go mod tidy`. The consumer gets only the types from the data-model module — no server dependencies (sqlx, pgx, chi, etc.) are pulled in.

### CI Configuration for Consumers

Add to CI environment variables:

```yaml
env:
  GOPROXY: "http://artifact-server.tail5ab057.ts.net/nexus/repository/go-internal/,https://proxy.golang.org,direct"
  GONOSUMDB: "github.com/sovraai/*"
  GONOSUMCHECK: "github.com/sovraai/*"
```
