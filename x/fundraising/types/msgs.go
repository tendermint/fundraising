package types

import (
	"time"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgCancelAuction(auctioneer string, auctionId uint64) *MsgCancelAuction {
	return &MsgCancelAuction{
		Auctioneer: auctioneer,
		AuctionId:  auctionId,
	}
}

func (msg MsgCancelAuction) Type() string {
	return sdk.MsgTypeURL(&MsgCancelAuction{})
}

func (msg MsgCancelAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Auctioneer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid auctioneer address: %v", err)
	}
	return nil
}

func NewMsgCreateBatchAuction(
	auctioneer string,
	startPrice math.LegacyDec,
	minBidPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate math.LegacyDec,
	startTime time.Time,
	endTime time.Time,
) *MsgCreateBatchAuction {
	return &MsgCreateBatchAuction{
		Auctioneer:        auctioneer,
		StartPrice:        startPrice,
		MinBidPrice:       minBidPrice,
		SellingCoin:       sellingCoin,
		PayingCoinDenom:   payingCoinDenom,
		VestingSchedules:  vestingSchedules,
		MaxExtendedRound:  maxExtendedRound,
		ExtendedRoundRate: extendedRoundRate,
		StartTime:         startTime,
		EndTime:           endTime,
	}
}

func (msg MsgCreateBatchAuction) Type() string {
	return sdk.MsgTypeURL(&MsgCreateBatchAuction{})
}

func (msg MsgCreateBatchAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Auctioneer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid auctioneer address: %v", err)
	}
	if !msg.StartPrice.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "start price must be positive")
	}
	if !msg.MinBidPrice.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "minimum price must be positive")
	}
	if err := msg.SellingCoin.Validate(); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid selling coin: %v", err)
	}
	if !msg.SellingCoin.Amount.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "selling coin amount must be positive")
	}
	if msg.SellingCoin.Denom == msg.PayingCoinDenom {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "selling coin denom must not be the same as paying coin denom")
	}
	if err := sdk.ValidateDenom(msg.PayingCoinDenom); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid paying coin denom: %v", err)
	}
	if !msg.EndTime.After(msg.StartTime) {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "end time must be set after start time")
	}
	if !msg.ExtendedRoundRate.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "extend rate must be positive")
	}
	return ValidateVestingSchedules(msg.VestingSchedules, msg.EndTime)
}

func NewMsgCreateFixedPriceAuction(
	auctioneer string,
	startPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	startTime time.Time,
	endTime time.Time,
) *MsgCreateFixedPriceAuction {
	return &MsgCreateFixedPriceAuction{
		Auctioneer:       auctioneer,
		StartPrice:       startPrice,
		SellingCoin:      sellingCoin,
		PayingCoinDenom:  payingCoinDenom,
		VestingSchedules: vestingSchedules,
		StartTime:        startTime,
		EndTime:          endTime,
	}
}

func (msg MsgCreateFixedPriceAuction) Type() string {
	return sdk.MsgTypeURL(&MsgCreateFixedPriceAuction{})
}

func (msg MsgCreateFixedPriceAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Auctioneer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid auctioneer address: %v", err)
	}
	if !msg.StartPrice.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "start price must be positive")
	}
	if err := msg.SellingCoin.Validate(); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid selling coin: %v", err)
	}
	if !msg.SellingCoin.Amount.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "selling coin amount must be positive")
	}
	if msg.SellingCoin.Denom == msg.PayingCoinDenom {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "selling coin denom must not be the same as paying coin denom")
	}
	if err := sdk.ValidateDenom(msg.PayingCoinDenom); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid paying coin denom: %v", err)
	}
	if !msg.EndTime.After(msg.StartTime) {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "end time must be set after start time")
	}
	return ValidateVestingSchedules(msg.VestingSchedules, msg.EndTime)
}

func NewMsgPlaceBid(
	auctionId uint64,
	bidder string,
	bidType BidType,
	price math.LegacyDec,
	coin sdk.Coin,
) *MsgPlaceBid {
	return &MsgPlaceBid{
		Bidder:    bidder,
		AuctionId: auctionId,
		BidType:   bidType,
		Price:     price,
		Coin:      coin,
	}
}

func (msg MsgPlaceBid) Type() string {
	return sdk.MsgTypeURL(&MsgPlaceBid{})
}

func (msg MsgPlaceBid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Bidder); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid bidder address: %v", err)
	}
	if !msg.Price.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "bid price must be positive value")
	}
	if err := msg.Coin.Validate(); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid bid coin: %v", err)
	}
	if !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid coin amount: %s", msg.Coin.Amount.String())
	}
	if msg.BidType != BidTypeFixedPrice && msg.BidType != BidTypeBatchWorth &&
		msg.BidType != BidTypeBatchMany {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid bid type: %T", msg.BidType.String())
	}
	return nil
}

func NewMsgAddAllowedBidder(
	auctionId uint64,
	allowedBidder AllowedBidder,
) *MsgAddAllowedBidder {
	return &MsgAddAllowedBidder{
		AuctionId:     auctionId,
		AllowedBidder: allowedBidder,
	}
}

func (msg *MsgAddAllowedBidder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.AllowedBidder.Bidder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg MsgAddAllowedBidder) Type() string {
	return sdk.MsgTypeURL(&MsgAddAllowedBidder{})
}

func (msg MsgAddAllowedBidder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.AllowedBidder.Bidder); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid bidder address: %v", err)
	}
	return nil
}

func NewMsgModifyBid(
	auctionId uint64,
	bidder string,
	bidId uint64,
	price math.LegacyDec,
	coin sdk.Coin,
) *MsgModifyBid {
	return &MsgModifyBid{
		Bidder:    bidder,
		AuctionId: auctionId,
		BidId:     bidId,
		Price:     price,
		Coin:      coin,
	}
}

func (msg MsgModifyBid) Type() string {
	return sdk.MsgTypeURL(&MsgModifyBid{})
}

func (msg MsgModifyBid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Bidder); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid bidder address: %v", err)
	}
	if !msg.Price.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "bid price must be positive value")
	}
	if err := msg.Coin.Validate(); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid bid coin: %v", err)
	}
	if !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "invalid coin amount: %s", msg.Coin.Amount.String())
	}
	return nil
}
