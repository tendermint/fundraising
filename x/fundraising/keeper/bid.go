package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetNextBidId increments bid id by one and set it.
func (k Keeper) GetNextBidIdWithUpdate(ctx sdk.Context, auctionId uint64) uint64 {
	seq := k.GetLastBidId(ctx, auctionId) + 1
	k.SetBidId(ctx, auctionId, seq)
	return seq
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
		return types.Bid{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", msg.AuctionId)
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return types.Bid{}, types.ErrInvalidAuctionStatus
	}

	if auction.GetPayingCoinDenom() != msg.BidCoin.Denom {
		return types.Bid{}, types.ErrInvalidPayingCoinDenom
	}

	allowedBiddersMap := make(map[string]sdk.Int) // map(bidder => maxBidAmount)
	for _, bidder := range auction.GetAllowedBidders() {
		allowedBiddersMap[bidder.GetBidder()] = bidder.MaxBidAmount
	}

	// The bidder must be in the allowed bidder list in order to bid
	maxBidAmt, found := allowedBiddersMap[msg.Bidder]
	if !found {
		return types.Bid{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bidder %s is not allowed to bid", msg.Bidder)
	}

	if err := k.ReservePayingCoin(ctx, auction.GetId(), msg.GetBidder(), msg.BidCoin); err != nil {
		return types.Bid{}, err
	}

	// Handle logics depending on auction type
	if auction.GetType() == types.AuctionTypeFixedPrice {
		if !msg.BidPrice.Equal(auction.GetStartPrice()) {
			return types.Bid{},
				sdkerrors.Wrapf(types.ErrInvalidStartPrice, "expected start price %s, got %s", auction.GetStartPrice(), msg.BidPrice)
		}

		receiveAmt := msg.BidCoin.Amount.ToDec().QuoTruncate(msg.BidPrice).TruncateInt()
		receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, receiveAmt)

		// The receive amount can't be greater than the bidder's maximum bid amount
		if receiveAmt.GT(maxBidAmt) {
			return types.Bid{},
				sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "bid amount %s, maximum bid amount %s", receiveAmt, maxBidAmt)
		}

		if auction.GetRemainingSellingCoin().IsLT(receiveCoin) {
			return types.Bid{},
				sdkerrors.Wrapf(types.ErrInsufficientRemainingAmount, "remaining coin amount %s", auction.GetRemainingSellingCoin())
		}

		remaining := auction.GetRemainingSellingCoin().Sub(receiveCoin)
		if err := auction.SetRemainingSellingCoin(remaining); err != nil {
			return types.Bid{}, err
		}
		k.SetAuction(ctx, auction)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypePlaceBid,
				sdk.NewAttribute(types.AttributeKeyBidAmount, receiveCoin.String()),
			),
		})
	} else {
		// TODO: implement English auction type
		return types.Bid{}, sdkerrors.Wrap(types.ErrInvalidAuctionType, "not supported auction type in this version")
	}

	seqId := k.GetNextBidIdWithUpdate(ctx, auction.GetId())

	bid := types.Bid{
		AuctionId: auction.GetId(),
		Id:        seqId,
		Bidder:    msg.Bidder,
		BidPrice:  msg.BidPrice,
		BidCoin:   msg.BidCoin,
		Height:    uint64(ctx.BlockHeader().Height),
		IsWinner:  true,
	}

	k.SetBid(ctx, bid.AuctionId, bid.Id, msg.GetBidder(), bid)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlaceBid,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyBidderAddress, msg.GetBidder().String()),
			sdk.NewAttribute(types.AttributeKeyBidPrice, msg.BidPrice.String()),
			sdk.NewAttribute(types.AttributeKeyBidCoin, msg.BidCoin.String()),
		),
	})

	return bid, nil
}
