# melange-forge

Lightweight, statically compiled Go binaries for distroless OCI images, packaged as apk via melange.

## Why

Distroless images have no shell, no `curl`, no `wget`. Some operations still require small binaries: healthchecks, probes, utilities. This repo provides minimal, statically compiled Go binaries packaged as apk packages via melange, ready to be embedded in any apko image.

Each binary is:
- Statically compiled (`CGO_ENABLED=0`)
- Signed and indexed as an apk package
- Tracked in the image SBOM via apko
- Scanned for vulnerabilities via Gosec CI

## Structure

```
melange-forge/
├── healthcheck-http/
│   ├── main.go                  # HTTP probe binary
│   ├── healthcheck-http.yaml    # melange build config
│   └── go.mod
├── healthcheck-sql/
│   ├── main.go                  # MariaDB/MySQL probe binary
│   ├── healthcheck-sql.yaml     # melange build config
│   ├── go.mod
│   └── go.sum
├── healthcheck-fcgi/
│   ├── main.go                  # FastCGI TCP probe binary
│   ├── healthcheck-fcgi.yaml    # melange build config
│   └── go.mod
└── issue-reporter/
    ├── main.go                  # CVE report formatter
    ├── grype.go                 # grype report structs
    ├── sbom.go                  # SPDX SBOM structs
    ├── report.go                # markdown body generator
    └── go.mod
```

## Binaries

### healthcheck-http

Performs an HTTP GET and exits 0 if the response code is between 200 and 399.

```
Usage: healthcheck-http --url=<url>

Flags:
  --url   URL to check (default: http://localhost/)
```

### healthcheck-sql

Performs a ping against a MariaDB/MySQL server and exits 0 if the server responds. Reads the password from a file (tmpfs/secret) and locks memory via `mlock` to prevent the password from being swapped to disk.

```
Usage: healthcheck-sql --host=<host> --port=<port> --user=<user> --secret=<path>

Flags:
  --host    MariaDB host (default: localhost)
  --port    MariaDB port (default: 3306)
  --user    MariaDB user (default: root)
  --secret  Path to password file (default: /run/secrets/mysql_password)
```

### healthcheck-fcgi

Opens a TCP connection to a FastCGI server (e.g. php-fpm) and exits 0 if the connection succeeds.

```
Usage: healthcheck-fcgi --addr=<host:port>

Flags:
  --addr  FastCGI address (default: localhost:9000)
```

### issue-reporter

Parses a grype JSON report and an optional apko SPDX SBOM, then outputs one JSON object per CVE on stdout. Designed to be piped into `gh issue create` in a GitHub Actions workflow.

Each output line contains:
- `title` — CVE ID + affected package + version
- `body` — formatted markdown with severity, CVSS score, EPSS, fix status, references, and upstream source from SBOM
- `severity` — lowercase severity for use as a GitHub label

```
Usage: issue-reporter --grype=<path> [--sbom=<path>] [--severity=<filter>]

Flags:
  --grype     Path to grype JSON report (default: grype-report.json)
  --sbom      Path to apko SPDX SBOM (optional, enriches output with upstream source)
  --severity  Comma-separated severity filter (default: high,critical)
```

Example output:

```json
{"title":"CVE-2026-8376 in perl:5.42.2-r2","body":"**Severity:** Critical\n...","severity":"critical"}
```

Example GHA usage:

```yaml
- name: Open issues if CVE found
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    ./issue-reporter --grype=grype-report.json --sbom=melange/mariadb/sbom-x86_64.spdx.json | \
    while read -r line; do
      title=$(echo "$line" | jq -r '.title')
      body=$(echo "$line" | jq -r '.body')
      severity=$(echo "$line" | jq -r '.severity')
      EXISTING=$(gh issue list --label security --state open --json title \
        | jq -r '.[].title' | grep "$title" || true)
      if [ -z "$EXISTING" ]; then
        gh issue create --title "$title" --body "$body" \
          --label "security" --label "severity:$severity"
      fi
    done
```

## Build

### Prerequisites

- [melange](https://github.com/chainguard-dev/melange)
- A signing key

```bash
melange keygen
```

### Build a package

```bash
melange build healthcheck-http/healthcheck-http.yaml \
  --arch amd64 \
  --signing-key melange.rsa
```

## Usage in apko.yaml

```yaml
contents:
  keyring:
    - https://packages.wolfi.dev/os/wolfi-signing.rsa.pub
    - ./melange.rsa.pub
  repositories:
    - https://packages.wolfi.dev/os
    - /path/to/melange-forge/healthcheck-http/packages
  packages:
    - wolfi-baselayout
    - your-service
    - healthcheck-http
```

Then in your docker-compose:

```yaml
healthcheck:
  test: ["CMD", "healthcheck-http", "--url=http://localhost:8080/status/"]
  interval: 10s
  timeout: 5s
  retries: 5
```

## CI

A GitHub Actions workflow runs Gosec on all `*.go` files on every push and pull request to `main`. The scan results are published in the job summary. The workflow blocks on any finding: no `|| true`.

## Roadmap

- Refactor Gosec CI to scan all `*.go` files in a single pass instead of a matrix per service
- GitHub Actions release workflow to publish `issue-reporter` as a binary artifact on each tag
- `issue-reporter`: include image name in issue title and body
- `issue-reporter`: CycloneDX SBOM support in addition to SPDX

## License

MIT
