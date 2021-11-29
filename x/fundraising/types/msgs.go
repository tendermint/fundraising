package types

import (
	time "time"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateFixedPriceAuction)(nil)
	_ sdk.Msg = (*MsgCreateEnglishAuction)(nil)
	_ sdk.Msg = (*MsgCancelFundraising)(nil)
	_ sdk.Msg = (*MsgPlaceBid)(nil)
)

// Message types for the fundraising module
const (
	TypeMsgCreateFixedPriceAuction = "create_fixed_price_auction"
	TypeMsgCreateEnglishAuction    = "create_english_auction"
	TypeMsgCancelFundraising       = "cancel_fundraising"
	TypeMsgPlaceBid                = "place_bid"
)

// NewMsgCreateFixedPriceAuction creates a new MsgCreateFixedPriceAuction.
func NewMsgCreateFixedPriceAuction(
	auctioneer string,
	startPrice sdk.Dec,
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

func (msg MsgCreateFixedPriceAuction) Route() string { return RouterKey }

func (msg MsgCreateFixedPriceAuction) Type() string { return TypeMsgCreateFixedPriceAuction }

func (msg MsgCreateFixedPriceAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Auctioneer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid auctioneer address: %v", err)
	}
	if !msg.StartPrice.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "start price must be positve")
	}
	if err := msg.SellingCoin.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid selling coin: %v", err)
	}
	if !msg.SellingCoin.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "selling coin amount must be positive")
	}
	if err := sdk.ValidateDenom(msg.PayingCoinDenom); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid paying coin denom: %v", err)
	}
	if !msg.EndTime.After(msg.StartTime) {
		return sdkerrors.Wrapf(ErrInvalidAuctionEndTime, "end time must be greater than start time")
	}
	if err := ValidateVestingSchedules(msg.VestingSchedules); err != nil {
		return err
	}
	return nil
}

func (msg MsgCreateFixedPriceAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateFixedPriceAuction) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateFixedPriceAuction) GetAuctioneer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCreateEnglishAuction creates a new MsgCreateEnglishAuction.
func NewMsgCreateEnglishAuction(
	auctionner string,
	startPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	maximumBidPrice sdk.Dec,
	extendRate sdk.Dec,
	startTime time.Time,
	endTime time.Time,
) *MsgCreateEnglishAuction {
	return &MsgCreateEnglishAuction{
		Auctioneer:       auctionner,
		StartPrice:       startPrice,
		SellingCoin:      sellingCoin,
		PayingCoinDenom:  payingCoinDenom,
		VestingSchedules: vestingSchedules,
		MaximumBidPrice:  maximumBidPrice,
		ExtendRate:       extendRate,
		StartTime:        startTime,
		EndTime:          endTime,
	}
}

func (msg MsgCreateEnglishAuction) Route() string { return RouterKey }

func (msg MsgCreateEnglishAuction) Type() string { return TypeMsgCreateEnglishAuction }

func (msg MsgCreateEnglishAuction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Auctioneer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid auctioneer address: %v", err)
	}
	if !msg.StartPrice.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "start price must be positve")
	}
	if err := msg.SellingCoin.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid selling coin: %v", err)
	}
	if !msg.SellingCoin.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "selling coin amount must be positive")
	}
	if err := sdk.ValidateDenom(msg.PayingCoinDenom); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid paying coin denom: %v", err)
	}
	if !msg.EndTime.After(msg.StartTime) {
		return sdkerrors.Wrapf(ErrInvalidAuctionEndTime, "end time must be greater than start time")
	}
	if err := ValidateVestingSchedules(msg.VestingSchedules); err != nil {
		return err
	}
	if !msg.MaximumBidPrice.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "maximum bid price must be positve")
	}
	if !msg.ExtendRate.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "extend rate must be positve")
	}
	return nil
}

func (msg MsgCreateEnglishAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateEnglishAuction) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateEnglishAuction) GetAuctioneer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCancelFundraising creates a new MsgCancelFundraising.
func NewMsgCancelFundraising(
	auctionner string,
	auctionId uint64,
) *MsgCancelFundraising {
	return &MsgCancelFundraising{
		Auctioneer: auctionner,
		AuctionId:  auctionId,
	}
}

func (msg MsgCancelFundraising) Route() string { return RouterKey }

func (msg MsgCancelFundraising) Type() string { return TypeMsgCancelFundraising }

func (msg MsgCancelFundraising) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Auctioneer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid auctioneer address %q: %v", msg.Auctioneer, err)
	}
	return nil
}

func (msg MsgCancelFundraising) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgCancelFundraising) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCancelFundraising) GetAuctioneer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Auctioneer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgPlaceBid creates a new MsgPlaceBid.
func NewMsgPlaceBid(
	id uint64,
	bidder string,
	price sdk.Dec,
	coin sdk.Coin,
) *MsgPlaceBid {
	return &MsgPlaceBid{
		AuctionId: id,
		Bidder:    bidder,
		Price:     price,
		Coin:      coin,
	}
}

func (msg MsgPlaceBid) Route() string { return RouterKey }

func (msg MsgPlaceBid) Type() string { return TypeMsgPlaceBid }

func (msg MsgPlaceBid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Bidder); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid bidder address: %v", err)
	}
	if !msg.Price.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "bid price must be positve value")
	}
	if err := msg.Coin.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid bid coin: %v", err)
	}
	return nil
}

func (msg MsgPlaceBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(&msg))
}

func (msg MsgPlaceBid) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgPlaceBid) GetBidder() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		panic(err)
	}
	return addr
}
