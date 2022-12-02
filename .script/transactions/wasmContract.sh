# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

# genesis account / chain -id is test, commands to store, instantiate, execute a contract. this eg- hackatom, github.com/CosmWasm/cosmwasm/contracts/hackatom

mantleNode tx wasm store /PATH_TO_WASM_COMTRACT/_.wasm --from test --gas 900000 -y --chain-id test

CODE_ID=$(mantleNode query wasm list-code --chain-id test | jq .[-1].id)

mantleNode keys add bob

INIT=$(jq -n --arg test $(mantleNode keys show -a test) --arg bob $(mantleNode keys show -a bob) '{"verifier":$test,"beneficiary":$bob}')

mantleNode tx wasm instantiate $CODE_ID "$INIT" --from test --amount=50000stake --label "escrow 1" -y --chain-id test

CONTRACT=$(mantleNode query wasm list-contract-by-code $CODE_ID --chain-id test | jq -r .[0].address)

MINT='{"asset_mint":{"properties":"test5:7, test89:76"}}'

mantleNode tx wasm execute $CONTRACT "$MINT" --from test -y --chain-id test

# issue asset normal
mantleNode tx assets mint --from test --properties test1:test1,test2:test2 --chain-id test
mantleNode q assets assets --chain-id test
