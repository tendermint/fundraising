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

	logger := k.Logger(ctx)

	// Get all auctions in the store and proceed operations depending on auction status.
	//
	// For AuctionStatusVesting, it first gets all vesting queues in the store and
	// look up the release time of each vesting queue to see if the module needs to distribute
	// the paying coin to the auctioneer.
	//
	// For AuctionStatusStandBy, it compares the current time and the start time of the auction and
	// updte the status if necessary.
	//
	// For AuctionStatusStarted, distribute the allocated paying coin for bidders for the auction and
	// set vesting schedules if they are defined.
	for _, auction := range k.GetAuctions(ctx) {
		if auction.GetType() == types.AuctionTypeFixedPrice {
			switch auction.GetStatus() {
			case types.AuctionStatusVesting:
				for _, vq := range k.GetVestingQueuesByAuctionId(ctx, auction.GetId()) {
					if !vq.GetReleaseTime().Before(ctx.BlockTime()) {
						if err := k.DistributePayingCoin(ctx, auction); err != nil {
							panic(err)
						}

						vq.Vested = true
						k.SetVestingQueue(ctx, auction.GetId(), vq.ReleaseTime, vq)
					}
				}

			case types.AuctionStatusStandBy:
				if types.IsAuctionStarted(auction.GetStartTime(), ctx.BlockTime()) {
					if err := auction.SetStatus(types.AuctionStatusStarted); err != nil {
						logger.Error("error is returned when setting auction status", "auction", auction)
					}
					k.SetAuction(ctx, auction)
				}

			case types.AuctionStatusStarted:
				if types.IsAuctionFinished(auction.GetEndTimes()[0], ctx.BlockTime()) {
					if err := k.DistributeSellingCoin(ctx, auction); err != nil {
						panic(err)
					}

					if err := k.SetVestingSchedules(ctx, auction); err != nil {
						panic(err)
					}
				}

			default:
				continue
			}
		}
	}
}
