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

func (k Keeper) ReleaseRemainingSellingCoin(ctx sdk.Context, auction types.AuctionI) error {
	sellingReserveAddr := auction.GetSellingReserveAddress()
	sellingCoinDenom := auction.GetSellingCoin().Denom

	// Send all the remaining selling coins back to the auctioneer
	spendableCoins := k.bankKeeper.SpendableCoins(ctx, sellingReserveAddr)
	releaseAmt := spendableCoins.AmountOf(sellingCoinDenom)
	releaseCoins := sdk.NewCoins(sdk.NewCoin(sellingCoinDenom, releaseAmt))

	if err := k.bankKeeper.SendCoins(ctx, sellingReserveAddr, auction.GetAuctioneer(), releaseCoins); err != nil {
		return err
	}
	return nil
}

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

// AllocateSellingCoin releases designated selling coin from the selling reserve account.
func (k Keeper) AllocateSellingCoin(ctx sdk.Context, auction types.AuctionI, mInfo MatchingInfo) error {
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

// AllocateVestingPayingCoin releases the selling coin from the vesting reserve account.
func (k Keeper) AllocateVestingPayingCoin(ctx sdk.Context, auction types.AuctionI) error {
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

// CancelAuction cancels the auction in an event when the auctioneer needs to modify the auction.
// However, it can only be canceled when the auction has not started yet.
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

	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
		),
	})

	return nil
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

	if err := types.ValidateAllowedBidders(bidders, auction.GetSellingCoin().Amount); err != nil {
		return err
	}

	// Append new bidders from the existing ones
	allowedBidders := auction.GetAllowedBidders()
	allowedBidders = append(allowedBidders, bidders...)

	_ = auction.SetAllowedBidders(allowedBidders)

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

	if err := k.ReleaseRemainingSellingCoin(ctx, auction); err != nil {
		panic(err)
	}

	if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
		panic(err)
	}
}

func (k Keeper) FinishBatchAuction(ctx sdk.Context, auction types.AuctionI) {
	ba := auction.(*types.BatchAuction)

	if ba.MaxExtendedRound+1 == uint32(len(auction.GetEndTimes())) {
		mInfo := k.CalculateBatchAllocation(ctx, auction)

		if err := k.AllocateSellingCoin(ctx, auction, mInfo); err != nil {
			panic(err)
		}

		if err := k.ReleaseRemainingSellingCoin(ctx, auction); err != nil {
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

	if err := k.ReleaseRemainingSellingCoin(ctx, auction); err != nil {
		panic(err)
	}

	if err := k.RefundPayingCoin(ctx, auction, mInfo); err != nil {
		panic(err)
	}

	if err := k.ApplyVestingSchedules(ctx, auction); err != nil {
		panic(err)
	}
}

// MatchingInfo holds information about a batch auction matching info.
type MatchingInfo struct {
	MatchedLen         int64              // the length of matched bids
	MatchedPrice       sdk.Dec            // the final matched price
	TotalMatchedAmount sdk.Int            // the total sold amount
	AllocationMap      map[string]sdk.Int // the allocation map that holds allocate amount information for each bidder
	ReservedMatchedMap map[string]sdk.Int
	RefundMap          map[string]sdk.Int // the refund map that holds refund amount information for each bidder
}

func (k Keeper) CalculateFixedPriceAllocation(ctx sdk.Context, auction types.AuctionI) MatchingInfo {
	mInfo := MatchingInfo{
		MatchedPrice:       sdk.ZeroDec(),
		TotalMatchedAmount: sdk.ZeroInt(),
		AllocationMap:      map[string]sdk.Int{},
	}

	totalMatchedAmt := sdk.ZeroInt()
	allocMap := map[string]sdk.Int{}

	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		bidAmt := b.GetBidSellingAmount(auction.GetPayingCoinDenom())

		// Accumulate bid amount if the bidder has other bid(s)
		if allocatedAmt, ok := allocMap[b.Bidder]; ok {
			allocMap[b.Bidder] = allocatedAmt.Add(bidAmt)
		} else {
			allocMap[b.Bidder] = bidAmt
		}
		totalMatchedAmt = totalMatchedAmt.Add(bidAmt)
		mInfo.MatchedLen = mInfo.MatchedLen + 1
	}

	mInfo.MatchedPrice = auction.GetStartPrice()
	mInfo.TotalMatchedAmount = totalMatchedAmt
	mInfo.AllocationMap = allocMap

	return mInfo
}

func (k Keeper) CalculateBatchAllocation(ctx sdk.Context, auction types.AuctionI) MatchingInfo {
	bids := k.GetBidsByAuctionId(ctx, auction.GetId())
	bids = types.SortByBidPrice(bids)

	mInfo := MatchingInfo{
		MatchedLen:         0,
		MatchedPrice:       sdk.ZeroDec(),
		TotalMatchedAmount: sdk.ZeroInt(),
		AllocationMap:      map[string]sdk.Int{},
		ReservedMatchedMap: map[string]sdk.Int{},
		RefundMap:          map[string]sdk.Int{},
	}
	allowedBiddersMap := auction.GetAllowedBiddersMap() // map(bidder => maxBidAmt)
	allocationMap := map[string]sdk.Int{}               // map(bidder => allocatedAmt)
	reservedMatchedMap := map[string]sdk.Int{}          // map(bidder => reservedMatchedAmt)
	reservedMap := map[string]sdk.Int{}                 // map(bidder => reservedAmt)
	refundMap := map[string]sdk.Int{}                   // map(bidder => refundAmt)

	// Initialize value for all maps
	for _, ab := range auction.GetAllowedBidders() {
		mInfo.AllocationMap[ab.Bidder] = sdk.ZeroInt()
		mInfo.ReservedMatchedMap[ab.Bidder] = sdk.ZeroInt()
		refundMap[ab.Bidder] = sdk.ZeroInt()
		reservedMap[ab.Bidder] = sdk.ZeroInt()
	}

	// Iterate from the highest matching bid price and stop until it finds
	// the matching information to store them into MatchingInfo object
	for _, bid := range bids {
		matchingPrice := bid.Price
		totalMatchedAmt := sdk.ZeroInt()

		// Add all allowed bidders for the matching price
		for _, ab := range auction.GetAllowedBidders() {
			allocationMap[ab.Bidder] = sdk.ZeroInt()
			reservedMatchedMap[ab.Bidder] = sdk.ZeroInt()
		}

		// Iterate all bids and execute the logics when the bid price is
		// higher than the current matching price
		for _, b := range bids {
			if b.Price.LT(matchingPrice) {
				continue
			}

			maxBidAmt := allowedBiddersMap[b.Bidder]
			allocateAmt := allocationMap[b.Bidder]

			// Uses minimum of the two amounts to prevent from exceeding the bidder's maximum bid amount
			if b.Type == types.BidTypeBatchWorth {
				bidAmt := b.Coin.Amount.ToDec().QuoTruncate(matchingPrice).TruncateInt()

				// MinInt(BidAmt, MaxBidAmt-AccumulatedBidAmt)
				matchingAmt := sdk.MinInt(bidAmt, maxBidAmt.Sub(allocateAmt))

				// Accumulate matching amount since a bidder can have multiple bids
				if alloc, ok := allocationMap[b.Bidder]; ok {
					allocationMap[b.Bidder] = alloc.Add(matchingAmt)
				}

				// Accumulate how much reserved paying coin amount is matched
				if reservedMatchedAmt, ok := reservedMatchedMap[b.Bidder]; ok {
					reserveAmt := sdk.ZeroInt()
					if matchingAmt.LT(bidAmt) {
						reserveAmt = matchingAmt.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
					} else {
						reserveAmt = b.Coin.Amount
					}
					reservedMatchedMap[b.Bidder] = reservedMatchedAmt.Add(reserveAmt)
				}

				totalMatchedAmt = totalMatchedAmt.Add(matchingAmt)
			} else {
				bidAmt := b.Coin.Amount

				// MinInt(BidAmt, MaxBidAmount-AccumulatedBidAmount)
				matchingAmt := sdk.MinInt(bidAmt, maxBidAmt.Sub(allocateAmt))

				// Accumulate matching amount since a bidder can have multiple bids
				if alloc, ok := allocationMap[b.Bidder]; ok {
					allocationMap[b.Bidder] = alloc.Add(matchingAmt)
				}

				// Accumulate how much reserved paying coin amount is matched
				if reservedMatchedAmt, ok := reservedMatchedMap[b.Bidder]; ok {
					reserveAmt := matchingAmt.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
					reservedMatchedMap[b.Bidder] = reservedMatchedAmt.Add(reserveAmt)
				}

				totalMatchedAmt = totalMatchedAmt.Add(matchingAmt)
			}
		}

		// Exit the iteration when the total matched amount is greater than the total selling coin amount
		if totalMatchedAmt.GT(auction.GetSellingCoin().Amount) {
			break
		}

		mInfo.MatchedLen = mInfo.MatchedLen + 1
		mInfo.MatchedPrice = matchingPrice
		mInfo.TotalMatchedAmount = totalMatchedAmt

		for _, ab := range auction.GetAllowedBidders() {
			mInfo.AllocationMap[ab.Bidder] = allocationMap[ab.Bidder]
			mInfo.ReservedMatchedMap[ab.Bidder] = reservedMatchedMap[ab.Bidder]
		}

		bid.SetMatched(true)
		k.SetBid(ctx, bid)
	}

	// Iterate all bids to get refund amount for each bidder
	// Calculate the refund amount by substracting allocate amount from
	// how much a bidder reserved to place a bid for the auction
	for _, b := range bids {
		if b.Type == types.BidTypeBatchWorth {
			reservedMap[b.Bidder] = reservedMap[b.Bidder].Add(b.Coin.Amount)
		} else {
			bidAmt := b.Coin.Amount.ToDec().Mul(b.Price).Ceil().TruncateInt()
			reservedMap[b.Bidder] = reservedMap[b.Bidder].Add(bidAmt)
		}
	}

	for bidder, reservedAmt := range reservedMap {
		reservedMatchedAmt, ok := mInfo.ReservedMatchedMap[bidder]
		if ok {
			refundMap[bidder] = reservedAmt.Sub(reservedMatchedAmt)
			continue
		}
		refundMap[bidder] = reservedAmt
	}

	mInfo.RefundMap = refundMap

	k.SetMatchedBidsLen(ctx, auction.GetId(), mInfo.MatchedLen)

	return mInfo
}
