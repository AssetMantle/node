package docs

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/AssetMantle/schema/x/data"
	baseData "github.com/AssetMantle/schema/x/data/base"
	errorConstants "github.com/AssetMantle/schema/x/errors/constants"
	baseIDs "github.com/AssetMantle/schema/x/ids/base"
	baseLists "github.com/AssetMantle/schema/x/lists/base"
	"github.com/AssetMantle/schema/x/properties"
	baseProperties "github.com/AssetMantle/schema/x/properties/base"
	"github.com/AssetMantle/schema/x/properties/constants"
	"github.com/AssetMantle/schema/x/qualified/base"
)

func nubIDHandler(context client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		transactionRequest := Prototype()
		if !rest.ReadRESTReq(responseWriter, httpRequest, context.LegacyAmino, &transactionRequest) {
			panic(errorConstants.IncorrectFormat)
		}

		if rest.CheckBadRequestError(responseWriter, transactionRequest.Validate()) {
			panic(errorConstants.IncorrectFormat)
		}

		req := transactionRequest.(request)

		nubID := baseIDs.NewStringID(req.NubID)
		immutables := base.NewImmutables(baseLists.NewPropertyList(baseProperties.NewMetaProperty(constants.NubIDProperty.GetKey(), baseData.NewIDData(nubID))))

		// TODO move to proper package
		var NubClassificationID = baseIDs.NewClassificationID(base.NewImmutables(baseLists.NewPropertyList(constants.NubIDProperty)), base.NewMutables(baseLists.NewPropertyList(constants.AuthenticationProperty)))

		rest.PostProcessResponse(responseWriter, context, newResponse(baseIDs.NewIdentityID(NubClassificationID, immutables).AsString(), "", nil))
	}
}

func identityIDHandler(context client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		classificationID, immutables, _ := ReadAndProcess(context, false, false, responseWriter, httpRequest)

		rest.PostProcessResponse(responseWriter, context, newResponse(baseIDs.NewIdentityID(classificationID, immutables).AsString(), "", nil))
	}
}

func identityClassificationHandler(context client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		id, immutables, _ := ReadAndProcess(context, true, true, responseWriter, httpRequest)
		rest.PostProcessResponse(responseWriter, context, newResponse(id.AsString(), immutables.GetProperty(constants.BondAmountProperty.GetID()).Get().(properties.MetaProperty).GetData().Get().(data.NumberData).AsString(), nil))
	}
}
