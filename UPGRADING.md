# Upgrading mantleNode

This guide provides instructions for upgrading to specific versions of mantleNode.

## v1.0.0 Upgrade

* The Yellow Spaceship Upgrade adds several new features including classifications, assets, identies, maintainers, metas, orders, and splits to the AssetMantle v1.0.0 codebase.

* The upgrade is scheduled for block `7135001`, August 18, 2023 at UTC 12:00.

* Validators should have a minimum of 64GB of RAM, or 64GB of swap, to prevent out of memory errors.

* In the case of any issues at upgrade time, coordination should take place via Discord.

### Go version

We recommend using go version 1.19

### Steps for upgrade

Verify you are currently running the correct version (v0.3.1) of the mantleNode

```bash
$ mantleNode version --long
name: AssetMantle
server_name: mantleNode
version: HEAD-5b2b0dcb37b107b0e0c1eaf9e907aa9f1a1992d9
commit: 5b2b0dcb37b107b0e0c1eaf9e907aa9f1a1992d9
build_tags: netgo,ledger
go: go version go1.18 linux/amd64
```

* Wait till the chain comes to a halt at the upgrade height.
* Validators/ nodes are requested to wait till chain halts. There will be a message that will pop up saying needs upgrade v1.0.0 after height xx.(Let the other nodes sync for a couple minutes) IMPORTANT: PLEASE WAIT FOR THE BINARY TO HALT ON ITS OWN. UPGRADE HEIGHT IS 7135001. Proceed when the chain halts at this height.
* Once the chain was halted, please stop the chain if its stopped.
* Download new binaries from https://github.com/assetmantle/node
* Replace the old binary with new binary. If you're running with systemd then stop the process and then replace the binary and start the systemd process again
* Verify if you are running new version of `mantleNode`

```bash

```

* Start your node with correct version
* Wait until 2/3+ of voting power has upgraded for the network to start producing blocks
* You can use the following commands to check peering status and state: you need jq installed

```bash
curl -sL http://127.0.0.1:26657/net_info | grep n_peers

curl -sL http://127.0.0.1:26657/consensus_state | jq -r .result.round_state.height_vote_set[].prevotes_bit_array
```
