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

	if err := k.ReservePayingCoin(ctx, auction.GetId(), msg.GetBidder(), msg.Coin); err != nil {
		return types.Bid{}, err
	}

	allowedBiddersMap := make(map[string]sdk.Int) // map(bidder => maxBidAmount)
	for _, ab := range auction.GetAllowedBidders() {
		allowedBiddersMap[ab.Bidder] = ab.MaxBidAmount
	}

	maxBidAmt, found := allowedBiddersMap[msg.Bidder]
	if !found {
		return types.Bid{}, types.ErrNotAllowedBidder
	}

	bidId := k.GetNextBidIdWithUpdate(ctx, auction.GetId())

	bid := types.Bid{
		AuctionId: auction.GetId(),
		Id:        bidId,
		Bidder:    msg.Bidder,
		Price:     msg.Price,
		Coin:      msg.Coin,
		Height:    uint64(ctx.BlockHeader().Height),
		IsWinner:  false,
	}

	// Place a bid depending on the auction and the bid types
	switch auction.GetType() {
	case types.AuctionTypeFixedPrice:
		if err := k.PlaceBidFixedPrice(ctx, msg, auction, maxBidAmt); err != nil {
			return types.Bid{}, err
		}

		bid.IsWinner = true

	case types.AuctionTypeBatch:
		if msg.BidType == types.BidTypeBatchWorth {
			if err := k.PlaceBidBatchWorth(ctx, msg, auction, maxBidAmt); err != nil {
				return types.Bid{}, err
			}
		} else if msg.BidType == types.BidTypeBatchMany {
			if err := k.PlaceBidBatchMany(ctx, msg, auction, maxBidAmt); err != nil {
				return types.Bid{}, err
			}
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

func (k Keeper) PlaceBidFixedPrice(ctx sdk.Context, msg *types.MsgPlaceBid, auction types.AuctionI, maxBidAmt sdk.Int) error {
	if msg.Coin.Denom != auction.GetPayingCoinDenom() {
		return types.ErrIncorrectCoinDenom
	}

	if !msg.Price.Equal(auction.GetStartPrice()) {
		return types.ErrInvalidStartPrice
	}

	// PayingCoinAmount / Price = ExchangedSellingCoinAmount
	exchangedSellingAmt := msg.Coin.Amount.ToDec().QuoTruncate(msg.Price).TruncateInt()
	exchangedSellingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, exchangedSellingAmt)

	// The bidder can't bid more than the remaining selling coin
	if auction.GetRemainingSellingCoin().IsLT(exchangedSellingCoin) {
		return sdkerrors.Wrapf(types.ErrInsufficientRemainingAmount, "remaining selling coin amount %s", auction.GetRemainingSellingCoin())
	}

	// Get the total bid amount by the bidder
	totalBidAmt := sdk.ZeroInt()
	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		if b.Bidder == msg.Bidder {
			exchangedSellingAmt := b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()
			totalBidAmt = totalBidAmt.Add(exchangedSellingAmt)
		}
	}

	// The sum of total bid amount and bid amount can't be more than the bidder's maximum bid amount
	totalBidAmt = totalBidAmt.Add(exchangedSellingAmt)
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

func (k Keeper) PlaceBidBatchWorth(ctx sdk.Context, msg *types.MsgPlaceBid, auction types.AuctionI, maxBidAmt sdk.Int) error {
	if msg.Coin.Denom != auction.GetPayingCoinDenom() {
		return types.ErrIncorrectCoinDenom
	}

	return nil
}

func (k Keeper) PlaceBidBatchMany(ctx sdk.Context, msg *types.MsgPlaceBid, auction types.AuctionI, maxBidAmt sdk.Int) error {
	if msg.Coin.Denom != auction.GetSellingCoin().Denom {
		return types.ErrIncorrectCoinDenom
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

	bid, found := k.GetBid(ctx, msg.AuctionId, msg.BidId)
	if !found {
		return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrNotFound, "bid not found")
	}

	// Not allowed to modify the bid type
	if bid.Coin.Denom != msg.Coin.Denom {
		return types.MsgModifyBid{}, types.ErrIncorrectCoinDenom
	}

	exchangedSellingAmtBefore := bid.Coin.Amount.ToDec().QuoTruncate(bid.Price).TruncateInt()
	exchangedSellingAmt := msg.Coin.Amount.ToDec().QuoTruncate(msg.Price).TruncateInt()

	// Either bid price or coin amount must be higher than the previous bid
	if exchangedSellingAmtBefore.LT(exchangedSellingAmt) {
		return types.MsgModifyBid{},
			sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "either bid price or coin amount must be higher than the previous bid")
	}

	bid.Price = msg.Price
	bid.Coin = msg.Coin
	bid.Height = uint64(ctx.BlockHeader().Height)

	k.SetBid(ctx, bid)

	return types.MsgModifyBid{}, nil
}
