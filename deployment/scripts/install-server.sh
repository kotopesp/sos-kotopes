#!/bin/bash
set -e

echo "=== SOS-KOTOPES SERVER INSTALLATION ==="
echo "Using existing deployment files from /deployment/"
echo "========================================"

echo "1. Updating system packages..."
sudo apt update
sudo apt upgrade -y

echo "2. Installing required packages..."
sudo apt install -y curl git nginx
sudo apt install -y python3-certbot-nginx
sudo apt install -y ufw

echo "3. Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

echo "4. Installing Docker Compose..."
DOCKER_COMPOSE_VERSION="v2.24.5"
sudo curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

echo "5. Configuring firewall..."
sudo ufw allow ssh
sudo ufw allow http
sudo ufw allow https
echo "y" | sudo ufw enable

PRODUCTION_BRAHCN='production'
echo "6. Cloning repository to /opt/..."
if [ ! -d "/opt" ]; then
    echo "Creating /opt directory..."
    sudo mkdir -p /opt
    sudo chown $(whoami):$(whoami) /opt 
fi

cd /opt || { echo "ERROR: Failed to cd to /opt"; exit 1; }

if [ ! -d "/opt/sos-kotopes" ]; then
    git clone https://github.com/kotopesp/sos-kotopes.git
    cd sos-kotopes || { echo "ERROR: Failed to cd to sos-kotopes"; exit 1; }
    echo "Switching to production branch..."
    git checkout ${PRODUCTION_BRAHCN}
else
    echo "Directory /opt/sos-kotopes already exists"
    cd sos-kotopes || { echo "ERROR: Failed to cd to sos-kotopes"; exit 1; }
    if [ ! -d ".git" ]; then
        echo "ERROR: /opt/sos-kotopes exists but is not a git repository!"
        exit 1
    fi
    echo "Updating existing repository and switching to production branch..."
    git fetch --all
    git checkout ${PRODUCTION_BRAHCN}
    git pull origin ${PRODUCTION_BRAHCN}
fi

echo "7. Verifying deployment files..."
DEPLOYMENT_FILES=(
    "deployment/nginx/sos-kotopes.conf"
    "deployment/scripts/deploy.sh"
    "deployment/scripts/setup-nginx.sh"
    "deployment/scripts/manage.sh"
)

for file in "${DEPLOYMENT_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo "ERROR: Required file $file not found!"
        exit 1
    fi
    echo "✓ Found: $file"
done

echo "8. Setting executable permissions..."
chmod +x deployment/scripts/*.sh

echo "9. Setting up Nginx using deployment/scripts/setup-nginx.sh..."
./deployment/scripts/setup-nginx.sh

echo "10. Creating command symlinks..."
sudo ln -sf /opt/sos-kotopes/deployment/scripts/deploy.sh /usr/local/bin/sos-deploy
sudo ln -sf /opt/sos-kotopes/deployment/scripts/manage.sh /usr/local/bin/sos-manage
sudo ln -sf /opt/sos-kotopes/deployment/scripts/setup-nginx.sh /usr/local/bin/sos-setup-nginx

echo "11. Setting up database backups..."
sudo ln -sf /opt/sos-kotopes/deployment/scripts/backup-db.sh /usr/local/bin/sos-backup
sudo ln -sf /opt/sos-kotopes/deployment/scripts/clean-backups.sh /usr/local/bin/sos-clean-backups
sudo ln -sf /opt/sos-kotopes/deployment/scripts/restore-db.sh /usr/local/bin/sos-restore-db

mkdir -p /opt/backups

echo "12. Setting up automatic backups..."
sudo cp /opt/sos-kotopes/deployment/systemd/sos-kotopes-backup.cron /etc/cron.d/
sudo chmod 644 /etc/cron.d/sos-kotopes-backup.cron
sudo touch /var/log/sos-kotopes-backup.log
sudo chmod 644 /var/log/sos-kotopes-backup.log



echo "13. Checking environment configuration..."
if [ ! -f .env ]; then
    PASSWORD=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
    cat > .env << 'ENV_EOF'
POSTGRES_USER=soskot_prod_user
POSTGRES_PASSWORD=${PASSWORD}
POSTGRES_DB=soskot_production
LOG_LEVEL=info
PORT=:8080
VK_CLIENT_ID=52010687
VK_CLIENT_SECRET=M6SWkM8KmIIJA60hTVx1
ENV_EOF
    chmod 600 .env
    echo ".env file created with secure passwords"
else
    echo ".env file already exists"
fi

echo "14. Running initial deployment using deployment/scripts/deploy.sh..."
./deployment/scripts/deploy.sh

echo "15 Enable HTTPS"
sudo certbot --nginx -d sos-kotopes.ru

echo "16. Final verification..."
echo "Waiting for services to start..."
sleep 5

echo "=== INSTALLATION COMPLETE ==="
echo ""
echo "Quick commands:"
echo "  sos-deploy          - Redeploy application"
echo "  sos-manage status   - Check service status"
echo "  sos-manage logs     - View application logs"
echo "  sos-manage restart  - Restart services"
echo ""
echo "Nginx commands:"
echo "  sos-setup-nginx    - Reconfigure nginx"
echo "  sudo nginx -t      - Test nginx configuration"
echo "  sudo systemctl reload nginx - Reload nginx"
echo ""
echo "Application URLs:"
echo "  Frontend: http://localhost:4200"
echo "  Backend:  http://localhost:8080"
echo "  Nginx:    http://localhost"
echo "========================================"
