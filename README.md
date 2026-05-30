# melange-forge

A lightweight Go toolkit providing minimal, statically compiled healthcheck binaries for distroless OCI images built with apko/melange.

## Why

Distroless images have no shell, no `curl`, no `wget`. Docker healthchecks need a binary. This repo provides minimal, statically compiled healthcheck binaries packaged as apk packages via melange, ready to be embedded in any apko image.

## Structure

```
melange-forge/
├── healthcheck-http/
│   ├── main.go                  # HTTP healthcheck binary
│   ├── healthcheck-http.yaml    # melange build config
│   └── go.mod
└── healthcheck-sql/
    ├── main.go                  # MariaDB/MySQL healthcheck binary
    ├── healthcheck-sql.yaml     # melange build config
    ├── go.mod
    └── go.sum
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

Performs a ping against a MariaDB/MySQL server and exits 0 if the server responds.

```
Usage: healthcheck-sql --host=<host> --port=<port> --user=<user> --password=<password>

Flags:
  --host      MariaDB host (default: localhost)
  --port      MariaDB port (default: 3306)
  --user      MariaDB user (default: root)
  --password  MariaDB password (default: "")
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

```bash
melange build healthcheck-sql/healthcheck-sql.yaml \
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

More healthcheck binaries will be added over time as new services require them.

## License

MIT
