# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

<<<<<<< HEAD
set -x

=======
>>>>>>> parent of b1a33d3 (tests: add shell script)
source ./.script/environment.sh

assetNode start >~/.AssetMantle/Node/log &
sleep 10
assetClient rest-server --chain-id "$AM_CHAIN_ID">~/.AssetMantle/Client/log &
echo "
Node and Client started up. For logs:
tail -f ~/.AssetMantle/Node/log
tail -f ~/.AssetMantle/Client/log
"
<<<<<<< HEAD

set +x
=======
>>>>>>> parent of b1a33d3 (tests: add shell script)
