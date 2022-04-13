module github.com/AssetMantle/node

go 1.17

require (
	github.com/cosmos/cosmos-sdk v0.45.1
	github.com/cosmos/ibc-go/v3 v3.0.0
	github.com/gorilla/mux v1.8.0
	github.com/gravity-devs/liquidity v1.5.0
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.3.0
	github.com/strangelove-ventures/packet-forward-middleware/v2 v2.1.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
)
replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tecbot/gorocksdb => github.com/cosmos/gorocksdb v1.2.0
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)


