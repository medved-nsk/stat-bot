#!/bin/zsh

docker build -t statbot .
docker image tag statbot warsong.me:5000/bots/statbot
docker push warsong.me:5000/bots/statbot