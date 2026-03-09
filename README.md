<p align="center">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original-wordmark.svg" width="120" alt="Go logo" />
</p>

<h1 align="center">ip-app-golang</h1>

<p align="center">A lightweight HTTP service that returns the host's IP address, health probes, metrics, and build info — built with Go's standard library.</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.22" />
  <img src="https://img.shields.io/badge/stdlib_only-FF6F00?style=for-the-badge&logo=go&logoColor=white" alt="stdlib only" />
  <img src="https://img.shields.io/badge/Docker-00B4D8?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
  <img src="https://img.shields.io/badge/Kubernetes-9B5DE5?style=for-the-badge&logo=kubernetes&logoColor=white" alt="Kubernetes" />
  <img src="https://img.shields.io/badge/License-MIT-06D6A0?style=for-the-badge" alt="MIT" />
</p>

---

## Endpoints

| Method | Path       | Description                              | Content-Type       |
|--------|------------|------------------------------------------|--------------------|
| GET    | `/`        | Hostname, outbound IP, welcome message   | application/json   |
| GET    | `/healthz` | Health check                             | application/json   |
| GET    | `/ready`   | Readiness probe                          | application/json   |
| GET    | `/live`    | Liveness probe                           | application/json   |
| GET    | `/status`  | Status check                             | application/json   |
| GET    | `/ping`    | Plain-text pong                          | text/plain         |
| GET    | `/metrics` | Prometheus-style metrics                 | text/plain         |
| GET    | `/info`    | Build metadata and runtime info          | application/json   |

---

## Quick Start

```bash
go run main.go
```

The server starts on port **8080** by default.

---

## Docker

Build and run using the multi-stage Dockerfile:

```bash
docker build -t ip-app-golang .
docker run -p 8080:8080 ip-app-golang
```

The multi-stage build compiles a static binary in `golang:1.22-alpine`, then copies it into a minimal `alpine:3.19` runtime image. The final image is small, contains no build tools, and runs as a non-root user.

---

## Test Endpoints

```bash
curl http://localhost:8080/
curl http://localhost:8080/healthz
curl http://localhost:8080/ready
curl http://localhost:8080/live
curl http://localhost:8080/status
curl http://localhost:8080/ping
curl http://localhost:8080/metrics
curl http://localhost:8080/info
```

---

<details>
<summary><strong>Sample Responses</strong></summary>

**GET /**
```json
{
  "hostname": "dev-machine",
  "ip": "192.168.1.42",
  "message": "Welcome to ip-app-golang — a simple IP discovery service."
}
```

**GET /healthz, /ready, /live, /status**
```json
{
  "status": "ok"
}
```

**GET /ping**
```
pong
```

**GET /metrics**
```
http_requests_total 47
cpu_usage_percent 0.0
memory_usage_bytes 1283072
uptime_seconds 312.45
```

**GET /info**
```json
{
  "app": "ip-app-golang",
  "version": "1.4.2",
  "commit": "a3f22b9",
  "buildTime": "2026-03-05",
  "go": "go1.22.0",
  "platform": "linux",
  "arch": "amd64",
  "env": "development"
}
```

</details>

---

## Project Structure

```
ip-app-golang/
  main.go          # HTTP server — all routes and handlers
  go.mod           # Module definition (github.com/itdefinedclass/ip-app-golang)
  Dockerfile       # Multi-stage build (build in golang:1.22-alpine, run in alpine:3.19)
  .gitignore       # Ignores the compiled binary, .env, and Windows executables
  README.md        # This file
```

---

## What is go.sum?

If a `go.sum` file appears in this project, here is what it does:

- **Tracks dependency checksums** — every module version you download gets a cryptographic hash recorded in `go.sum` so builds are reproducible and tamper-proof.
- **Equivalent to `package-lock.json` (Node) or a `requirements.txt` lock (Python)** — it pins the exact bytes of each dependency.
- **Safe to regenerate** — delete it and run `go mod tidy`; Go will re-download dependencies and recreate the file.
- **Should be committed** — check it into version control so every developer and CI pipeline uses identical dependency hashes.

This project currently uses only the standard library, so no `go.sum` file is generated.

---

## Environment Variables

| Variable  | Default       | Description                        |
|-----------|---------------|------------------------------------|
| `PORT`    | `8080`        | TCP port the server listens on     |
| `APP_ENV` | `development` | Environment label shown in `/info` |

---

## Go Build and Packaging

**Static binary** — `go build` produces a single, statically linked executable with no runtime dependencies:

```bash
go build -o ip-app-golang main.go
./ip-app-golang
```

**Cross-compilation** — target any OS/architecture combination with two environment variables:

```bash
GOOS=linux   GOARCH=amd64 go build -o ip-app-golang main.go   # Linux x86_64
GOOS=linux   GOARCH=arm64 go build -o ip-app-golang main.go   # Linux ARM64
GOOS=darwin  GOARCH=arm64 go build -o ip-app-golang main.go   # macOS Apple Silicon
GOOS=windows GOARCH=amd64 go build -o ip-app-golang.exe main.go  # Windows
```

**Docker for production** — the included multi-stage Dockerfile builds with `CGO_ENABLED=0` for a fully static binary, then packages it in a minimal Alpine image. This is the recommended approach for Kubernetes deployments.

---

<p align="center"><sub>Built for the DevOps training class — keep shipping.</sub></p>
