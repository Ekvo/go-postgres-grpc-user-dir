#!bin/bash

until pg_isready -h db -p 5432 -U user-manager --dbname=user-store; do
  echo "PostgreSQL no connect..."
  sleep 1
done

echo "PostgreSQL connect."

/usr/src/app/user &
user_pid=$!

trap 'echo "SIG recived, kill user_pid"; kill $user_pid' SIGINT SIGTERM SIGKILL

wait $user_pid