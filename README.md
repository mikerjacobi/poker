### echomongo 
This repo implements a simple authentication system and web interface using react/echo/mongo.  Users can create accounts, login, and logout.  A session id is stored as a cookie and in mongo next to the account who owns the session.  For page loads, the session gets passed as a standard cookie.  For ajax requests, the session id gets pulled into a header variable for CORS reasons.  User passwords are stored bcrypted.

#### building docker container
1) run: `./go.py`

#### restarting container for development
1) run: `go build && docker restart echomongo_echo_1 && docker logs --tail=1 -f echomongo_echo_1`
