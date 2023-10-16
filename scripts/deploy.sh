set -e

SSH_USER=anton
SSH_HOST=amiselaytes.com
DEPLOY_PATH=/var/www/tatadata

ssh "$SSH_USER@$SSH_HOST" <<EOF
  cd "$DEPLOY_PATH"
  git pull origin main
  docker compose down
  docker compose up -d --remove-orphans postgres mosquitto 
  docker compose up -d --remove-orphans --build server
EOF