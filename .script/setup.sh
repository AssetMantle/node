# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

.script/reset.sh

source ./.script/environment.sh

mantleNode init "$AM_NODE_NAME_1" --chain-id "$AM_CHAIN_ID"
mantleNode config keyring-backend "$AM_KEYRING_BACKEND"
mantleNode config chain-id "$AM_CHAIN_ID"
mantleNode config output "$AM_OUTPUT"
mantleNode keys add "$AM_NAME_1" --keyring-backend "$AM_KEYRING_BACKEND" <<<"n" >/dev/null 2>&1 0>&1
mantleNode keys add "$AM_NAME_1" --recover --keyring-backend "$AM_KEYRING_BACKEND" <<<"y
$AM_SEED_1"
mantleNode genesis add-genesis-account "$AM_NAME_1" "$AM_GENESIS_BALANCE_1$AM_STAKE_DENOMINATION" --keyring-backend "$AM_KEYRING_BACKEND"
mantleNode genesis gentx "$AM_NAME_1" "$AM_GENESIS_STAKE_1$AM_STAKE_DENOMINATION" --chain-id "$AM_CHAIN_ID" --keyring-backend "$AM_KEYRING_BACKEND"
mantleNode genesis collect-gentxs
# Update the voting period parameter in genesis.json to 10s
jq '.app_state.gov.params.voting_period = "10s"' $HOME/.mantleNode/config/genesis.json > $HOME/.mantleNode/config/genesis_tmp.json && mv $HOME/.mantleNode/config/genesis_tmp.json $HOME/.mantleNode/config/genesis.json