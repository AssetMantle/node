package docs

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	errorConstants "github.com/AssetMantle/schema/x/errors/constants"
	"github.com/AssetMantle/schema/x/ids"
	baseIDs "github.com/AssetMantle/schema/x/ids/base"
)

func splitIDHandler(context client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		transactionRequest := Prototype()
		if !rest.ReadRESTReq(responseWriter, httpRequest, context.LegacyAmino, &transactionRequest) {
			panic(errorConstants.IncorrectFormat)
		}

		if rest.CheckBadRequestError(responseWriter, transactionRequest.Validate()) {
			panic(errorConstants.IncorrectFormat)
		}

		req := transactionRequest.(request)

		fromID, _ := baseIDs.ReadIdentityID(req.FromID)

		coins, _ := sdkTypes.ParseCoinsNormalized(req.Coins)

		var coinID ids.CoinID
		for _, coin := range coins {
			coinID = baseIDs.NewCoinID(baseIDs.NewStringID(coin.Denom))
		}

		rest.PostProcessResponse(responseWriter, context, newResponse(baseIDs.NewSplitID(fromID, coinID).AsString(), "", nil))
	}
}
