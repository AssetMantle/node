# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

source ./.script/reset.sh
source ./.script/environment.sh

mantleNode keys add $AM_ACC01 $KEYRING
mantleNode keys add $AM_ACC02 $KEYRING

mantleNode init "$AM_NODE" --chain-id "$AM_CHAIN"
mantleNode add-genesis-account "$AM_ACC01" "$AM_GENESIS_BAL$AM_DENOM" $KEYRING
mantleNode gentx "$AM_ACC01" "$AM_GENESIS_STAKE$AM_DENOM" $KEYRING --chain-id "$AM_CHAIN"
mantleNode collect-gentxs
