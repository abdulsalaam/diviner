#!/bin/bash
 docker-compose -f docker-compose-cli.yaml down
 docker rm $(docker ps -aq)
 docker rmi $(docker images --filter=reference='dev*' -q)
