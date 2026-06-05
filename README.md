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
└── healthcheck-fcgi/
    ├── main.go                  # FastCGI TCP probe binary
    ├── healthcheck-fcgi.yaml    # melange build config
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

## Roadmap

More binaries will be added over time as new services require them.

## License

MIT
