
docker build --target test -t tailscale-ssl-proxy_build-test .
docker compose -f docker-compose.test.yml up
docker compose -f docker-compose.test.yml down
