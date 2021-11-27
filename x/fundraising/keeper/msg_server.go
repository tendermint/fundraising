package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"

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
	// TODO: not implemented yet
	return &types.MsgCreateFixedPriceAuctionResponse{}, nil
}

// CreateEnglishAuction defines a method to create english auction.
func (k msgServer) CreateEnglishAuction(goCtx context.Context, msg *types.MsgCreateEnglishAuction) (*types.MsgCreateEnglishAuctionResponse, error) {
	// TODO: not implemented yet
	return &types.MsgCreateEnglishAuctionResponse{}, nil
}

// CancelFundraising defines a method to cancel fundraising.
func (k msgServer) CancelFundraising(goCtx context.Context, msg *types.MsgCancelFundraising) (*types.MsgCancelFundraisingResponse, error) {
	// TODO: not implemented yet
	return &types.MsgCancelFundraisingResponse{}, nil
}

// PlaceBid defines a method to cancel fundraising.
func (k msgServer) PlaceBid(goCtx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	// TODO: not implemented yet
	return &types.MsgPlaceBidResponse{}, nil
}
