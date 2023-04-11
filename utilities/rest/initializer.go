package rest

import (
	"node/utilities/rest/keys/add"
	"node/utilities/rest/sign"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func RegisterRESTRoutes(context client.Context, router *mux.Router) {
	add.RegisterRESTRoutes(context, router)
	sign.RegisterRESTRoutes(context, router)
}
