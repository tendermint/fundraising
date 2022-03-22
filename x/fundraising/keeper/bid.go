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

// PlaceBid places a bid for the auction.
func (k Keeper) PlaceBid(ctx sdk.Context, msg *types.MsgPlaceBid) (types.Bid, error) {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return types.Bid{}, sdkerrors.Wrap(sdkerrors.ErrNotFound, "auction not found")
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return types.Bid{}, types.ErrInvalidAuctionStatus
	}

	if auction.GetType() == types.AuctionTypeBatch {
		if msg.Price.LT(auction.(*types.BatchAuction).MinBidPrice) {
			return types.Bid{}, types.ErrInsufficientMinBidPrice
		}
	}

	_, found = auction.GetAllowedBiddersMap()[msg.Bidder]
	if !found {
		return types.Bid{}, types.ErrNotAllowedBidder
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
		IsMatched: false,
	}

	// Place a bid depending on the auction type and the bid type
	switch bid.Type {
	case types.BidTypeFixedPrice:
		if err := k.ValidateFixedPriceBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}

		bidPayingAmt := bid.GetBidPayingAmount(auction.GetPayingCoinDenom())
		bidPayingCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), bidPayingAmt)
		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), bidPayingCoin); err != nil {
			return types.Bid{}, err
		}

		bidSellingAmt := bid.GetBidSellingAmount(auction.GetPayingCoinDenom())
		bidSellingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, bidSellingAmt)
		remaining := auction.GetRemainingSellingCoin().Sub(bidSellingCoin)

		_ = auction.SetRemainingSellingCoin(remaining)
		k.SetAuction(ctx, auction)
		bid.SetMatched(true)

	case types.BidTypeBatchWorth:
		if err := k.ValidateBatchWorthBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), msg.Coin); err != nil {
			return types.Bid{}, err
		}

	case types.BidTypeBatchMany:
		if err := k.ValidateBatchManyBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}

		reserveAmt := msg.Coin.Amount.ToDec().Mul(msg.Price).Ceil().TruncateInt()
		reserveCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), reserveAmt)

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), reserveCoin); err != nil {
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

func (k Keeper) ValidateFixedPriceBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if bid.Coin.Denom != auction.GetPayingCoinDenom() && bid.Coin.Denom != auction.GetSellingCoin().Denom {
		return types.ErrIncorrectCoinDenom
	}

	if !bid.Price.Equal(auction.GetStartPrice()) {
		return sdkerrors.Wrap(types.ErrInvalidStartPrice, "start price must be equal to the auction start price")
	}

	// PayingCoinAmount / Price = BidSellingCoinAmount
	bidSellingAmt := bid.GetBidSellingAmount(auction.GetPayingCoinDenom())
	bidSellingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, bidSellingAmt)

	// The bidder can't bid more than the remaining selling coin
	remainingSellingCoin := auction.GetRemainingSellingCoin()
	if remainingSellingCoin.IsLT(bidSellingCoin) {
		return sdkerrors.Wrapf(types.ErrInsufficientRemainingAmount, "remaining selling coin amount %s", remainingSellingCoin)
	}

	// Get the total bid amount by the bidder
	totalBidAmt := sdk.ZeroInt()
	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		if b.Bidder == bid.Bidder {
			totalBidAmt = totalBidAmt.Add(b.GetBidSellingAmount(auction.GetPayingCoinDenom()))
		}
	}

	totalBidAmt = totalBidAmt.Add(bidSellingAmt)
	maxBidAmt := auction.GetMaxBidAmount(bid.Bidder)

	// The  sum of total bid amount can't be more than the bidder's maximum bid amount
	if totalBidAmt.GT(maxBidAmt) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

func (k Keeper) ValidateBatchWorthBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if bid.Coin.Denom != auction.GetPayingCoinDenom() {
		return types.ErrIncorrectCoinDenom
	}

	// Get the bid amount by the bidder with bid_price
	bidAmt := bid.Coin.Amount.ToDec().QuoTruncate(bid.Price).TruncateInt()
	maxBidAmt := auction.GetMaxBidAmount(bid.Bidder)

	// The  sum of total bid amount can't be more than the bidder's maximum bid amount
	if bidAmt.GT(maxBidAmt) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

func (k Keeper) ValidateBatchManyBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if bid.Coin.Denom != auction.GetSellingCoin().Denom {
		return types.ErrIncorrectCoinDenom
	}

	// Get the bid amount by the bidder with bid_price
	bidAmt := bid.Coin.Amount
	maxBidAmt := auction.GetMaxBidAmount(bid.Bidder)

	// The  sum of total bid amount can't be more than the bidder's maximum bid amount
	if bidAmt.GT(maxBidAmt) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

// ModifyBid handles types.MsgModifyBid and stores the modified bid.
// A bidder must provide either greater bid price or coin amount.
// They are not permitted to modify with less bid price or coin amount.
func (k Keeper) ModifyBid(ctx sdk.Context, msg *types.MsgModifyBid) error {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return sdkerrors.Wrap(sdkerrors.ErrNotFound, "auction not found")
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return types.ErrInvalidAuctionStatus
	}

	if auction.GetType() != types.AuctionTypeBatch {
		return types.ErrIncorrectAuctionType
	}

	bid, found := k.GetBid(ctx, msg.AuctionId, msg.BidId)
	if !found {
		return sdkerrors.Wrap(sdkerrors.ErrNotFound, "bid not found")
	}

	if !bid.GetBidder().Equals(msg.GetBidder()) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only the bid creator can modify the bid")
	}

	if msg.Price.LT(auction.(*types.BatchAuction).MinBidPrice) {
		return types.ErrInsufficientMinBidPrice
	}

	// Modifying bid type is not allowed
	if bid.Coin.Denom != msg.Coin.Denom {
		return sdkerrors.Wrap(types.ErrIncorrectCoinDenom, "modifying bid type is not allowed")
	}

	// Either bid price or coin amount must be greater than the modifying values
	if msg.Price.LT(bid.Price) || msg.Coin.Amount.LT(bid.Coin.Amount) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "bid price or coin amount cannot be lower")
	}

	if msg.Price.Equal(bid.Price) && msg.Coin.Amount.Equal(bid.Coin.Amount) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "bid price and coin amount must be changed")
	}

	// Reserve bid amount difference
	switch bid.Type {
	case types.BidTypeBatchWorth:
		diffBidCoin := msg.Coin.Sub(bid.Coin)
		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), diffBidCoin); err != nil {
			return err
		}
	case types.BidTypeBatchMany:
		prevBidAmt := msg.Coin.Amount.ToDec().Mul(msg.Price)
		currBidAmt := bid.Coin.Amount.ToDec().Mul(bid.Price)
		diffBidAmt := prevBidAmt.Sub(currBidAmt).TruncateInt()
		diffBidCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), diffBidAmt)

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), diffBidCoin); err != nil {
			return err
		}
	}

	bid.Price = msg.Price
	bid.Coin = msg.Coin
	bid.Height = uint64(ctx.BlockHeader().Height)

	k.SetBid(ctx, bid)

	return nil
}
