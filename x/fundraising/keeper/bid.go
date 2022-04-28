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

// PlaceBid places a bid for the selling coin of the auction.
func (k Keeper) PlaceBid(ctx sdk.Context, msg *types.MsgPlaceBid) (types.Bid, error) {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return types.Bid{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d not found", msg.AuctionId)
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return types.Bid{}, types.ErrInvalidAuctionStatus
	}

	if auction.GetType() == types.AuctionTypeBatch {
		if msg.Price.LT(auction.(*types.BatchAuction).MinBidPrice) {
			return types.Bid{}, types.ErrInsufficientMinBidPrice
		}
	}

	allowedBidders := k.GetAllowedBiddersByAuction(ctx, auction.GetId())
	allowedBiddersMap := types.GetAllowedBiddersMap(allowedBidders)
	_, found = allowedBiddersMap[msg.Bidder]
	if !found {
		return types.Bid{}, types.ErrNotAllowedBidder
	}

	bid := types.Bid{
		AuctionId: msg.AuctionId,
		Id:        k.GetNextBidIdWithUpdate(ctx, auction.GetId()),
		Bidder:    msg.Bidder,
		Type:      msg.BidType,
		Price:     msg.Price,
		Coin:      msg.Coin,
		IsMatched: false,
	}

	payingCoinDenom := auction.GetPayingCoinDenom()

	// Place a bid depending on the bid type
	switch bid.Type {
	case types.BidTypeFixedPrice:
		if err := k.ValidateFixedPriceBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}

		fa := auction.(*types.FixedPriceAuction)

		// Reserve bid amount
		bidPayingAmt := bid.ConvertToPayingAmount(payingCoinDenom)
		bidPayingCoin := sdk.NewCoin(payingCoinDenom, bidPayingAmt)
		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), bidPayingCoin); err != nil {
			return types.Bid{}, err
		}

		// Subtract bid amount from the remaining
		bidSellingAmt := bid.ConvertToSellingAmount(payingCoinDenom)
		bidSellingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, bidSellingAmt)
		fa.RemainingSellingCoin = fa.RemainingSellingCoin.Sub(bidSellingCoin)

		k.SetAuction(ctx, fa)
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

		reserveAmt := bid.ConvertToPayingAmount(payingCoinDenom)
		reserveCoin := sdk.NewCoin(payingCoinDenom, reserveAmt)

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), reserveCoin); err != nil {
			return types.Bid{}, err
		}
	}

	// Call before bid placed hook
	k.BeforeBidPlaced(ctx, bid.AuctionId, bid.Bidder, bid.Type, bid.Price, bid.Coin)

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

// ValidateFixedPriceBid validates a fixed price bid type.
func (k Keeper) ValidateFixedPriceBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if auction.GetType() != types.AuctionTypeFixedPrice {
		return types.ErrIncorrectAuctionType
	}

	if bid.Coin.Denom != auction.GetPayingCoinDenom() &&
		bid.Coin.Denom != auction.GetSellingCoin().Denom {
		return types.ErrIncorrectCoinDenom
	}

	if !bid.Price.Equal(auction.GetStartPrice()) {
		return sdkerrors.Wrap(types.ErrInvalidStartPrice, "start price must be equal to the start price of the auction")
	}

	// For remaining coin validation, convert bid amount in selling coin denom
	bidAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())
	bidCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, bidAmt)
	remainingCoin := auction.(*types.FixedPriceAuction).RemainingSellingCoin

	if remainingCoin.IsLT(bidCoin) {
		return sdkerrors.Wrapf(types.ErrInsufficientRemainingAmount, "remaining selling coin amount %s", remainingCoin)
	}

	// Get the total bid amount by the bidder
	totalBidAmt := sdk.ZeroInt()
	for _, bid := range k.GetBidsByBidder(ctx, bid.GetBidder()) {
		if bid.AuctionId == auction.GetId() {
			bidSellingAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())
			totalBidAmt = totalBidAmt.Add(bidSellingAmt)
		}
	}

	allowedBidder, found := k.GetAllowedBidder(ctx, bid.AuctionId, bid.GetBidder())
	if !found {
		return sdkerrors.Wrap(sdkerrors.ErrNotFound, "bidder is not found in allowed bidder list")
	}

	totalBidAmt = totalBidAmt.Add(bidAmt)

	// The total bid amount can't be greater than the bidder's maximum bid amount
	if totalBidAmt.GT(allowedBidder.MaxBidAmount) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

// ValidateBatchWorthBid validates a batch worth bid type.
func (k Keeper) ValidateBatchWorthBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if auction.GetType() != types.AuctionTypeBatch {
		return types.ErrIncorrectAuctionType
	}

	if bid.Coin.Denom != auction.GetPayingCoinDenom() {
		return types.ErrIncorrectCoinDenom
	}

	allowedBidder, found := k.GetAllowedBidder(ctx, bid.AuctionId, bid.GetBidder())
	if !found {
		return sdkerrors.Wrap(sdkerrors.ErrNotFound, "bidder is not found in allowed bidder list")
	}

	bidAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())

	// The total bid amount can't be greater than the bidder's maximum bid amount
	if bidAmt.GT(allowedBidder.MaxBidAmount) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

// ValidateBatchManyBid validates a batch many bid type.
func (k Keeper) ValidateBatchManyBid(ctx sdk.Context, auction types.AuctionI, bid types.Bid) error {
	if auction.GetType() != types.AuctionTypeBatch {
		return types.ErrIncorrectAuctionType
	}

	if bid.Coin.Denom != auction.GetSellingCoin().Denom {
		return types.ErrIncorrectCoinDenom
	}

	allowedBidder, found := k.GetAllowedBidder(ctx, bid.AuctionId, bid.GetBidder())
	if !found {
		return sdkerrors.Wrap(sdkerrors.ErrNotFound, "bidder is not found in allowed bidder list")
	}

	bidAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())

	// The total bid amount can't be greater than the bidder's maximum bid amount
	if bidAmt.GT(allowedBidder.MaxBidAmount) {
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

	if bid.Coin.Denom != msg.Coin.Denom {
		return types.ErrIncorrectCoinDenom
	}

	if msg.Price.LT(bid.Price) || msg.Coin.Amount.LT(bid.Coin.Amount) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bid price or coin amount cannot be lower")
	}

	if msg.Price.Equal(bid.Price) && msg.Coin.Amount.Equal(bid.Coin.Amount) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bid price and coin amount must be changed")
	}

	// Reserve bid amount difference
	switch bid.Type {
	case types.BidTypeBatchWorth:
		diffReserveCoin := msg.Coin.Sub(bid.Coin)
		if diffReserveCoin.IsPositive() {
			if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), diffReserveCoin); err != nil {
				return err
			}
		}
	case types.BidTypeBatchMany:
		prevReserveAmt := bid.Coin.Amount.ToDec().Mul(bid.Price).Ceil()
		currReserveAmt := msg.Coin.Amount.ToDec().Mul(msg.Price).Ceil()
		diffReserveAmt := currReserveAmt.Sub(prevReserveAmt).TruncateInt()
		diffReserveCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), diffReserveAmt)
		if diffReserveCoin.IsPositive() {
			if err := k.ReservePayingCoin(ctx, msg.AuctionId, msg.GetBidder(), diffReserveCoin); err != nil {
				return err
			}
		}
	}

	bid.Price = msg.Price
	bid.Coin = msg.Coin

	// Call the before mid modified hook
	k.BeforeBidModified(ctx, bid.AuctionId, bid.Bidder, bid.Type, bid.Price, bid.Coin)

	k.SetBid(ctx, bid)

	return nil
}
