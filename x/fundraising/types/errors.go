package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/fundraising module sentinel errors
var (
	ErrInvalidAuctionEndTime = sdkerrors.Register(ModuleName, 4, "invalid auction end time")
)
