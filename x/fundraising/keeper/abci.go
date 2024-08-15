package keeper

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (k Keeper) BeginBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Get all auctions from the store and execute operations depending on auction status.
	auctions, err := k.Auctions(ctx)
	if err != nil {
		return err
	}
	for _, auction := range auctions {
		switch auction.GetStatus() {
		case types.AuctionStatusStandBy:
			err = k.ExecuteStandByStatus(ctx, auction)
		case types.AuctionStatusStarted:
			err = k.ExecuteStartedStatus(ctx, auction)
		case types.AuctionStatusVesting:
			err = k.ExecuteVestingStatus(ctx, auction)
		default:
			err = fmt.Errorf("invalid auction status %s", auction.GetStatus())
		}
	}
	return err
}
