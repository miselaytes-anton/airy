set -e

SSH_USER=airy
SSH_HOST=amiselaytes.com
DEPLOY_PATH=/var/www/airy

IMAGE_FILE=airy-backend.tar.gz
IMAGE_TAG=airy-backend:latest

# Build docker image and push it to VM
docker build .  --tag "$IMAGE_TAG" --platform linux/x86_64
docker image save "$IMAGE_TAG" | gzip > "$IMAGE_FILE"
scp ./"$IMAGE_FILE" "$SSH_USER@$SSH_HOST":"$DEPLOY_PATH"
rm -rf ./"$IMAGE_FILE"
 
# SSH into VM, load docker image and restart containers
ssh "$SSH_USER@$SSH_HOST" <<EOF
    set -e
    cd "$DEPLOY_PATH"
    git reset --hard main
    git pull origin main
    docker load < "$IMAGE_FILE"
    make docker-down
    make docker-prod
    docker image prune --force
    rm -rf "$IMAGE_FILE"
EOF