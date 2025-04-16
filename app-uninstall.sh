#!/bin/bash

set -e

export BACKEND_PORT=${BACKEND_PORT:-5000}
export FRONTEND_PORT=${FRONTEND_PORT:-3000}

echo "### Checking and killing processes running on BACKEND_PORT=$BACKEND_PORT and FRONTEND_PORT=$FRONTEND_PORT"

BACKEND_PID=$(sudo lsof -t -i :$BACKEND_PORT || true)
if [ -n "$BACKEND_PID" ]; then
    echo "Process ID for BACKEND_PORT=$BACKEND_PORT: $BACKEND_PID"
    echo "Killing process $BACKEND_PID running on BACKEND_PORT=$BACKEND_PORT"
    sudo kill -9 $BACKEND_PID || echo "Failed to kill process on BACKEND_PORT=$BACKEND_PORT"
else
    echo "No process found running on BACKEND_PORT=$BACKEND_PORT"
fi

FRONTEND_PID=$(sudo lsof -t -i :$FRONTEND_PORT || true)
if [ -n "$FRONTEND_PID" ]; then
    echo "Process ID for FRONTEND_PORT=$FRONTEND_PORT: $FRONTEND_PID"
    echo "Killing process $FRONTEND_PID running on FRONTEND_PORT=$FRONTEND_PORT"
    sudo kill -9 $FRONTEND_PID || echo "Failed to kill process on FRONTEND_PORT=$FRONTEND_PORT"
else
    echo "No process found running on FRONTEND_PORT=$FRONTEND_PORT"
fi

echo "### Stopping backend and frontend processes"
pkill -f "go run main.go" || echo "No backend process found"
pkill -f "react-scripts start" || echo "No frontend process found"
pkill -f "npx react-scripts start" || echo "No setsid frontend process found"

echo "### Removing Go installation"
sudo rm -rf /usr/local/go
sudo rm -f /etc/profile.d/go.sh

echo "### Removing Node.js and npm"
sudo yum remove -y nodejs

echo "### Removing application files"
rm -rf /opt/polling-app
rm -rf ~/PollApp

echo "### Cleaning up temporary files"
rm -rf /tmp/go*.linux-amd64.tar.gz

# echo "### Cleaning up environment files and logs"
# rm -f ~/PollApp/backend/PollApp/backend.log || echo "No backend log found"
# rm -f ~/PollApp/frontend/frontend.log || echo "No frontend log found"
# rm -f ~/PollApp/frontend/.env || echo "No .env file found"

echo "### Uninstall complete"
