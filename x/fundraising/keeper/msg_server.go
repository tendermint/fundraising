package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
func (m msgServer) CreateFixedPriceAuction(goCtx context.Context, msg *types.MsgCreateFixedPriceAuction) (*types.MsgCreateFixedPriceAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.CreateFixedPriceAuction(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateFixedPriceAuctionResponse{}, nil
}

// CreateEnglishAuction defines a method to create english auction.
func (m msgServer) CreateEnglishAuction(goCtx context.Context, msg *types.MsgCreateEnglishAuction) (*types.MsgCreateEnglishAuctionResponse, error) {
	// TODO: not implemented yet
	return &types.MsgCreateEnglishAuctionResponse{}, nil
}

// CancelAuction defines a method to cancel auction.
func (m msgServer) CancelAuction(goCtx context.Context, msg *types.MsgCancelAuction) (*types.MsgCancelAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.CancelAuction(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCancelAuctionResponse{}, nil
}

// PlaceBid defines a method to place bid for the auction.
func (m msgServer) PlaceBid(goCtx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.PlaceBid(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgPlaceBidResponse{}, nil
}

// AddAllowedBidder defines a method to add an allowed bidder.
// This message is created for testing purpose and it must not be used in mainnet.
func (m msgServer) AddAllowedBidder(goCtx context.Context, msg *types.MsgAddAllowedBidder) (*types.MsgAddAllowedBidderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if EnableAddAllowedBidder {
		if err := m.Keeper.AddAllowedBidders(ctx, msg.AuctionId, []types.AllowedBidder{msg.AllowedBidder}); err != nil {
			return nil, err
		}
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "EnableAddAllowedBidder is disabled")
	}

	return &types.MsgAddAllowedBidderResponse{}, nil
}
