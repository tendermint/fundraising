package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

// CreateFixedPriceAuction defines a method to create fixed price auction.
func (k msgServer) CreateFixedPriceAuction(goCtx context.Context, msg *types.MsgCreateFixedPriceAuction) (*types.MsgCreateFixedPriceAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.CreateFixedPriceAuction(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateFixedPriceAuctionResponse{}, nil
}

// CreateEnglishAuction defines a method to create english auction.
func (k msgServer) CreateEnglishAuction(goCtx context.Context, msg *types.MsgCreateEnglishAuction) (*types.MsgCreateEnglishAuctionResponse, error) {
	// TODO: not implemented yet
	return &types.MsgCreateEnglishAuctionResponse{}, nil
}

// CancelAuction defines a method to cancel auction.
func (k msgServer) CancelAuction(goCtx context.Context, msg *types.MsgCancelAuction) (*types.MsgCancelAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.CancelAuction(ctx, msg.AuctionId); err != nil {
		return nil, err
	}

	return &types.MsgCancelAuctionResponse{}, nil
}

// PlaceBid defines a method to place bid for the auction.
func (k msgServer) PlaceBid(goCtx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.PlaceBid(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgPlaceBidResponse{}, nil
}
