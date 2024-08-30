package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) AddAllowedBidder(ctx context.Context, msg *types.MsgAddAllowedBidder) (*types.MsgAddAllowedBidderResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.AllowedBidder.Bidder); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if !EnableAddAllowedBidder {
		return nil, sdkerrors.Wrap(errors.ErrInvalidRequest, "EnableAddAllowedBidder is disabled")
	}

	if err := k.Keeper.AddAllowedBidders(ctx, msg.AuctionId, []types.AllowedBidder{msg.AllowedBidder}); err != nil {
		return nil, err
	}

	return &types.MsgAddAllowedBidderResponse{}, nil
}

func (k msgServer) CancelAuction(ctx context.Context, msg *types.MsgCancelAuction) (*types.MsgCancelAuctionResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Auctioneer); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if err := k.Keeper.CancelAuction(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCancelAuctionResponse{}, nil
}

func (k msgServer) CreateBatchAuction(ctx context.Context, msg *types.MsgCreateBatchAuction) (*types.MsgCreateBatchAuctionResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Auctioneer); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if _, err := k.Keeper.CreateBatchAuction(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateBatchAuctionResponse{}, nil
}

func (k msgServer) CreateFixedPriceAuction(ctx context.Context, msg *types.MsgCreateFixedPriceAuction) (*types.MsgCreateFixedPriceAuctionResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Auctioneer); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if _, err := k.Keeper.CreateFixedPriceAuction(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateFixedPriceAuctionResponse{}, nil
}

func (k msgServer) ModifyBid(ctx context.Context, msg *types.MsgModifyBid) (*types.MsgModifyBidResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Bidder); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if err := k.Keeper.ModifyBid(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgModifyBidResponse{}, nil
}

func (k msgServer) PlaceBid(ctx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Bidder); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid authority address")
	}

	if _, err := k.Keeper.PlaceBid(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgPlaceBidResponse{}, nil
}
