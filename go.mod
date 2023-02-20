module github.com/AssetMantle/node

go 1.14

require (
	bou.ke/monkey v1.0.1 // indirect
	github.com/AssetMantle/modules v0.3.2-0.20230220095921-efdf481b47e1
	github.com/CosmWasm/go-cosmwasm v0.10.0 // indirect
	github.com/CosmWasm/wasmd v0.29.2
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d // indirect
	github.com/cosmos/cosmos-sdk v0.45.9
	github.com/cosmos/ibc-go/v3 v3.3.0
	github.com/spf13/cobra v1.5.0
	github.com/swaggo/http-swagger v1.3.3 // indirect
	github.com/tendermint/iavl v0.14.3 // indirect
	github.com/tendermint/tendermint v0.34.21
)

replace (
	github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.0
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
