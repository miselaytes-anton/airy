set -e
 
ssh "$SSH_USER@$SSH_HOST" <<EOF
    cd "$DEPLOY_PATH"
    git reset --hard main
    git pull origin main
    make docker-down
    make docker-prod
    docker image prune --force
EOF