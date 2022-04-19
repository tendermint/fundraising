package simulation

// DONTCOVER

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Simulation parameter constants.
const (
	AuctionCreationFee = "auction_creation_fee"
	ExtendedPeriod     = "extended_period"
)

// GenAuctionCreationFee return randomized auction creation fee.
func GenAuctionCreationFee(r *rand.Rand) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(simulation.RandIntBetween(r, 0, 100_000_000))))
}

// GenExtendedPeriod return default extended period.
func GenExtendedPeriod(r *rand.Rand) uint32 {
	return uint32(simulation.RandIntBetween(r, int(types.DefaultExtendedPeriod), 10))
}

// RandomizedGenState generates a random GenesisState.
func RandomizedGenState(simState *module.SimulationState) {
	var auctionCreationFee sdk.Coins
	simState.AppParams.GetOrGenerate(
		simState.Cdc, AuctionCreationFee, &auctionCreationFee, simState.Rand,
		func(r *rand.Rand) { auctionCreationFee = GenAuctionCreationFee(r) },
	)

	var extendedPeriod uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ExtendedPeriod, &extendedPeriod, simState.Rand,
		func(r *rand.Rand) { extendedPeriod = GenExtendedPeriod(r) },
	)

	genState := types.GenesisState{
		Params: types.Params{
			AuctionCreationFee: auctionCreationFee,
			ExtendedPeriod:     extendedPeriod,
		},
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&genState)
}
