package keeper

import (
	"strconv"
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetAuctionId returns the global auction ID counter.
func (k Keeper) GetAuctionId(ctx sdk.Context) uint64 {
	var id uint64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AuctionIdKey)
	if bz == nil {
		id = 0 // initialize the auction id
	} else {
		val := gogotypes.UInt64Value{}
		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}
		id = val.GetValue()
	}
	return id
}

// GetNextAuctionIdWithUpdate increments auction id by one and set it.
func (k Keeper) GetNextAuctionIdWithUpdate(ctx sdk.Context) uint64 {
	id := k.GetAuctionId(ctx) + 1
	k.SetAuctionId(ctx, id)
	return id
}

// SetAuctionId sets the global auction ID counter.
func (k Keeper) SetAuctionId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.AuctionIdKey, bz)
}

// GetAuction returns an auction for a given auction id.
func (k Keeper) GetAuction(ctx sdk.Context, id uint64) (auction types.AuctionI, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAuctionKey(id))
	if bz == nil {
		return auction, false
	}
	return k.decodeAuction(bz), true
}

// GetAuctions returns all auctions in the store.
func (k Keeper) GetAuctions(ctx sdk.Context) (auctions []types.AuctionI) {
	k.IterateAuctions(ctx, func(auction types.AuctionI) (stop bool) {
		auctions = append(auctions, auction)
		return false
	})

	return auctions
}

// SetAuction sets an auction with the given auction id.
func (k Keeper) SetAuction(ctx sdk.Context, auction types.AuctionI) {
	id := auction.GetId()
	store := ctx.KVStore(k.storeKey)

	bz, err := k.MarshalAuction(auction)
	if err != nil {
		panic(err)
	}

	store.Set(types.GetAuctionKey(id), bz)
}

// RemoveAuction removes the auction from the store
func (k Keeper) RemoveAuction(ctx sdk.Context, auction types.AuctionI) {
	id := auction.GetId()
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAuctionKey(id))
}

// IterateAuctions iterates over all the stored auctions and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAuctions(ctx sdk.Context, cb func(auction types.AuctionI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AuctionKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		auction := k.decodeAuction(iterator.Value())

		if cb(auction) {
			break
		}
	}
}

func (k Keeper) decodeAuction(bz []byte) types.AuctionI {
	acc, err := k.UnmarshalAuction(bz)
	if err != nil {
		panic(err)
	}

	return acc
}

// MarshalAuction serializes an auction.
func (k Keeper) MarshalAuction(auction types.AuctionI) ([]byte, error) { // nolint:interfacer
	return k.cdc.MarshalInterface(auction)
}

// UnmarshalAuction returns an auction from raw serialized
// bytes of a Proto-based Auction type.
func (k Keeper) UnmarshalAuction(bz []byte) (auction types.AuctionI, err error) {
	return auction, k.cdc.UnmarshalInterface(bz, &auction)
}

// CreateFixedPriceAuction sets fixed price auction.
func (k Keeper) CreateFixedPriceAuction(ctx sdk.Context, msg *types.MsgCreateFixedPriceAuction) error {
	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	auctioneerAcc, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		return err
	}

	// escrow the selling coin to the selling reserve account
	sellingReserveAcc := types.SellingReserveAcc(msg.SellingCoin.Denom)
	if err := k.bankKeeper.SendCoins(ctx, auctioneerAcc, sellingReserveAcc, sdk.NewCoins(msg.SellingCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow selling coin to selling reserve account")
	}

	payingReserveAcc := types.PayingReserveAcc(msg.SellingCoin.Denom)
	vestingReserveAcc := types.VestingReserveAcc(msg.SellingCoin.Denom)

	baseAuction := types.NewBaseAuction(
		nextId,
		types.AuctionTypeFixedPrice,
		auctioneerAcc.String(),
		sellingReserveAcc.String(),
		payingReserveAcc.String(),
		msg.StartPrice,
		msg.SellingCoin,
		msg.PayingCoinDenom,
		vestingReserveAcc.String(),
		msg.VestingSchedules,
		sdk.ZeroDec(),
		msg.SellingCoin, // add selling coin to total selling coin
		msg.StartTime,
		[]time.Time{msg.EndTime},
		types.AuctionStatusStandBy,
	)

	fixedPriceAuction := types.NewFixedPriceAuction(baseAuction)

	k.SetAuction(ctx, fixedPriceAuction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateFixedPriceAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyAuctioneerAddress, auctioneerAcc.String()),
			sdk.NewAttribute(types.AttributeKeyStartPrice, msg.StartPrice.String()),
			sdk.NewAttribute(types.AttributeKeySellingPoolAddress, sellingReserveAcc.String()),
			sdk.NewAttribute(types.AttributeKeyPayingPoolAddress, payingReserveAcc.String()),
			sdk.NewAttribute(types.AttributeKeyVestingPoolAddress, vestingReserveAcc.String()),
			sdk.NewAttribute(types.AttributeKeySellingCoin, msg.SellingCoin.String()),
			sdk.NewAttribute(types.AttributeKeyPayingCoinDenom, msg.PayingCoinDenom),
			// sdk.NewAttribute(types.AttributeKeyVestingSchedules, msg.VestingSchedules), // TODO: stringtify
			sdk.NewAttribute(types.AttributeKeyStartTime, msg.StartTime.String()),
			sdk.NewAttribute(types.AttributeKeyEndTime, msg.EndTime.String()),
			sdk.NewAttribute(types.AttributeKeyAuctionStatus, types.AuctionStatusStandBy.String()),
		),
	})

	return nil
}

// CancelAuction cancels the auction in an event of modification for the auction.
// The auctioneer can only delete it when it is not already started.
func (k Keeper) CancelAuction(ctx sdk.Context, id uint64) error {
	auction, found := k.GetAuction(ctx, id)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", id)
	}

	if auction.GetStatus() != types.AuctionStatusStandBy {
		return sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "invalid auction status")
	}

	// TODO: consider if we want the auction to be deleted or leave history
	k.RemoveAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
		),
	})

	return nil
}

// PlaceBid places bid for the auction.
func (k Keeper) PlaceBid(ctx sdk.Context, msg *types.MsgPlaceBid) error {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", msg.AuctionId)
	}

	if !msg.Price.Equal(auction.GetStartPrice()) {
		return sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to start price")
	}

	// substract total selling coin from the request amount of coin
	auction.SetTotalSellingCoin(auction.GetTotalSellingCoin().Sub(msg.Coin))

	k.SetAuction(ctx, auction)

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
