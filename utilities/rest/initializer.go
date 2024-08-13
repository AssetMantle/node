package rest

import (
	documentIDGetters "github.com/AssetMantle/modules/utilities/rest/id_getters/docs"
	"github.com/AssetMantle/node/utilities/rest/faucet"
	"github.com/AssetMantle/node/utilities/rest/keys/add"
	"github.com/AssetMantle/node/utilities/rest/sign"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func RegisterRESTRoutes(context client.Context, router *mux.Router) {
	add.RegisterRESTRoutes(context, router)
	faucet.RegisterRESTRoutes(context, router)
	sign.RegisterRESTRoutes(context, router)
	documentIDGetters.RegisterRESTRoutes(context, router)
}
