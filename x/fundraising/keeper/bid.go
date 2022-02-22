package keeper

import (
	"fmt"
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

	if auction.GetPayingCoinDenom() != msg.Coin.Denom {
		return types.Bid{}, types.ErrInvalidPayingCoinDenom
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

	// Place a bid depending on the auction type and the bid type
	switch auction.GetType() {
	case types.AuctionTypeFixedPrice:
		if err := k.PlaceFixedPriceBid(ctx, msg, auction, maxBidAmt); err != nil {
			return types.Bid{}, err
		}
	case types.AuctionTypeBatch:
		if msg.BidType == types.BidTypeBatchWorth {
			if err := k.PlaceBatchWorthBid(ctx, msg, auction, maxBidAmt); err != nil {
				return types.Bid{}, err
			}
		} else if msg.BidType == types.BidTypeBatchMany {
			if err := k.PlaceBatchManyBid(ctx, msg, auction, maxBidAmt); err != nil {
				return types.Bid{}, err
			}
		}
	}

	bidId := k.GetNextBidIdWithUpdate(ctx, auction.GetId())
	height := uint64(ctx.BlockHeader().Height)

	bid := types.Bid{
		AuctionId: auction.GetId(),
		Id:        bidId,
		Bidder:    msg.Bidder,
		Price:     msg.Price,
		Coin:      msg.Coin,
		Height:    height,
		IsWinner:  true,
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

func (k Keeper) PlaceFixedPriceBid(ctx sdk.Context, msg *types.MsgPlaceBid, auction types.AuctionI, maxBidAmt sdk.Int) error {
	if !msg.Price.Equal(auction.GetStartPrice()) {
		return types.ErrInvalidStartPrice
	}

	// TODO: better to handle this calculation logic in types?
	receiveAmt := msg.Coin.Amount.ToDec().QuoTruncate(msg.Price).TruncateInt()
	receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, receiveAmt)

	totalBidAmt := sdk.ZeroInt()
	for _, b := range k.GetBidsByBidder(ctx, msg.GetBidder()) {
		if b.Type == types.BidTypeFixedPrice {
			totalBidAmt = totalBidAmt.Add(b.Coin.Amount)
		}
	}

	// The bidder can't bid more than the sum of total bid amount
	if totalBidAmt.GT(receiveAmt) {
		return types.ErrOverMaxBidAmountLimit
	}

	// The bidder can't bid more than the remaining selling coin
	if auction.GetRemainingSellingCoin().IsLT(receiveCoin) {
		return sdkerrors.Wrapf(types.ErrInsufficientRemainingAmount, "remaining selling coin amount %s", auction.GetRemainingSellingCoin())
	}

	remaining := auction.GetRemainingSellingCoin().Sub(receiveCoin)
	_ = auction.SetRemainingSellingCoin(remaining)

	k.SetAuction(ctx, auction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePlaceBid,
			sdk.NewAttribute(types.AttributeKeyBidAmount, receiveCoin.String()),
		),
	})

	return nil
}

func (k Keeper) PlaceBatchWorthBid(ctx sdk.Context, msg *types.MsgPlaceBid, auction types.AuctionI, maxBidAmt sdk.Int) error {

	return nil
}

func (k Keeper) PlaceBatchManyBid(ctx sdk.Context, msg *types.MsgPlaceBid, auction types.AuctionI, maxBidAmt sdk.Int) error {
	return nil
}

// ModifyBid modifies the auctioneer's bid
func (k Keeper) ModifyBid(ctx sdk.Context, msg *types.MsgModifyBid) (types.MsgModifyBid, error) {

	// TODO: not implemented yet
	// 2. bid_id must be one of the existing bids in the auction with auction_id
	// 3. bid_price must be higher of the price of bid_id and/or coin amount must be higher of the coin amount of bid_id

	auction, found := k.GetAuction(ctx, msg.AuctionId)
	if !found {
		return types.MsgModifyBid{}, sdkerrors.Wrap(sdkerrors.ErrNotFound, "auction not found")
	}

	bid := &types.Bid{}
	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		if b.Bidder == msg.Bidder {
			bid = &b
		}
	}

	// if bid == nil {
	// 	return
	// }

	fmt.Println("bid: ", bid)

	return types.MsgModifyBid{}, nil
}
