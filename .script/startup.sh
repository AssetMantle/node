# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

set -x

source ./environment.sh

assetNode start >~/.AssetMantle/Node/log &
sleep 10
assetClient rest-server --laddr tcp://127.0.0.1:1317 --chain-id "$AM_CHAIN" --trust-node>~/.AssetMantle/Client/log &
sleep 10
echo "
Node and Client started up. For logs:
tail -f ~/.AssetMantle/Node/log
tail -f ~/.AssetMantle/Client/log
"

set +x