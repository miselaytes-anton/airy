set -e

ssh "$SSH_USER@$SSH_HOST" <<EOF
  cd "$DEPLOY_PATH"
  git pull origin main
  docker compose down
  docker compose up -d --remove-orphans postgres mosquitto 
  docker compose up -d --remove-orphans --build server
EOF