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

	for _, auction := range k.GetAuctions(ctx) {
		if auction.GetType() == types.AuctionTypeFixedPrice {
			switch auction.GetStatus() {
			case types.AuctionStatusVesting:
				// Look up the release time of vesting queue for the auction and
				// see if the module needs to be distribute the vested amount of coin to the auctioneer

				// TODO: get all vesting queues (auctionId -> ProtocolBuffer(VestingQueue))

			case types.AuctionStatusStandBy:
				// Update the status to AUCTION_STATUS_STARTED if the start time is passed over the current time
				if types.IsAuctionStarted(auction, ctx.BlockTime()) {
					auction.SetStatus(types.AuctionStatusStarted)
				}
			case types.AuctionStatusFinished:
				// Calculate allocation for each bidder of the auction and distribute them to the bidders
				// Lastly, store vesting queue if the auction has any vesting schedules and set status to AuctionStatusVesting
				for _, bid := range k.GetBids(ctx, auction.GetId()) {
					logger.Info("Bid information", "bid", bid)
				}

			default:
				continue
			}
		}
	}
}
