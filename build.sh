
docker build . -t ssl-proxy_build-release
docker compose -f docker-compose.build.yml up
docker compose -f docker-compose.build.yml down
