package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"

	"github.com/AssetMantle/node/utilities/rest/keys/add"
	"github.com/AssetMantle/node/utilities/rest/sign"

	"github.com/AssetMantle/node/utilities/rest/faucet"
)

func RegisterRESTRoutes(context client.Context, router *mux.Router) {
	add.RegisterRESTRoutes(context, router)
	faucet.RegisterRESTRoutes(context, router)
	sign.RegisterRESTRoutes(context, router)
}
