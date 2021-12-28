package keeper

import (
	"strconv"
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

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

// DistributeSellingCoin releases designated selling coin from the selling reserve account.
func (k Keeper) DistributeSellingCoin(ctx sdk.Context, auction types.AuctionI) error {
	sellingReserveAcc := types.SellingReserveAcc(auction.GetId())

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	totalBidCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, sdk.ZeroInt())

	// distribute coins to all bidders from the selling reserve account
	for _, bid := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		receiveAmt := bid.Coin.Amount.ToDec().QuoTruncate(bid.Price).TruncateInt()
		receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, receiveAmt)

		bidderAcc, err := sdk.AccAddressFromBech32(bid.GetBidder())
		if err != nil {
			return err
		}

		inputs = append(inputs, banktypes.NewInput(sellingReserveAcc, sdk.NewCoins(receiveCoin)))
		outputs = append(outputs, banktypes.NewOutput(bidderAcc, sdk.NewCoins(receiveCoin)))

		totalBidCoin = totalBidCoin.Add(receiveCoin)
	}

	reserveBalance := k.bankKeeper.GetBalance(ctx, sellingReserveAcc, auction.GetSellingCoin().Denom)
	remainingCoin := reserveBalance.Sub(totalBidCoin)

	// send remaining coin to the auctioneer
	inputs = append(inputs, banktypes.NewInput(sellingReserveAcc, sdk.NewCoins(remainingCoin)))
	outputs = append(outputs, banktypes.NewOutput(auction.GetAuctioneer(), sdk.NewCoins(remainingCoin)))

	// send all at once
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}

	return nil
}

// DistributePayingCoin releases the selling coin from the vesting reserve account.
func (k Keeper) DistributePayingCoin(ctx sdk.Context, auction types.AuctionI) error {
	for _, vq := range k.GetVestingQueuesByAuctionId(ctx, auction.GetId()) {
		if vq.IsVestingReleasable(ctx.BlockTime()) {
			vestingReserveAcc := auction.GetVestingReserveAddress()

			if err := k.bankKeeper.SendCoins(ctx, vestingReserveAcc, auction.GetAuctioneer(), sdk.NewCoins(vq.PayingCoin)); err != nil {
				return sdkerrors.Wrap(err, "failed to release paying coin to the auctioneer")
			}

			vq.Released = true
			k.SetVestingQueue(ctx, auction.GetId(), vq.ReleaseTime, vq)
		}
	}

	return nil
}

// ReserveSellingCoin reserves the selling coin to the selling reserve account.
func (k Keeper) ReserveSellingCoin(ctx sdk.Context, auctionId uint64, auctioneerAcc sdk.AccAddress, sellingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, auctioneerAcc, types.SellingReserveAcc(auctionId), sdk.NewCoins(sellingCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to reserve selling coin")
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

	if err := k.ReserveSellingCoin(ctx, nextId, auctioneerAcc, msg.SellingCoin); err != nil {
		return err
	}

	sellingReserveAcc := types.SellingReserveAcc(nextId)
	payingReserveAcc := types.PayingReserveAcc(nextId)
	vestingReserveAcc := types.VestingReserveAcc(nextId)

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
			sdk.NewAttribute(types.AttributeKeySellingReserveAddress, sellingReserveAcc.String()),
			sdk.NewAttribute(types.AttributeKeyPayingReserveAddress, payingReserveAcc.String()),
			sdk.NewAttribute(types.AttributeKeyVestingReserveAddress, vestingReserveAcc.String()),
			sdk.NewAttribute(types.AttributeKeySellingCoin, msg.SellingCoin.String()),
			sdk.NewAttribute(types.AttributeKeyPayingCoinDenom, msg.PayingCoinDenom),
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

	sellingReserveAcc := auction.GetSellingReserveAddress()
	reserveBalance := k.bankKeeper.GetBalance(ctx, sellingReserveAcc, auction.GetSellingCoin().Denom)

	if err := k.bankKeeper.SendCoins(ctx, sellingReserveAcc, auction.GetAuctioneer(), sdk.NewCoins(reserveBalance)); err != nil {
		return sdkerrors.Wrap(err, "failed to release selling coin")
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
