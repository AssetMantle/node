#!/usr/bin/env bash

NONCE="$RANDOM"
SLEEP="15"
ACCOUNT1="test1"
ACCOUNT2="test2"
KEYRING="--keyring-backend test"
CHAIN_ID="--chain-id load-test-1"
NUB_ID="nubID$NONCE"

# send coin
assetClient tx send $ACCOUNT1 $(assetClient keys show test2 --address --keyring-backend test) 1000000stake $KEYRING $CHAIN_ID --yes ; sleep $SLEEP

# send tx nub
assetClient tx identities nub --from $ACCOUNT1 --nubID $NUB_ID $KEYRING $CHAIN_ID --yes ; sleep $SLEEP

NUB_CLASS_ID=$(assetClient query identities identities | grep -i -A 1 'classificationid' | grep 'idstring:' | cut -d ':' -f2)
NUB_HASH_ID=$(assetClient query identities identities | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2)
NUBACCOUNT="$NUB_CLASS_ID|$NUB_HASH_ID"

# send define nub
assetClient tx identities define --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --from $ACCOUNT1 --fromID $NUBACCOUNT $KEYRING $CHAIN_ID --yes ; sleep $SLEEP

DEFINE_HASH_ID=$(assetClient query classifications classifications | grep -i -A 1 'hashid' | grep 'idstring:' | cut -d ':' -f2)
CLASS_ID="$CHAIN_ID.$DEFINE_HASH_ID"

assetClient tx identities issue --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" --from $ACCOUNT1 --fromID $NUBACCOUNT --to $(assetClient keys show test2 --address --keyring-backend test) --classificationID "$CLASS_ID" $CHAIN_ID $KEYRING --yes ; sleep $SLEEP

assetClient tx identities provision --from $ACCOUNT1 --to $(assetClient keys show test2 --address --keyring-backend test) --identityID $NUBACCOUNT $CHAIN_ID $KEYRING --yes ; sleep $SLEEP

assetClient tx identities unprovision --from $ACCOUNT1 --to $(assetClient keys show test2 --address --keyring-backend test) --identityID $NUBACCOUNT $CHAIN_ID $KEYRING --yes ; sleep $SLEEP

assetClient tx identities provision --from $ACCOUNT1 --to $(assetClient keys show test1 --address --keyring-backend test) --identityID $NUBACCOUNT $CHAIN_ID $KEYRING --yes ; sleep $SLEEP

assetClient tx identities unprovision --from $ACCOUNT1 --to $(assetClient keys show test1 --address --keyring-backend test) --identityID $NUBACCOUNT $CHAIN_ID $KEYRING --yes ; sleep $SLEEP

assetClient tx metas reveal --from $ACCOUNT1 --data "S|test1" $CHAIN_ID $KEYRING --yes ; sleep $SLEEP

assetClient tx assets define --from $ACCOUNT1 --fromID $NUBACCOUNT --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" $CHAIN_ID $KEYRING --yes ; sleep $SLEEP

assetClient tx assets mint --from $ACCOUNT1 --fromID $NUBACCOUNT --classificationID "$CLASS_ID" --toID $NUBACCOUNT --immutableProperties "testImmutable1:S|num1" --immutableMetaProperties "testImmutableMeta2:S|num2" --mutableProperties "testMutable3:S|num3" --mutableMetaProperties "testMutableMeta4:S|num4" $CHAIN_ID $KEYRING --yes ; sleep $SLEEP

