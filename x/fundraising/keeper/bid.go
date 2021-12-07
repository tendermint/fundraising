package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetBid returns a bid for the given auction id and sequence number.
// A bidder can have as many bids as they want, so sequence is required to get the bid.
func (k Keeper) GetBid(ctx sdk.Context, auctionID uint64, sequence uint64) (bid types.Bid, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBidKey(auctionID, sequence))
	if bz == nil {
		return bid, false
	}
	k.cdc.MustUnmarshal(bz, &bid)
	return bid, true
}

// SetBid sets a bid with the given arguments.
func (k Keeper) SetBid(ctx sdk.Context, auctionID uint64, sequence uint64, bidderAcc sdk.AccAddress, bid types.Bid) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&bid)
	store.Set(types.GetBidKey(auctionID, sequence), bz)
	store.Set(types.GetBidIndexKey(bidderAcc, auctionID, sequence), []byte{})
}

// GetBids returns all bids registered in the store.
func (k Keeper) GetBids(ctx sdk.Context, auctionID uint64) []types.Bid {
	bids := []types.Bid{}
	k.IterateBids(ctx, auctionID, func(auctionID uint64, bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})
	return bids
}

// GetBidsByBidder returns all bids that are created by a bidder.
func (k Keeper) GetBidsByBidder(ctx sdk.Context, bidderAcc sdk.AccAddress) []types.Bid {
	bids := []types.Bid{}
	k.IterateBidsByBidder(ctx, bidderAcc, func(bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})
	return bids
}

// IterateBids iterates through all bids stored in the store
// and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateBids(ctx sdk.Context, auctionID uint64, cb func(auctionID uint64, bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidAuctionIDKey(auctionID))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		auctionID, _ := types.ParseBidKey(iter.Key())
		if cb(auctionID, bid) {
			break
		}
	}
}

// IterateBidsByBidder iterates through all bids by a bidder stored in the store
// and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateBidsByBidder(ctx sdk.Context, bidderAcc sdk.AccAddress, cb func(bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidIndexByBidderPrefix(bidderAcc))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		auctionID, sequence := types.ParseBidIndexKey(iter.Key())
		bid, _ := k.GetBid(ctx, auctionID, sequence)
		if cb(bid) {
			break
		}
	}
}

// PlaceBid places bid for the auction.
func (k Keeper) PlaceBid(ctx sdk.Context, msg *types.MsgPlaceBid) error {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", msg.AuctionId)
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return sdkerrors.Wrapf(types.ErrInvalidAuctionStatus, "unable to bid because the auction is in %s", auction.GetStatus().String())
	}

	// Bid amount must have greater than or equal to the amount of coin they want to bid
	requireAmt := msg.Price.Mul(msg.Coin.Amount.ToDec()).TruncateInt()
	balance := k.bankKeeper.GetBalance(ctx, msg.GetBidder(), auction.GetPayingCoinDenom())
	if balance.Amount.Sub(requireAmt).IsNegative() {
		return sdkerrors.ErrInsufficientFunds
	}

	if !auction.GetTotalSellingCoin().Sub(msg.Coin).IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "request coin must be lower than or equal to the remaining total selling coin")
	}

	if auction.GetType() == types.AuctionTypeFixedPrice {
		if !msg.Price.Equal(auction.GetStartPrice()) {
			return sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to the start price of the auction")
		}

		// Bidder cannot bid more than the total selling coin
		remaining := auction.GetTotalSellingCoin().Sub(msg.Coin)
		if err := auction.SetTotalSellingCoin(remaining); err != nil {
			return err
		}

		k.SetAuction(ctx, auction)
	}

	sequenceId := k.GetNextSequenceWithUpdate(ctx)

	bid := types.Bid{
		AuctionId: auction.GetId(),
		Sequence:  sequenceId,
		Bidder:    msg.Bidder,
		Price:     msg.Price,
		Coin:      msg.Coin,
		Height:    uint64(ctx.BlockHeader().Height),
		IsWinner:  false, // it becomes true when a bidder receives succesfully during distribution in endblocker
	}

	k.SetBid(ctx, bid.AuctionId, bid.Sequence, msg.GetBidder(), bid)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlaceBid,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyBidderAddress, msg.GetBidder().String()),
			sdk.NewAttribute(types.AttributeKeyBidPrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeKeyBidCoin, msg.Coin.String()),
		),
	})

	return nil
}
