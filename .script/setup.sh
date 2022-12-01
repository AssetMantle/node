# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

source ./.script/reset.sh
source ./.script/environment.sh

assetClient keys add $AM_ACC01 $KEYRING
assetClient keys add $AM_ACC02 $KEYRING

sleep $SLEEP

assetNode init "$AM_NODE" --chain-id "$AM_CHAIN"
assetNode add-genesis-account "$AM_ACC01" "$AM_GENESIS_BAL$AM_DENOM" $KEYRING
assetNode gentx --name "$AM_ACC01" --amount "$AM_GENESIS_STAKE$AM_DENOM" $KEYRING
assetNode collect-gentxs
