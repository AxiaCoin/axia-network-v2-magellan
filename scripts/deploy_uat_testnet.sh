#!/bin/bash
echo ====checkout to testnet-uat====
git checkout testnet-uat
echo ====docker login====
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 667299798306.dkr.ecr.us-east-1.amazonaws.com
echo ====build image====
docker build -t 667299798306.dkr.ecr.us-east-1.amazonaws.com/axia-indexer-testnet-infra-magellan3ba08964-h5eca3o7tqfq:latest .
echo ====push image====
docker push 667299798306.dkr.ecr.us-east-1.amazonaws.com/axia-indexer-testnet-infra-magellan3ba08964-h5eca3o7tqfq:latest
