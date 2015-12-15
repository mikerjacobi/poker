### poker 

#### building poker app from scratch
* cd client
* npm install
* cd semantic
* gulp build
* yes to RTL, otherwise take all the defaults 
* cd ../..
* go build && docker-compose up -d
* insert echo.math({count:0})

#### restarting container for development
1) run: `go build && docker restart poker_echo_1 && docker logs --tail=1 -f poker_echo_1`
