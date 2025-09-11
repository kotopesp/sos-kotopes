#!/bin/bash
set -e

echo "=== Setting up Nginx reverse proxy ==="

sudo cp /opt/sos-kotopes/deployment/nginx/sos-kotopes.conf /etc/nginx/sites-available/sos-kotopes

sudo ln -sf /etc/nginx/sites-available/sos-kotopes /etc/nginx/sites-enabled/

sudo rm -f /etc/nginx/sites-enabled/default

sudo nginx -t
sudo systemctl reload nginx

echo "=== Nginx setup completed ==="
echo "Config file: /etc/nginx/sites-available/sos-kotopes"
