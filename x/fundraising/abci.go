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

	// 1. AuctionStatusVesting: look up the release time of vesting queue for the auction and
	//	  see if the module needs to be release the vested amount of coin to the auctioneer
	//
	// 2. AuctionStatusStandBy: update the status to AUCTION_STATUS_STARTED
	//    if the start time is passed over the current time
	//
	// 3. AuctionStatusStarted: proceed to calculate allocation for each bidder
	//    if the start time is on time or has passed the current block time.
	//
	for _, auction := range k.GetAuctions(ctx) {
		if auction.GetType() == types.AuctionTypeFixedPrice {
			switch auction.GetStatus() {
			case types.AuctionStatusVesting:
				if err := k.DistributePayingCoin(ctx, auction); err != nil {
					panic(err)
				}

			case types.AuctionStatusStandBy:
				if types.IsAuctionStarted(auction.GetStartTime(), ctx.BlockTime()) {
					if err := auction.SetStatus(types.AuctionStatusStarted); err != nil {
						logger.Error("error is returned when setting auction status", "auction", auction)
					}
				}

			case types.AuctionStatusStarted:
				if !auction.GetStartTime().Before(ctx.BlockTime()) {
					if err := k.DistributeSellingCoin(ctx, auction); err != nil {
						panic(err)
					}

					// Set vesting schedules and change status
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
