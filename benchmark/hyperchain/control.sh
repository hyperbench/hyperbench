#!/bin/bash

set -e

Usage="bash control.sh {stop|start|restart}"

case "$1" in
    start)
        echo "Start the Docker containers."
        docker-compose -f docker-compose.yml up -d
        echo "Started"
        ;;
    stop)
        echo "Stop the Docker containers."
        docker-compose -f docker-compose.yml stop
        echo "Stopped"
        ;;
    restart)
        echo "Shut down the Docker containers for the system tests."
        docker-compose -f docker-compose.yml kill && docker-compose -f docker-compose.yml down

        docker-compose -f docker-compose.yml up -d
        echo "Restarted"
        ;;
    *)
        echo "Usage: $Usage"
        ;;
esac

