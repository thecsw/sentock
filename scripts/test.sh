#!/bin/sh
docker-compose up -d postgres
echo -n "Sleeping until DB is operational ... "
sleep 4
echo -e '\033[32mdone\033[0m'
go test -v -count=1 ./...
docker-compose down
