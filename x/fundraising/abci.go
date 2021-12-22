package fundraising

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Get all auctions in the store and proceed operations depending on auction status.
	//
	// For AuctionStatusVesting, it first gets all vesting queues in the store and
	// look up the release time of each vesting queue to see if the module needs to distribute
	// the paying coin to the auctioneer.
	//
	// For AuctionStatusStandBy, it compares the current and start time of the auction and
	// update the status if the auction status is ready to be updated.
	//
	// For AuctionStatusStarted, distribute the allocated paying coin to bidders for the auction and
	// set vesting schedules if they are defined.
	for _, auction := range k.GetAuctions(ctx) {
		switch auction.GetStatus() {
		case types.AuctionStatusStandBy:
			if types.IsAuctionStarted(auction.GetStartTime(), ctx.BlockTime()) {
				_ = auction.SetStatus(types.AuctionStatusStarted)
				k.SetAuction(ctx, auction)
			}

		case types.AuctionStatusStarted:
			if auction.GetType() == types.AuctionTypeFixedPrice {
				if types.IsAuctionFinished(auction.GetEndTimes()[0], ctx.BlockTime()) {
					if err := k.DistributeSellingCoin(ctx, auction); err != nil {
						panic(err)
					}

					if err := k.SetVestingSchedules(ctx, auction); err != nil {
						panic(err)
					}
				}
			}

		case types.AuctionStatusVesting:
			if err := k.DistributePayingCoin(ctx, auction); err != nil {
				panic(err)
			}

		default:
			continue
		}
	}
}
