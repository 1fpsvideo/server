#!/bin/bash
# 1fps Deployment Script

# This script deploys the 1fps server to the remote machine.
# It builds the server locally, stops the service on the remote machine,
# updates Docker components if necessary, copies the binary, and restarts the service.

# To set up the Ubuntu service (do this manually, once):
# 1. Create a systemd service file: sudo nano /etc/systemd/system/1fps.service
# 2. Add the following content:
#    [Unit]
#    Description=1fps Server
#    After=network.target
#
#    [Service]
#    ExecStart=/home/ubuntu/work/1fps/server
#    WorkingDirectory=/home/ubuntu/work/1fps
#    User=ubuntu
#    Restart=always
#
#    [Install]
#    WantedBy=multi-user.target
#
# 3. Save the file and exit
# 4. Reload systemd: sudo systemctl daemon-reload
# 5. Enable the service: sudo systemctl enable 1fps.service

# Remote server details
REMOTE_HOST="51.81.245.182"
REMOTE_DIR="/home/ubuntu/work/1fps"

# Build the server locally for Linux
echo "Building server locally for Linux..."
GOOS=linux GOARCH=amd64 go build -o server.linux server.go
if [ $? -ne 0 ]; then
    echo "Build failed. Aborting deployment."
    exit 1
fi

# Stop the service on the remote machine
echo "Stopping 1fps service on remote machine..."
ssh $REMOTE_HOST "sudo systemctl stop 1fps.service && sudo systemctl status 1fps.service"

# Check if docker-compose.yml needs to be updated
echo "Checking if docker-compose.yml needs to be updated..."
if ! ssh $REMOTE_HOST "test -e $REMOTE_DIR/docker-compose.yml && diff <(cat $REMOTE_DIR/docker-compose.yml) <(cat -)" < docker-compose.yml; then
    echo "docker-compose.yml needs updating. Copying file and restarting containers..."
    scp docker-compose.yml $REMOTE_HOST:$REMOTE_DIR/
    ssh $REMOTE_HOST "cd $REMOTE_DIR && docker-compose down && docker-compose up -d"
else
    echo "docker-compose.yml is up to date."
fi

# Copy .env file to remote host
echo "Copying .env file to remote machine..."
scp .env $REMOTE_HOST:$REMOTE_DIR/

# Parse .env file and test Redis connection on remote server
REDIS_HOST=$(grep REDIS_HOST .env | cut -d '=' -f2)
REDIS_PORT=$(grep REDIS_PORT .env | cut -d '=' -f2)
echo "Testing Redis connection on remote server (${REDIS_HOST}:${REDIS_PORT})..."
ssh $REMOTE_HOST "
    while ! nc -z ${REDIS_HOST} ${REDIS_PORT}; do
        echo 'Redis is unavailable - sleeping for 5 seconds'
        sleep 5
    done
    echo 'Redis is available'
"

# Copy the server binary
echo "Copying server binary to remote machine..."
scp server.linux $REMOTE_HOST:$REMOTE_DIR/server

# Remove the local Linux binary
echo "Removing local Linux binary..."
rm server.linux

# Start the service on the remote machine
echo "Starting 1fps service on remote machine..."
ssh $REMOTE_HOST "sudo systemctl start 1fps.service && sudo systemctl status 1fps.service"

echo "Deployment completed successfully."
