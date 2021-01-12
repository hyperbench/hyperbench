#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -e

echo "Shut down the Docker containers for the system tests."
docker-compose -f docker-compose.yml kill && docker-compose -f docker-compose.yml down

echo "remove the local state"
rm -f ~/.hfc-key-store/*

echo "remove fabric containers"
num=$(docker ps -a | awk '{if($2~/^((hyperledger)|(dev)).*/)print $1}' | wc -l)
if [ $num -ne 0 ];then
        docker rm $(docker ps -a | awk '{if($2~/^((hyperledger)|(dev)).*/)print $1}')
fi

echo "remove chaincode images"
chaincode_num=$(docker images dev-* -q | wc -l)
if [ $chaincode_num -ne 0 ];then
        docker rmi $(docker images dev-* -q)
fi
echo "Your system is now clean."
echo "Restarting..."

./start.sh

echo "Restarted"