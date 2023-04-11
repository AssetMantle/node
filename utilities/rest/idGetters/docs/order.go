package docs

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/AssetMantle/schema/x/data"
	baseData "github.com/AssetMantle/schema/x/data/base"
	baseIDs "github.com/AssetMantle/schema/x/ids/base"
	"github.com/AssetMantle/schema/x/properties"
	baseProperties "github.com/AssetMantle/schema/x/properties/base"
	"github.com/AssetMantle/schema/x/properties/constants"
	propertiesUtilities "github.com/AssetMantle/schema/x/properties/utilities"
	"github.com/AssetMantle/schema/x/qualified/base"
	baseTypes "github.com/AssetMantle/schema/x/types/base"
)

func orderIDHandler(context client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		req, classificationID, ImmutableMetaProperties, ImmutableProperties, _, _ := Read(context, responseWriter, httpRequest)
		makerOwnableSplit, _ := sdkTypes.NewDecFromStr(req.MakerOwnableSplit)

		takerOwnableSplit, _ := sdkTypes.NewDecFromStr(req.TakerOwnableSplit)

		fromID, _ := baseIDs.ReadIdentityID(req.FromID)

		takerID, _ := baseIDs.ReadIdentityID(req.TakerID)

		makerOwnableID, _ := baseIDs.ReadOwnableID(req.MakerOwnableID)
		height, _ := strconv.Atoi(req.Height)
		takerOwnableID, _ := baseIDs.ReadOwnableID(req.TakerOwnableID)

		immutableMetaProperties := ImmutableMetaProperties.
			Add(baseProperties.NewMetaProperty(constants.ExchangeRateProperty.GetKey(), baseData.NewDecData(takerOwnableSplit.QuoTruncate(sdkTypes.SmallestDec()).QuoTruncate(makerOwnableSplit)))).
			Add(baseProperties.NewMetaProperty(constants.CreationHeightProperty.GetKey(), baseData.NewHeightData(baseTypes.NewHeight(int64(height))))).
			Add(baseProperties.NewMetaProperty(constants.MakerOwnableIDProperty.GetKey(), baseData.NewIDData(makerOwnableID))).
			Add(baseProperties.NewMetaProperty(constants.TakerOwnableIDProperty.GetKey(), baseData.NewIDData(takerOwnableID))).
			Add(baseProperties.NewMetaProperty(constants.MakerIDProperty.GetKey(), baseData.NewIDData(fromID))).
			Add(baseProperties.NewMetaProperty(constants.TakerIDProperty.GetKey(), baseData.NewIDData(takerID)))

		Immutables := base.NewImmutables(immutableMetaProperties.Add(propertiesUtilities.AnyPropertyListToPropertyList(ImmutableProperties.GetList()...)...))

		// Mutables := base.NewMutables(baseLists.NewPropertyList(propertiesUtilities.AnyPropertyListToPropertyList(append(mutableMetaProperties.GetList(), mutableProperties.GetList()...)...)...))

		// Immutables := base.NewImmutables(immutables.GetImmutablePropertyList().Add(baseProperties.NewMetaProperty(constants.BondAmountProperty.GetKey(), baseData.NewDecData(GetTotalWeight(immutables, Mutables).Mul(sdkTypes.NewDec(1))))))
		rest.PostProcessResponse(responseWriter, context, newResponse(baseIDs.NewOrderID(classificationID, Immutables).AsString(), "", nil))
	}
}
func orderClassificationHandler(context client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		_, _, immutableMetaProperties, immutableProperties, mutableMetaProperties, mutableProperties := Read(context, responseWriter, httpRequest)
		immutables := base.NewImmutables(
			immutableMetaProperties.Add(
				propertiesUtilities.AnyPropertyListToPropertyList(
					immutableProperties.Add(
						constants.ExchangeRateProperty.ToAnyProperty(),
						constants.CreationHeightProperty.ToAnyProperty(),
						constants.MakerOwnableIDProperty.ToAnyProperty(),
						constants.TakerOwnableIDProperty.ToAnyProperty(),
						constants.MakerIDProperty.ToAnyProperty(),
						constants.TakerIDProperty.ToAnyProperty(),
					).GetList()...,
				)...,
			),
		)
		mutables := base.NewMutables(
			mutableMetaProperties.Add(
				propertiesUtilities.AnyPropertyListToPropertyList(
					mutableProperties.Add(
						constants.ExpiryHeightProperty.ToAnyProperty(),
						constants.MakerOwnableSplitProperty.ToAnyProperty(),
					).GetList()...,
				)...,
			),
		)
		Immutables := base.NewImmutables(immutables.GetImmutablePropertyList().Add(baseProperties.NewMetaProperty(constants.BondAmountProperty.GetKey(), baseData.NewNumberData(GetTotalWeight(immutables, mutables)*baseData.NewNumberData(1).Get()))))
		rest.PostProcessResponse(responseWriter, context, newResponse(baseIDs.NewClassificationID(Immutables, mutables).AsString(), Immutables.GetProperty(constants.BondAmountProperty.GetID()).Get().(properties.MetaProperty).GetData().Get().(data.NumberData).AsString(), nil))
	}
}
