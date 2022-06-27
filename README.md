# Magellon

A data processing pipeline for the [AXIA network](https://axc.network).

## Features

- Maintains a persistent log of all consensus events and decisions made on the Axia network.
- Indexes SwapChain, CoreChain, and AXChain transactions.
- An API allowing easy exploration of the index.

## Prerequisite

https://docs.docker.com/engine/install/ubuntu/

https://docs.docker.com/compose/install/

## Quick Start with Standalone Mode on Test (testnet) network

The easiest way to get started is to try out the standalone mode.

```shell script
git clone https://github.com/axiacoin/axia-network-v2-magellan.git $GOPATH/github.com/axiacoin/axia-network-v2-magellan
cd $GOPATH/github.com/axiacoin/axia-network-v2-magellan
make dev_env_start
make standalone_run
```

## [Production Deployment](docs/deployment.md)

