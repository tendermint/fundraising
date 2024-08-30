package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/fundraising module sentinel errors
var (
	ErrInvalidSigner               = sdkerrors.Register(ModuleName, 1101, "expected gov account as only signer for proposal message")
	ErrInvalidAuctionType          = sdkerrors.Register(ModuleName, 1102, "invalid auction type")
	ErrInvalidStartPrice           = sdkerrors.Register(ModuleName, 1103, "invalid start price")
	ErrInvalidVestingSchedules     = sdkerrors.Register(ModuleName, 1104, "invalid vesting schedules")
	ErrInvalidAuctionStatus        = sdkerrors.Register(ModuleName, 1105, "invalid auction status")
	ErrInvalidMaxBidAmount         = sdkerrors.Register(ModuleName, 1106, "invalid maximum bid amount")
	ErrIncorrectAuctionType        = sdkerrors.Register(ModuleName, 1107, "incorrect auction type")
	ErrIncorrectCoinDenom          = sdkerrors.Register(ModuleName, 1108, "incorrect coin denom")
	ErrEmptyAllowedBidders         = sdkerrors.Register(ModuleName, 1109, "empty bidders")
	ErrNotAllowedBidder            = sdkerrors.Register(ModuleName, 1110, "not allowed bidder")
	ErrOverMaxBidAmountLimit       = sdkerrors.Register(ModuleName, 1111, "over maximum bid amount limit")
	ErrInsufficientRemainingAmount = sdkerrors.Register(ModuleName, 1112, "insufficient remaining amount")
	ErrInsufficientMinBidPrice     = sdkerrors.Register(ModuleName, 1113, "insufficient bid price")
)
