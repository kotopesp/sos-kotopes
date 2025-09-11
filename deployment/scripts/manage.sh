#!/bin/bash

case "$1" in
    start)
        echo "Starting services..."
        cd /opt/sos-kotopes
        docker-compose up -d
        ;;
    stop)
        echo "Stopping services..."
        cd /opt/sos-kotopes
        docker-compose down
        ;;
    restart)
        echo "Restarting services..."
        cd /opt/sos-kotopes
        docker-compose restart
        ;;
    logs)
        echo "Showing logs..."
        cd /opt/sos-kotopes
        docker-compose logs -f
        ;;
    status)
        echo "Services status:"
        cd /opt/sos-kotopes
        docker-compose ps
        ;;
    update)
        echo "Updating from git..."
        cd /opt/sos-kotopes
        git pull origin production
        /opt/sos-kotopes/deployment/scripts/deploy.sh
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|logs|status|update}"
        exit 1
        ;;
esac
