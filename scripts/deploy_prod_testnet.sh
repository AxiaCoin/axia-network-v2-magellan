#!/bin/bash
echo ====checkout to testnet-prod====
git checkout testnet-prod
echo ====docker login====
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 606361176142.dkr.ecr.us-east-1.amazonaws.com
echo ====build image====
docker build -t 606361176142.dkr.ecr.us-east-1.amazonaws.com/axia-indexer-testnet-infra-magellan3ba08964-lpj0j6lyors3:latest .
echo ====push image====
docker push 606361176142.dkr.ecr.us-east-1.amazonaws.com/axia-indexer-testnet-infra-magellan3ba08964-lpj0j6lyors3:latest