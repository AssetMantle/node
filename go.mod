module github.com/AssetMantle/node

go 1.14

require (
	github.com/AssetMantle/modules v0.4.0
	github.com/CosmWasm/wasmd v0.29.2
	github.com/cosmos/cosmos-sdk v0.45.9
	github.com/cosmos/ibc-go/v3 v3.3.0
	github.com/spf13/cobra v1.5.0
	github.com/tendermint/tendermint v0.34.21
)

replace (
	github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.0
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
