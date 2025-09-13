#!/bin/bash

case "$1" in
    start|stop|restart|logs|status|update|backup)
        if [ ! -d "/opt/sos-kotopes" ]; then
            echo "ERROR: Project directory /opt/sos-kotopes not found!"
            exit 1
        fi
        cd /opt/sos-kotopes || { echo "ERROR: Cannot cd to /opt/sos-kotopes"; exit 1; }
        ;;
esac

case "$1" in
    start)
        echo "Starting services..."
        docker-compose up -d
        ;;
    stop)
        echo "Stopping services..."
        docker-compose down
        ;;
    restart)
        echo "Restarting services..."
        docker-compose restart
        ;;
    logs)
        echo "Showing logs..."
        docker-compose logs -f
        ;;
    status)
        echo "Services status:"
        docker-compose ps
        ;;
    update)
        echo "Updating from git..."
        
        if ! git fetch origin production; then
            echo "ERROR: Cannot fetch from git"
            exit 1
        fi
        
        CURRENT_BRANCH=$(git branch --show-current)
        if [ "$CURRENT_BRANCH" != "production" ]; then
            echo "Warning: Not on production branch (current: $CURRENT_BRANCH)"
        fi
        # !!!MAY BE CHANGED TO MASTER BRANCH!!!
        git pull origin production
        
        if [ ! -f "deployment/scripts/deploy.sh" ]; then
            echo "ERROR: deploy.sh not found!"
            exit 1
        fi
        
        ./deployment/scripts/deploy.sh
        ;;
    backup)
        echo "Creating database backup..."
        if [ -f "deployment/scripts/backup-db.sh" ]; then
            ./deployment/scripts/backup-db.sh
        else
            echo "ERROR: backup-db.sh not found!"
            exit 1
        fi
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|logs|status|update|backup}"
        echo ""
        echo "Commands:"
        echo "  start     - Start services"
        echo "  stop      - Stop services"
        echo "  restart   - Restart services"
        echo "  logs      - Show logs (follow mode)"
        echo "  status    - Show services status"
        echo "  update    - Update from git and deploy"
        echo "  backup    - Create database backup"
        exit 1
        ;;
esac
