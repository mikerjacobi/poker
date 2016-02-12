#script to quickly reboot the selenium stuff

rm -rf poker-tests
docker rm -f $(docker ps -aq)
cd ../server
docker-compose up -d && docker exec server_mongo_1 mongo /fixtures/fixtures.js
cd ../test
