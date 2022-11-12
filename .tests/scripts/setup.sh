#!/bin/bash

rm -rf $HOME/.AssetMantle/Client/
rm -rf $HOME/.AssetMantle/Node/

assetClient keys add test1 --keyring-backend test
assetClient keys add test2 --keyring-backend test

assetNode init test --chain-id load-test-1
assetNode add-genesis-account test1 10000000000000000000000stake --keyring-backend test
assetNode gentx --name test1 --amount 1000000000stake --keyring-backend test
assetNode collect-gentxs
