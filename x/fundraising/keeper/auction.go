package keeper

import (
	"fmt"
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

// ReserveCreationFees reserves the auction creation fee to the fee collector account.
func (k Keeper) ReserveCreationFees(ctx sdk.Context, auctioneerAddr sdk.AccAddress) error {
	params := k.GetParams(ctx)
	feeCollectorAddr, err := sdk.AccAddressFromBech32(params.FeeCollectorAddress)
	if err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoins(ctx, auctioneerAddr, feeCollectorAddr, params.AuctionCreationFee); err != nil {
		return sdkerrors.Wrap(err, "failed to reserve auction creation fee")
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

	if err := k.ReserveCreationFees(ctx, auctioneerAddr); err != nil {
		return &types.FixedPriceAuction{}, err
	}

	if err := k.ReserveSellingCoin(ctx, nextId, auctioneerAddr, msg.SellingCoin); err != nil {
		return &types.FixedPriceAuction{}, err
	}

	baseAuction := types.NewBaseAuction(
		nextId,
		types.AuctionTypeFixedPrice,
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

	// Update status if the start time is already passed the current time
	if baseAuction.IsAuctionStarted(ctx.BlockTime()) {
		baseAuction.Status = types.AuctionStatusStarted
	}

	auction := types.NewFixedPriceAuction(baseAuction)
	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedPriceAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyAuctioneerAddress, auction.GetAuctioneer().String()),
			sdk.NewAttribute(types.AttributeKeyStartPrice, auction.GetStartPrice().String()),
			sdk.NewAttribute(types.AttributeKeySellingReserveAddress, auction.GetSellingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyPayingReserveAddress, auction.GetPayingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyVestingReserveAddress, auction.GetVestingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeySellingCoin, auction.GetSellingCoin().String()),
			sdk.NewAttribute(types.AttributeKeyPayingCoinDenom, auction.GetPayingCoinDenom()),
			sdk.NewAttribute(types.AttributeKeyStartTime, auction.GetStartTime().String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyAuctionStatus, auction.GetStatus().String()),
		),
	})

	return auction, nil
}

// CreateEnglishAuction sets english auction.
func (k Keeper) CreateEnglishAuction(ctx sdk.Context, msg *types.MsgCreateEnglishAuction) (*types.EnglishAuction, error) {
	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	auctioneerAddr, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		return &types.EnglishAuction{}, err
	}

	if ctx.BlockTime().After(msg.EndTime) {
		return &types.EnglishAuction{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be prior to current time")
	}

	if err := k.ReserveCreationFees(ctx, auctioneerAddr); err != nil {
		return &types.EnglishAuction{}, err
	}

	if err := k.ReserveSellingCoin(ctx, nextId, auctioneerAddr, msg.SellingCoin); err != nil {
		return &types.EnglishAuction{}, err
	}

	baseAuction := types.NewBaseAuction(
		nextId,
		types.AuctionTypeEnglish,
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

	// Update status if the start time is already passed the current time
	if baseAuction.IsAuctionStarted(ctx.BlockTime()) {
		baseAuction.Status = types.AuctionStatusStarted
	}

	auction := types.NewEnglishAuction(
		baseAuction,
		msg.MaximumBidPrice,
		msg.Extended,
		msg.ExtendRate,
	)
	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedPriceAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyAuctioneerAddress, auction.GetAuctioneer().String()),
			sdk.NewAttribute(types.AttributeKeyStartPrice, auction.GetStartPrice().String()),
			sdk.NewAttribute(types.AttributeKeySellingReserveAddress, auction.GetSellingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyPayingReserveAddress, auction.GetPayingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyVestingReserveAddress, auction.GetVestingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeySellingCoin, auction.GetSellingCoin().String()),
			sdk.NewAttribute(types.AttributeKeyPayingCoinDenom, auction.GetPayingCoinDenom()),
			sdk.NewAttribute(types.AttributeKeyStartTime, auction.GetStartTime().String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyAuctionStatus, auction.GetStatus().String()),
			sdk.NewAttribute(types.AttributeKeyMaximumBidPrice, msg.MaximumBidPrice.String()),
			sdk.NewAttribute(types.AttributeKeyExtended, strconv.FormatUint(uint64(msg.Extended), 10)),
			sdk.NewAttribute(types.AttributeKeyExtendRate, msg.ExtendRate.String()),
		),
	})

	return auction, nil
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

	// Release the selling coin to the auctioneer.
	sellingReserveAddr := auction.GetSellingReserveAddress()
	auctioneerAddr := auction.GetAuctioneer()
	reserveCoin := k.bankKeeper.GetBalance(ctx, sellingReserveAddr, auction.GetSellingCoin().Denom)

	if err := k.bankKeeper.SendCoins(ctx, sellingReserveAddr, auctioneerAddr, sdk.NewCoins(reserveCoin)); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to release selling coin")
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

// CalculateWinners ...
// 1. Calculate the winners and if the extend rate is exceeded, extend the end time
// 2. Distribute allocation to the winners
// 3. Transfer the paying coin to the auctioneer with the given vesting schedules
// By default, extended auction round is triggered once for all english auctions
func (k Keeper) CalculateWinners(ctx sdk.Context, auction types.AuctionI) error {
	bids := k.GetBidsByAuctionId(ctx, auction.GetId())
	bids = types.SanitizeReverseBids(bids)

	// first round needs to calculate the winning price
	if len(auction.GetEndTimes()) == 1 {
		totalSellingAmt := sdk.ZeroDec()
		totalCoinAmt := sdk.ZeroDec()
		remainingAmt := auction.GetRemainingCoin().Amount

		for _, bid := range bids {
			totalCoinAmt = totalCoinAmt.Add(bid.Coin.Amount.ToDec())
			totalSellingAmt = totalCoinAmt.QuoTruncate(bid.Price)
		}

		remainingAmt = remainingAmt.Sub(totalSellingAmt.TruncateInt())
		remainingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, remainingAmt)

		_ = auction.SetRemainingCoin(remainingCoin)

		// TODO: fillPrice, store winning bids list, and set second last time (current block time)

	} else {
		// TODO
		fmt.Println("")
	}

	// TODO: distribution and transferring the paying coin to the auctioneer

	return nil
}
