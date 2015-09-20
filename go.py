#!/usr/bin/python
import os

cmd = "docker ps -a | grep echo | awk '{print $1}' | xargs docker rm -f"
os.system(cmd)
cmd = "docker images | grep echo | awk '{print $3}' | xargs docker rmi -f"
os.system(cmd)
cmd = "rm echomongo"
os.system(cmd)
cmd = "go build && docker-compose up -d"
os.system(cmd)
