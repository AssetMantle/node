# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

.script/reset.sh

source ./.script/environment.sh

mantleNode init "$AM_NODE_NAME_1" --chain-id "$AM_CHAIN_ID"
mantleNode keys add "$AM_NAME_1" <<<"n" >/dev/null 2>&1 0>&1
mantleNode keys add "$AM_NAME_1" --recover <<<"y
$AM_SEED_1"
mantleNode add-genesis-account "$AM_NAME_1" "$AM_GENESIS_BALANCE_1$AM_STAKE_DENOMINATION"
mantleNode gentx "$AM_NAME_1" "$AM_GENESIS_STAKE_1$AM_STAKE_DENOMINATION" --chain-id "$AM_CHAIN_ID"
mantleNode collect-gentxs
