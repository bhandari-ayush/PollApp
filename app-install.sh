#!/bin/bash

set -e

GO_VERSION="1.24.1"
APP_DIR="/opt/polling-app/backend"  
MAIN_FILE="main.go"                 
APP_IP=$(hostname -I | awk '{print $1}') 

echo "### Updating system packages"
sudo yum update -y

echo "### Installing dependencies for APP"
if ! command -v wget &> /dev/null; then
    sudo yum install -y wget
else
    echo "wget is already installed"
fi

if ! command -v git &> /dev/null; then
    sudo yum install -y git
else
    echo "git is already installed"
fi

if ! command -v tar &> /dev/null; then
    sudo yum install -y tar
else
    echo "tar is already installed"
fi

if ! rpm -q bash-completion &> /dev/null; then
    sudo yum install -y bash-completion
else
    echo "bash-completion is already installed"
fi

echo "### Checking and installing lsof if not already installed"
if ! command -v lsof &> /dev/null; then
    sudo yum install -y lsof
else
    echo "lsof is already installed"
fi

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

if [ -d "PollApp" ]; then
    echo "### Removing existing PollApp directory"
    rm -rf PollApp
fi

echo "### Cloning PollApp repository"
git clone https://github.com/bhandari-ayush/PollApp.git
cd PollApp/backend/PollApp

echo "### Initializing Go modules and downloading dependencies"
go mod tidy

export DB_HOST=@@{DB.address}@@
export DB_USER=root
export BACKEND_PORT=5000
export FRONTEND_PORT=3000
export ENV=prod

echo "### Running the backend server"
nohup go run $MAIN_FILE > backend.log 2>&1 &

sleep 60

echo "### Checking if the backend server is up"
BACKEND_HEALTH_URL="http://localhost:$BACKEND_PORT/v1/health"
for i in {1..10}; do
    RESPONSE=$(curl -s --head --request GET $BACKEND_HEALTH_URL)
    echo "Response from backend (attempt $i):"
    echo "$RESPONSE"
    if echo "$RESPONSE" | grep "200 OK" > /dev/null; then
        echo "Backend server is up and running at $BACKEND_HEALTH_URL"
        break
    else
        echo "Waiting for backend server to be up... ($i/10)"
        sleep 2
    fi
done

if ! curl -s --head --request GET $BACKEND_HEALTH_URL | grep "200 OK" > /dev/null; then
    echo "Backend server failed to start. Please check the logs."
    exit 1
fi

echo "Installing Node.js and npm..."
curl -sL https://rpm.nodesource.com/setup_16.x | sudo bash -
    sudo yum install -y nodejs

echo "Verifying Node.js and npm installation..."
node -v
npm -v

cd ../../frontend

echo "REACT_APP_BACKEND_URL=http://$APP_IP:$BACKEND_PORT" > .env

echo "Installing dependencies..."
npm install

echo "Starting the React development server..."
setsid env BROWSER=none HOST=0.0.0.0 PORT=$FRONTEND_PORT npx react-scripts start > frontend.log 2>&1 < /dev/null &

sleep 20

echo "### Setup complete"
echo "Backend is running on http://$APP_IP:$BACKEND_PORT"
echo "Frontend is running on http://$APP_IP:$FRONTEND_PORT"