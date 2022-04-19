package fundraising

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Get all auctions from the store and execute operations depending on auction status.
	for _, auction := range k.GetAuctions(ctx) {
		switch auction.GetStatus() {
		case types.AuctionStatusStandBy:
			k.ExecuteStandByStatus(ctx, auction)

		case types.AuctionStatusStarted:
			k.ExecuteStartedStatus(ctx, auction)

		case types.AuctionStatusVesting:
			k.ExecuteVestingStatus(ctx, auction)

		default:
			continue
		}
	}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// // Get all auctions from the store and execute operations depending on auction status.
	// for _, auction := range k.GetAuctions(ctx) {
	// 	switch auction.GetStatus() {
	// 	// case types.AuctionStatusStandBy:
	// 	// 	k.ExecuteStandByStatus(ctx, auction) 10 EndBlocker Started -> 11 Endblocker

	// 	case types.AuctionStatusStarted:
	// 		k.ExecuteStartedStatus(ctx, auction)

	// 	case types.AuctionStatusVesting:
	// 		k.ExecuteVestingStatus(ctx, auction)

	// 	default:
	// 		continue
	// 	}
	// }
}
