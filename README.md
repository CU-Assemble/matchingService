This is matching service which is responsible for contain parcipants for each activity in this app

This microservice use rabbitMQ (message broker) to communicate between activity microservice and this one

before using this service please run the following command to create docker container

docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
