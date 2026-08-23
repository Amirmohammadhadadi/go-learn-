// ============================================================================
// FILE: docker_guide.go
// TITLE: راهنمای کامل Docker و Docker Compose برای Go - Multi-Stage Build
// HOW TO RUN: go run docker_guide.go (فایل توضیحی - دستورات در ترمینال اجرا شوند)
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Docker و Multi-Stage Build چیست؟
// ============================================================================
//
// Docker: پلتفرمی برای打包 و اجرای برنامه‌ها در کانتینرهای ایزوله
//
// Multi-Stage Build: تکنیکی برای ساخت imageهای کوچک و امن
//
// مزایای Multi-Stage Build:
// 1. کاهش حجم image (از ~1GB به ~10MB)
// 2. عدم inclusion ابزارهای build در image نهایی
// 3. امنیت بیشتر (بدون compiler در image نهایی)
// 4. لایه‌های کمتر و build سریع‌تر
//
// مراحل Multi-Stage Build:
// مرحله 1: Builder - کامپایل برنامه با تمام وابستگی‌ها
// مرحله 2: Runtime - اجرای باینری کامپایل شده
//
// قانون طلایی:
// "از Alpine Linux برای runtime image استفاده کن (حجم کم).
//  از .dockerignore برای排除 فایل‌های غیرضروری استفاده کن.
//  برنامه را statically compile کن تا وابستگی به libc نداشته باشی."
// ============================================================================

package __docker

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🐳 DOCKER & DOCKER COMPOSE FOR GO")
	fmt.Println("Multi-Stage Build | Optimization | Production Ready")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: Dockerfile بهینه (Multi-Stage)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 SECTION 1: OPTIMIZED DOCKERFILE (Multi-Stage)")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# Dockerfile (Multi-Stage Build)
# ============================================================================

# مرحله 1: Builder (کامپایل برنامه)
# ============================================================================
FROM golang:1.24-alpine AS builder

# نصب ابزارهای مورد نیاز (اختیاری)
RUN apk add --no-cache git ca-certificates tzdata

# تنظیم دایرکتوری کاری
WORKDIR /app

# کپی فایل‌های go.mod و go.sum (برای کش لایه‌ها)
COPY go.mod go.sum ./

# دانلود وابستگی‌ها
RUN go mod download && go mod verify

# کپی کل کد منبع
COPY . .

# آپشن‌های کامپایل برای image کوچک‌تر
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOAMD64=v2

# کامپایل برنامه
RUN go build -ldflags="-s -w" -o /app/bin/app ./cmd/server

# مرحله 2: Runtime (اجرای برنامه)
# ============================================================================
FROM alpine:3.19 AS runtime

# نصب گواهی‌های SSL و timezone data
RUN apk add --no-cache ca-certificates tzdata

# ایجاد کاربر غیر-root برای اجرای امن
RUN adduser -D -g '' appuser

# تنظیم timezone (اختیاری)
ENV TZ=Asia/Tehran

# کپی باینری از مرحله builder
COPY --from=builder /app/bin/app /app

# کپی فایل‌های استاتیک (در صورت وجود)
# COPY --from=builder /app/static /static

# استفاده از کاربر غیر-root
USER appuser

# پورت برنامه
EXPOSE 8080

# نقطه ورود
ENTRYPOINT ["/app"]

# CMD (قابل override)
CMD ["--config", "/etc/app/config.yaml"]
`)

	// ============================================================================
	// بخش 2: Dockerfile با Debug Support
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 2: DOCKERFILE WITH DEBUG SUPPORT")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# Dockerfile.debug (برای دیباگ با Delve)
# ============================================================================

# مرحله 1: Builder با ابزارهای دیباگ
# ============================================================================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# کامپایل با پشتیبانی دیباگ (بدون بهینه‌سازی)
RUN go build -gcflags="all=-N -l" -o /app/bin/app ./cmd/server

# مرحله 2: Runtime با دیباگر
# ============================================================================
FROM alpine:3.19 AS debug

RUN apk add --no-cache ca-certificates tzdata curl

# نصب Delve (دیباگر Go)
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY --from=builder /app/bin/app /app

# پورت دیباگ
EXPOSE 8080 40000

# اجرا با Delve
ENTRYPOINT ["dlv", "exec", "/app", "--headless", "--listen=:40000", "--api-version=2", "--accept-multiclient"]
`)

	// ============================================================================
	// بخش 3: .dockerignore (فایلهای exclude)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚫 SECTION 3: .dockerignore")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# .dockerignore (بهبود سرعت build و امنیت)
# ============================================================================

# Version control
.git/
.gitignore
.gitattributes

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# Documentation
docs/
*.md
README.md
LICENSE

# Build artifacts
bin/
dist/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Testing
coverage.out
*.test
*.prof
*.pprof
*.cov
*.out

# Dependencies
vendor/

# Environment files
.env
.env.*
*.env
*.env.local

# Logs
*.log
logs/

# Temporary files
tmp/
temp/
*.tmp

# OS files
.DS_Store
Thumbs.db
desktop.ini

# Kubernetes
kubernetes/
helm/
*.yaml
*.yml

# Docker
Dockerfile*
.dockerignore
docker-compose*.yml

# CI/CD
.github/
.gitlab-ci.yml
Jenkinsfile
`)

	// ============================================================================
	// بخش 4: Docker Compose برای توسعه
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚙️ SECTION 4: DOCKER COMPOSE FOR DEVELOPMENT")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# docker-compose.yml (توسعه)
# ============================================================================

version: '3.8'

services:
  # اپلیکیشن Go
  app:
    build:
      context: .
      dockerfile: Dockerfile
      target: runtime  # یا debug برای دیباگ
    container_name: go_app
    ports:
      - "8080:8080"
      - "40000:40000"  # برای دیباگ
    environment:
      - APP_ENV=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=appdb
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    volumes:
      - ./configs:/etc/app/configs:ro
      - ./logs:/var/log/app
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    networks:
      - app_network

  # PostgreSQL
  postgres:
    image: postgres:16-alpine
    container_name: postgres_db
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=appdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app_network

  # Redis
  redis:
    image: redis:7-alpine
    container_name: redis_cache
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app_network

volumes:
  postgres_data:
  redis_data:

networks:
  app_network:
    driver: bridge
`)

	// ============================================================================
	// بخش 5: Docker Compose برای پروداکشن
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🏭 SECTION 5: DOCKER COMPOSE FOR PRODUCTION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# docker-compose.prod.yml (پروداکشن)
# ============================================================================

version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
      target: runtime
    image: registry.example.com/app:${VERSION:-latest}
    container_name: go_app_prod
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
    volumes:
      - ./configs/prod:/etc/app/configs:ro
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    networks:
      - app_network

  postgres:
    image: postgres:16-alpine
    container_name: postgres_prod
    environment:
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 30s
      timeout: 10s
      retries: 5
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
    networks:
      - app_network

  redis:
    image: redis:7-alpine
    container_name: redis_prod
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
    networks:
      - app_network

  nginx:
    image: nginx:alpine
    container_name: nginx_proxy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/prod.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - app
    networks:
      - app_network

volumes:
  postgres_data:
  redis_data:

networks:
  app_network:
    driver: overlay
`)

	// ============================================================================
	// بخش 6: Makefile برای خودکارسازی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 6: MAKEFILE FOR AUTOMATION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# Makefile (خودکارسازی دستورات Docker)
# ============================================================================

.PHONY: help build run stop clean dev prod logs shell test

help:  ## نمایش کمک
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build:  ## ساخت image
	@echo "Building Docker image..."
	docker build -t myapp:latest -f Dockerfile .

build-prod:  ## ساخت image پروداکشن
	@echo "Building production image..."
	docker build --target runtime -t myapp:prod -f Dockerfile .

run:  ## اجرای کانتینر
	@echo "Running container..."
	docker-compose up -d

stop:  ## توقف کانتینر
	@echo "Stopping containers..."
	docker-compose down

clean:  ## پاک کردن کانتینرها و imageها
	@echo "Cleaning up..."
	docker-compose down -v
	docker system prune -f

dev:  ## اجرا در حالت توسعه (با hot-reload)
	@echo "Starting development environment..."
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

prod:  ## اجرا در حالت پروداکشن
	@echo "Starting production environment..."
	docker-compose -f docker-compose.prod.yml up -d

logs:  ## نمایش لاگ‌ها
	docker-compose logs -f

shell:  ## ورود به شل کانتینر
	docker-compose exec app sh

test:  ## اجرای تست‌ها در کانتینر
	docker-compose run --rm app go test -v ./...

build-multi:  ## ساخت image برای معماری‌های مختلف
	@echo "Building multi-architecture image..."
	docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest --push .

size:  ## نمایش حجم image
	docker images | grep myapp

scan:  ## اسکن image برای vulnerabilities
	docker scan myapp:latest

save:  ## ذخیره image به فایل
	docker save myapp:latest -o myapp.tar

load:  ## بارگذاری image از فایل
	docker load -i myapp.tar
`)

	// ============================================================================
	// بخش 7: Optimizations (بهبودهای عملکردی)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 7: DOCKER OPTIMIZATIONS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. IMAGE SIZE OPTIMIZATION                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use Alpine Linux (5MB vs 100MB)                                      │
│    • Multi-stage builds                                                   │
│    • Static compilation (CGO_ENABLED=0)                                   │
│    • Strip debug symbols (-ldflags="-s -w")                               │
│    • Use distroless images for ultra-small images                         │
│                                                                           │
│    Example: FROM gcr.io/distroless/static                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. BUILD PERFORMANCE                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Order layers by frequency of change                                  │
│    • Use Docker BuildKit (DOCKER_BUILDKIT=1)                              │
│    • Cache go mod download                                                │
│    • Use --cache-from for CI/CD pipelines                                 │
│    • Parallel builds with --parallel                                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. LAYER CACHING                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    # بد (کمتر caching)                                                   │
│    COPY . .                                                               │
│    RUN go build                                                           │
│                                                                           │
│    # خوب (بهترین caching)                                                │
│    COPY go.mod go.sum ./                                                  │
│    RUN go mod download                                                    │
│    COPY . .                                                               │
│    RUN go build                                                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. SECURITY                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Run as non-root user                                                 │
│    • Use specific version tags (not 'latest')                             │
│    • Regularly update base images                                         │
│    • Scan images for vulnerabilities                                      │
│    • Use secrets management (Docker secrets)                              │
│    • Read-only root filesystem                                            │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	// ============================================================================
	// بخش 8: Common Commands
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 SECTION 8: COMMON DOCKER COMMANDS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ IMAGE COMMANDS                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│ docker build -t myapp:latest .                                           │
│ docker build --target runtime -t myapp:latest .                          │
│ docker build --no-cache -t myapp:latest .                                │
│ docker build --build-arg VERSION=1.0 -t myapp:1.0 .                      │
│ docker images                                                             │
│ docker rmi myapp:latest                                                   │
│ docker save myapp:latest -o myapp.tar                                     │
│ docker load -i myapp.tar                                                  │
│ docker tag myapp:latest myregistry/myapp:latest                           │
│ docker push myregistry/myapp:latest                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CONTAINER COMMANDS                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│ docker run -d --name myapp -p 8080:8080 myapp:latest                     │
│ docker run -it --rm myapp:latest sh                                       │
│ docker ps                                                                 │
│ docker ps -a                                                              │
│ docker stop myapp                                                         │
│ docker start myapp                                                        │
│ docker restart myapp                                                      │
│ docker rm myapp                                                           │
│ docker logs -f myapp                                                      │
│ docker exec -it myapp sh                                                  │
│ docker cp myapp:/app/config.yaml ./                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ DOCKER COMPOSE COMMANDS                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│ docker-compose up -d                                                      │
│ docker-compose down                                                       │
│ docker-compose down -v (حذف volumes)                                      │
│ docker-compose restart                                                    │
│ docker-compose logs -f                                                    │
│ docker-compose exec app sh                                                │
│ docker-compose ps                                                         │
│ docker-compose build --no-cache                                           │
│ docker-compose pull                                                       │
│ docker-compose config                                                     │
│ docker-compose -f docker-compose.prod.yml up -d                           │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	// ============================================================================
	// بخش 9: Best Practices
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 DOCKER BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. DOCKERFILE                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Use specific base image tags (not 'latest')                         │
│    ✔ Use multi-stage builds                                               │
│    ✔ Order layers by frequency of change                                  │
│    ✔ Use .dockerignore                                                    │
│    ✔ Run as non-root user                                                 │
│    ✔ Use COPY instead of ADD (unless needed)                              │
│    ✔ Combine RUN commands (apt-get update && apt-get install)             │
│    ✔ Remove temporary files in same RUN layer                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. SECURITY                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Use specific version tags                                            │
│    ✔ Regularly update base images                                         │
│    ✔ Scan images with docker scan                                         │
│    ✔ Use secrets management, not env vars for secrets                     │
│    ✔ Run containers as non-root                                           │
│    ✔ Use read-only root filesystem                                        │
│    ✔ Drop unnecessary capabilities                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. MONITORING                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Set resource limits (--memory, --cpus)                              │
│    ✔ Configure logging driver and limits                                  │
│    ✔ Use health checks                                                   │
│    ✔ Monitor container metrics                                           │
│    ✔ Use labels for metadata                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. CI/CD PIPELINE                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Cache Docker layers                                                 │
│    ✔ Use multi-stage builds in CI                                        │
│    ✔ Tag images with commit SHA and version                              │
│    ✔ Push to registry                                                    │
│    ✔ Scan for vulnerabilities                                            │
│    ✔ Test images before deployment                                       │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	// ============================================================================
	// بخش 10: Example Workflow
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 10: COMPLETE WORKFLOW")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ DEVELOPMENT WORKFLOW                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. Development with hot-reload                                            │
│     $ docker-compose -f docker-compose.yml -f docker-compose.dev.yml up    │
│                                                                             │
│  2. Run tests                                                              │
│     $ docker-compose run --rm app go test -v ./...                         │
│                                                                             │
│  3. Build production image                                                 │
│     $ docker build --target runtime -t myapp:${VERSION} .                 │
│                                                                             │
│  4. Test image                                                             │
│     $ docker run --rm myapp:${VERSION} go test -v ./...                    │
│                                                                             │
│  5. Scan for vulnerabilities                                               │
│     $ docker scan myapp:${VERSION}                                         │
│                                                                             │
│  6. Push to registry                                                       │
│     $ docker tag myapp:${VERSION} registry.example.com/myapp:${VERSION}   │
│     $ docker push registry.example.com/myapp:${VERSION}                   │
│                                                                             │
│  7. Deploy to production                                                   │
│     $ docker-compose -f docker-compose.prod.yml pull                       │
│     $ docker-compose -f docker-compose.prod.yml up -d                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

📝 ENVIRONMENT VARIABLES (.env):

   # Development
   APP_ENV=development
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=appdb

   # Production
   APP_ENV=production
   DB_HOST=prod-db.example.com
   DB_PORT=5432
   DB_USER=appuser
   DB_PASSWORD=securepassword123
   DB_NAME=appdb_prod

🔗 USEFUL RESOURCES:

   • Docker Hub: https://hub.docker.com
   • Official Go image: https://hub.docker.com/_/golang
   • Official Alpine image: https://hub.docker.com/_/alpine
   • Multi-stage builds: https://docs.docker.com/build/building/multi-stage/
   • Docker Compose: https://docs.docker.com/compose/
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🐳 DOCKER & DOCKER COMPOSE - COMPLETE")
	fmt.Println("Ready to containerize your Go applications!")
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
