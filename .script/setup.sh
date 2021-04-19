.script/reset.sh

source ./.script/environment.sh

assetNode init "$AM_NODE_NAME_1" --chain-id "$AM_CHAIN_ID"
assetClient keys add "$AM_NAME_1" <<<"n" >/dev/null 2>&1 0>&1
assetClient keys add "$AM_NAME_1" --recover <<<"y
$AM_SEED_1"
assetNode add-genesis-account "$AM_NAME_1" "$AM_GENESIS_BALANCE_1$AM_STAKE_DENOMINATION"
assetNode gentx --name "$AM_NAME_1" --amount "$AM_GENESIS_STAKE_1$AM_STAKE_DENOMINATION"
assetNode collect-gentxs
