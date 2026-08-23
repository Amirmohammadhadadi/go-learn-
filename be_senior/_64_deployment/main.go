// ============================================================================
// FILE: deployment_guide.go
// TITLE: راهنمای کامل استقرار Go روی VPS و سرویس‌های ابری
// HOW TO RUN: go run deployment_guide.go (فایل توضیحی - دستورات در ترمینال اجرا شوند)
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - گزینه‌های استقرار Go
// ============================================================================
//
// گزینه‌های مختلف برای استقرار برنامه‌های Go:
//
// 1. VPS (Virtual Private Server)
//    - DigitalOcean Droplet, Linode, Vultr, Hetzner
//    - کنترل کامل روی سرور
//    - نیاز به مدیریت سرور (security updates, monitoring)
//
// 2. AWS EC2
//    - مقیاس‌پذیر، قابلیت‌های پیشرفته
//    - ELB, Auto-scaling, CloudWatch
//    - پیچیدگی بیشتر
//
// 3. Heroku / Platform-as-a-Service
//    - ساده‌ترین روش (git push)
//    - مدیریت خودکار (scaling, logging, SSL)
//    - محدودیت‌ها و هزینه بالاتر
//
// 4. Container-based (Docker + Registry)
//    - یکسان در همه محیط‌ها
//    - آسان برای orchestration (Kubernetes, ECS)
//
// قانون طلایی:
// "برای start سرپا، از VPS یا Heroku استفاده کن.
//  برای traffic بالا و نیاز به auto-scaling، از AWS استفاده کن.
//  همیشه از reverse proxy (Nginx/Caddy) در مقابل برنامه استفاده کن.
//  از systemd یا supervisor برای مدیریت process استفاده کن."
// ============================================================================

package __deployment

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🚀 DEPLOYMENT GUIDE FOR GO")
	fmt.Println("VPS | DigitalOcean | AWS EC2 | Heroku")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: آماده‌سازی برنامه برای استقرار
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 1: PREPARING THE APPLICATION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# 1.1 Create optimized binary
# ============================================================================
$ cd /path/to/your/project
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o myapp ./cmd/server

# Flags explanation:
# - CGO_ENABLED=0   : Static binary (no libc dependency)
# - GOOS=linux      : Target OS
# - GOARCH=amd64    : Target architecture
# - -ldflags="-s -w": Strip debug symbols (reduces binary size)

# 1.2 Check binary size and dependencies
$ ls -lh myapp
$ file myapp
$ ldd myapp  # should say "not a dynamic executable"

# 1.3 Create .env file (for configuration)
$ cat > .env << EOF
APP_ENV=production
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
EOF
`)

	// ============================================================================
	// بخش 2: استقرار روی VPS (DigitalOcean)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🖥️ SECTION 2: VPS DEPLOYMENT (DigitalOcean)")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# 2.1 Create DigitalOcean Droplet
# ============================================================================
# 1. Login to DigitalOcean
# 2. Create Droplet:
#    - Ubuntu 22.04 LTS
#    - Basic plan ($6-12/month)
#    - Add SSH key
# 3. Note the IP address

# 2.2 Initial server setup
# ============================================================================
$ ssh root@<YOUR_DROPLET_IP>

# Update system
$ apt update && apt upgrade -y

# Create deploy user
$ adduser deploy
$ usermod -aG sudo deploy
$ su - deploy

# Add SSH key for deploy user
$ mkdir -p ~/.ssh
$ echo "your-public-ssh-key" >> ~/.ssh/authorized_keys
$ chmod 600 ~/.ssh/authorized_keys

# 2.3 Install Go (if building on server)
# ============================================================================
$ wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
$ echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
$ echo 'export GOPATH=$HOME/go' >> ~/.profile
$ source ~/.profile

# 2.4 Deploy application
# ============================================================================
# From local machine
$ scp -r ./myapp deploy@<SERVER_IP>:~/
$ scp ./.env deploy@<SERVER_IP>:~/

# Or using rsync
$ rsync -avz --delete ./myapp deploy@<SERVER_IP>:~/app/

# 2.5 Create systemd service
# ============================================================================
$ sudo cat > /etc/systemd/system/myapp.service << EOF
[Unit]
Description=My Go Application
After=network.target

[Service]
Type=simple
User=deploy
WorkingDirectory=/home/deploy
EnvironmentFile=/home/deploy/.env
ExecStart=/home/deploy/myapp
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
$ sudo systemctl daemon-reload
$ sudo systemctl enable myapp
$ sudo systemctl start myapp

# Check status
$ sudo systemctl status myapp
$ sudo journalctl -u myapp -f

# 2.6 Install and configure Nginx (reverse proxy)
# ============================================================================
$ sudo apt install nginx -y

$ sudo cat > /etc/nginx/sites-available/myapp << EOF
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

$ sudo ln -s /etc/nginx/sites-available/myapp /etc/nginx/sites-enabled/
$ sudo rm /etc/nginx/sites-enabled/default
$ sudo nginx -t
$ sudo systemctl restart nginx

# 2.7 Setup HTTPS with Certbot (Let's Encrypt)
# ============================================================================
$ sudo apt install certbot python3-certbot-nginx -y
$ sudo certbot --nginx -d your-domain.com
$ sudo certbot renew --dry-run

# 2.8 Setup firewall
# ============================================================================
$ sudo ufw allow 22/tcp
$ sudo ufw allow 80/tcp
$ sudo ufw allow 443/tcp
$ sudo ufw enable
`)

	// ============================================================================
	// بخش 3: اتوماسیون با Script (deploy.sh)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🤖 SECTION 3: AUTOMATION WITH DEPLOY SCRIPT")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
#!/bin/bash
# ============================================================================
# deploy.sh - Script خودکار استقرار
# ============================================================================

set -e

# Variables
APP_NAME="myapp"
SERVER_USER="deploy"
SERVER_IP="YOUR_SERVER_IP"
SERVER_PATH="/home/$SERVER_USER/$APP_NAME"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Starting deployment...${NC}"

# 1. Build application
echo "Building application..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $APP_NAME ./cmd/server

# 2. Run tests
echo "Running tests..."
go test -v ./...

# 3. Create deployment package
echo "Creating deployment package..."
tar -czf $APP_NAME.tar.gz $APP_NAME .env

# 4. Upload to server
echo "Uploading to server..."
scp $APP_NAME.tar.gz $SERVER_USER@$SERVER_IP:$SERVER_PATH/

# 5. Deploy on server
echo "Deploying on server..."
ssh $SERVER_USER@$SERVER_IP << EOF
    cd $SERVER_PATH
    tar -xzf $APP_NAME.tar.gz
    rm $APP_NAME.tar.gz
    
    # Restart service
    sudo systemctl restart $APP_NAME
    
    # Check status
    sudo systemctl status $APP_NAME --no-pager
EOF

# 6. Cleanup
rm $APP_NAME.tar.gz
rm $APP_NAME

echo -e "${GREEN}Deployment completed!${NC}"
`)

	// ============================================================================
	// بخش 4: استقرار روی AWS EC2
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("☁️ SECTION 4: AWS EC2 DEPLOYMENT")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# 4.1 Install AWS CLI
# ============================================================================
$ pip install awscli
$ aws configure
# Enter: Access Key ID, Secret Key, Region (us-east-1), output format (json)

# 4.2 Create EC2 instance
# ============================================================================
$ aws ec2 run-instances \
    --image-id ami-0c55b159cbfafe1f0 \
    --instance-type t2.micro \
    --key-name my-key \
    --security-group-ids sg-12345678 \
    --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=myapp}]'

# 4.3 Get instance IP
# ============================================================================
$ aws ec2 describe-instances --filters "Name=tag:Name,Values=myapp" \
    --query "Reservations[].Instances[].PublicIpAddress"

# 4.4 Create Elastic IP (for static IP)
# ============================================================================
$ aws ec2 allocate-address
$ aws ec2 associate-address --instance-id <INSTANCE_ID> --public-ip <ELASTIC_IP>

# 4.5 Setup Application Load Balancer (for auto-scaling)
# ============================================================================
$ aws elbv2 create-load-balancer \
    --name myapp-lb \
    --subnets subnet-12345678 subnet-87654321 \
    --security-groups sg-12345678

$ aws elbv2 create-target-group \
    --name myapp-tg \
    --protocol HTTP \
    --port 8080 \
    --vpc-id vpc-12345678

$ aws elbv2 create-listener \
    --load-balancer-arn <LB_ARN> \
    --protocol HTTP \
    --port 80 \
    --default-actions Type=forward,TargetGroupArn=<TG_ARN>

# 4.6 Create Auto Scaling Group
# ============================================================================
$ aws autoscaling create-auto-scaling-group \
    --auto-scaling-group-name myapp-asg \
    --launch-template LaunchTemplateName=myapp-template \
    --min-size 2 \
    --max-size 10 \
    --desired-capacity 2 \
    --vpc-zone-identifier "subnet-12345678,subnet-87654321"

$ aws autoscaling put-scaling-policy \
    --auto-scaling-group-name myapp-asg \
    --policy-name myapp-scale-out \
    --policy-type TargetTrackingScaling \
    --target-tracking-configuration file://scaling-config.json

# 4.7 CloudWatch monitoring
# ============================================================================
$ aws cloudwatch put-metric-alarm \
    --alarm-name myapp-high-cpu \
    --comparison-operator GreaterThanThreshold \
    --evaluation-periods 2 \
    --metric-name CPUUtilization \
    --namespace AWS/EC2 \
    --period 60 \
    --statistic Average \
    --threshold 70 \
    --actions-enabled \
    --alarm-actions arn:aws:sns:us-east-1:123456789012:myapp-topic
`)

	// ============================================================================
	// بخش 5: استقرار روی Heroku (PaaS)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 5: HEROKU DEPLOYMENT")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# 5.1 Install Heroku CLI
# ============================================================================
$ curl https://cli-assets.heroku.com/install.sh | sh

# 5.2 Login to Heroku
# ============================================================================
$ heroku login

# 5.3 Create Procfile
# ============================================================================
$ cat > Procfile << EOF
web: ./myapp
EOF

# 5.4 Create app
# ============================================================================
$ heroku create myapp-name

# 5.5 Set environment variables
# ============================================================================
$ heroku config:set APP_ENV=production
$ heroku config:set DB_URL=postgres://...
$ heroku config:set JWT_SECRET=your-secret

# 5.6 Add add-ons
# ============================================================================
$ heroku addons:create heroku-postgresql:hobby-dev
$ heroku addons:create heroku-redis:hobby-dev

# 5.7 Deploy
# ============================================================================
$ git push heroku main

# 5.8 Scale dynos
# ============================================================================
$ heroku ps:scale web=2

# 5.9 Open app
# ============================================================================
$ heroku open

# 5.10 View logs
# ============================================================================
$ heroku logs --tail

# 5.11 Run one-off commands
# ============================================================================
$ heroku run go run cmd/migrate/main.go

# 5.12 Create pipeline (for staging/production)
# ============================================================================
$ heroku pipelines:create myapp-pipeline --app=myapp-staging --stage=staging
$ heroku pipelines:add myapp-pipeline --app=myapp-production --stage=production
`)

	// ============================================================================
	// بخش 6: Docker-based Deployment (برای همه پلتفرم‌ها)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🐳 SECTION 6: DOCKER DEPLOYMENT")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# 6.1 Multi-stage Dockerfile
# ============================================================================
$ cat > Dockerfile.prod << 'EOF'
# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/server

# Stage 2: Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Tehran

COPY --from=builder /app/app /app
COPY --from=builder /app/.env /app/.env

EXPOSE 8080

CMD ["/app"]
EOF

# 6.2 Build and push to registry
# ============================================================================
$ docker build -f Dockerfile.prod -t myapp:latest .
$ docker tag myapp:latest myregistry/myapp:latest
$ docker push myregistry/myapp:latest

# 6.3 Deploy with docker-compose on VPS
# ============================================================================
$ cat > docker-compose.yml << EOF
version: '3.8'

services:
  app:
    image: myregistry/myapp:latest
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: \${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  pgdata:
EOF

$ docker-compose up -d
`)

	// ============================================================================
	// بخش 7: Continuous Deployment (GitHub Actions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 7: CONTINUOUS DEPLOYMENT")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# .github/workflows/deploy.yml
# ============================================================================

name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Build
        run: |
          CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myapp ./cmd/server

      - name: Deploy to VPS
        uses: appleboy/scp-action@master
        with:
          host: \${{ secrets.SERVER_HOST }}
          username: \${{ secrets.SERVER_USER }}
          key: \${{ secrets.SSH_KEY }}
          source: "myapp,.env"
          target: "/home/deploy/myapp"

      - name: Restart service
        uses: appleboy/ssh-action@master
        with:
          host: \${{ secrets.SERVER_HOST }}
          username: \${{ secrets.SERVER_USER }}
          key: \${{ secrets.SSH_KEY }}
          script: |
            sudo systemctl restart myapp

  deploy-heroku:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: akhileshns/heroku-deploy@v3.13.15
        with:
          heroku_api_key: \${{ secrets.HEROKU_API_KEY }}
          heroku_app_name: \${{ secrets.HEROKU_APP_NAME }}
          heroku_email: \${{ secrets.HEROKU_EMAIL }}
`)

	// ============================================================================
	// بخش 8: Monitoring & Logging بعد از استقرار
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 8: POST-DEPLOYMENT MONITORING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# 8.1 Install monitoring tools
# ============================================================================
# Prometheus + Grafana
$ docker run -d --name prometheus -p 9090:9090 prom/prometheus
$ docker run -d --name grafana -p 3000:3000 grafana/grafana

# 8.2 Setup log aggregation (Loki)
# ============================================================================
$ docker run -d --name loki -p 3100:3100 grafana/loki

# 8.3 Setup uptime monitoring (UptimeRobot)
# ============================================================================
# 1. Sign up at uptimerobot.com
# 2. Add monitor for https://your-domain.com/health
# 3. Configure alert contacts

# 8.4 Setup alerting
# ============================================================================
# Configure Prometheus alerts
$ cat > alerts.yml << EOF
groups:
  - name: myapp
    rules:
      - alert: InstanceDown
        expr: up{job="myapp"} == 0
        for: 1m
        annotations:
          summary: "Application is down"
EOF
`)

	// ============================================================================
	// بخش 9: Troubleshooting
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 9: TROUBLESHOOTING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ COMMON ISSUES AND SOLUTIONS                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Issue 1: "port already in use"                                            │
│ Solution: $ sudo lsof -i :8080                                            │
│           $ kill -9 <PID>                                                 │
│                                                                             │
│ Issue 2: "connection refused"                                             │
│ Solution: Check firewall: sudo ufw status                                 │
│           Check service: sudo systemctl status myapp                      │
│                                                                             │
│ Issue 3: "too many open files"                                            │
│ Solution: $ ulimit -n 65536                                               │
│           Add to /etc/security/limits.conf:                               │
│           * soft nofile 65536                                              │
│           * hard nofile 65536                                              │
│                                                                             │
│ Issue 4: "out of memory"                                                  │
│ Solution: Check memory usage: free -h                                     │
│           Add swap: fallocate -l 2G /swapfile                             │
│                                                                             │
│ Issue 5: SSL certificate expired                                          │
│ Solution: $ sudo certbot renew                                            │
│           Add to crontab: 0 0 * * * certbot renew --quiet                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	// ============================================================================
	// بخش 10: Best Practices
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 DEPLOYMENT BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. SECURITY                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use SSH keys (disable password authentication)                       │
│    • Run services as non-root user                                        │
│    • Use firewall (ufw)                                                   │
│    • Keep system updated                                                  │
│    • Use environment variables for secrets                                │
│    • Enable HTTPS with Let's Encrypt                                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. MONITORING                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Set up health checks (/health)                                       │
│    • Monitor resources (CPU, Memory, Disk)                                │
│    • Configure alerts                                                     │
│    • Centralize logs                                                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. BACKUP                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Backup database daily                                                │
│    • Backup configuration files                                           │
│    • Store backups off-server                                             │
│    • Test restore procedure                                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. PERFORMANCE                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use reverse proxy (Nginx)                                            │
│    • Enable gzip compression                                              │
│    • Set up CDN for static assets                                         │
│    • Use connection pooling                                               │
│    • Implement caching strategy                                           │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	// ============================================================================
	// بخش 11: جمع‌بندی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ PLATFORM        │ COMMAND / METHOD              │ WHEN TO USE              │
├─────────────────┼───────────────────────────────┼──────────────────────────┤
│ VPS (DigitalOcean)│ systemd + Nginx             │ Full control, budget     │
│ AWS EC2         │ Auto Scaling + LB             │ High traffic, enterprise │
│ Heroku          │ git push heroku              │ Fast deployment, simple  │
│ Docker          │ docker-compose               │ Portability, microservices│
└─────────────────┴───────────────────────────────┴──────────────────────────┘

🔗 USEFUL LINKS:

   • DigitalOcean: https://digitalocean.com
   • AWS Free Tier: https://aws.amazon.com/free
   • Heroku: https://heroku.com
   • Let's Encrypt: https://letsencrypt.org
   • Nginx: https://nginx.org
   • Prometheus: https://prometheus.io
   • Grafana: https://grafana.com
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 DEPLOYMENT GUIDE - COMPLETE")
	fmt.Println("Your Go application is ready for production!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
