#!/bin/bash

ASSETCLIENT="~/node/build/assetClient"


# IDENTITIES NUB TX:
ID_NUB_TX=`$ASSETCLIENT tx identities nub --from test --nubID "deepanshutr" --chain-id "test-chain-1" --keyring-backend test --output json -y | jq -r '.txhash'`
sleep 6.5; $ASSETCLIENT query tx $ID_NUB_TX --chain-id "test-chain-1" --trust-node --output json | jq

# IDENTITIES DEFINE TX:
ID_DEFINE_TX=`$ASSETCLIENT tx identities define --from test --fromID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --immutableMetaProperties "A:S|a" --immutableProperties "B:S|" --mutableMetaProperties "C:S|" --mutableProperties "D:S|" --chain-id "test-chain-1" --keyring-backend test --output json -y | jq -r '.txhash'`
sleep 6.5; $ASSETCLIENT query tx $ID_DEFINE_TX --chain-id "test-chain-1" --trust-node --output json | jq

# IDENTITIES ISSUE TX:
ID_ISSUE_TX=`$ASSETCLIENT tx identities issue --from test --fromID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --to "mantle1pkkayn066msg6kn33wnl5srhdt3tnu2vuet86j" --classificationID "qt24iLAgD4h9V0_ST1-OKn8kpILZQ5BONEc9ILhnq5k=" --immutableMetaProperties "A:S|a" --immutableProperties "B:S|b" --mutableMetaProperties "C:S|c" --mutableProperties "D:S|d" --chain-id "test-chain-1" --keyring-backend test --output json -y | jq -r '.txhash'`
sleep 6.5; $ASSETCLIENT query tx $ID_ISSUE_TX --chain-id "test-chain-1" --trust-node --output json | jq

# ASSETS DEFINE TX:
ASSETS_DEFINE_TX=`$ASSETCLIENT tx assets define --from test --fromID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --immutableMetaProperties "A:S|a" --immutableProperties "B:S|" --mutableMetaProperties "C:S|,supply:D|" --mutableProperties "D:S|" --chain-id "test-chain-1" --keyring-backend test --output json -y | jq -r '.txhash'`
sleep 6.5; $ASSETCLIENT query tx $ASSETS_DEFINE_TX --chain-id "test-chain-1" --trust-node --output json | jq

# ASSETS MINT TX:
ASSETS_MINT_TX=`$ASSETCLIENT tx assets mint --from test --fromID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --toID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --classificationID "xuyPBH-VaAtiY-crFwTDPjoQU56Cx28AcVHbte4K4NY=" --immutableMetaProperties "A:S|a" --immutableProperties "B:S|b" --mutableMetaProperties "C:S|c,supply:D|100" --mutableProperties "D:S|d" --chain-id "test-chain-1" --keyring-backend test --output json -y | jq -r '.txhash'`
sleep 6.5; $ASSETCLIENT query tx $ASSETS_MINT_TX --chain-id "test-chain-1" --trust-node --output json | jq

# SPLITS WRAP TX:
SPLITS_WRAP_TX=`$ASSETCLIENT tx splits wrap --from test --fromID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --coins 900000000000stake --chain-id "test-chain-1" --keyring-backend test --output json -y | jq -r '.txhash'`
sleep 6.5; $ASSETCLIENT query tx $SPLITS_WRAP_TX --chain-id "test-chain-1" --trust-node --output json | jq

# ORDERS DEFINE TX:
ORDERS_DEFINE_TX=`$ASSETCLIENT tx orders define --from test --fromID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --immutableMetaProperties "E:S|e" --immutableProperties "F:S|f" --mutableMetaProperties "G:S|g" --mutableProperties "H:S|h" --chain-id "test-chain-1" --keyring-backend test --output json -y | jq -r '.txhash'`
sleep 6.5; $ASSETCLIENT query tx $ASSETS_DEFINE_TX --chain-id "test-chain-1" --trust-node --output json | jq

# ORDERS MAKE TX:
ORDERS_MAKE_TX=`$ASSETCLIENT tx orders make --from test --fromID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --classificationID "XJOYgfgnyIIGa-qnZMun-FDECTY_fkbZ89wVkuLzSKU=" --takerID "Lkj8sS7M1GMTZT65lJ6K_1vfK5uT8C39icnvHMfQ5RA=" --makerOwnableID "mbil64D0TJZScClVWQLe8JlrmtqkTO6pMul9zfLFlBE=" --takerOwnableID stake --makerOwnableSplit 1 --takerOwnableSplit 1 --expiresIn 100000 --immutableMetaProperties "E:S|e" --immutableProperties "F:S|f" --mutableMetaProperties "G:S|g" --mutableProperties "H:S|h" --chain-id "test-chain-1" --keyring-backend test --output json -y | jq -r '.txhash'`
sleep 6.5; $ASSETCLIENT query tx $ORDERS_MAKE_TX --chain-id "test-chain-1" --trust-node --output json | jq

