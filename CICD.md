# CI/CD Pipeline — ip-app-golang

## Pipeline Overview

```
Checkout → Lint → Test → Build → Docker Build → Docker Push → Deploy Staging → Deploy Production
```

| Stage | Tool | Purpose |
|---|---|---|
| Checkout | Jenkins SCM | Pull source code |
| Lint | golangci-lint | Static analysis & code style enforcement |
| Test | go test | Unit tests with race detection & coverage |
| Build | go build | Compile binary artifact |
| Docker Build | docker | Build container image tagged with `<git-sha>-<build-num>` |
| Docker Push | docker | Push image to registry (main/master/release branches only) |
| Deploy Staging | kubectl | Auto-deploy to staging after push |
| Deploy Production | kubectl + manual gate | Manual approval required before production deploy |

---

## Jenkins Setup

### Required Plugins

| Plugin | Purpose |
|---|---|
| Pipeline | Declarative pipeline support |
| Docker Pipeline | Docker build/push in pipeline |
| HTML Publisher | Publish Go coverage report |
| Warnings Next Generation | Display golangci-lint results |
| Credentials Binding | Inject secrets into steps |
| Kubernetes CLI | kubectl access via kubeconfig credentials |

### Credentials to Configure

Go to **Manage Jenkins > Credentials > System > Global credentials**:

| Credential ID | Type | Description |
|---|---|---|
| `docker-registry-credentials` | Username/Password | Docker registry login |
| `kubeconfig-staging` | Secret file | kubeconfig for staging cluster |
| `kubeconfig-production` | Secret file | kubeconfig for production cluster |

### Tool Configuration

Go to **Manage Jenkins > Tools**:
- Add **Go installation** named `go` with version `1.22`

---

## Prerequisites — Local Setup

### macOS

```bash
# Install Go
brew install go

# Install golangci-lint
brew install golangci-lint

# Install Docker Desktop
brew install --cask docker

# Install kubectl
brew install kubectl
```

### Ubuntu

```bash
# Install Go
sudo apt update
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b $(go env GOPATH)/bin v1.57.2

# Install Docker
sudo apt install -y docker.io
sudo systemctl enable --now docker
sudo usermod -aG docker $USER  # logout and back in after this

# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -sL https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```

### Windows (PowerShell as Administrator)

```powershell
# Install Chocolatey (if not already installed)
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install tools
choco install golang -y
choco install docker-desktop -y
choco install kubernetes-cli -y

# Install golangci-lint
choco install golangci-lint -y
# or via PowerShell:
# Invoke-WebRequest -Uri https://github.com/golangci/golangci-lint/releases/download/v1.57.2/golangci-lint-1.57.2-windows-amd64.zip -OutFile golangci-lint.zip
# Expand-Archive golangci-lint.zip -DestinationPath C:\tools\
```

---

## Customization

### Change Docker Registry

Edit the `DOCKER_REGISTRY` and `IMAGE_REPO` values in the `environment` block:

```groovy
DOCKER_REGISTRY = 'ghcr.io'                    // GitHub Container Registry
// or
DOCKER_REGISTRY = '123456789.dkr.ecr.us-east-1.amazonaws.com'  // AWS ECR
IMAGE_REPO      = "${DOCKER_REGISTRY}/your-org/${APP_NAME}"
```

### Change Deployment Strategy

The deploy stages use `kubectl set image` (rolling update). To switch to Helm:

```groovy
sh "helm upgrade --install ${APP_NAME} ./helm-chart --set image.tag=${IMAGE_TAG} -n staging"
```

### Add Slack Notifications

Uncomment and configure the `slackSend` line in the `post { failure { ... } }` block.

### Branch Strategy

- `main` / `master` — triggers full pipeline including push + deploy
- `release/*` — triggers build and push only (no deploy)
- Feature branches — runs lint, test, build only
