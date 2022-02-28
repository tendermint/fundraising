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
		IsWinner:  false,
	}

	// Place a bid depending on the auction type and the bid type
	switch bid.Type {
	case types.BidTypeFixedPrice:
		if err := k.ValidateFixedPriceBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), msg.Coin); err != nil {
			return types.Bid{}, err
		}

		exchangedSellingAmt := bid.GetExchangedSellingAmount()
		exchangedSellingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, exchangedSellingAmt)
		remaining := auction.GetRemainingSellingCoin().Sub(exchangedSellingCoin)

		_ = auction.SetRemainingSellingCoin(remaining)
		k.SetAuction(ctx, auction)
		bid.SetWinner(true)

	case types.BidTypeBatchWorth:
		if bid.Coin.Denom != auction.GetPayingCoinDenom() {
			return types.Bid{}, types.ErrIncorrectCoinDenom
		}

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), msg.Coin); err != nil {
			return types.Bid{}, err
		}

	case types.BidTypeBatchMany:
		if bid.Coin.Denom != auction.GetSellingCoin().Denom {
			return types.Bid{}, types.ErrIncorrectCoinDenom
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
	if bid.Coin.Denom != auction.GetPayingCoinDenom() {
		return types.ErrIncorrectCoinDenom
	}

	if !bid.Price.Equal(auction.GetStartPrice()) {
		return sdkerrors.Wrap(types.ErrInvalidStartPrice, "start price must be equal to the auction start price")
	}

	// PayingCoinAmount / Price = ExchangedSellingCoinAmount
	exchangedSellingAmt := bid.GetExchangedSellingAmount()
	exchangedSellingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, exchangedSellingAmt)

	// The bidder can't bid more than the remaining selling coin
	remainingSellingCoin := auction.GetRemainingSellingCoin()
	if remainingSellingCoin.IsLT(exchangedSellingCoin) {
		return sdkerrors.Wrapf(types.ErrInsufficientRemainingAmount, "remaining selling coin amount %s", remainingSellingCoin)
	}

	// Get the total bid amount by the bidder
	totalBidAmt := sdk.ZeroInt()
	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		if b.Bidder == bid.Bidder {
			totalBidAmt = totalBidAmt.Add(b.GetExchangedSellingAmount())
		}
	}

	totalBidAmt = totalBidAmt.Add(exchangedSellingAmt)
	maxBidAmt := auction.GetMaxBidAmount(bid.Bidder)

	// The  sum of total bid amount can't be more than the bidder's maximum bid amount
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

	if auction.GetStatus() != types.AuctionStatusStarted {
		return types.MsgModifyBid{}, types.ErrInvalidAuctionStatus
	}

	if auction.GetType() != types.AuctionTypeBatch {
		return types.MsgModifyBid{}, types.ErrIncorrectAuctionType
	}

	bid, found := k.GetBid(ctx, msg.AuctionId, msg.BidId)
	if !found {
		return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrNotFound, "bid not found")
	}

	if bid.Bidder != msg.Bidder {
		return types.MsgModifyBid{}, types.ErrIncorrectOwner
	}

	// Modifying bid type is not allowed
	if bid.Coin.Denom != msg.Coin.Denom {
		return types.MsgModifyBid{}, types.ErrIncorrectCoinDenom
	}

	// TODO: coin amount varies depending on bid type?
	// The bid price or coin amount must be higher than the modifying bid one
	if msg.Price.LTE(bid.Price) && msg.Coin.Amount.LTE(bid.Coin.Amount) {
		return types.MsgModifyBid{},
			sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bid price or coin amount cannot be lower")
	}

	// Reserve bid amount difference
	switch bid.Type {
	case types.BidTypeBatchWorth:
		diffBidCoin := msg.Coin.Sub(bid.Coin)

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), diffBidCoin); err != nil {
			return types.MsgModifyBid{}, err
		}

	case types.BidTypeBatchMany:
		prevBidAmt := msg.Coin.Amount.ToDec().Mul(msg.Price)
		currBidAmt := bid.Coin.Amount.ToDec().Mul(bid.Price)
		diffBidAmt := prevBidAmt.Sub(currBidAmt).TruncateInt()
		diffBidCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), diffBidAmt)

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), diffBidCoin); err != nil {
			return types.MsgModifyBid{}, err
		}
	}

	bid.Price = msg.Price
	bid.Coin = msg.Coin
	bid.Height = uint64(ctx.BlockHeader().Height)

	k.SetBid(ctx, bid)

	return types.MsgModifyBid{}, nil
}
