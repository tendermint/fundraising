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

// CreateFixedPriceAuction sets a fixed price auction.
func (k Keeper) CreateFixedPriceAuction(ctx sdk.Context, msg *types.MsgCreateFixedPriceAuction) (types.AuctionI, error) {
	if ctx.BlockTime().After(msg.EndTime) { // EndTime < CurrentTime
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be set after the current time")
	}

	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	if err := k.ReserveCreationFee(ctx, msg.GetAuctioneer()); err != nil {
		return nil, err
	}

	if err := k.ReserveSellingCoin(ctx, nextId, msg.GetAuctioneer(), msg.SellingCoin); err != nil {
		return nil, err
	}

	// Allowed bidder list is empty when an auction is created
	// The module is fundamentally designed to delegate authorization
	// to an external module to add allowed bidder list for an auction
	allowedBidders := []types.AllowedBidder{}
	endTimes := []time.Time{msg.EndTime} // it is an array data type to handle BatchAuction

	ba := types.NewBaseAuction(
		nextId,
		types.AuctionTypeFixedPrice,
		allowedBidders,
		msg.Auctioneer,
		types.SellingReserveAddress(nextId).String(),
		types.PayingReserveAddress(nextId).String(),
		msg.StartPrice,
		msg.SellingCoin,
		msg.PayingCoinDenom,
		types.VestingReserveAddress(nextId).String(),
		msg.VestingSchedules,
		msg.SellingCoin,
		msg.StartTime,
		endTimes,
		types.AuctionStatusStandBy,
	)

	// Update status if the start time is already passed over the current time
	if ba.ShouldAuctionStarted(ctx.BlockTime()) {
		ba.Status = types.AuctionStatusStarted
	}

	auction := types.NewFixedPriceAuction(ba)

	// Call the before auction created hook
	k.BeforeFixedPriceAuctionCreated(
		ctx,
		auction.Auctioneer,
		auction.StartPrice,
		auction.SellingCoin,
		auction.PayingCoinDenom,
		auction.VestingSchedules,
		auction.StartTime,
		auction.EndTimes[0],
	)

	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedPriceAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyAuctioneerAddress, auction.GetAuctioneer().String()),
			sdk.NewAttribute(types.AttributeKeySellingReserveAddress, auction.GetSellingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyPayingReserveAddress, auction.GetPayingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyStartPrice, auction.GetStartPrice().String()),
			sdk.NewAttribute(types.AttributeKeySellingCoin, auction.GetSellingCoin().String()),
			sdk.NewAttribute(types.AttributeKeyPayingCoinDenom, auction.GetPayingCoinDenom()),
			sdk.NewAttribute(types.AttributeKeyVestingReserveAddress, auction.GetVestingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyRemainingSellingCoin, auction.GetRemainingSellingCoin().String()),
			sdk.NewAttribute(types.AttributeKeyStartTime, auction.GetStartTime().String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyAuctionStatus, auction.GetStatus().String()),
		),
	})

	return auction, nil
}

// CreateBatchAuction sets batch auction.
func (k Keeper) CreateBatchAuction(ctx sdk.Context, msg *types.MsgCreateBatchAuction) (types.AuctionI, error) {
	if ctx.BlockTime().After(msg.EndTime) { // EndTime < CurrentTime
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be set after the current time")
	}

	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	if err := k.ReserveCreationFee(ctx, msg.GetAuctioneer()); err != nil {
		return nil, err
	}

	if err := k.ReserveSellingCoin(ctx, nextId, msg.GetAuctioneer(), msg.SellingCoin); err != nil {
		return nil, err
	}

	// Allowed bidder list is empty when an auction is created
	// The module is fundamentally designed to delegate authorization
	// to an external module to add allowed bidder list for an auction
	allowedBidders := []types.AllowedBidder{}
	endTimes := []time.Time{msg.EndTime} // it is an array data type to handle BatchAuction

	baseAuction := types.NewBaseAuction(
		nextId,
		types.AuctionTypeBatch,
		allowedBidders,
		msg.Auctioneer,
		types.SellingReserveAddress(nextId).String(),
		types.PayingReserveAddress(nextId).String(),
		msg.StartPrice,
		msg.SellingCoin,
		msg.PayingCoinDenom,
		types.VestingReserveAddress(nextId).String(),
		msg.VestingSchedules,
		msg.SellingCoin,
		msg.StartTime,
		endTimes,
		types.AuctionStatusStandBy,
	)

	// Update status if the start time is already passed the current time
	if baseAuction.ShouldAuctionStarted(ctx.BlockTime()) {
		baseAuction.Status = types.AuctionStatusStarted
	}

	auction := types.NewBatchAuction(
		baseAuction,
		msg.MinBidPrice,
		sdk.ZeroDec(),
		msg.MaxExtendedRound,
		msg.ExtendedRoundRate,
	)

	// Call the before auction created hook
	k.BeforeBatchAuctionCreated(
		ctx,
		auction.Auctioneer,
		auction.StartPrice,
		auction.MinBidPrice,
		auction.SellingCoin,
		auction.PayingCoinDenom,
		auction.VestingSchedules,
		auction.MaxExtendedRound,
		auction.ExtendedRoundRate,
		auction.StartTime,
		auction.EndTimes[0],
	)

	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateBatchAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyAuctioneerAddress, auction.GetAuctioneer().String()),
			sdk.NewAttribute(types.AttributeKeySellingReserveAddress, auction.GetSellingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyPayingReserveAddress, auction.GetPayingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyStartPrice, auction.GetStartPrice().String()),
			sdk.NewAttribute(types.AttributeKeySellingCoin, auction.GetSellingCoin().String()),
			sdk.NewAttribute(types.AttributeKeyPayingCoinDenom, auction.GetPayingCoinDenom()),
			sdk.NewAttribute(types.AttributeKeyVestingReserveAddress, auction.GetVestingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyRemainingSellingCoin, auction.GetRemainingSellingCoin().String()),
			sdk.NewAttribute(types.AttributeKeyStartTime, auction.GetStartTime().String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyAuctionStatus, auction.GetStatus().String()),
			sdk.NewAttribute(types.AttributeKeyMinBidPrice, auction.MinBidPrice.String()),
			sdk.NewAttribute(types.AttributeKeyMaxExtendedRound, fmt.Sprint(auction.MaxExtendedRound)),
			sdk.NewAttribute(types.AttributeKeyExtendedRoundRate, auction.ExtendedRoundRate.String()),
		),
	})

	return auction, nil
}

// CancelAuction cancels the auction. It can only be canceled when the auction has not started yet.
func (k Keeper) CancelAuction(ctx sdk.Context, msg *types.MsgCancelAuction) error {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d not found", msg.AuctionId)
	}

	if auction.GetAuctioneer().String() != msg.Auctioneer {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only the autioneer can cancel the auction")
	}

	if auction.GetStatus() != types.AuctionStatusStandBy {
		return sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "only the stand by auction can be cancelled")
	}

	sellingReserveAddr := auction.GetSellingReserveAddress()
	sellingCoinDenom := auction.GetSellingCoin().Denom
	spendableCoins := k.bankKeeper.SpendableCoins(ctx, sellingReserveAddr)
	releaseCoin := sdk.NewCoin(sellingCoinDenom, spendableCoins.AmountOf(sellingCoinDenom))

	// Release the selling coin back to the auctioneer
	if err := k.bankKeeper.SendCoins(ctx, sellingReserveAddr, auction.GetAuctioneer(), sdk.NewCoins(releaseCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to release the selling coin")
	}

	_ = auction.SetRemainingSellingCoin(sdk.NewCoin(sellingCoinDenom, sdk.ZeroInt()))
	_ = auction.SetStatus(types.AuctionStatusCancelled)

	// Call the before auction canceled hook
	k.BeforeAuctionCanceled(ctx, msg.AuctionId, msg.Auctioneer)

	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
		),
	})

	return nil
}

// AddAllowedBidders is a function that is implemented for an external module.
// An external module uses this function to add allowed bidders in the auction's allowed bidders list.
// It doesn't look up the bidder's previous maximum bid amount. Instead, it overlaps.
// It doesn't have any auctioneer's verficiation logic because the module is fundamentally designed
// to delegate full authorization to an external module.
// It is up to an external module to freely add necessary verficiation and operations depending on their use cases.
func (k Keeper) AddAllowedBidders(ctx sdk.Context, auctionId uint64, bidders []types.AllowedBidder) error {
	auction, found := k.GetAuction(ctx, auctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", auctionId)
	}

	if len(bidders) == 0 {
		return types.ErrEmptyAllowedBidders
	}

	if err := types.ValidateAllowedBidders(bidders, auction.GetSellingCoin().Amount); err != nil {
		return err
	}

	// Append new bidders from the existing ones
	allowedBidders := auction.GetAllowedBidders()
	allowedBidders = append(allowedBidders, bidders...)

	if err := auction.SetAllowedBidders(allowedBidders); err != nil {
		return err
	}

	// Call the before allowed bidders added hook
	k.BeforeAllowedBiddersAdded(ctx, auctionId, bidders)

	k.SetAuction(ctx, auction)

	return nil
}

// UpdateAllowedBidder is a function that is implemented for an external module.
// An external module uses this function to update maximum bid amount of particular allowed bidder in the auction.
// It doesn't have any auctioneer's verficiation logic because the module is fundamentally designed
// to delegate full authorization to an external module.
// It is up to an external module to freely add necessary verficiation and operations depending on their use cases.
func (k Keeper) UpdateAllowedBidder(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, maxBidAmount sdk.Int) error {
	auction, found := k.GetAuction(ctx, auctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", auctionId)
	}

	if maxBidAmount.IsNil() {
		return types.ErrInvalidMaxBidAmount
	}

	if !maxBidAmount.IsPositive() {
		return types.ErrInvalidMaxBidAmount
	}

	if _, found := auction.GetAllowedBiddersMap()[bidder.String()]; !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bidder %s is not found", bidder.String())
	}

	_ = auction.SetMaxBidAmount(bidder.String(), maxBidAmount)

	// Call the before allowed bidders updated hook
	k.BeforeAllowedBidderUpdated(ctx, auctionId, bidder, maxBidAmount)

	k.SetAuction(ctx, auction)

	return nil
}

// AllocateSellingCoin allocates allocated selling coin for all matched bids in MatchingInfo and
// releases them from the selling reserve account.
func (k Keeper) AllocateSellingCoin(ctx sdk.Context, auction types.AuctionI, mInfo MatchingInfo) error {
	// Call the before seling coin distributed hook
	k.BeforeSellingCoinsAllocated(ctx, auction.GetId(), mInfo.AllocationMap, mInfo.RefundMap)

	sellingReserveAddr := auction.GetSellingReserveAddress()
	sellingCoinDenom := auction.GetSellingCoin().Denom

	inputs := []banktypes.Input{}
	outputs := []banktypes.Output{}

	// Allocate coins to all matched bidders in AllocationMap and
	// set the amounts in trasnaction inputs and outputs from the selling reserve account
	for bidder, allocAmt := range mInfo.AllocationMap {
		if allocAmt.IsZero() {
			continue
		}
		allocateCoins := sdk.NewCoins(sdk.NewCoin(sellingCoinDenom, allocAmt))
		bidderAddr, _ := sdk.AccAddressFromBech32(bidder)

		inputs = append(inputs, banktypes.NewInput(sellingReserveAddr, allocateCoins))
		outputs = append(outputs, banktypes.NewOutput(bidderAddr, allocateCoins))
	}

	// Send all at once
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	return nil
}

// ReleaseVestingPayingCoin releases the vested selling coin to the auctioneer from the vesting reserve account.
func (k Keeper) ReleaseVestingPayingCoin(ctx sdk.Context, auction types.AuctionI) error {
	vqs := k.GetVestingQueuesByAuctionId(ctx, auction.GetId())
	vqsLen := len(vqs)

	for i, vq := range vqs {
		if vq.ShouldRelease(ctx.BlockTime()) {
			vestingReserveAddr := auction.GetVestingReserveAddress()
			auctioneerAddr := auction.GetAuctioneer()
			payingCoins := sdk.NewCoins(vq.PayingCoin)

			if err := k.bankKeeper.SendCoins(ctx, vestingReserveAddr, auctioneerAddr, payingCoins); err != nil {
				return sdkerrors.Wrap(err, "failed to release paying coin to the auctioneer")
			}

			vq.SetReleased(true)
			k.SetVestingQueue(ctx, vq)

			// Update status to AuctionStatusFinished when all the amounts are released
			if i == vqsLen-1 {
				_ = auction.SetStatus(types.AuctionStatusFinished)
				k.SetAuction(ctx, auction)
			}
		}
	}

	return nil
}

// RefundRemainingSellingCoin refunds the remaining selling coin back to the auctioneer.
// This function is called right after the selling coin is sold.
func (k Keeper) RefundRemainingSellingCoin(ctx sdk.Context, auction types.AuctionI) error {
	sellingReserveAddr := auction.GetSellingReserveAddress()
	sellingCoinDenom := auction.GetSellingCoin().Denom
	spendableCoins := k.bankKeeper.SpendableCoins(ctx, sellingReserveAddr)
	releaseCoins := sdk.NewCoins(sdk.NewCoin(sellingCoinDenom, spendableCoins.AmountOf(sellingCoinDenom)))

	if err := k.bankKeeper.SendCoins(ctx, sellingReserveAddr, auction.GetAuctioneer(), releaseCoins); err != nil {
		return err
	}
	return nil
}

// RefundPayingCoin refunds paying coin back to the bidders.
func (k Keeper) RefundPayingCoin(ctx sdk.Context, auction types.AuctionI, mInfo MatchingInfo) error {
	payingReserveAddr := auction.GetPayingReserveAddress()
	payingCoinDenom := auction.GetPayingCoinDenom()

	inputs := []banktypes.Input{}
	outputs := []banktypes.Output{}

	// Refund the unmatched bid amount back to the bidder
	for bidder, refundAmt := range mInfo.RefundMap {
		if refundAmt.IsZero() {
			continue
		}

		bidderAddr, err := sdk.AccAddressFromBech32(bidder)
		if err != nil {
			return err
		}
		refundCoins := sdk.NewCoins(sdk.NewCoin(payingCoinDenom, refundAmt))

		inputs = append(inputs, banktypes.NewInput(payingReserveAddr, refundCoins))
		outputs = append(outputs, banktypes.NewOutput(bidderAddr, refundCoins))
	}

	// Send all at once
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	return nil
}

// ExtendRound extends another round of ExtendedPeriod value for the auction.
func (k Keeper) ExtendRound(ctx sdk.Context, ba *types.BatchAuction) {
	params := k.GetParams(ctx)
	extendedPeriod := ctx.BlockTime().AddDate(0, 0, int(params.ExtendedPeriod))

	endTimes := ba.GetEndTimes()
	endTimes = append(endTimes, extendedPeriod)

	_ = ba.SetEndTimes(endTimes)
	k.SetAuction(ctx, ba)
}

// FinishFixedPriceAuction finishes a fixed price auction.
func (k Keeper) FinishFixedPriceAuction(ctx sdk.Context, auction types.AuctionI) {
	mInfo := k.CalculateFixedPriceAllocation(ctx, auction)

	if err := k.AllocateSellingCoin(ctx, auction, mInfo); err != nil {
		panic(err)
	}

	if err := k.RefundRemainingSellingCoin(ctx, auction); err != nil {
		panic(err)
	}

	if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
		panic(err)
	}
}

// FinishBatchAuction finishes a batch auction.
func (k Keeper) FinishBatchAuction(ctx sdk.Context, auction types.AuctionI) {
	ba := auction.(*types.BatchAuction)

	if ba.MaxExtendedRound+1 == uint32(len(auction.GetEndTimes())) {
		mInfo := k.CalculateBatchAllocation(ctx, auction)

		if err := k.AllocateSellingCoin(ctx, auction, mInfo); err != nil {
			panic(err)
		}

		if err := k.RefundRemainingSellingCoin(ctx, auction); err != nil {
			panic(err)
		}

		if err := k.RefundPayingCoin(ctx, auction, mInfo); err != nil {
			panic(err)
		}

		if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
			panic(err)
		}
	}

	// Extend round since there is no last matched length to compare with
	lastMatchedLen := k.GetLastMatchedBidsLen(ctx, ba.GetId())
	if lastMatchedLen == 0 {
		k.CalculateBatchAllocation(ctx, auction)
		k.ExtendRound(ctx, ba)
		return
	}

	mInfo := k.CalculateBatchAllocation(ctx, auction)

	currDec := sdk.NewDec(mInfo.MatchedLen)
	lastDec := sdk.NewDec(lastMatchedLen)
	diff := sdk.OneDec().Sub(currDec.Quo(lastDec)) // 1 - (CurrentMatchedLenDec / LastMatchedLenDec)

	// To prevent from auction sniping technique, compare the extended round rate with
	// the current and the last length of matched bids to determine
	// if the auction needs another extended round
	if diff.GTE(ba.ExtendedRoundRate) {
		k.ExtendRound(ctx, ba)
		return
	}

	if err := k.AllocateSellingCoin(ctx, auction, mInfo); err != nil {
		panic(err)
	}

	if err := k.RefundRemainingSellingCoin(ctx, auction); err != nil {
		panic(err)
	}

	if err := k.RefundPayingCoin(ctx, auction, mInfo); err != nil {
		panic(err)
	}

	if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
		panic(err)
	}
}
