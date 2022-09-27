<div align="center">
  <h1>Node</h1>

![GitHub Workflow Status](https://github.com/assetmantle/node/actions/workflows/test.yml/badge.svg)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=AssetMantle_modules&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=AssetMantle_modules)
[![Go Report Card](https://goreportcard.com/badge/github.com/assetmantle/node)](https://goreportcard.com/report/github.com/assetmantle/node)
[![License: Apache-2.0](https://img.shields.io/github/license/assetmantle/node.svg)](https://github.com/assetmantle/node/blob/main/LICENSE)
[![Lines Of Code](https://img.shields.io/tokei/lines/github/assetmantle/node)](https://github.com/assetmantle/node)
[![Version](https://img.shields.io/github/tag/assetmantle/node.svg?cacheSeconds=3600)](https://github.com/assetmantle/node/latest)

</div>

Application implementing the minimum clique of AssetMantle modules enabling interNFT definition, issuance, ownership
transfer and decentralized exchange.

## Hardware Requirements

* **Minimal**
  * 1 GB RAM
  * 50 GB HDD
  * 1.4 GHz CPU
* **Recommended**
  * 2 GB RAM
  * 100 GB HDD
  * 2.0 GHz x2 CPU

> NOTE: SSDs have limited TBW before non-catastrophic data errors. Running a full node requires a TB+ writes per day, causing rapid deterioration of SSDs over HDDs of comparable quality.

## Operating System

* Linux/Windows/MacOS(x86)
* **Recommended**
* Linux(x86_64)

## Installation Steps

>Prerequisite: go1.14+ required. [ref](https://golang.org/doc/install)
>Prerequisite: git. [ref](https://github.com/git/git)
>Optional requirement: GNU make. [ref](https://www.gnu.org/software/make/manual/html_node/index.html)

* Clone git repository

```shell
git clone https://github.com/AssetMantle/node.git
```

* Checkout release tag

```shell
git fetch --tags
git checkout [vX.X.X]
```

* Install

```shell
cd assetMantle
make all
```

### Generate keys

`assetClient keys add [key_name]`

or

`assetMantle keys add [key_name] --recover` to regenerate keys with your [BIP39](https://github.com/bitcoin/bips/tree/master/bip-0039) mnemonic

### Connect to a chain and start node

* [Install](#installation-steps) assetMantle application
* Initialize node

```shell
assetNode init [NODE_NAME]
```

* Replace `${HOME}/.assetNode/config/genesis.json` with the genesis file of the chain.
* Add `persistent_peers` or `seeds` in `${HOME}/.assetNode/config/config.toml`
* Start node

```shell
assetNode start
```

### Initialize a new chain and start node

* Initialize: `assetNode init [node_name] --chain-id [chain_name]`
* Add key for genesis account `assetClient keys add [genesis_key_name]`
* Add genesis account `assetNode add-genesis-account [genesis_key_name] 10000000000000000000stake`
* Create a validator at genesis `assetNode gentx --name [genesis_key_name] --amount 10000000stake`
* Collect genesis transactions `assetNode collect-gentxs`
* Start node `assetNode start`
* To start api server `assetClient rest-server`

### Reset chain

```shell
rm -rf ~/.assetNode
```

### Shutdown node

```shell
killall assetNode
```

### Check version

```shell
assetNode version
```

## Containeirzed environment

> Make sure you have latest docker version, Docker for mac can be [slow](https://twitter.com/pratikbin/status/1570722135571861504). Recommend using linux

Start node and client

```shell
# port 26657 and 1317 are exposed
make docker-compose
```

### clean

```shell
make docker-clean
```

## Contributing

If you want to contribute to AssetMantle Modules, please read the instructions in [CONTRIBUTING.md](CONTRIBUTING.md).

<div align="center">

[![Discord](https://dcbadge.vercel.app/api/server/8tSZ2NPSnS)](https://discord.gg/8tSZ2NPSnS)
[![Twitter](https://img.shields.io/twitter/follow/AssetMantle?color=blue&label=Twitter&style=for-the-badge&cacheSeconds=3600&logo=twitter)](https://twitter.com/AssetMantle)
[![Reddit](https://img.shields.io/reddit/subreddit-subscribers/AssetMantle?style=for-the-badge&cacheSeconds=3600&logo=reddit&label=Reddit%20r/assetmantle&logoColor=white)](https://www.reddit.com/r/AssetMantle/)
[![Twitter](https://img.shields.io/youtube/channel/subscribers/UCQkov-0kol99KGMxyXc-a6Q?label=YouTube&cacheSeconds=3600&logoColor=red&style=for-the-badge&logo=YouTube)](https://twitter.com/AssetMantle)

</div>

<div align="center">
    <div style="display:flex; justify-content:space-around;">
        <h3 style="margin:-5px 10px 5px;">Contributors</h3>
        <hr align="left" width="20%">
    </div>
    <img src="https://contrib.rocks/image?repo=assetmantle/node&columns=80" style="width:700;"/>
</div>
