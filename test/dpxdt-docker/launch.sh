#! /bin/bash

echo "Running DB Migrations"
echo 'server.db.create_all()
' | ./run_shell.sh

echo "Launching Dpxdt server"
./run_combined.sh
