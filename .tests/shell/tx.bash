#!/bin/bash

set -x

NONCE="$RANDOM"
SLEEP="7"
ACC01="test1"
ACC02="test2"
CHAIN="load-test-1"
KEYRING="--keyring-backend test"
CHAINID="--chain-id load-test-1"
NUBID="nubID$NONCE"
OUT="--output json"
ADDR01="$(assetClient keys show test1 --address --keyring-backend test)"
ADDR02="$(assetClient keys show test2 --address --keyring-backend test)"

sleep $SLEEP

SEND_TX=$(assetClient tx send $ACC01 $ADDR02 1000000stake $KEYRING $CHAINID $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $SEND_TX --chain-id load-test-1 --trust-node
IDENTITIES_NUB_TX=$(assetClient tx identities nub --from $ACC01 --nubID $NUBID $KEYRING $CHAINID $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $IDENTITIES_NUB_TX --chain-id load-test-1 --trust-node
NUB_CLASS=$(assetClient query identities identities | grep -i -A 1 'classificationid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g')
NUB_HASH=$(assetClient query identities identities | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g')
NUB_ACC="$NUB_CLASS|$NUB_HASH"

echo $NUB_ACC

IDENTITIES_DEFINE_TX=$(assetClient tx identities define --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --from $ACC01 --fromID $NUB_ACC $KEYRING $CHAINID $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $IDENTITIES_DEFINE_TX --chain-id load-test-1 --trust-node
DEFINE_HASH=$(assetClient query classifications classifications | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | tail -1)
DEFINE_ACC="$CHAIN.$DEFINE_HASH"

echo $DEFINE_ACC

IDENTITIES_ISSUE_TX=$(assetClient tx identities issue --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --from $ACC01 --fromID $NUB_ACC --to $ADDR01 --classificationID "$DEFINE_ACC" $CHAINID $KEYRING $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $IDENTITIES_ISSUE_TX --chain-id load-test-1 --trust-node

ASSETS_DEFINE_TX=$(assetClient tx assets define --from $ACC01 --fromID $NUB_ACC --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" $CHAINID $KEYRING $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $ASSETS_DEFINE_TX --chain-id load-test-1 --trust-node
CLASS_ASSET=$(assetClient query classifications classifications | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | tail -1)

ASSETS_MINT=$(assetClient tx assets mint --from $ACC01 --fromID $NUB_ACC --toID $NUB_ACC --classificationID "$CHAIN.$CLASS_ASSET" --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" $CHAINID $KEYRING $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $ASSETS_MINT --chain-id load-test-1 --trust-node
ASSET_CLASS=$(assetClient query assets assets | grep -i -A 1 'classificationid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | tail -1)
ASSET_HASH=$(assetClient query assets assets | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | tail -1)
ASSET_ACC="$ASSET_CLASS|$ASSET_HASH"

ASSETS_MUTATE=$(assetClient tx assets mutate --from $ACC01 --fromID $NUB_ACC --assetID $ASSET_ACC --mutableProperties "testMutable3:S|num31" --mutableMetaProperties "testMutableMeta4:S|num42" $KEYRING $CHAINID $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $ASSETS_MUTATE --chain-id load-test-1 --trust-node

SPLITS_WRAP=$(assetClient tx splits wrap --from $ACC01 --fromID $NUB_ACC --coins 20stake $KEYRING $CHAINID $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $SPLITS_WRAP --chain-id load-test-1 --trust-node

ORDERS_DEFINE=$(assetClient tx orders define --from $ACC01 --fromID $NUB_ACC --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" $KEYRING $CHAINID $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $ORDERS_DEFINE --chain-id load-test-1 --trust-node
ORDERS_HASH=$(assetClient query classifications classifications | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2 | sed 's/ //g' | head -3 | tail -1)
ORDERS_CLASS="$CHAIN.$ORDERS_HASH"

ORDERS_MAKE=$(assetClient tx orders make --from $ACC01 --fromID $NUB_ACC --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --classificationID $ORDERS_CLASS --makerOwnableID $ASSET_ACC --takerOwnableID stake --makerOwnableSplit "0.001" --takerOwnableSplit "0.0001" $KEYRING $CHAINID $OUT -y | jq -r '.txhash')
sleep $SLEEP
assetClient query tx $ORDERS_MAKE --chain-id load-test-1 --trust-node

# assetClient tx orders cancel 
# assetClient tx orders take

set +x
