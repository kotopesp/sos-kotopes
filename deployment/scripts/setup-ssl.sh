#!/bin/bash
set -e

echo "=== Setting up SSL certificates ==="

sudo mkdir -p /var/www/certbot
sudo chown -R www-data:www-data /var/www/certbot

sudo certbot certonly --webroot --non-interactive --agree-tos \
    --email admin@sos-kotopes.ru \
    -w /var/www/certbot \
    -d sos-kotopes.ru \
    -d www.sos-kotopes.ru

if [ ! -f "/etc/letsencrypt/live/sos-kotopes.ru/fullchain.pem" ]; then
    echo "ERROR: SSL certificates were not created!"
    exit 1
fi

echo "=== SSL setup completed ==="