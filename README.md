### echomongo 
This repo implements a simple authentication system and web interface using react/echo/mongo.  Users can create accounts, login, and logout.

#### building docker container
1) open go.py and replace file system paths accordingly

2) run: `./go.py`

#### restarting container for development
1) run: `go build && docker restart echomongo_echo_1 && docker logs --tail=1 -f echomongo_echo_1`
