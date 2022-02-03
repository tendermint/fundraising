package keeper

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetNextAuctionIdWithUpdate increments auction id by one and set it.
func (k Keeper) GetNextAuctionIdWithUpdate(ctx sdk.Context) uint64 {
	id := k.GetLastAuctionId(ctx) + 1
	k.SetAuctionId(ctx, id)
	return id
}

// DistributeSellingCoin releases designated selling coin from the selling reserve account.
func (k Keeper) DistributeSellingCoin(ctx sdk.Context, auction types.AuctionI) error {
	sellingReserveAddress := auction.GetSellingReserveAddress()

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	totalBidCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, sdk.ZeroInt())

	// Distribute coins to all bidders from the selling reserve account
	for _, bid := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		receiveAmt := bid.Coin.Amount.ToDec().QuoTruncate(bid.Price).TruncateInt()
		receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, receiveAmt)

		bidderAddr, err := sdk.AccAddressFromBech32(bid.GetBidder())
		if err != nil {
			return err
		}

		inputs = append(inputs, banktypes.NewInput(sellingReserveAddress, sdk.NewCoins(receiveCoin)))
		outputs = append(outputs, banktypes.NewOutput(bidderAddr, sdk.NewCoins(receiveCoin)))

		totalBidCoin = totalBidCoin.Add(receiveCoin)
	}

	reserveBalance := k.bankKeeper.GetBalance(ctx, sellingReserveAddress, auction.GetSellingCoin().Denom)
	remainingCoin := reserveBalance.Sub(totalBidCoin)

	// Send remaining coin to the auctioneer
	inputs = append(inputs, banktypes.NewInput(sellingReserveAddress, sdk.NewCoins(remainingCoin)))
	outputs = append(outputs, banktypes.NewOutput(auction.GetAuctioneer(), sdk.NewCoins(remainingCoin)))

	// Send all at once
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	return nil
}

// DistributePayingCoin releases the selling coin from the vesting reserve account.
func (k Keeper) DistributePayingCoin(ctx sdk.Context, auction types.AuctionI) error {
	lenVestingQueue := len(k.GetVestingQueuesByAuctionId(ctx, auction.GetId()))

	for i, vq := range k.GetVestingQueuesByAuctionId(ctx, auction.GetId()) {
		if vq.IsVestingReleasable(ctx.BlockTime()) {
			vestingReserveAddress := auction.GetVestingReserveAddress()

			if err := k.bankKeeper.SendCoins(ctx, vestingReserveAddress, auction.GetAuctioneer(), sdk.NewCoins(vq.PayingCoin)); err != nil {
				return sdkerrors.Wrap(err, "failed to release paying coin to the auctioneer")
			}

			vq.Released = true
			k.SetVestingQueue(ctx, auction.GetId(), vq.ReleaseTime, vq)

			// Set finished status when vesting schedule is ended
			if i == lenVestingQueue-1 {
				if err := auction.SetStatus(types.AuctionStatusFinished); err != nil {
					return err
				}

				k.SetAuction(ctx, auction)
			}
		}
	}

	return nil
}

// ReserveSellingCoin reserves the selling coin to the selling reserve account.
func (k Keeper) ReserveSellingCoin(ctx sdk.Context, auctionId uint64, auctioneerAddr sdk.AccAddress, sellingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, auctioneerAddr, types.SellingReserveAddress(auctionId), sdk.NewCoins(sellingCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to reserve selling coin")
	}
	return nil
}

// ReleaseSellingCoin releases the selling coin to the auctioneer.
func (k Keeper) ReleaseSellingCoin(ctx sdk.Context, auction types.AuctionI) error {
	sellingReserveAddress := auction.GetSellingReserveAddress()
	auctioneerAddr := auction.GetAuctioneer()

	reserveBalance := k.bankKeeper.GetBalance(ctx, sellingReserveAddress, auction.GetSellingCoin().Denom)

	if err := k.bankKeeper.SendCoins(ctx, sellingReserveAddress, auctioneerAddr, sdk.NewCoins(reserveBalance)); err != nil {
		return sdkerrors.Wrap(err, "failed to release selling coin")
	}
	return nil
}

// CreateFixedPriceAuction sets fixed price auction.
func (k Keeper) CreateFixedPriceAuction(ctx sdk.Context, msg *types.MsgCreateFixedPriceAuction) (*types.FixedPriceAuction, error) {
	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	auctioneerAddr, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		return &types.FixedPriceAuction{}, err
	}

	if ctx.BlockTime().After(msg.EndTime) {
		return &types.FixedPriceAuction{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be prior to current time")
	}

	params := k.GetParams(ctx)
	feeCollectorAddr, err := sdk.AccAddressFromBech32(params.FeeCollectorAddress)
	if err != nil {
		return &types.FixedPriceAuction{}, err
	}

	if err := k.bankKeeper.SendCoins(ctx, msg.GetAuctioneer(), feeCollectorAddr, params.AuctionCreationFee); err != nil {
		return &types.FixedPriceAuction{}, sdkerrors.Wrap(err, "failed to pay auction creation fee")
	}

	if err := k.ReserveSellingCoin(ctx, nextId, auctioneerAddr, msg.SellingCoin); err != nil {
		return &types.FixedPriceAuction{}, err
	}

	baseAuction := types.NewBaseAuction(
		nextId,
		types.AuctionTypeFixedPrice,
		[]types.AllowedBidder{}, // it is always nil when an auction is created
		auctioneerAddr.String(),
		types.SellingReserveAddress(nextId).String(),
		types.PayingReserveAddress(nextId).String(),
		msg.StartPrice,
		msg.SellingCoin,
		msg.PayingCoinDenom,
		types.VestingReserveAddress(nextId).String(),
		msg.VestingSchedules,
		sdk.ZeroDec(),
		msg.SellingCoin, // add selling coin to remaining coin
		msg.StartTime,
		[]time.Time{msg.EndTime},
		types.AuctionStatusStandBy,
	)

	// Update status if the start time is already passed over the current time
	if baseAuction.IsAuctionStarted(ctx.BlockTime()) {
		baseAuction.Status = types.AuctionStatusStarted
	}

	auction := types.NewFixedPriceAuction(baseAuction)
	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedPriceAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyAuctioneerAddress, auctioneerAddr.String()),
			sdk.NewAttribute(types.AttributeKeyStartPrice, msg.StartPrice.String()),
			sdk.NewAttribute(types.AttributeKeySellingReserveAddress, auction.GetSellingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyPayingReserveAddress, auction.GetPayingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyVestingReserveAddress, auction.GetVestingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeySellingCoin, msg.SellingCoin.String()),
			sdk.NewAttribute(types.AttributeKeyPayingCoinDenom, msg.PayingCoinDenom),
			sdk.NewAttribute(types.AttributeKeyStartTime, msg.StartTime.String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyAuctionStatus, types.AuctionStatusStandBy.String()),
		),
	})

	return auction, nil
}

// CreateEnglishAuction sets english auction.
func (k Keeper) CreateEnglishAuction(ctx sdk.Context, msg *types.MsgCreateEnglishAuction) error {
	// TODO: not implemented yet
	return nil
}

// CancelAuction cancels the auction in an event when the auctioneer needs to modify the auction.
// However, it can only be canceled when the auction has not started yet.
func (k Keeper) CancelAuction(ctx sdk.Context, msg *types.MsgCancelAuction) (types.AuctionI, error) {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", msg.AuctionId)
	}

	if auction.GetAuctioneer().String() != msg.Auctioneer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify ownership of the auction")
	}

	if auction.GetStatus() != types.AuctionStatusStandBy {
		return nil, sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "auction cannot be canceled due to current status")
	}

	if err := k.ReleaseSellingCoin(ctx, auction); err != nil {
		return nil, err
	}

	_ = auction.SetRemainingCoin(sdk.NewCoin(auction.GetSellingCoin().Denom, sdk.ZeroInt()))
	_ = auction.SetStatus(types.AuctionStatusCancelled)

	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
		),
	})

	return auction, nil
}

// AddAllowedBidders is a function for an external module and it simply adds new bidder(s) to AllowedBidder list.
// Note that it doesn't do auctioneer verification because the module is generalized for broader use cases.
// It is designed to delegate to an external module to add necessary verification and logics depending on their use case.
func (k Keeper) AddAllowedBidders(ctx sdk.Context, auctionId uint64, bidders []types.AllowedBidder) error {
	auction, found := k.GetAuction(ctx, auctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", auctionId)
	}

	if len(bidders) == 0 {
		return types.ErrEmptyAllowedBidders
	}

	for _, bidder := range bidders {
		if bidder.MaxBidAmount.IsNegative() {
			return types.ErrInvalidMaxBidAmount
		}
	}

	// Add one or more bidders to allowed bidder list
	if err := auction.SetAllowedBidders(bidders); err != nil {
		return err
	}
	k.SetAuction(ctx, auction)

	return nil
}

// UpdateAllowedBidder is a function for an external module and it simply updates the bidder's maximum bid amount.
// Note that it doesn't do auctioneer verification because the module is generalized for broader use cases.
// It is designed to delegate to an external module to add necessary verification and logics depending on their use case.
func (k Keeper) UpdateAllowedBidder(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, maxBidAmount sdk.Int) error {
	auction, found := k.GetAuction(ctx, auctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", auctionId)
	}

	// Check the bidder existence in the allowed bidder list
	allowedBidders := make(map[string]sdk.Int)
	for _, bidder := range auction.GetAllowedBidders() {
		allowedBidders[bidder.GetBidder()] = bidder.MaxBidAmount
	}

	if _, found := allowedBidders[bidder.String()]; !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bidder %s is not found", bidder.String())
	}

	if !maxBidAmount.IsPositive() {
		return types.ErrInvalidMaxBidAmount
	}

	// Update the bidder's maximum bid amount
	if err := auction.SetAllowedBidders([]types.AllowedBidder{{
		Bidder:       bidder.String(),
		MaxBidAmount: maxBidAmount,
	}}); err != nil {
		return err
	}
	k.SetAuction(ctx, auction)

	return nil
}
