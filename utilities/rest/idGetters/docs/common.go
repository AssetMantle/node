package docs

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	codecUtilities "github.com/AssetMantle/schema/utilities/codec"
	baseData "github.com/AssetMantle/schema/x/data/base"
	"github.com/AssetMantle/schema/x/ids"
	baseIDs "github.com/AssetMantle/schema/x/ids/base"
	"github.com/AssetMantle/schema/x/lists"
	"github.com/AssetMantle/schema/x/lists/utilities"
	baseProperties "github.com/AssetMantle/schema/x/properties/base"
	"github.com/AssetMantle/schema/x/properties/constants"
	propertiesUtilities "github.com/AssetMantle/schema/x/properties/utilities"
	"github.com/AssetMantle/schema/x/qualified"
	"github.com/AssetMantle/schema/x/qualified/base"
)

func RegisterLegacyAminoCodec(legacyAmino *codec.LegacyAmino) {
	codecUtilities.RegisterModuleConcrete(legacyAmino, request{})
}

func GetTotalWeight(immutables qualified.Immutables, mutables qualified.Mutables) int64 {
	totalWeight := int64(0)
	for _, property := range append(immutables.GetImmutablePropertyList().GetList(), mutables.GetMutablePropertyList().GetList()...) {
		totalWeight += property.Get().GetBondWeight()
	}
	return totalWeight
}

func ReadAndProcess(context client.Context, addAuth bool, addBond bool, responseWriter http.ResponseWriter, httpRequest *http.Request) (ids.ClassificationID, qualified.Immutables, qualified.Mutables) {
	_, classificationID, immutableMetaProperties, immutableProperties, mutableMetaProperties, mutableProperties := Read(context, responseWriter, httpRequest)
	Immutables, Mutables := Process(immutableMetaProperties, immutableProperties, mutableMetaProperties, mutableProperties, addAuth, addBond)
	if len(classificationID.Bytes()) != 0 {
		return classificationID, Immutables, Mutables
	}
	return baseIDs.NewClassificationID(Immutables, Mutables), Immutables, Mutables
}

func Read(context client.Context, responseWriter http.ResponseWriter, httpRequest *http.Request) (request, ids.ClassificationID, lists.PropertyList, lists.PropertyList, lists.PropertyList, lists.PropertyList) {
	transactionRequest := Prototype()
	if !rest.ReadRESTReq(responseWriter, httpRequest, context.LegacyAmino, &transactionRequest) {
		return request{}, nil, nil, nil, nil, nil
	}

	if rest.CheckBadRequestError(responseWriter, transactionRequest.Validate()) {
		return request{}, nil, nil, nil, nil, nil
	}

	req := transactionRequest.(request)

	immutableMetaProperties, _ := utilities.ReadMetaPropertyList(req.ImmutableMetaProperties)

	immutableProperties, _ := utilities.ReadMetaPropertyList(req.ImmutableProperties)

	immutableProperties = immutableProperties.ScrubData()

	mutableMetaProperties, _ := utilities.ReadMetaPropertyList(req.MutableMetaProperties)

	mutableProperties, _ := utilities.ReadMetaPropertyList(req.MutableProperties)

	mutableProperties = mutableProperties.ScrubData()

	classificationID, _ := baseIDs.ReadClassificationID(req.ClassificationID)
	return req, classificationID, immutableMetaProperties, immutableProperties, mutableMetaProperties, mutableProperties
}

func Process(immutableMetaProperties, immutableProperties, mutableMetaProperties, mutableProperties lists.PropertyList, addAuth bool, addBond bool) (qualified.Immutables, qualified.Mutables) {
	immutables := base.NewImmutables(immutableMetaProperties.Add(propertiesUtilities.AnyPropertyListToPropertyList(immutableProperties.GetList()...)...))
	var Mutables qualified.Mutables
	if addAuth {
		Mutables = base.NewMutables(mutableMetaProperties.Add(propertiesUtilities.AnyPropertyListToPropertyList(mutableProperties.Add(constants.AuthenticationProperty.ToAnyProperty()).GetList()...)...))
	} else {
		Mutables = base.NewMutables(mutableMetaProperties.Add(propertiesUtilities.AnyPropertyListToPropertyList(mutableProperties.GetList()...)...))
	}
	var Immutables qualified.Immutables
	if addBond {
		Immutables = base.NewImmutables(immutables.GetImmutablePropertyList().Add(baseProperties.NewMetaProperty(constants.BondAmountProperty.GetKey(), baseData.NewNumberData(GetTotalWeight(immutables, Mutables)*baseData.NewNumberData(1).Get()))))
	} else {
		Immutables = immutables
	}
	return Immutables, Mutables
}

func RegisterRESTRoutes(context client.Context, router *mux.Router) {
	router.HandleFunc("/get/classification/identity", identityClassificationHandler(context)).Methods("POST")
	router.HandleFunc("/get/classification/asset", assetClassificationHandler(context)).Methods("POST")
	router.HandleFunc("/get/classification/order", orderClassificationHandler(context)).Methods("POST")
	router.HandleFunc("/get/document/assetID", assetIDHandler(context)).Methods("POST")
	router.HandleFunc("/get/document/identityID", identityIDHandler(context)).Methods("POST")
	router.HandleFunc("/get/document/orderID", orderIDHandler(context)).Methods("POST")
	router.HandleFunc("/get/document/nubID", nubIDHandler(context)).Methods("POST")
	router.HandleFunc("/get/document/splitID", splitIDHandler(context)).Methods("POST")
}
