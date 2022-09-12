# node

Application implementing the minimum clique of AssetMantle modules enabling interNFT definition, issuance, ownership
transfer and decentralized exchange.

[![LoC](https://tokei.rs/b1/github/persistenceOne/assetMantle)](https://github.com/persistenceOne/assetMantle)

## Talk to us!
*   [Twitter](https://twitter.com/PersistenceOne)
*   [Telegram](https://t.me/PersistenceOneChat)
*   [Discord](https://discord.com/channels/796174129077813248)

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
>Prerequisite: go1.15+ required. [ref](https://golang.org/doc/install)

>Prerequisite: git. [ref](https://github.com/git/git)

>Optional requirement: GNU make. [ref](https://www.gnu.org/software/make/manual/html_node/index.html)


* Clone git repository
```shell
git clone https://github.com/persistenceOne/assetMantle.git
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