package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/fundraising module sentinel errors
var (
	ErrInvalidAuctionType          = sdkerrors.Register(ModuleName, 2, "invalid auction type")
	ErrInvalidStartPrice           = sdkerrors.Register(ModuleName, 3, "invalid start price")
	ErrInvalidVestingSchedules     = sdkerrors.Register(ModuleName, 4, "invalid vesting schedules")
	ErrInvalidAuctionStatus        = sdkerrors.Register(ModuleName, 5, "invalid auction status")
	ErrIncorrectAuctionType        = sdkerrors.Register(ModuleName, 6, "incorrect auction type")
	ErrIncorrectCoinDenom          = sdkerrors.Register(ModuleName, 7, "incorrect coin denom")
	ErrEmptyAllowedBidders         = sdkerrors.Register(ModuleName, 8, "empty bidders")
	ErrInvalidMaxBidAmount         = sdkerrors.Register(ModuleName, 9, "invalid maximum bid amount")
	ErrOverMaxBidAmountLimit       = sdkerrors.Register(ModuleName, 10, "over maximum bid amount limit")
	ErrInsufficientRemainingAmount = sdkerrors.Register(ModuleName, 11, "insufficient remaining amount")
	ErrNotAllowedBidder            = sdkerrors.Register(ModuleName, 12, "not allowed bidder")
	ErrIncorrectOwner              = sdkerrors.Register(ModuleName, 13, "incorrect owner")
	ErrInvalidBidId                = sdkerrors.Register(ModuleName, 14, "incorrect bid id")
	ErrInvalidModify               = sdkerrors.Register(ModuleName, 15, "bid price or coin amount cannot be lower: invalid request")
)
