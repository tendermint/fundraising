package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetNextSequence increments sequence number by one and set it.
func (k Keeper) GetNextSequenceWithUpdate(ctx sdk.Context, auctionId uint64) uint64 {
	seq := k.GetLastSequence(ctx, auctionId) + 1
	k.SetSequence(ctx, auctionId, seq)
	return seq
}

// ReservePayingCoin reserves paying coin to the paying reserve account.
func (k Keeper) ReservePayingCoin(ctx sdk.Context, auctionId uint64, bidderAcc sdk.AccAddress, payingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, bidderAcc, types.PayingReserveAcc(auctionId), sdk.NewCoins(payingCoin)); err != nil {
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
		return types.Bid{}, sdkerrors.Wrapf(types.ErrInvalidAuctionStatus, "unable to bid because the auction is in %s", auction.GetStatus().String())
	}

	if auction.GetPayingCoinDenom() != msg.Coin.Denom {
		return types.Bid{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "coin denom must match with the paying coin denom")
	}

	if err := k.ReservePayingCoin(ctx, auction.GetId(), msg.GetBidder(), msg.Coin); err != nil {
		return types.Bid{}, err
	}

	receiveAmt := msg.Coin.Amount.ToDec().QuoTruncate(msg.Price).TruncateInt()
	receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, receiveAmt)

	// The bidder cannot bid more than the remaining coin
	remaining := auction.GetRemainingCoin().Sub(receiveCoin)
	if remaining.IsNegative() {
		return types.Bid{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "request coin must be lower than or equal to the remaining coin")
	}

	if auction.GetType() == types.AuctionTypeFixedPrice {
		if !msg.Price.Equal(auction.GetStartPrice()) {
			return types.Bid{}, sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to the start price of the auction")
		}

		if err := auction.SetRemainingCoin(remaining); err != nil {
			return types.Bid{}, err
		}

		k.SetAuction(ctx, auction)
	} else if auction.GetType() == types.AuctionTypeEnglish {
		bids := k.GetBidsByBidder(ctx, msg.GetBidder())
		bidsLen := len(bids)

		for i, bid := range bids {
			if i == bidsLen-1 {
				if msg.Price.LTE(bid.Price) {
					return types.Bid{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bid price can only increase for English auction")
				}
			}
		}

		if err := auction.SetRemainingCoin(remaining); err != nil {
			return types.Bid{}, err
		}

		k.SetAuction(ctx, auction)
	}

	sequenceId := k.GetNextSequenceWithUpdate(ctx, auction.GetId())

	bid := types.Bid{
		AuctionId: auction.GetId(),
		Sequence:  sequenceId,
		Bidder:    msg.Bidder,
		Price:     msg.Price,
		Coin:      msg.Coin,
		Height:    uint64(ctx.BlockHeader().Height),
		Eligible:  true,
	}

	k.SetBid(ctx, bid.AuctionId, bid.Sequence, msg.GetBidder(), bid)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlaceBid,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyBidderAddress, msg.GetBidder().String()),
			sdk.NewAttribute(types.AttributeKeyBidPrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeKeyBidCoin, msg.Coin.String()),
			sdk.NewAttribute(types.AttributeKeyBidAmount, receiveCoin.String()),
		),
	})

	return bid, nil
}
