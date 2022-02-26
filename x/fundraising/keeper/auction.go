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
func (k Keeper) DistributeSellingCoin(ctx sdk.Context, auction types.AuctionI, bids []types.Bid) error {
	sellingReserveAddr := auction.GetSellingReserveAddress()
	sellingCoinDenom := auction.GetSellingCoin().Denom

	totalBidCoin := sdk.NewCoin(sellingCoinDenom, sdk.ZeroInt())
	inputs := []banktypes.Input{}
	outputs := []banktypes.Output{}

	// Loop through all bids and set the allocated coins to all bidders from the selling reserve account
	for _, b := range bids {
		switch b.Type {
		case types.BidTypeFixedPrice:
			exchangedSellingAmt := b.GetExchangedSellingAmount()
			exchangedSellingCoin := sdk.NewCoin(sellingCoinDenom, exchangedSellingAmt)

			inputs = append(inputs, banktypes.NewInput(sellingReserveAddr, sdk.NewCoins(exchangedSellingCoin)))
			outputs = append(outputs, banktypes.NewOutput(b.GetBidder(), sdk.NewCoins(exchangedSellingCoin)))

			totalBidCoin = totalBidCoin.Add(exchangedSellingCoin)

		case types.BidTypeBatchWorth:

		case types.BidTypeBatchMany:

		}
	}

	reserveCoin := k.bankKeeper.GetBalance(ctx, sellingReserveAddr, sellingCoinDenom)
	remainingCoin := reserveCoin.Sub(totalBidCoin)

	// Send remaining coin to the auctioneer
	inputs = append(inputs, banktypes.NewInput(sellingReserveAddr, sdk.NewCoins(remainingCoin)))
	outputs = append(outputs, banktypes.NewOutput(auction.GetAuctioneer(), sdk.NewCoins(remainingCoin)))

	// Send all at once
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	return nil
}

// DistributePayingCoin releases the selling coin from the vesting reserve account.
func (k Keeper) DistributePayingCoin(ctx sdk.Context, auction types.AuctionI) error {
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

// CreateFixedPriceAuction sets fixed price auction.
func (k Keeper) CreateFixedPriceAuction(ctx sdk.Context, msg *types.MsgCreateFixedPriceAuction) (*types.FixedPriceAuction, error) {
	if ctx.BlockTime().After(msg.EndTime) {
		return &types.FixedPriceAuction{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be set prior to the current time")
	}

	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	if err := k.ReserveCreationFee(ctx, msg.GetAuctioneer()); err != nil {
		return &types.FixedPriceAuction{}, err
	}

	if err := k.ReserveSellingCoin(ctx, nextId, msg.GetAuctioneer(), msg.SellingCoin); err != nil {
		return &types.FixedPriceAuction{}, err
	}

	allowedBidders := []types.AllowedBidder{} // it is nil when an auction is created
	winningPrice := msg.StartPrice            // it is start price
	numWinningBidders := uint64(0)            // initial value is 0
	remainingSellingCoin := msg.SellingCoin   // it is starting with selling coin amount
	endTimes := []time.Time{msg.EndTime}      // it is an array data type to handle BatchAuction

	baseAuction := types.NewBaseAuction(
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
		winningPrice,
		numWinningBidders,
		remainingSellingCoin,
		msg.StartTime,
		endTimes,
		types.AuctionStatusStandBy,
	)

	// Update status if the start time is already passed over the current time
	if baseAuction.ShouldAuctionStarted(ctx.BlockTime()) {
		baseAuction.Status = types.AuctionStatusStarted
	}

	auction := types.NewFixedPriceAuction(baseAuction)
	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedPriceAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyAuctioneerAddress, msg.Auctioneer),
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

// CreateBatchAuction sets batch auction.
func (k Keeper) CreateBatchAuction(ctx sdk.Context, msg *types.MsgCreateBatchAuction) (*types.BatchAuction, error) {
	if ctx.BlockTime().After(msg.EndTime) {
		return &types.BatchAuction{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be set prior to the current time")
	}

	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	if err := k.ReserveCreationFee(ctx, msg.GetAuctioneer()); err != nil {
		return &types.BatchAuction{}, err
	}

	if err := k.ReserveSellingCoin(ctx, nextId, msg.GetAuctioneer(), msg.SellingCoin); err != nil {
		return &types.BatchAuction{}, err
	}

	allowedBidders := []types.AllowedBidder{} // it is nil when an auction is created
	winningPrice := sdk.ZeroDec()             // TODO: makes sense to have start price?
	numWinningBidders := uint64(0)            // initial value is 0
	remainingSellingCoin := msg.SellingCoin   // it is starting with selling coin amount
	endTimes := []time.Time{msg.EndTime}      // it is an array data type to handle BatchAuction

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
		winningPrice,
		numWinningBidders,
		remainingSellingCoin,
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
		msg.MaxExtendedRound,
		msg.ExtendedRoundRate,
	)
	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedPriceAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyAuctioneerAddress, msg.Auctioneer),
			sdk.NewAttribute(types.AttributeKeyStartPrice, auction.GetStartPrice().String()),
			sdk.NewAttribute(types.AttributeKeySellingReserveAddress, auction.GetSellingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyPayingReserveAddress, auction.GetPayingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeyVestingReserveAddress, auction.GetVestingReserveAddress().String()),
			sdk.NewAttribute(types.AttributeKeySellingCoin, auction.GetSellingCoin().String()),
			sdk.NewAttribute(types.AttributeKeyPayingCoinDenom, auction.GetPayingCoinDenom()),
			sdk.NewAttribute(types.AttributeKeyStartTime, auction.GetStartTime().String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyAuctionStatus, auction.GetStatus().String()),
			sdk.NewAttribute(types.AttributeKeyMaxExtendedRound, fmt.Sprint(msg.MaxExtendedRound)),
			sdk.NewAttribute(types.AttributeKeyExtendedRoundRate, msg.ExtendedRoundRate.String()),
		),
	})

	return auction, nil
}

// CancelAuction cancels the auction in an event when the auctioneer needs to modify the auction.
// However, it can only be canceled when the auction has not started yet.
func (k Keeper) CancelAuction(ctx sdk.Context, msg *types.MsgCancelAuction) (types.AuctionI, error) {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "auction not found")
	}

	if auction.GetAuctioneer().String() != msg.Auctioneer {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify ownership of the auction")
	}

	if auction.GetStatus() != types.AuctionStatusStandBy {
		return nil, sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "auction cannot be canceled due to current status")
	}

	sellingReserveAddr := auction.GetSellingReserveAddress()
	auctioneerAddr := auction.GetAuctioneer()
	releaseCoin := k.bankKeeper.GetBalance(ctx, sellingReserveAddr, auction.GetSellingCoin().Denom)

	// Release the selling coin back to the auctioneer
	if err := k.bankKeeper.SendCoins(ctx, sellingReserveAddr, auctioneerAddr, sdk.NewCoins(releaseCoin)); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to release the selling coin")
	}

	_ = auction.SetRemainingSellingCoin(sdk.NewCoin(auction.GetSellingCoin().Denom, sdk.ZeroInt()))
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

	if err := types.ValidateAllowedBidders(bidders); err != nil {
		return err
	}

	// Append new bidders from the existing ones
	allowedBidders := auction.GetAllowedBidders()
	allowedBidders = append(allowedBidders, bidders...)

	if err := auction.SetAllowedBidders(allowedBidders); err != nil {
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

	k.SetAuction(ctx, auction)

	return nil
}

func (k Keeper) ExtendRound(ctx sdk.Context, ba *types.BatchAuction) {
	endTimes := ba.GetEndTimes()
	endTimes = append(endTimes, ctx.BlockTime())

	_ = ba.SetEndTimes(endTimes)
	k.SetAuction(ctx, ba)
}

// AllocationInfo holds information about a batch auction allocation.
type AllocationInfo struct {
	MatchedBids  []types.Bid // the matched bids
	MatchedPrice sdk.Dec     // the final matched price
	SoldAmount   sdk.Int     // the total sold amount
}

func (k Keeper) CalculateAllocation(ctx sdk.Context, auction types.AuctionI) AllocationInfo {
	bids := k.GetBidsByAuctionId(ctx, auction.GetId())
	bids = types.SortByBidPrice(bids)

	allocInfo := AllocationInfo{
		MatchedBids:  []types.Bid{},
		MatchedPrice: bids[0].Price,
		SoldAmount:   sdk.ZeroInt(),
	}

	allowedBidders := auction.GetAllowedBidders()
	allowedBiddersMap := auction.GetAllowedBiddersMap() // map(bidder => maxBidAmt)
	accumulatedMap := make(map[string]sdk.Int)          // map(bidder => accumulatedAmt)
	totalSellingAmt := auction.GetSellingCoin().Amount

	// Iterate from the highest bid price and find the last matching bid price and total sold amount
	// It doesn't concern about a partial amount of coins for the bid after the last matching bid
	for _, b := range bids {
		matchingPrice := b.Price
		accumulatedAmt := sdk.ZeroInt()

		// Add all allowed bidders to accumulatedMap for every matching price
		for _, ab := range allowedBidders {
			accumulatedMap[ab.Bidder] = sdk.ZeroInt()
		}

		for _, b := range bids {
			if b.Price.LT(matchingPrice) {
				continue
			}

			if b.Type == types.BidTypeBatchWorth {
				// Min(ExchangedSellingAmt, MaxBidAmt-AccumulatedBidAmt)
				matchingAmt := b.Coin.Amount.ToDec().QuoTruncate(matchingPrice).TruncateInt()
				matchingAmt = sdk.MinInt(matchingAmt, allowedBiddersMap[b.Bidder].Sub(accumulatedMap[b.Bidder]))

				accumulatedAmt = accumulatedAmt.Add(matchingAmt)
				accumulatedMap[b.Bidder] = accumulatedMap[b.Bidder].Add(matchingAmt)
			} else {
				// Min(SellingAmt, MaxBidAmount-AccumulatedBidAmount)
				diffAmt := allowedBiddersMap[b.Bidder].Sub(accumulatedMap[b.Bidder])
				matchingAmt := sdk.MinInt(b.Coin.Amount, diffAmt)

				accumulatedAmt = accumulatedAmt.Add(matchingAmt)
				accumulatedMap[b.Bidder] = accumulatedMap[b.Bidder].Add(matchingAmt)
			}
		}

		if accumulatedAmt.GT(totalSellingAmt) {
			break
		}

		b.SetWinner(true)
		k.SetBid(ctx, b)

		allocInfo.MatchedBids = append(allocInfo.MatchedBids, b)
		allocInfo.MatchedPrice = matchingPrice
		allocInfo.SoldAmount = accumulatedAmt
	}

	k.SetMatchedBidsLen(ctx, auction.GetId(), len(allocInfo.MatchedBids))

	return allocInfo
}

func (k Keeper) FinishFixedPriceAuction(ctx sdk.Context, auction types.AuctionI) {
	bids := k.GetBidsByAuctionId(ctx, auction.GetId())

	if err := k.DistributeSellingCoin(ctx, auction, bids); err != nil {
		panic(err)
	}

	if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
		panic(err)
	}
}

func (k Keeper) FinishBatchAuction(ctx sdk.Context, auction types.AuctionI) {
	allocInfo := k.CalculateAllocation(ctx, auction)

	ba := auction.(*types.BatchAuction)
	if ba.MaxExtendedRound == 0 {
		if err := k.DistributeSellingCoin(ctx, auction, allocInfo.MatchedBids); err != nil {
			panic(err)
		}

		if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
			panic(err)
		}

	} else {
		// Compare with the last matched bids length and
		// determine if it needs another round
		currMatchedLen := int64(len(allocInfo.MatchedBids))
		lastMatchedLen := k.GetMatchedBidsLen(ctx, ba.Id)
		currDec := sdk.NewDec(currMatchedLen)
		lastDec := sdk.NewDec(lastMatchedLen)

		// 1 - (currentMatchedLenDec / lastMatchedLenDec)
		diff := sdk.OneDec().Sub(currDec.Quo(lastDec))

		// Extend another round if the diff is greater than or equal to the extended round rate
		if diff.GTE(ba.ExtendedRoundRate) {
			k.ExtendRound(ctx, ba)
		} else {
			if err := k.DistributeSellingCoin(ctx, auction, allocInfo.MatchedBids); err != nil {
				panic(err)
			}

			if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
				panic(err)
			}
		}
	}
}
