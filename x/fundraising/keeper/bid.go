package keeper

import (
	"context"
	"errors"
	"strconv"

	"cosmossdk.io/collections"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errcode "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetNextBidIdWithUpdate increments bid id by one and set it.
func (k Keeper) GetNextBidIdWithUpdate(ctx context.Context, auctionId uint64) (uint64, error) {
	seq, err := k.BidSeq.Get(ctx, auctionId)
	if errors.Is(err, collections.ErrNotFound) {
		seq = 0
	} else if err != nil {
		return 0, err
	}
	seq++
	return seq, k.BidSeq.Set(ctx, auctionId, seq)
}

// GetBidsByAuctionId returns all bids associated with the auction id that are registered in the store.
func (k Keeper) GetBidsByAuctionId(ctx context.Context, auctionId uint64) ([]types.Bid, error) {
	bids := make([]types.Bid, 0)
	rng := collections.NewPrefixedPairRange[uint64, uint64](auctionId)
	err := k.Bid.Walk(ctx, rng, func(key collections.Pair[uint64, uint64], bid types.Bid) (bool, error) {
		bids = append(bids, bid)
		return false, nil
	})
	return bids, err
}

// GetBidsByBidder returns all bids associated with the bidder that are registered in the store.
func (k Keeper) GetBidsByBidder(ctx context.Context, bidderAddr sdk.AccAddress) ([]types.Bid, error) {
	bids := make([]types.Bid, 0)
	// TODO find a way to store by bidder id to avoid read all store
	// rng := collections.NewPrefixedPairRange[uint64, uint64](bidderAddr)
	err := k.Bid.Walk(ctx, nil, func(key collections.Pair[uint64, uint64], bid types.Bid) (bool, error) {
		if bid.Bidder == bidderAddr.String() {
			bids = append(bids, bid)
		}
		return false, nil
	})
	return bids, err
}

// Bids returns all Bid.
func (k Keeper) Bids(ctx context.Context) ([]types.Bid, error) {
	bids := make([]types.Bid, 0)
	err := k.IterateBids(ctx, func(_ collections.Pair[uint64, uint64], bid types.Bid) (bool, error) {
		bids = append(bids, bid)
		return false, nil
	})
	return bids, err
}

// IterateBids iterates over all the Bids and performs a callback function.
func (k Keeper) IterateBids(ctx context.Context, cb func(collections.Pair[uint64, uint64], types.Bid) (bool, error)) error {
	err := k.Bid.Walk(ctx, nil, cb)
	if err != nil {
		return err
	}
	return nil
}

// PlaceBid places a bid for the selling coin of the auction.
func (k Keeper) PlaceBid(ctx context.Context, msg *types.MsgPlaceBid) (types.Bid, error) {
	auction, err := k.Auction.Get(ctx, msg.AuctionId)
	if err != nil {
		return types.Bid{}, err
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return types.Bid{}, types.ErrInvalidAuctionStatus
	}

	if auction.GetType() == types.AuctionTypeBatch {
		if msg.Price.LT(auction.(*types.BatchAuction).MinBidPrice) {
			return types.Bid{}, types.ErrInsufficientMinBidPrice
		}
	}

	bidder, err := sdk.AccAddressFromBech32(msg.GetBidder())
	if err != nil {
		return types.Bid{}, err
	}

	_, err = k.AllowedBidder.Get(ctx, collections.Join(auction.GetId(), bidder))
	if err != nil {
		return types.Bid{}, sdkerrors.Wrap(types.ErrNotAllowedBidder, err.Error())
	}

	if err := k.PayPlaceBidFee(ctx, bidder); err != nil {
		return types.Bid{}, sdkerrors.Wrap(err, "failed to pay place bid fee")
	}

	bidID, err := k.GetNextBidIdWithUpdate(ctx, auction.GetId())
	if err != nil {
		return types.Bid{}, sdkerrors.Wrap(err, "failed to get next bid id")
	}
	bid := types.Bid{
		AuctionId: msg.AuctionId,
		Id:        bidID,
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
		if err := k.ReservePayingCoin(ctx, msg.AuctionId, bidder, bidPayingCoin); err != nil {
			return types.Bid{}, sdkerrors.Wrap(err, "failed to reserve paying coin")
		}

		// Subtract bid amount from the remaining
		bidSellingAmt := bid.ConvertToSellingAmount(payingCoinDenom)
		bidSellingCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, bidSellingAmt)
		fa.RemainingSellingCoin = fa.RemainingSellingCoin.Sub(bidSellingCoin)

		if err := k.Auction.Set(ctx, fa.GetId(), fa); err != nil {
			return types.Bid{}, err
		}

		bid.SetMatched(true)

	case types.BidTypeBatchWorth:
		if err := k.ValidateBatchWorthBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, bidder, msg.Coin); err != nil {
			return types.Bid{}, sdkerrors.Wrap(err, "failed to reserve paying coin")
		}

	case types.BidTypeBatchMany:
		if err := k.ValidateBatchManyBid(ctx, auction, bid); err != nil {
			return types.Bid{}, err
		}

		reserveAmt := bid.ConvertToPayingAmount(payingCoinDenom)
		reserveCoin := sdk.NewCoin(payingCoinDenom, reserveAmt)

		if err := k.ReservePayingCoin(ctx, msg.AuctionId, bidder, reserveCoin); err != nil {
			return types.Bid{}, sdkerrors.Wrap(err, "failed to reserve paying coin")
		}
	}

	// Call before bid placed hook
	if err := k.BeforeBidPlaced(ctx, bid.AuctionId, bid.Id, bid.Bidder, bid.Type, bid.Price, bid.Coin); err != nil {
		return types.Bid{}, err
	}

	if err := k.Bid.Set(ctx, collections.Join(bid.AuctionId, bid.Id), bid); err != nil {
		return types.Bid{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlaceBid,
			sdk.NewAttribute(types.AttributeKeyAuctionId, strconv.FormatUint(auction.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyBidderAddress, msg.GetBidder()),
			sdk.NewAttribute(types.AttributeKeyBidPrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeKeyBidCoin, msg.Coin.String()),
		),
	})

	return bid, nil
}

// ValidateFixedPriceBid validates a fixed price bid type.
func (k Keeper) ValidateFixedPriceBid(ctx context.Context, auction types.AuctionI, bid types.Bid) error {
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
	bids, err := k.GetBidsByBidder(ctx, bid.GetBidder())
	if err != nil {
		return err
	}
	totalBidAmt := math.ZeroInt()
	for _, bid := range bids {
		if bid.AuctionId == auction.GetId() {
			bidSellingAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())
			totalBidAmt = totalBidAmt.Add(bidSellingAmt)
		}
	}

	allowedBidder, err := k.AllowedBidder.Get(ctx, collections.Join(bid.AuctionId, bid.GetBidder()))
	if err != nil {
		return sdkerrors.Wrap(err, "bidder is not found in allowed bidder list")
	}

	totalBidAmt = totalBidAmt.Add(bidAmt)

	// The total bid amount can't be greater than the bidder's maximum bid amount
	if totalBidAmt.GT(allowedBidder.MaxBidAmount) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

// ValidateBatchWorthBid validates a batch worth bid type.
func (k Keeper) ValidateBatchWorthBid(ctx context.Context, auction types.AuctionI, bid types.Bid) error {
	if auction.GetType() != types.AuctionTypeBatch {
		return types.ErrIncorrectAuctionType
	}

	if bid.Coin.Denom != auction.GetPayingCoinDenom() {
		return types.ErrIncorrectCoinDenom
	}

	allowedBidder, err := k.AllowedBidder.Get(ctx, collections.Join(bid.AuctionId, bid.GetBidder()))
	if err != nil {
		return sdkerrors.Wrap(err, "bidder is not found in allowed bidder list")
	}

	bidAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())

	// The total bid amount can't be greater than the bidder's maximum bid amount
	if bidAmt.GT(allowedBidder.MaxBidAmount) {
		return types.ErrOverMaxBidAmountLimit
	}

	return nil
}

// ValidateBatchManyBid validates a batch many bid type.
func (k Keeper) ValidateBatchManyBid(ctx context.Context, auction types.AuctionI, bid types.Bid) error {
	if auction.GetType() != types.AuctionTypeBatch {
		return types.ErrIncorrectAuctionType
	}

	if bid.Coin.Denom != auction.GetSellingCoin().Denom {
		return types.ErrIncorrectCoinDenom
	}

	allowedBidder, err := k.AllowedBidder.Get(ctx, collections.Join(bid.AuctionId, bid.GetBidder()))
	if err != nil {
		return sdkerrors.Wrap(err, "bidder is not found in allowed bidder list")
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
func (k Keeper) ModifyBid(ctx context.Context, msg *types.MsgModifyBid) error {
	auction, err := k.Auction.Get(ctx, msg.AuctionId)
	if err != nil {
		return err
	}

	if auction.GetStatus() != types.AuctionStatusStarted {
		return types.ErrInvalidAuctionStatus
	}

	if auction.GetType() != types.AuctionTypeBatch {
		return types.ErrIncorrectAuctionType
	}

	bid, err := k.Bid.Get(ctx, collections.Join(msg.AuctionId, msg.BidId))
	if err != nil {
		return err
	}

	bidder, err := sdk.AccAddressFromBech32(msg.GetBidder())
	if err != nil {
		return err
	}

	if !bid.GetBidder().Equals(bidder) {
		return sdkerrors.Wrap(errcode.ErrUnauthorized, "only the bid creator can modify the bid")
	}

	if msg.Price.LT(auction.(*types.BatchAuction).MinBidPrice) {
		return types.ErrInsufficientMinBidPrice
	}

	if bid.Coin.Denom != msg.Coin.Denom {
		return types.ErrIncorrectCoinDenom
	}

	if msg.Price.LT(bid.Price) || msg.Coin.Amount.LT(bid.Coin.Amount) {
		return sdkerrors.Wrap(errcode.ErrInvalidRequest, "bid price or coin amount cannot be lower")
	}

	if msg.Price.Equal(bid.Price) && msg.Coin.Amount.Equal(bid.Coin.Amount) {
		return sdkerrors.Wrap(errcode.ErrInvalidRequest, "bid price and coin amount must be changed")
	}

	// Reserve bid amount difference
	switch bid.Type {
	case types.BidTypeBatchWorth:
		diffReserveCoin := msg.Coin.Sub(bid.Coin)
		if diffReserveCoin.IsPositive() {
			if err := k.ReservePayingCoin(ctx, msg.AuctionId, bidder, diffReserveCoin); err != nil {
				return sdkerrors.Wrap(err, "failed to reserve paying coin")
			}
		}
	case types.BidTypeBatchMany:
		prevReserveAmt := math.LegacyNewDecFromInt(bid.Coin.Amount).Mul(bid.Price).Ceil()
		currReserveAmt := math.LegacyNewDecFromInt(msg.Coin.Amount).Mul(msg.Price).Ceil()
		diffReserveAmt := currReserveAmt.Sub(prevReserveAmt).TruncateInt()
		diffReserveCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), diffReserveAmt)
		if diffReserveCoin.IsPositive() {
			if err := k.ReservePayingCoin(ctx, msg.AuctionId, bidder, diffReserveCoin); err != nil {
				return sdkerrors.Wrap(err, "failed to reserve paying coin")
			}
		}
	}

	bid.Price = msg.Price
	bid.Coin = msg.Coin

	// Call the before mid modified hook
	if err := k.BeforeBidModified(ctx, bid.AuctionId, bid.Id, bid.Bidder, bid.Type, bid.Price, bid.Coin); err != nil {
		return err
	}

	if err := k.Bid.Set(ctx, collections.Join(bid.AuctionId, bid.Id), bid); err != nil {
		return err
	}
	return nil
}
