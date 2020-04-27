#!/bin/sh
docker-compose build --no-cache
docker-compose up -d josh twippy caddy
docker-compose logs -f --tail 100
