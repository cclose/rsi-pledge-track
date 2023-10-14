#!/bin/bash

# Wait for the database service to be ready
until psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME} -c '\q'; do
  echo "Waiting for the database to become available..."
  sleep 2
done

# Run your application after the database is ready
exec ./rsiAPI
