# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

source ./.script/reset.sh
source ./.script/environment.sh

mantleNode keys add $AM_ACC01 --recover
mantleNode keys add $AM_ACC02

mantleNode init "$AM_NODE" --chain-id "$AM_CHAIN"
mantleNode add-genesis-account "$AM_ACC01" "$AM_GENESIS_BAL$AM_DENOM"
mantleNode gentx "$AM_ACC01" "$AM_GENESIS_STAKE$AM_DENOM" --chain-id "$AM_CHAIN"
mantleNode collect-gentxs

# dirt arrest crawl trumpet enemy attend soap absurd clip learn debate maximum hawk pattern rich isolate creek echo peasant close appear south advance mean
