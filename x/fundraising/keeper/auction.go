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

// GetSequence returns the last sequence number of the bid.
func (k Keeper) GetSequence(ctx sdk.Context) uint64 {
	var seq uint64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SequenceKey)
	if bz == nil {
		seq = 0 // initialize the sequence
	} else {
		val := gogotypes.UInt64Value{}
		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}
		seq = val.GetValue()
	}
	return seq
}

// GetNextSequence increments sequence number by one and set it.
func (k Keeper) GetNextSequenceWithUpdate(ctx sdk.Context) uint64 {
	seq := k.GetSequence(ctx) + 1
	k.SetSequence(ctx, seq)
	return seq
}

// SetSequence sets the sequence number of the bid.
func (k Keeper) SetSequence(ctx sdk.Context, seq uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: seq})
	store.Set(types.SequenceKey, bz)
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

// DistributeSellingCoin releases designated selling coin from the selling reserve module account.
func (k Keeper) DistributeSellingCoin(ctx sdk.Context, auction types.AuctionI) error {
	sellingReserveAcc := types.SellingReserveAcc(auction.GetId())

	// distribute coin to all bidders who bid for the auction
	for _, bid := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		bidAmt := bid.Coin.Amount.ToDec().Quo(bid.Price).TruncateInt()
		bidCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, bidAmt)

		bidderAcc, err := sdk.AccAddressFromBech32(bid.GetBidder())
		if err != nil {
			return err
		}

		if err := k.bankKeeper.SendCoins(ctx, sellingReserveAcc, bidderAcc, sdk.NewCoins(bidCoin)); err != nil {
			return err
		}
	}

	reserveBalance := k.bankKeeper.GetBalance(ctx, sellingReserveAcc, auction.GetSellingCoin().Denom)

	// Send remaining selling coin to the auctioneer
	if err := k.bankKeeper.SendCoins(ctx, sellingReserveAcc, auction.GetAuctioneer(), sdk.NewCoins(reserveBalance)); err != nil {
		return err
	}

	return nil
}

// DistributePayingCoin releases vested selling coin from the vesting reserve module account.
func (k Keeper) DistributePayingCoin(ctx sdk.Context, auction types.AuctionI) error {
	for _, vq := range k.GetVestingQueuesByAuctionId(ctx, auction.GetId()) {
		if types.IsVested(vq.GetReleaseTime(), ctx.BlockTime()) {
			vestingReserveAcc := types.VestingReserveAcc(auction.GetId())
			reserveBalance := k.bankKeeper.GetBalance(ctx, vestingReserveAcc, auction.GetPayingCoinDenom())

			if err := k.bankKeeper.SendCoins(ctx, vestingReserveAcc, auction.GetAuctioneer(), sdk.NewCoins(reserveBalance)); err != nil {
				return err
			}

			vq.Vested = true
			k.SetVestingQueue(ctx, auction.GetId(), vq.ReleaseTime, vq)
		}
	}

	return nil
}

// CreateFixedPriceAuction sets fixed price auction.
func (k Keeper) CreateFixedPriceAuction(ctx sdk.Context, msg *types.MsgCreateFixedPriceAuction) error {
	nextId := k.GetNextAuctionIdWithUpdate(ctx)

	auctioneerAcc, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		return err
	}

	if ctx.BlockTime().After(msg.EndTime) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "end time must be prior to current time")
	}

	params := k.GetParams(ctx)
	auctionFeeCollectorAcc, err := sdk.AccAddressFromBech32(params.AuctionFeeCollector)
	if err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoins(ctx, msg.GetAuctioneer(), auctionFeeCollectorAcc, params.AuctionCreationFee); err != nil {
		return sdkerrors.Wrap(err, "failed to pay auction creation fee")
	}

	sellingReserveAcc := types.SellingReserveAcc(nextId)
	payingReserveAcc := types.PayingReserveAcc(nextId)
	vestingReserveAcc := types.VestingReserveAcc(nextId)

	// reserve the selling coin to the selling reserve account
	if err := k.bankKeeper.SendCoins(ctx, auctioneerAcc, sellingReserveAcc, sdk.NewCoins(msg.SellingCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow selling coin to selling reserve account")
	}

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
		msg.SellingCoin, // add selling coin to remaining coin
		msg.StartTime,
		[]time.Time{msg.EndTime},
		types.AuctionStatusStandBy,
	)

	// updates status if the start time is already passed over the current time
	if types.IsAuctionStarted(baseAuction.GetStartTime(), ctx.BlockTime()) {
		baseAuction.Status = types.AuctionStatusStarted
	}

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

// CreateEnglishAuction sets english auction.
func (k Keeper) CreateEnglishAuction(ctx sdk.Context, msg *types.MsgCreateEnglishAuction) error {
	// TODO: not implemented yet
	return nil
}

// CancelAuction cancels the auction in an event when the auctioneer needs to modify the auction.
// However, it can only be canceled when the auction has not started yet.
func (k Keeper) CancelAuction(ctx sdk.Context, msg *types.MsgCancelAuction) error {
	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", msg.AuctionId)
	}

	if auction.GetAuctioneer().String() != msg.Auctioneer {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify ownership of the auction")
	}

	if auction.GetStatus() != types.AuctionStatusStandBy {
		return sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "auction cannot be canceled due to current status")
	}

	if err := auction.SetStatus(types.AuctionStatusCancelled); err != nil {
		return err
	}

	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAuction,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
		),
	})

	return nil
}
