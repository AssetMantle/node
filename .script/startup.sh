# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

source ./.script/environment.sh

mantleNode start >~/.mantleNode/log &
echo "
Node and Client started up. For logs:
tail -f ~/.mantleNode/log
"
