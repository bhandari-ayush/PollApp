#!/bin/bash

set -e


GO_VERSION="1.24.1"
APP_DIR="/opt/polling-app/backend"  
MAIN_FILE="main.go"                 
BACKEND_IP=$(hostname -I | awk '{print $1}') 

echo "### Updating system packages"
sudo yum update -y

echo "### Installing dependencies for Go and Git"
sudo yum install -y wget git tar bash-completion

echo "### Removing any existing Go installation"
sudo rm -rf /usr/local/go

echo "### Downloading Go $GO_VERSION"
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -P /tmp

echo "### Installing Go"
sudo tar -C /usr/local -xzf /tmp/go${GO_VERSION}.linux-amd64.tar.gz

echo "### Setting up Go environment variables"
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh
export PATH=$PATH:/usr/local/go/bin

echo "### Verifying Go installation"
go version

git clone https://github.com/bhandari-ayush/PollApp.git
cd PollApp/backend/PollApp

echo "### Initializing Go modules and downloading dependencies"
go mod tidy

export DB_HOST=@@{Postgres.address}@@
export DB_USER=root
export BACKEND_PORT=8080
export FRONTEND_PORT=3000
export ENV=prod

echo "### Running the backend server"
nohup go run $MAIN_FILE > backend.log 2>&1 &


echo "Installing Node.js and npm..."
curl -sL https://rpm.nodesource.com/setup_16.x | sudo bash -
sudo yum install -y nodejs

echo "Verifying Node.js and npm installation..."
node -v
npm -v

cd ../ui

echo "REACT_APP_BACKEND_URL=http://$BACKEND_IP:$BACKEND_PORT" > .env

echo "Installing dependencies..."
npm install

echo "Starting the React development server..."
setsid env BROWSER=none HOST=0.0.0.0 PORT=3000 npx react-scripts start > frontend.log 2>&1 < /dev/null &

sleep 20

echo "### Setup complete"
echo "Backend is running on http://$BACKEND_IP:$BACKEND_PORT"
echo "Frontend is running on http://$BACKEND_IP:$FRONTEND_PORT"