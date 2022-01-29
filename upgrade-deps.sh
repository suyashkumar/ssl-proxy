
docker build --target build -t tailscale-ssl-proxy_upgrade-deps .docker build .
docker compose -f docker-compose.upgrade-deps.yml up
docker compose -f docker-compose.upgrade-deps.yml down
