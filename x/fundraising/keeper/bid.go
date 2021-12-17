package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetBid returns a bid for the given auction id and sequence number.
// A bidder can have as many bids as they want, so sequence is required to get the bid.
func (k Keeper) GetBid(ctx sdk.Context, auctionId uint64, sequence uint64) (bid types.Bid, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBidKey(auctionId, sequence))
	if bz == nil {
		return bid, false
	}
	k.cdc.MustUnmarshal(bz, &bid)
	return bid, true
}

// SetBid sets a bid with the given arguments.
func (k Keeper) SetBid(ctx sdk.Context, auctionId uint64, sequence uint64, bidderAcc sdk.AccAddress, bid types.Bid) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&bid)
	store.Set(types.GetBidKey(auctionId, sequence), bz)
	store.Set(types.GetBidIndexKey(bidderAcc, auctionId, sequence), []byte{})
}

// GetBids returns all bids registered in the store.
func (k Keeper) GetBids(ctx sdk.Context) []types.Bid {
	bids := []types.Bid{}
	k.IterateBids(ctx, func(bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})
	return bids
}

// GetBidsByAuctionId returns all bids associated with the auction id that are registered in the store.
func (k Keeper) GetBidsByAuctionId(ctx sdk.Context, auctionId uint64) []types.Bid {
	bids := []types.Bid{}
	k.IterateBidsByAuctionId(ctx, auctionId, func(bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})
	return bids
}

// GetBidsByBidder returns all bids associated with the bidder that are registered in the store.
func (k Keeper) GetBidsByBidder(ctx sdk.Context, bidderAcc sdk.AccAddress) []types.Bid {
	bids := []types.Bid{}
	k.IterateBidsByBidder(ctx, bidderAcc, func(bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})
	return bids
}

// IterateBids iterates through all bids stored in the store and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateBids(ctx sdk.Context, cb func(bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.BidKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		if cb(bid) {
			break
		}
	}
}

// IterateBidsByAuctionId iterates through all bids associated with the auction id stored in the store
// and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateBidsByAuctionId(ctx sdk.Context, auctionId uint64, cb func(bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidAuctionIDKey(auctionId))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		if cb(bid) {
			break
		}
	}
}

// IterateBidsByBidder iterates through all bids associated with the bidder stored in the store
// and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateBidsByBidder(ctx sdk.Context, bidderAcc sdk.AccAddress, cb func(bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidIndexByBidderPrefix(bidderAcc))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		auctionId, sequence := types.ParseBidIndexKey(iter.Key())
		bid, _ := k.GetBid(ctx, auctionId, sequence)
		if cb(bid) {
			break
		}
	}
}

// ReservePayingCoin reserves paying coin to the paying reserve account.
func (k Keeper) ReservePayingCoin(ctx sdk.Context, auctionId uint64, bidderAcc sdk.AccAddress, payingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, bidderAcc, types.PayingReserveAcc(auctionId), sdk.NewCoins(payingCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to reserve paying coin")
	}
	return nil
}

// PlaceBid places a bid for the auction.
func (k Keeper) PlaceBid(ctx sdk.Context, msg *types.MsgPlaceBid) error {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", msg.AuctionId)
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return sdkerrors.Wrapf(types.ErrInvalidAuctionStatus, "unable to bid because the auction is in %s", auction.GetStatus().String())
	}

	if auction.GetPayingCoinDenom() != msg.Coin.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "coin denom must match with the paying coin denom")
	}

	bidAmt := msg.Coin.Amount.ToDec().Quo(msg.Price).TruncateInt()
	receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, bidAmt)
	balanceAmt := k.bankKeeper.GetBalance(ctx, msg.GetBidder(), auction.GetPayingCoinDenom()).Amount

	// the bidder must have greater than or equal to the bid amount
	if balanceAmt.Sub(bidAmt).IsNegative() {
		return sdkerrors.ErrInsufficientFunds
	}

	// the bidder cannot bid more than the remaining coin
	remaining := auction.GetRemainingCoin().Sub(receiveCoin)
	if remaining.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "request coin must be lower than or equal to the remaining total selling coin")
	}

	if auction.GetType() == types.AuctionTypeFixedPrice {
		if !msg.Price.Equal(auction.GetStartPrice()) {
			return sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to the start price of the auction")
		}

		if err := auction.SetRemainingCoin(remaining); err != nil {
			return err
		}

		k.SetAuction(ctx, auction)

		if err := k.ReservePayingCoin(ctx, auction.GetId(), msg.GetBidder(), msg.Coin); err != nil {
			return err
		}

	} else {
		// TODO: implement English auction type
		return sdkerrors.Wrap(types.ErrInvalidAuctionType, "not supported auction type in this version")
	}

	sequenceId := k.GetNextSequenceWithUpdate(ctx)

	bid := types.Bid{
		AuctionId: auction.GetId(),
		Sequence:  sequenceId,
		Bidder:    msg.Bidder,
		Price:     msg.Price,
		Coin:      msg.Coin,
		Height:    uint64(ctx.BlockHeader().Height),
		Eligible:  false, // it becomes true when a bidder receives succesfully during distribution in endblocker
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

	return nil
}
