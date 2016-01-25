### poker 

#### building from scratch
* cd client
* npm install
* cd semantic
* gulp build
* yes to RTL, otherwise take all defaults 
* cd ../../server
* go build && docker-compose up -d
* insert echo.math({count:0})

#### development commands
restart the server, from server/
* go build && docker restart server_echo_1 && docker logs --tail=1 -f server_echo_1

recompile javascript, from client/
* webpack --watch

#### testing
* make sure phantomjs is installed globally: `sudo npm install -g phantomjs`
* `cd test`
* `behave`
