#!/bin/bash

set -x

API="http://127.0.0.1:1317"
DELAY="7"
TEMP="$PWD/temp"
NONCE="$RANDOM"
FIRST="first_$NONCE"
CHAIN="load-test-1"
COINS="1000stake"
NUB="nubID$NONCE"
MAIN=$(assetClient keys show test1 --address --keyring-backend test)

source ./json.bash

# setup() {
#   rm -rf /root/.AssetMantle/Node/
#   assetNode init test --chain-id load-test-1
#   assetNode add-genesis-account test1 10000000000000000000000stake --keyring-backend test
#   assetNode gentx --name test1 --amount 1000000000stake --keyring-backend test
#   assetNode collect-gentxs
#   systemctl restart mantle-node.service
#   systemctl restart mantle-client.service
# }

mkdir -p $TEMP;
sleep $DELAY;

# keys add

curl -sL -X POST $API/keys/add \
  -H "Content-Type: application/json" \
  -d "$(keys_add_data)" | jq '.' | awk 'NF' > $TEMP/keys_add.log.json; sleep $DELAY;

ACCOUNT_FIRST=$(cat $TEMP/keys_add.log.json | jq -r '.result.keyOutput.address');

curl -sL -X POST $API/sign \
  -H "Content-Type: application/json" \
  -d "$(sign_tx_data)" | jq '.' | awk 'NF' > $TEMP/sign_tx.log.json; sleep $DELAY;

ACCOUNT_FIRST_PUB_TYPE=$(cat $TEMP/sign_tx.log.json | jq -r '.result.tx.signatures[].pub_key.type');
ACCOUNT_FIRST_PUB_VALUE=$(cat $TEMP/sign_tx.log.json | jq -r '.result.tx.signatures[].pub_key.value');
ACCOUNT_FIRST_SIG=$(cat $TEMP/sign_tx.log.json | jq -r '.result.tx.signatures[].signature');

curl -sL -X POST $API/txs \
  -H "Content-Type: application/json" \
  -d "$(broadcast_tx_data)" | jq '.' | awk 'NF' > $TEMP/broadcast_tx.log.json; sleep $DELAY;

curl -sL -X POST $API/identities/nub \
  -H "Content-Type: application/json" \
  -d "$(identities_nub_tx_data)" | jq '.' | awk 'NF' > $TEMP/identities_nub_tx_data.log.json; sleep $DELAY;

IDENTITIES_NUB_TX_HASH=$(cat $TEMP/identities_nub_tx_data.log.json | jq -r '.txhash');

curl -sL -X GET $API/txs/$IDENTITIES_NUB_TX_HASH \
  -H "Content-Type: application/json" | jq '.' | awk 'NF'

curl -sL -X GET $API/identities/identities/test \
  -H "Content-Type: application/json" | jq '.' | awk 'NF' > $TEMP/query_identities.log.json;

ID_CLASSIFICATION=$(cat $TEMP/query_identities.log.json | jq -r '.result.list[0].value.Document.id.value.classificationID.value.idString');
ID_HASH=$(cat $TEMP/query_identities.log.json | jq -r '.result.list[0].value.Document.id.value.hashID.value.idString');

curl -sL -X POST $API/identities/define \
  -H "Content-Type: application/json" \
  -d "$(identities_define_tx_data)" | jq '.' | awk 'NF' > $TEMP/identities_define_tx_data.log.json; sleep $DELAY;

IDENTITIES_DEFINE_TX_HASH=$(cat $TEMP/identities_define_tx_data.log.json | jq -r '.txhash');

curl -sL -X GET $API/txs/$IDENTITIES_DEFINE_TX_HASH \
  -H "Content-Type: application/json" | jq '.' | awk 'NF'

curl -sL -X GET $API/classifications/classifications/test \
  -H "Content-Type: application/json" | jq '.' | awk 'NF' > $TEMP/query_classifications.log.json;

ID_DEFINE_HASH=$(cat $TEMP/query_classifications.log.json | jq -r '.result.list[1].value.Document.id.value.HashID.value.idString' | tail -1);

curl -sL -X POST $API/identities/issue \
  -H "Content-Type: application/json" \
  -d "$(identities_issue_tx_data)" | jq '.' | awk 'NF' > $TEMP/identities_issue_tx_data.log.json; sleep $DELAY;

IDENTITIES_ISSUE_TX_HASH=$(cat $TEMP/identities_issue_tx_data.log.json | jq -r '.txhash');

curl -sL -X GET $API/txs/$IDENTITIES_ISSUE_TX_HASH \
  -H "Content-Type: application/json" | jq '.' | awk 'NF'

curl -sL -X POST $API/assets/define \
  -H "Content-Type: application/json" \
  -d "$(assets_define_tx_data)" | jq '.' | awk 'NF' > $TEMP/assets_define_tx_data.log.json; sleep $DELAY;

ASSETS_DEFINE_TX_HASH=$(cat $TEMP/assets_define_tx_data.log.json | jq -r '.txhash');

curl -sL -X GET $API/txs/$ASSETS_DEFINE_TX_HASH \
  -H "Content-Type: application/json" | jq '.' | awk 'NF'



rm -rf $TEMP;

set +x 
