# building docker container
1) open go.py and replace filesys paths accordingly
2) then run: ./go.py

# restarting container for development
1) go build && docker restart echomongo_echo_1 && docker logs --tail=1 -f echomongo_echo_1
