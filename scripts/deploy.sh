set -e

SSH_USER=airy
SSH_HOST=amiselaytes.com
DEPLOY_PATH=/var/www/airy

ssh "$SSH_USER@$SSH_HOST" <<EOF
  cd "$DEPLOY_PATH"
  git pull origin main
  make docker-prod
EOF