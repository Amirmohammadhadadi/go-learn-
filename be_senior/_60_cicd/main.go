// ============================================================================
// FILE: cicd_guide.go
// TITLE: راهنمای کامل CI/CD در Go با GitLab CI و GitHub Actions
// HOW TO RUN: go run cicd_guide.go (فایل توضیحی - دستورات در ترمینال اجرا شوند)
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - CI/CD چیست؟
// ============================================================================
//
// CI/CD (Continuous Integration/Continuous Deployment) مجموعه‌ای از practiceهاست
// که به شما امکان می‌دهد کد را به صورت خودکار build، test و deploy کنید.
//
// Continuous Integration (CI):
// - ادغام مکرر کد در مخزن اصلی
// - اجرای خودکار تست‌ها
// - شناسایی سریع مشکلات
//
// Continuous Deployment (CD):
// - استقرار خودکار پس از گذراندن تست‌ها
// - کاهش زمان بین نوشتن کد و رسیدن به پروداکشن
//
// قانون طلایی:
// "هر commit باید build و test شود. هرگز با تست شکست خورده merge نکن.
//  از cache برای بهبود سرعت استفاده کن. Secrets را هرگز در فایل‌های CI ذخیره نکن."
// ============================================================================

package __cicd

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🚀 CI/CD FOR GO")
	fmt.Println("GitLab CI | GitHub Actions")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: GitLab CI (پایپلاین کامل)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 1: GITLAB CI COMPLETE PIPELINE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# .gitlab-ci.yml (پایپلاین کامل GitLab CI برای Go)
# ============================================================================

# مراحل (stages)
stages:
  - test
  - build
  - scan
  - package
  - deploy

# متغیرهای عمومی (می‌تواند در UI تنظیم شود)
variables:
  GO_VERSION: "1.24"
  GOPROXY: "https://proxy.golang.org,direct"
  GOSUMDB: "sum.golang.org"
  CGO_ENABLED: "0"
  GOOS: "linux"
  GOARCH: "amd64"

# Cache برای وابستگی‌ها
cache:
  key: ${CI_COMMIT_REF_SLUG}
  paths:
    - .go/pkg/mod
    - .go-cache

# قبل از هر job (تنظیمات اولیه)
before_script:
  - mkdir -p .go
  - export GOPATH="$CI_PROJECT_DIR/.go"
  - export GOCACHE="$CI_PROJECT_DIR/.go-cache"
  - go version
  - go env

# ============================================================================
# مرحله 1: تست (Test)
# ============================================================================

# Lint (بررسی کیفیت کد)
lint:
  stage: test
  image: golangci/golangci-lint:latest
  script:
    - golangci-lint run --timeout=5m
  allow_failure: false

# Unit Tests (تست‌های واحد)
unit-tests:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go mod download
    - go test -v -race -coverprofile=coverage.out ./...
    - go tool cover -func=coverage.out
    - go tool cover -html=coverage.out -o coverage.html
  artifacts:
    paths:
      - coverage.html
      - coverage.out
    expire_in: 7 days
  coverage: '/total:.*?(\d+\.\d+)%/'

# Integration Tests (تست‌های یکپارچه)
integration-tests:
  stage: test
  image: golang:${GO_VERSION}
  services:
    - postgres:15-alpine
    - redis:7-alpine
  variables:
    POSTGRES_DB: testdb
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: password
    REDIS_HOST: redis
  script:
    - go test -v -tags=integration ./...
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event"'
    - if: '$CI_COMMIT_BRANCH == "main"'
    - if: '$CI_COMMIT_BRANCH == "develop"'

# Benchmark Tests (تست‌های عملکرد)
benchmark:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go test -bench=. -benchmem ./... > benchmark.txt
  artifacts:
    paths:
      - benchmark.txt
    expire_in: 7 days
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'
    - if: '$CI_COMMIT_BRANCH == "develop"'

# ============================================================================
# مرحله 2: Build (کامپایل)
# ============================================================================

# Build برای معماری‌های مختلف
build-multi-arch:
  stage: build
  image: golang:${GO_VERSION}
  script:
    # Build for Linux amd64
    - GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/app-linux-amd64 ./cmd/server
    # Build for Linux arm64
    - GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/app-linux-arm64 ./cmd/server
    # Build for Windows amd64
    - GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/app-windows-amd64.exe ./cmd/server
    # Build for macOS amd64
    - GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/app-darwin-amd64 ./cmd/server
  artifacts:
    paths:
      - bin/
    expire_in: 7 days
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'
    - if: '$CI_COMMIT_TAG'

# ============================================================================
# مرحله 3: Scan (امنیت)
# ============================================================================

# Security Scan با govulncheck
govulncheck:
  stage: scan
  image: golang:${GO_VERSION}
  script:
    - go install golang.org/x/vuln/cmd/govulncheck@latest
    - govulncheck ./...
  allow_failure: true

# Dependency Scanning
dependency-scan:
  stage: scan
  image: aquasec/trivy:latest
  script:
    - trivy fs --severity HIGH,CRITICAL --no-progress .
  allow_failure: true

# ============================================================================
# مرحله 4: Package (ساخت Docker Image)
# ============================================================================

# Build & Push Docker Image
docker-build:
  stage: package
  image: docker:24
  services:
    - docker:dind
  variables:
    DOCKER_DRIVER: overlay2
    DOCKER_TLS_CERTDIR: "/certs"
  script:
    # Login to Registry
    - echo "$CI_REGISTRY_PASSWORD" | docker login -u "$CI_REGISTRY_USER" --password-stdin $CI_REGISTRY
    # Build image
    - docker build --target runtime -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA -f Dockerfile .
    - docker tag $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA $CI_REGISTRY_IMAGE:latest
    # Push images
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA
    - docker push $CI_REGISTRY_IMAGE:latest
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'
    - if: '$CI_COMMIT_TAG'

# ============================================================================
# مرحله 5: Deploy (استقرار)
# ============================================================================

# Deploy to Staging
deploy-staging:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache openssh-client
    - mkdir -p ~/.ssh
    - echo "$SSH_PRIVATE_KEY" > ~/.ssh/id_rsa
    - chmod 600 ~/.ssh/id_rsa
  script:
    - ssh -o StrictHostKeyChecking=no $STAGING_USER@$STAGING_HOST "
        docker pull $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA &&
        docker-compose -f /app/docker-compose.yml up -d --no-deps app
      "
  environment:
    name: staging
    url: https://staging.example.com
  rules:
    - if: '$CI_COMMIT_BRANCH == "develop"'

# Deploy to Production
deploy-production:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache openssh-client
    - mkdir -p ~/.ssh
    - echo "$SSH_PRIVATE_KEY_PROD" > ~/.ssh/id_rsa
    - chmod 600 ~/.ssh/id_rsa
  script:
    - ssh -o StrictHostKeyChecking=no $PROD_USER@$PROD_HOST "
        docker pull $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA &&
        docker-compose -f /app/docker-compose.yml up -d --no-deps app
      "
  environment:
    name: production
    url: https://example.com
  rules:
    - if: '$CI_COMMIT_TAG'
  when: manual
`)

	// ============================================================================
	// بخش 2: GitHub Actions (پایپلاین کامل)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🐙 SECTION 2: GITHUB ACTIONS COMPLETE PIPELINE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# .github/workflows/ci.yml (پایپلاین کامل GitHub Actions برای Go)
# ============================================================================

name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
    tags: ['v*']
  pull_request:
    branches: [main, develop]
  workflow_dispatch:  # اجرای دستی

env:
  GO_VERSION: '1.24'
  GOPROXY: 'https://proxy.golang.org,direct'
  CGO_ENABLED: '0'

jobs:
  # ============================================================================
  # Job 1: Lint
  # ============================================================================
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: --timeout=5m

  # ============================================================================
  # Job 2: Test
  # ============================================================================
  test:
    name: Test
    runs-on: ubuntu-latest
    needs: lint
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Download dependencies
        run: go mod download

      - name: Run unit tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
          flags: unittests

      - name: Run benchmark
        run: go test -bench=. -benchmem ./...

  # ============================================================================
  # Job 3: Security
  # ============================================================================
  security:
    name: Security Scan
    runs-on: ubuntu-latest
    permissions:
      security-events: write
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Run govulncheck
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...

      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'

      - name: Upload Trivy results to GitHub Security tab
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'

  # ============================================================================
  # Job 4: Build
  # ============================================================================
  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [test, security]
    strategy:
      matrix:
        os: [linux, windows, darwin]
        arch: [amd64, arm64]
        exclude:
          - os: windows
            arch: arm64
          - os: darwin
            arch: arm64
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Build
        env:
          GOOS: ${{ matrix.os }}
          GOARCH: ${{ matrix.arch }}
        run: |
          OUTPUT_NAME=app-${{ matrix.os }}-${{ matrix.arch }}
          if [ "${{ matrix.os }}" = "windows" ]; then OUTPUT_NAME="$OUTPUT_NAME.exe"; fi
          go build -ldflags="-s -w" -o bin/$OUTPUT_NAME ./cmd/server

      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: app-${{ matrix.os }}-${{ matrix.arch }}
          path: bin/

  # ============================================================================
  # Job 5: Docker Build
  # ============================================================================
  docker:
    name: Docker Build & Push
    runs-on: ubuntu-latest
    needs: build
    if: github.event_name == 'push' && (github.ref == 'refs/heads/main' || startsWith(github.ref, 'refs/tags/v'))
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to DockerHub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_TOKEN }}

      - name: Login to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: |
            ${{ secrets.DOCKER_USERNAME }}/myapp
            ghcr.io/${{ github.repository }}
          tags: |
            type=ref,event=branch
            type=ref,event=tag
            type=sha,format=short
            type=raw,value=latest,enable={{is_default_branch}}

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          target: runtime
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  # ============================================================================
  # Job 6: Deploy
  # ============================================================================
  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    needs: docker
    if: github.ref == 'refs/heads/develop'
    environment:
      name: staging
      url: https://staging.example.com
    steps:
      - name: Deploy to Staging Server
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.STAGING_HOST }}
          username: ${{ secrets.STAGING_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            docker pull ghcr.io/${{ github.repository }}:sha-${GITHUB_SHA::7}
            docker-compose -f /app/docker-compose.yml up -d --no-deps app
            docker image prune -f

  deploy-production:
    name: Deploy to Production
    runs-on: ubuntu-latest
    needs: docker
    if: startsWith(github.ref, 'refs/tags/v')
    environment:
      name: production
      url: https://example.com
    steps:
      - name: Deploy to Production Server
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.PROD_HOST }}
          username: ${{ secrets.PROD_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY_PROD }}
          script: |
            docker pull ghcr.io/${{ github.repository }}:${{ github.ref_name }}
            docker-compose -f /app/docker-compose-prod.yml up -d --no-deps app
            docker image prune -f

  # ============================================================================
  # Job 7: Release
  # ============================================================================
  release:
    name: Create Release
    runs-on: ubuntu-latest
    needs: [build, docker]
    if: startsWith(github.ref, 'refs/tags/v')
    permissions:
      contents: write
    steps:
      - name: Download artifacts
        uses: actions/download-artifact@v4
        with:
          path: artifacts

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: artifacts/**/*
          generate_release_notes: true
`)

	// ============================================================================
	// بخش 3: GitLab CI برای Monorepo
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📁 SECTION 3: GITLAB CI FOR MONOREPO")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# .gitlab-ci.yml (Monorepo با چند سرویس)
# ============================================================================

# تنظیمات برای monorepo با چند سرویس
variables:
  GO_VERSION: "1.24"

# مرحله اول: فقط تغییرات مربوط به service1 را بررسی کن
.service1-changes: &service1-changes
  rules:
    - changes:
        - "services/service1/**/*"
        - "pkg/**/*"
    - if: '$CI_MERGE_REQUEST_TARGET_BRANCH_NAME == "main"'

# هر سرویس به صورت جداگانه test و build می‌شود
service1:test:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - cd services/service1
    - go mod download
    - go test -v ./...
  <<: *service1-changes

service1:build:
  stage: build
  image: golang:${GO_VERSION}
  script:
    - cd services/service1
    - go build -o bin/service1 ./cmd/server
  artifacts:
    paths:
      - services/service1/bin/
  <<: *service1-changes

service1:docker:
  stage: package
  image: docker:24
  services:
    - docker:dind
  script:
    - docker build -t service1:latest -f services/service1/Dockerfile .
  <<: *service1-changes
`)

	// ============================================================================
	// بخش 4: GitHub Actions برای Monorepo
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📁 SECTION 4: GITHUB ACTIONS FOR MONOREPO")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# .github/workflows/monorepo.yml (Monorepo با Path Filters)
# ============================================================================

name: Monorepo CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

# استفاده از path filters برای کاهش اجرای غیرضروری
jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      service1: ${{ steps.filter.outputs.service1 }}
      service2: ${{ steps.filter.outputs.service2 }}
    steps:
      - uses: actions/checkout@v4
      - uses: dorny/paths-filter@v3
        id: filter
        with:
          filters: |
            service1:
              - 'services/service1/**'
              - 'pkg/shared/**'
            service2:
              - 'services/service2/**'
              - 'pkg/shared/**'

  service1-ci:
    needs: detect-changes
    if: needs.detect-changes.outputs.service1 == 'true'
    uses: ./.github/workflows/service-ci.yml
    with:
      service: service1
    secrets: inherit

  service2-ci:
    needs: detect-changes
    if: needs.detect-changes.outputs.service2 == 'true'
    uses: ./.github/workflows/service-ci.yml
    with:
      service: service2
    secrets: inherit
`)

	// ============================================================================
	// بخش 5: Secrets Management
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔐 SECTION 5: SECRETS MANAGEMENT")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ GITLAB CI SECRETS                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│   1. Go to Settings → CI/CD → Variables                                   │
│   2. Add variables:                                                        │
│      - REGISTRY_PASSWORD (Masked)                                          │
│      - SSH_PRIVATE_KEY (Masked, Protected)                                 │
│      - DOCKER_PASSWORD (Masked)                                            │
│      - AWS_ACCESS_KEY_ID (Masked)                                          │
│                                                                             │
│   Usage in .gitlab-ci.yml:                                                 │
│   - $CI_REGISTRY_PASSWORD                                                  │
│   - $SSH_PRIVATE_KEY                                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ GITHUB ACTIONS SECRETS                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│   1. Go to Settings → Secrets and variables → Actions                     │
│   2. Add secrets:                                                          │
│      - DOCKER_USERNAME                                                     │
│      - DOCKER_TOKEN                                                        │
│      - STAGING_HOST                                                        │
│      - STAGING_USER                                                        │
│      - SSH_PRIVATE_KEY                                                     │
│      - GITHUB_TOKEN (automatically provided)                               │
│                                                                             │
│   Usage in workflow:                                                       │
│   - ${{ secrets.DOCKER_USERNAME }}                                         │
│   - ${{ secrets.SSH_PRIVATE_KEY }}                                         │
└─────────────────────────────────────────────────────────────────────────────┘

⚠️  NEVER:
   • Store secrets in code or config files
   • Print secrets in logs
   • Commit .env files
   • Use secrets in MRs from forks without review
`)

	// ============================================================================
	// بخش 6: Cache Optimization
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 6: CACHE OPTIMIZATION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ GITLAB CI CACHE                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│ cache:                                                                    │
│   key: ${CI_COMMIT_REF_SLUG}                                              │
│   paths:                                                                  │
│     - .go/pkg/mod                                                         │
│     - .go-cache                                                           │
│                                                                             │
│ # Cache بین branches مختلف                                                │
│ cache:                                                                    │
│   key: ${CI_JOB_NAME}                                                     │
│   paths:                                                                  │
│     - .go/pkg/mod                                                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ GITHUB ACTIONS CACHE                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│ - name: Cache Go modules                                                  │
│   uses: actions/cache@v4                                                  │
│   with:                                                                   │
│     path: |                                                               │
│       ~/.cache/go-build                                                   │
│       ~/go/pkg/mod                                                        │
│     key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}                │
│     restore-keys: |                                                       │
│       ${{ runner.os }}-go-                                                │
│                                                                             │
│ # Built-in Go cache in setup-go action                                    │
│ - uses: actions/setup-go@v5                                               │
│   with:                                                                   │
│     go-version: '1.24'                                                    │
│     cache: true                                                           │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	// ============================================================================
	// بخش 7: Best Practices
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 CI/CD BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. PIPELINE STRUCTURE                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Stages: test → build → scan → package → deploy                       │
│    ✔ Fail fast: lint and test first                                       │
│    ✔ Parallel jobs for speed                                              │
│    ✔ Manual approval for production deployment                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. TESTING                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Run unit tests with race detector                                    │
│    ✔ Run integration tests with service containers                        │
│    ✔ Generate and store coverage reports                                  │
│    ✔ Run benchmarks for performance regression                            │
│    ✔ Test multiple Go versions                                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. SECURITY                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Run SAST (Static Application Security Testing)                       │
│    ✔ Scan dependencies for vulnerabilities                                │
│    ✔ Scan Docker images                                                  │
│    ✔ Use secrets management                                               │
│    ✔ Sign your commits and tags                                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. PERFORMANCE                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Cache Go modules and build artifacts                                 │
│    ✔ Use parallel jobs                                                    │
│    ✔ Use smaller runner images                                            │
│    ✔ Optimize Docker layers                                               │
│    ✔ Use buildx for cross-platform builds                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. DEPLOYMENT                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Use environment protection (branch restrictions)                     │
│    ✔ Implement blue-green or canary deployments                           │
│    ✔ Health checks after deployment                                       │
│    ✔ Automatic rollback on failure                                        │
│    ✔ Audit logging of deployments                                         │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	// ============================================================================
	// بخش 8: Troubleshooting
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 8: TROUBLESHOOTING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ COMMON ISSUES AND SOLUTIONS                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Issue 1: "cannot find package"                                            │
│ Solution: Run 'go mod download' before build                              │
│                                                                             │
│ Issue 2: Cache not working                                                │
│ Solution: Check cache key and paths                                       │
│                                                                             │
│ Issue 3: Tests pass locally but fail in CI                                │
│ Solution: Check environment variables, timeouts, test ordering            │
│                                                                             │
│ Issue 4: Docker build slow                                                │
│ Solution: Use multi-stage builds, layer caching, BuildKit                 │
│                                                                             │
│ Issue 5: Secrets not available                                            │
│ Solution: Check secret scope (protected branches only)                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

🔗 USEFUL LINKS:

   • GitLab CI Documentation: https://docs.gitlab.com/ee/ci/
   • GitHub Actions Documentation: https://docs.github.com/en/actions
   • Go Testing: https://pkg.go.dev/testing
   • Docker Buildx: https://docs.docker.com/buildx/working-with-buildx/
   • Trivy Scanner: https://github.com/aquasecurity/trivy
   • golangci-lint: https://golangci-lint.run/
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 CI/CD - COMPLETE")
	fmt.Println("Ready to automate your Go application delivery!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
