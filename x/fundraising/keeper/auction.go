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

// AllocateSellingCoin releases designated selling coin from the selling reserve account.
func (k Keeper) AllocateSellingCoin(ctx sdk.Context, auction types.AuctionI, mInfo MatchingInfo) error {
	sellingReserveAddr := auction.GetSellingReserveAddress()
	sellingCoinDenom := auction.GetSellingCoin().Denom

	inputs := []banktypes.Input{}
	outputs := []banktypes.Output{}

	// Loop through all bids and set the allocated coins in transaction inputs and outputs
	// from the selling reserve account
	for _, alloc := range mInfo.Allocations {
		switch auction.GetType() {
		case types.AuctionTypeFixedPrice:
			allocateCoins := sdk.NewCoins(sdk.NewCoin(sellingCoinDenom, alloc.AllocateAmount))
			bidderAddr, _ := sdk.AccAddressFromBech32(alloc.Bidder)

			inputs = append(inputs, banktypes.NewInput(sellingReserveAddr, allocateCoins))
			outputs = append(outputs, banktypes.NewOutput(bidderAddr, allocateCoins))

		case types.AuctionTypeBatch:
			allocateCoins := sdk.NewCoins(sdk.NewCoin(sellingCoinDenom, alloc.AllocateAmount))
			bidderAddr, _ := sdk.AccAddressFromBech32(alloc.Bidder)

			inputs = append(inputs, banktypes.NewInput(sellingReserveAddr, allocateCoins))
			outputs = append(outputs, banktypes.NewOutput(bidderAddr, allocateCoins))

			// Refund the unmatched bid amount back to the bidder
			if alloc.ReserveAmount.GT(alloc.AllocateAmount) {
				refundAmt := alloc.ReserveAmount.Sub(alloc.AllocateAmount)
				refundCoins := sdk.NewCoins(sdk.NewCoin(sellingCoinDenom, refundAmt))

				inputs = append(inputs, banktypes.NewInput(sellingReserveAddr, refundCoins))
				outputs = append(outputs, banktypes.NewOutput(bidderAddr, refundCoins))
			}
		}
	}

	// Refund all the remaining coins back to the auctioneer
	if auction.GetType() == types.AuctionTypeFixedPrice {
		reserveCoin := k.bankKeeper.GetBalance(ctx, sellingReserveAddr, sellingCoinDenom)
		totalAmt := sdk.NewCoin(sellingCoinDenom, mInfo.TotalMatchedAmount)
		refundAmt := reserveCoin.Sub(totalAmt)

		inputs = append(inputs, banktypes.NewInput(sellingReserveAddr, sdk.NewCoins(refundAmt)))
		outputs = append(outputs, banktypes.NewOutput(auction.GetAuctioneer(), sdk.NewCoins(refundAmt)))
	}

	// Send all at once
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	return nil
}

// AllocatePayingCoin releases the selling coin from the vesting reserve account.
func (k Keeper) AllocatePayingCoin(ctx sdk.Context, auction types.AuctionI) error {
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
func (k Keeper) CreateFixedPriceAuction(ctx sdk.Context, msg *types.MsgCreateFixedPriceAuction) (types.AuctionI, error) {
	if ctx.BlockTime().After(msg.EndTime) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be set prior to the current time")
	}

	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	if err := k.ReserveCreationFee(ctx, msg.GetAuctioneer()); err != nil {
		return nil, err
	}

	if err := k.ReserveSellingCoin(ctx, nextId, msg.GetAuctioneer(), msg.SellingCoin); err != nil {
		return nil, err
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
func (k Keeper) CreateBatchAuction(ctx sdk.Context, msg *types.MsgCreateBatchAuction) (types.AuctionI, error) {
	if ctx.BlockTime().After(msg.EndTime) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be set prior to the current time")
	}

	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	if err := k.ReserveCreationFee(ctx, msg.GetAuctioneer()); err != nil {
		return nil, err
	}

	if err := k.ReserveSellingCoin(ctx, nextId, msg.GetAuctioneer(), msg.SellingCoin); err != nil {
		return nil, err
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
	params := k.GetParams(ctx)
	extendedPeriod := ctx.BlockTime().AddDate(0, 0, int(params.ExtendedPeriod))

	endTimes := ba.GetEndTimes()
	endTimes = append(endTimes, extendedPeriod)

	_ = ba.SetEndTimes(endTimes)
	k.SetAuction(ctx, ba)
}

func (k Keeper) FinishFixedPriceAuction(ctx sdk.Context, auction types.AuctionI) {
	mInfo := k.CalculateFixedPriceAllocation(ctx, auction)

	if err := k.AllocateSellingCoin(ctx, auction, mInfo); err != nil {
		panic(err)
	}

	if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
		panic(err)
	}
}

func (k Keeper) FinishBatchAuction(ctx sdk.Context, auction types.AuctionI) {
	mInfo := k.CalculateBatchAllocation(ctx, auction)

	ba := auction.(*types.BatchAuction)
	if ba.MaxExtendedRound == 0 {
		if err := k.AllocateSellingCoin(ctx, auction, mInfo); err != nil {
			panic(err)
		}

		if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
			panic(err)
		}
		return
	}

	currMatchedLen := mInfo.MatchedLen
	lastMatchedLen := k.GetMatchedBidsLen(ctx, ba.Id)

	// Extend round since there is no last matched length to compare with
	if lastMatchedLen == 0 {
		k.ExtendRound(ctx, ba)
		return
	}

	// To prevent from auction sniping technique, compare the extended round rate with
	// the current and the last length of matched bids to determine
	// if the auction needs another extended round
	currDec := sdk.NewDec(currMatchedLen)
	lastDec := sdk.NewDec(lastMatchedLen)
	diff := sdk.OneDec().Sub(currDec.Quo(lastDec)) // 1 - (CurrentMatchedLenDec / LastMatchedLenDec)

	if diff.GTE(ba.ExtendedRoundRate) {
		k.ExtendRound(ctx, ba)
		return
	}

	if err := k.AllocateSellingCoin(ctx, auction, mInfo); err != nil {
		panic(err)
	}

	if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
		panic(err)
	}
}

// MatchingInfo holds information about a batch auction matching info.
type MatchingInfo struct {
	MatchedLen         int64   // the length of matched bids
	MatchedPrice       sdk.Dec // the final matched price
	TotalMatchedAmount sdk.Int // the total sold amount
	Allocations        []AllocationInfo
}

// AllocationInfo holds information about a bidder's allocation.
type AllocationInfo struct {
	Bidder         string
	AllocateAmount sdk.Int
	ReserveAmount  sdk.Int
}

func (k Keeper) CalculateFixedPriceAllocation(ctx sdk.Context, auction types.AuctionI) MatchingInfo {
	mInfo := MatchingInfo{
		MatchedLen:         0,
		MatchedPrice:       sdk.ZeroDec(),
		TotalMatchedAmount: sdk.ZeroInt(),
		Allocations:        []AllocationInfo{},
	}

	totalMatchedAmt := sdk.ZeroInt()
	allocations := []AllocationInfo{}

	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		bidAmt := b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()

		allocations = append(allocations, AllocationInfo{
			Bidder:         b.Bidder,
			AllocateAmount: bidAmt,
		})

		totalMatchedAmt = totalMatchedAmt.Add(bidAmt)
		mInfo.MatchedLen = mInfo.MatchedLen + 1
	}

	mInfo.MatchedPrice = auction.GetStartPrice()
	mInfo.TotalMatchedAmount = totalMatchedAmt
	mInfo.Allocations = allocations

	return mInfo
}

func (k Keeper) CalculateBatchAllocation(ctx sdk.Context, auction types.AuctionI) MatchingInfo {
	bids := k.GetBidsByAuctionId(ctx, auction.GetId())
	bids = types.SortByBidPrice(bids)

	allowedBidders := auction.GetAllowedBidders()
	if len(allowedBidders) == 0 {
		panic("This can't occur since the allowed bidders must exist")
	}

	mInfo := MatchingInfo{
		MatchedLen:         0,
		MatchedPrice:       sdk.ZeroDec(),
		TotalMatchedAmount: sdk.ZeroInt(),
		Allocations:        []AllocationInfo{},
	}
	allowedBiddersMap := auction.GetAllowedBiddersMap() // map(bidder => maxBidAmt)
	accumulatedMap := make(map[string]sdk.Int)          // map(bidder => accumulatedAmt)
	reservedMap := make(map[string]sdk.Int)             // map(bidder => reservedAmt)

	// Iterate from the highest matching bid price and find the matching information
	// to store them in MatchingInfo object
	for _, b := range bids {
		matchingPrice := b.Price
		totalMatchedAmt := sdk.ZeroInt()

		// Add all allowed bidders to the maps for the matching price
		// Accumulated and reserved amounts must be initialized per matching price
		for _, ab := range allowedBidders {
			accumulatedMap[ab.Bidder] = sdk.ZeroInt()
			reservedMap[ab.Bidder] = sdk.ZeroInt()
		}

		for _, b := range bids {
			if b.Price.LT(matchingPrice) {
				continue
			}

			// Uses minimum of the two amounts to prevent from exceeding the bidder's maximum bid amount
			if b.Type == types.BidTypeBatchWorth {
				maxBidAmt := allowedBiddersMap[b.Bidder]
				accumulatedAmt := accumulatedMap[b.Bidder]
				bidAmt := b.Coin.Amount.ToDec().QuoTruncate(matchingPrice).TruncateInt()

				// MinInt(bidAmt, MaxBidAmt-AccumulatedBidAmt)
				matchingAmt := sdk.MinInt(bidAmt, maxBidAmt.Sub(accumulatedAmt))

				reservedMap[b.Bidder] = reservedMap[b.Bidder].Add(b.Coin.Amount)
				accumulatedMap[b.Bidder] = accumulatedMap[b.Bidder].Add(matchingAmt)
				totalMatchedAmt = totalMatchedAmt.Add(matchingAmt)
			} else {
				maxBidAmt := allowedBiddersMap[b.Bidder]
				accumulatedAmt := accumulatedMap[b.Bidder]
				bidAmt := b.Coin.Amount

				// MinInt(BidAmt, MaxBidAmount-AccumulatedBidAmount)
				matchingAmt := sdk.MinInt(bidAmt, maxBidAmt.Sub(accumulatedAmt))

				bidAmt = b.Coin.Amount.ToDec().Mul(b.Price).Ceil().TruncateInt()
				reservedMap[b.Bidder] = reservedMap[b.Bidder].Add(bidAmt)
				accumulatedMap[b.Bidder] = accumulatedMap[b.Bidder].Add(matchingAmt)
				totalMatchedAmt = totalMatchedAmt.Add(matchingAmt)
			}
		}

		if totalMatchedAmt.GT(auction.GetSellingCoin().Amount) {
			break
		}

		b.SetWinner(true)
		k.SetBid(ctx, b)

		mInfo.MatchedLen = mInfo.MatchedLen + 1
		mInfo.MatchedPrice = matchingPrice
		mInfo.TotalMatchedAmount = totalMatchedAmt
	}

	// Store allocation info from the maps
	allocsInfo := []AllocationInfo{}
	for bidder, accumulatedAmt := range accumulatedMap {
		allocsInfo = append(allocsInfo, AllocationInfo{
			Bidder:         bidder,
			AllocateAmount: accumulatedAmt,
			ReserveAmount:  reservedMap[bidder],
		})
	}
	mInfo.Allocations = allocsInfo

	// 100  000000
	// 3250 000000
	// 1000 000000
	// 2500 000000
	// 100  000000
	// Total: 3211 111111

	for _, alloc := range mInfo.Allocations {
		fmt.Println("MatchedPrice: ", mInfo.MatchedPrice)
		fmt.Println("Bidder: ", alloc.Bidder)
		fmt.Println("AllocateAmt: ", alloc.AllocateAmount)
		fmt.Println("ReservedAmt: ", alloc.ReserveAmount)
		fmt.Println("")
	}

	k.SetMatchedBidsLen(ctx, auction.GetId(), mInfo.MatchedLen)

	return mInfo
}
