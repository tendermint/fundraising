package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetNextBidId increments bid id by one and set it.
func (k Keeper) GetNextBidIdWithUpdate(ctx sdk.Context, auctionId uint64) uint64 {
	id := k.GetLastBidId(ctx, auctionId) + 1
	k.SetBidId(ctx, auctionId, id)
	return id
}

// ReservePayingCoin reserves paying coin to the paying reserve account.
func (k Keeper) ReservePayingCoin(ctx sdk.Context, auctionId uint64, bidderAddr sdk.AccAddress, payingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, bidderAddr, types.PayingReserveAddress(auctionId), sdk.NewCoins(payingCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to reserve paying coin")
	}
	return nil
}

// PlaceBid places a bid for the auction.
func (k Keeper) PlaceBid(ctx sdk.Context, msg *types.MsgPlaceBid) (types.Bid, error) {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return types.Bid{}, sdkerrors.Wrap(sdkerrors.ErrNotFound, "auction not found")
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return types.Bid{}, types.ErrInvalidAuctionStatus
	}

	_, found = auction.GetAllowedBiddersMap()[msg.Bidder]
	if !found {
		return types.Bid{}, types.ErrNotAllowedBidder
	}

	// !! For BidTypeBatchMany, msg.Coin * msg.Price must be reserved.
	if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), msg.Coin); err != nil {
		return types.Bid{}, err
	}

	nextBidId := k.GetNextBidIdWithUpdate(ctx, auction.GetId())

	bid := types.Bid{
		AuctionId: msg.AuctionId,
		Id:        nextBidId,
		Bidder:    msg.Bidder,
		Type:      msg.BidType,
		Price:     msg.Price,
		Coin:      msg.Coin,
		Height:    uint64(ctx.BlockHeader().Height),
		IsWinner:  false,
	}

	// Place a bid depending on the auction and the bid types
	switch bid.Type {
	case types.BidTypeFixedPrice:
		if err := k.HandleFixedPriceBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}
		bid.IsWinner = true

	case types.BidTypeBatchWorth:
		if err := k.HandleBatchWorthBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}

	case types.BidTypeBatchMany:
		if err := k.HandleBatchManyBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}
	}

	k.SetBid(ctx, bid)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlaceBid,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyBidderAddress, msg.GetBidder().String()),
			sdk.NewAttribute(types.AttributeKeyBidPrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeKeyBidCoin, msg.Coin.String()),
		),
	})

	return bid, nil
}

func (k Keeper) HandleFixedPriceBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if bid.Coin.Denom != auction.GetPayingCoinDenom() {
		return types.ErrIncorrectCoinDenom
	}

	if !bid.Price.Equal(auction.GetStartPrice()) {
		return sdkerrors.Wrap(types.ErrInvalidStartPrice, "start price must be equal to the auction start price")
	}

	// PayingCoinAmount / Price = ExchangedSellingCoinAmount
	exchangedSellingAmt := bid.Coin.Amount.ToDec().QuoTruncate(bid.Price).TruncateInt()
	exchangedSellingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, exchangedSellingAmt)

	// The bidder can't bid more than the remaining selling coin
	if auction.GetRemainingSellingCoin().IsLT(exchangedSellingCoin) {
		return sdkerrors.Wrapf(types.ErrInsufficientRemainingAmount, "remaining selling coin amount %s", auction.GetRemainingSellingCoin())
	}

	// Get the total bid amount by the bidder
	totalBidAmt := sdk.ZeroInt()
	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		if b.Bidder == bid.Bidder {
			exchangedSellingAmt := b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()
			totalBidAmt = totalBidAmt.Add(exchangedSellingAmt)
		}
	}

	totalBidAmt = totalBidAmt.Add(exchangedSellingAmt)
	maxBidAmt := auction.GetMaxBidAmount(bid.Bidder)

	// The sum of total bid amount and bid amount can't be more than the bidder's maximum bid amount
	if totalBidAmt.GT(maxBidAmt) {
		return types.ErrOverMaxBidAmountLimit
	}

	remaining := auction.GetRemainingSellingCoin().Sub(exchangedSellingCoin)
	_ = auction.SetRemainingSellingCoin(remaining)

	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlaceBid,
			sdk.NewAttribute(types.AttributeKeyBidAmount, exchangedSellingCoin.String()),
		),
	})

	return nil
}

func (k Keeper) HandleBatchWorthBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if bid.Coin.Denom != auction.GetPayingCoinDenom() {
		return types.ErrIncorrectCoinDenom
	}

	// Get the total bid amount by the bidder
	totalBidAmt := sdk.ZeroInt()
	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		if b.Bidder == bid.Bidder {
			if b.Type == types.BidTypeBatchMany {
				totalBidAmt = totalBidAmt.Add(b.Coin.Amount)
			} else if b.Type == types.BidTypeBatchWorth {
				exchangedSellingAmt := b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()
				totalBidAmt = totalBidAmt.Add(exchangedSellingAmt)
			}
		}
	}

	exchangedSellingAmt := bid.Coin.Amount.ToDec().Quo(bid.Price).Ceil().TruncateInt()
	totalBidAmt = totalBidAmt.Add(exchangedSellingAmt)
	maxBidAmt := auction.GetMaxBidAmount(bid.Bidder)

	if totalBidAmt.GT(maxBidAmt) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

func (k Keeper) HandleBatchManyBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if bid.Coin.Denom != auction.GetSellingCoin().Denom {
		return types.ErrIncorrectCoinDenom
	}

	// Get the total bid amount by the bidder
	totalBidAmt := sdk.ZeroInt()
	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		if b.Bidder == bid.Bidder {
			if b.Type == types.BidTypeBatchMany {
				totalBidAmt = totalBidAmt.Add(b.Coin.Amount)
			} else if b.Type == types.BidTypeBatchWorth {
				exchangedSellingAmt := b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()
				totalBidAmt = totalBidAmt.Add(exchangedSellingAmt)
			}
		}
	}

	totalBidAmt = totalBidAmt.Add(bid.Coin.Amount)
	maxBidAmt := auction.GetMaxBidAmount(bid.Bidder)

	if totalBidAmt.GT(maxBidAmt) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

// ModifyBid modifies the auctioneer's bid
func (k Keeper) ModifyBid(ctx sdk.Context, msg *types.MsgModifyBid) (types.MsgModifyBid, error) {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrNotFound, "auction not found")
	}

	if auction.GetType() != types.AuctionTypeBatch {
		return types.MsgModifyBid{}, types.ErrIncorrectAuctionType
	}

	// !! Here seems to need modification: A same BidId can be in different bidders.
	// need to check only within bid list of this bidder.
	bid, found := k.GetBid(ctx, msg.AuctionId, msg.BidId)
	if !found {
		return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrNotFound, "bid not found")
	}

	// Modifying bid type is not allowed
	if bid.Coin.Denom != msg.Coin.Denom {
		return types.MsgModifyBid{}, types.ErrIncorrectCoinDenom
	}

	// Modified by Jeongho
	// Either bid price or coin amount must be higher than the previous bid
	if msg.Price.LT(bid.Price) {
		return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bid price cannot be lower")
	}
	if msg.Coin.IsLT(bid.Coin) {
		return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bid coin amount cannot be lower")
	}
	if msg.Price.Equal(bid.Price) && msg.Coin.IsEqual(bid.Coin) {
		return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Either bid price or coin amount must be higher")
	}

	// OR can be implemented as below
	//if !(msg.Price.GTE(bid.Price) && msg.Coin.IsGTE(bid.Coin)) {
	//	return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Either bid price or coin amount must be higher")
	//}

	// !! Here seems to need modification: This needs to handle differently according to Bid Type
	// If BidTypeBatchWorth --> just subtract msg.Coin
	// If BidTypeBatchMany --> subtract msg.Coin * msg.Price from bid.Coin * bid.Price
	// Reserve the bid amount difference
	diffBidAmt := msg.Coin.Sub(bid.Coin) // Suggest to change the name from diffBidAmt --> diffBidPayingAmt

	// !! Need to check if total Bid Amount does not exceed maxBidAmount.
	// HERE

	if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), diffBidAmt); err != nil {
		return types.MsgModifyBid{}, err
	}

	bid.Price = msg.Price
	bid.Coin = msg.Coin
	bid.Height = uint64(ctx.BlockHeader().Height)

	k.SetBid(ctx, bid)

	return types.MsgModifyBid{}, nil
}
