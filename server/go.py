#!/usr/bin/python
import os

cmd = "docker ps -a | grep server | awk '{print $1}' | xargs docker rm -f"
os.system(cmd)
cmd = "docker images | grep server | awk '{print $3}' | xargs docker rmi -f"
os.system(cmd)
cmd = "rm poker"
os.system(cmd)
cmd = "go build -o poker && docker-compose up -d"
os.system(cmd)
