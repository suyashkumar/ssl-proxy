
docker build --target release -t tailscale-ssl-proxy_build-release .
docker compose --env-file .auth/github.env -f docker-compose.build.yml up
docker compose -f docker-compose.build.yml down
