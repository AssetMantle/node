#!/bin/bash

keys_add_data() 
{
cat <<EOF
{
  "name": "$FIRST"
}
EOF
}

sign_tx_data()
{
cat <<EOF
{
  "baseReq":{
    "from":"${MAIN}",
    "chain_id": "${CHAIN}"
  },
  "type": "cosmos-sdk/StdTx",
  "value": {
    "msg": [
      {
        "type": "cosmos-sdk/MsgSend",
        "value": {
          "from_address": "${MAIN}",
          "to_address": "${ACCOUNT_FIRST}",
          "amount": [
            {
              "denom": "stake",
              "amount": "10000"
            }
          ]
        }
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
    },
    "signatures": null,
    "memo": ""
  }
}
EOF
}

broadcast_tx_data()
{
cat <<EOF
{
       "tx": {
            "msg": [
                {
                    "type": "cosmos-sdk/MsgSend",
                    "value": {
                        "from_address": "${MAIN}",
                        "to_address": "${ACCOUNT_FIRST}",
                        "amount": [
                            {
                                "denom": "stake",
                                "amount": "10000"
                            }
                        ]
                    }
                }
            ],
            "fee": {
                "amount": [],
                "gas": "200000"
            },
            "signatures": [
                {
                    "pub_key": {
                    	"type": "${ACCOUNT_FIRST_PUB_TYPE}",
                    	"value": "${ACCOUNT_FIRST_PUB_VALUE}"
                    	},
                    "signature": "${ACCOUNT_FIRST_SIG}"
                }
              ],
            "memo": ""
        },
    "mode": "block"
}
EOF
}

identities_nub_tx_data()
{
cat <<EOF
{
"type": "github.com/AssetMantle/modules/modules/identities/internal/transactions/nub/transactionRequest",
 "value": {
    "baseReq": {
      "from": "${ACCOUNT_FIRST}", 
      "chain_id": "${CHAIN}"
    },
    "nubID": "${NUB}"
  }
}
EOF
}

identities_define_tx_data()
{
cat <<EOF
{
"type": "github.com/AssetMantle/modules/modules/identities/internal/transactions/define/transactionRequest",
 "value": {
    "baseReq": {
      "from": "${ACCOUNT_FIRST}",
      "chain_id": "${CHAIN}"
    },
    "fromID": "${ID_CLASSIFICATION}|${ID_HASH}",
    "immutableMetaProperties": "test1:S|ImmutableMeta",
    "immutableProperties": "testMeta2:S|Immutable",
    "mutableMetaProperties": "nubDefine3:S|MutableMeta",
    "mutableProperties": "nubDefine4:S|Mutable"
  }
}
EOF
}

identities_issue_tx_data()
{
cat <<EOF
{
"type": "github.com/AssetMantle/modules/modules/identities/internal/transactions/issue/transactionRequest",
 "value": {
    "baseReq": {
      "from": "${ACCOUNT_FIRST}",
      "chain_id": "${CHAIN}"
    },
    "to": "${ACCOUNT_FIRST}",
    "fromID": "${ID_CLASSIFICATION}|${ID_HASH}",
    "classificationID": "${CHAIN}.${ID_DEFINE_HASH}",
    "immutableMetaProperties": "test1:S|ImmutableMeta",
    "immutableProperties": "testMeta2:S|Immutable",
    "mutableMetaProperties": "nubIssue3:S|MutableMeta",
    "mutableProperties": "nubIssue4:S|Mutable"
  }
}
EOF
}

assets_define_tx_data()
{
cat <<EOF
{
"type": "github.com/AssetMantle/modules/modules/assets/internal/transactions/define/transactionRequest",
 "value": {
    "baseReq": {
      "from": "${ACCOUNT_FIRST}",
      "chain_id": "${CHAIN}"
    },
    "fromID": "${ID_CLASSIFICATION}|${ID_HASH}",
    "immutableMetaProperties": "test1:S|ImmutableMeta",
    "immutableProperties": "testMeta2:S|Immutable",
    "mutableMetaProperties": "assetDefine3:S|MutableMeta",
    "mutableProperties": "assetDefine4:S|Mutable"
	}
}
EOF
}