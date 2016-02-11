### poker 

#### building from scratch
* cd client
* npm install
* cd semantic
* gulp build
* yes to RTL, otherwise take all defaults 
* cd ../../server
* go build && docker-compose up -d

#### development commands
restart the server, from server/
* go build && docker restart server_echo_1 && docker logs --tail=1 -f server_echo_1

recompile javascript, from client/
* webpack 

#### testing
* make sure phantomjs is installed globally: `sudo npm install -g phantomjs`
* `cd test`
* `npm install`
* `npm test`


add local ip hostname to selenium nodes
  from server_selenium_chrome_1: echo "192.168.8.88    dev" >> /etc/hosts

