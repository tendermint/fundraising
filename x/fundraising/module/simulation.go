package fundraising

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/tendermint/fundraising/testutil/sample"

	fundraisingsimulation "github.com/tendermint/fundraising/x/fundraising/simulation"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// avoid unused import issue
var _ = sample.AccAddress

const (
	opWeightMsgCreateFixedPriceAuction          = "op_weight_msg_create_fixed_price_auction"
	defaultWeightMsgCreateFixedPriceAuction int = 20

	opWeightMsgCreateBatchAuction          = "op_weight_msg_create_batch_auction"
	defaultWeightMsgCreateBatchAuction int = 20

	opWeightMsgCancelAuction          = "op_weight_msg_cancel_auction"
	defaultWeightMsgCancelAuction int = 15

	opWeightMsgPlaceBid          = "op_weight_msg_place_bid"
	defaultWeightMsgPlaceBid int = 80

	opWeightMsgModifyBid          = "op_weight_msg_modify_bid"
	defaultWeightMsgModifyBid int = 15

	// this line is used by starport scaffolding # simapp/module/const
)

const (
	AuctionCreationFee = "auction_creation_fee"
	ExtendedPeriod     = "extended_period"
)

// GenAuctionCreationFee return randomized auction creation fee.
func GenAuctionCreationFee(r *rand.Rand) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simtypes.RandIntBetween(r, 0, 100_000_000))))
}

// GenExtendedPeriod return default extended period.
func GenExtendedPeriod(r *rand.Rand) uint32 {
	return uint32(simtypes.RandIntBetween(r, int(types.DefaultExtendedPeriod), 10))
}

// RandomizedGenState generates a random GenesisState.
func RandomizedGenState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	fundraisingGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}

	var auctionCreationFee sdk.Coins
	simState.AppParams.GetOrGenerate(AuctionCreationFee, &auctionCreationFee, simState.Rand, func(r *rand.Rand) {
		auctionCreationFee = GenAuctionCreationFee(r)
	})

	var extendedPeriod uint32
	simState.AppParams.GetOrGenerate(ExtendedPeriod, &extendedPeriod, simState.Rand, func(r *rand.Rand) {
		extendedPeriod = GenExtendedPeriod(r)
	})

	fundraisingGenesis.Params.AuctionCreationFee = auctionCreationFee
	fundraisingGenesis.Params.ExtendedPeriod = extendedPeriod
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&fundraisingGenesis)
}

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateFixedPriceAuction int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateFixedPriceAuction, &weightMsgCreateFixedPriceAuction, nil,
		func(_ *rand.Rand) {
			weightMsgCreateFixedPriceAuction = defaultWeightMsgCreateFixedPriceAuction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateFixedPriceAuction,
		fundraisingsimulation.SimulateMsgCreateFixedPriceAuction(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgCreateBatchAuction int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateBatchAuction, &weightMsgCreateBatchAuction, nil,
		func(_ *rand.Rand) {
			weightMsgCreateBatchAuction = defaultWeightMsgCreateBatchAuction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateBatchAuction,
		fundraisingsimulation.SimulateMsgCreateBatchAuction(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgCancelAuction int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelAuction, &weightMsgCancelAuction, nil,
		func(_ *rand.Rand) {
			weightMsgCancelAuction = defaultWeightMsgCancelAuction
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelAuction,
		fundraisingsimulation.SimulateMsgCancelAuction(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgPlaceBid int
	simState.AppParams.GetOrGenerate(opWeightMsgPlaceBid, &weightMsgPlaceBid, nil,
		func(_ *rand.Rand) {
			weightMsgPlaceBid = defaultWeightMsgPlaceBid
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgPlaceBid,
		fundraisingsimulation.SimulateMsgPlaceBid(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	var weightMsgModifyBid int
	simState.AppParams.GetOrGenerate(opWeightMsgModifyBid, &weightMsgModifyBid, nil,
		func(_ *rand.Rand) {
			weightMsgModifyBid = defaultWeightMsgModifyBid
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgModifyBid,
		fundraisingsimulation.SimulateMsgModifyBid(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateFixedPriceAuction,
			defaultWeightMsgCreateFixedPriceAuction,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				fundraisingsimulation.SimulateMsgCreateFixedPriceAuction(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateBatchAuction,
			defaultWeightMsgCreateBatchAuction,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				fundraisingsimulation.SimulateMsgCreateBatchAuction(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCancelAuction,
			defaultWeightMsgCancelAuction,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				fundraisingsimulation.SimulateMsgCancelAuction(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgPlaceBid,
			defaultWeightMsgPlaceBid,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				fundraisingsimulation.SimulateMsgPlaceBid(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgModifyBid,
			defaultWeightMsgModifyBid,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				fundraisingsimulation.SimulateMsgModifyBid(am.accountKeeper, am.bankKeeper, am.keeper, simState.TxConfig)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
