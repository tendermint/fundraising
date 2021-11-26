package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/fundraising module sentinel errors
var (
	ErrInvalidAuctionType = sdkerrors.Register(ModuleName, 2, "invalid auction type")
	ErrInvalidStartPrice  = sdkerrors.Register(ModuleName, 3, "invalid start price")
)
