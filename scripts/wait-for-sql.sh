#!/bin/sh
until COMPOSE_PROJECT_NAME=test docker-compose exec postgres postgresadmin ping -P 3306 -u root --password=testing-password | grep "alive" ; do
  >&2 echo "PostgresSQL is unavailable - waiting for it... 😴"
  sleep 1
done
