#!/bin/bash
set -e

echo "=== Setting up Nginx reverse proxy ==="

# Copy nginx config
sudo cp /opt/sos-kotopes/deployment/nginx/sos-kotopes.conf /etc/nginx/sites-available/sos-kotopes

# Enable site
sudo ln -sf /etc/nginx/sites-available/sos-kotopes /etc/nginx/sites-enabled/

# Remove default nginx site
sudo rm -f /etc/nginx/sites-enabled/default

# Test and reload nginx
sudo nginx -t
sudo systemctl reload nginx

echo "=== Nginx setup completed ==="
echo "Config file: /etc/nginx/sites-available/sos-kotopes"
