# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

set -x
source environment.sh

NONCE="$RANDOM"
SLEEP="7"
NUBID="nubID$NONCE"
OUT="--output json"
AM_ADDR01="$(mantleNode keys show $AM_ACC01 --address $KEYRING)"
AM_ADDR02="$(mantleNode keys show $AM_ACC02 --address $KEYRING)"

sleep $SLEEP

SEND_TX=$(mantleNode tx send $AM_ACC01 $AM_ADDR02 1000000$AM_DENOM $KEYRING --chain-id $AM_CHAIN $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $SEND_TX --chain-id $AM_CHAIN --trust-node
IDENTITIES_NUB_TX=$(mantleNode tx identities nub --from $AM_ACC01 --nubID $NUBID $KEYRING --chain-id $AM_CHAIN $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $IDENTITIES_NUB_TX --chain-id $AM_CHAIN --trust-node
NUB_CLASS=$(mantleNode query identities identities | grep -i -A 1 'classificationid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g')
NUB_HASH=$(mantleNode query identities identities | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g')
NUB_ACC="$NUB_CLASS|$NUB_HASH"

echo $NUB_ACC

IDENTITIES_DEFINE_TX=$(mantleNode tx identities define --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --from $AM_ACC01 --fromID $NUB_ACC $KEYRING --chain-id $AM_CHAIN $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $IDENTITIES_DEFINE_TX --chain-id $AM_CHAIN --trust-node
DEFINE_HASH=$(mantleNode query classifications classifications | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | tail -1)
DEFINE_ACC="$AM_CHAIN.$DEFINE_HASH"

echo $DEFINE_ACC

IDENTITIES_ISSUE_TX=$(mantleNode tx identities issue --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --from $AM_ACC01 --fromID $NUB_ACC --to $AM_ADDR01 --classificationID "$DEFINE_ACC" --chain-id $AM_CHAIN $KEYRING $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $IDENTITIES_ISSUE_TX --chain-id $AM_CHAIN --trust-node

ASSETS_DEFINE_TX=$(mantleNode tx assets define --from $AM_ACC01 --fromID $NUB_ACC --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --chain-id $AM_CHAIN $KEYRING $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $ASSETS_DEFINE_TX --chain-id $AM_CHAIN --trust-node
CLASS_ASSET=$(mantleNode query classifications classifications | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | tail -1)

ASSETS_MINT=$(mantleNode tx assets mint --from $AM_ACC01 --fromID $NUB_ACC --toID $NUB_ACC --classificationID "$AM_CHAIN.$CLASS_ASSET" --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --chain-id $AM_CHAIN $KEYRING $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $ASSETS_MINT --chain-id $AM_CHAIN --trust-node
ASSET_CLASS=$(mantleNode query assets assets | grep -i -A 1 'classificationid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | tail -1)
ASSET_HASH=$(mantleNode query assets assets | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | tail -1)
ASSET_ACC="$ASSET_CLASS|$ASSET_HASH"

ASSETS_MUTATE=$(mantleNode tx assets mutate --from $AM_ACC01 --fromID $NUB_ACC --assetID $ASSET_ACC --mutableProperties "testMutable3:S|num31" --mutableMetaProperties "testMutableMeta4:S|num42" $KEYRING --chain-id $AM_CHAIN $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $ASSETS_MUTATE --chain-id $AM_CHAIN --trust-node

SPLITS_WRAP=$(mantleNode tx splits wrap --from $AM_ACC01 --fromID $NUB_ACC --coins 20stake $KEYRING --chain-id $AM_CHAIN $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $SPLITS_WRAP --chain-id $AM_CHAIN --trust-node

SPLITS_UNWRAP=$(mantleNode tx splits unwrap --from $AM_ACC01 --fromID $NUB_ACC --ownableID $AM_DENOM --split 1 $KEYRING)
sleep $SLEEP
mantleNode query tx $SPLITS_UNWRAP --chain-id $AM_CHAIN --trust-node

# mantleNode tx splits unwrap -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --ownableID stake --split 1 $KEYRING $MODE


ORDERS_DEFINE=$(mantleNode tx orders define --from $AM_ACC01 --fromID $NUB_ACC --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" $KEYRING --chain-id $AM_CHAIN $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $ORDERS_DEFINE --chain-id $AM_CHAIN --trust-node
ORDERS_HASH=$(mantleNode query classifications classifications | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | head -3 | tail -1)
ORDERS_CLASS="$AM_CHAIN.$ORDERS_HASH"
ORDERS_MAKE=$(mantleNode tx orders make --from $AM_ACC01 --fromID $NUB_ACC --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --classificationID $ORDERS_CLASS --makerOwnableID $ASSET_ACC --takerOwnableID stake --makerOwnableSplit "0.000000000000000001" --takerOwnableSplit "0.000000000000000001" $KEYRING --chain-id $AM_CHAIN $OUT -y | jq -r '.txhash')
sleep $SLEEP
mantleNode query tx $ORDERS_MAKE --chain-id $AM_CHAIN --trust-node

mantleNode q orders orders
# mantleNode tx orders take

set +x
