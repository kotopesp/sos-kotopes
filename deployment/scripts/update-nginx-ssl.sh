#!/bin/bash
set -e

echo "=== Updating Nginx with SSL configuration ==="

sudo cp /opt/sos-kotopes/deployment/nginx/sos-kotopes-ssl.conf /etc/nginx/sites-available/sos-kotopes

sudo nginx -t
sudo systemctl reload nginx

echo "=== Nginx SSL configuration updated ==="