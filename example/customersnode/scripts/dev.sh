#!/usr/bin/env bash
set -e

cleanup() {
    docker-compose -f  docker-compose.yml down
    trap '' EXIT INT TERM
    exit $err
}

trap cleanup SIGINT EXIT

# Make sure docker-compose is installed
if ! hash docker-compose 2>/dev/null; then
  echo -e '\033[0;31mPlease install docker-compose\033[0m'
  exit 1
fi

COMPOSE_HTTP_TIMEOUT=120 docker-compose -f docker-compose.yml up -d --force-recreate

npm run dev-server
