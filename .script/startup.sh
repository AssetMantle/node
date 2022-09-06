# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

source ./.script/environment.sh

assetNode start >~/.AssetMantle/Node/log &
sleep 10
assetClient rest-server --chain-id "$AM_CHAIN_ID">~/.AssetMantle/Client/log &
echo "
Node and Client started up. For logs:
tail -f ~/.AssetMantle/Node/log
tail -f ~/.AssetMantle/Client/log
"
