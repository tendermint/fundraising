package types

import (
	time "time"

	proto "github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	SellingReserveAccPrefix string = "SellingReserveAcc"
	PayingReserveAccPrefix  string = "PayingReserveAcc"
	VestingReserveAccPrefix string = "VestingReserveAcc"
	AccNameSplitter         string = "|"

	// ReserveAddressType is an address type of reserve pool for selling or paying.
	// The module uses the address type of 32 bytes length, but it can be changed depending on Cosmos SDK's direction.
	ReserveAddressType = AddressType32Bytes
)

var (
	_ AuctionI = (*FixedPriceAuction)(nil)
	_ AuctionI = (*EnglishAuction)(nil)
)

// NewBaseAuction creates a new BaseAuction object
//nolint:interfacer
func NewBaseAuction(
	id uint64, typ AuctionType, auctioneerAddr string, sellingPoolAddr string,
	payingPoolAddr string, startPrice sdk.Dec, sellingCoin sdk.Coin,
	payingCoinDenom string, vestingAddr string, vestingSchedules []VestingSchedule,
	winningPrice sdk.Dec, totalSellingCoin sdk.Coin, startTime time.Time,
	endTimes []time.Time, status AuctionStatus,
) *BaseAuction {
	return &BaseAuction{
		Id:                 id,
		Type:               typ,
		Auctioneer:         auctioneerAddr,
		SellingPoolAddress: sellingPoolAddr,
		PayingPoolAddress:  payingPoolAddr,
		StartPrice:         startPrice,
		SellingCoin:        sellingCoin,
		PayingCoinDenom:    payingCoinDenom,
		VestingAddress:     vestingAddr,
		VestingSchedules:   vestingSchedules,
		WinningPrice:       winningPrice,
		TotalSellingCoin:   totalSellingCoin,
		StartTime:          startTime,
		EndTimes:           endTimes,
		Status:             status,
	}
}

func (ba BaseAuction) GetId() uint64 { //nolint:golint
	return ba.Id
}

func (ba *BaseAuction) SetId(id uint64) error { //nolint:golint
	ba.Id = id
	return nil
}

func (ba BaseAuction) GetType() AuctionType {
	return ba.Type
}

func (ba *BaseAuction) SetType(typ AuctionType) error {
	ba.Type = typ
	return nil
}

func (ba BaseAuction) GetAuctioneer() string {
	return ba.Auctioneer
}

func (ba *BaseAuction) SetAuctioneer(addr string) error {
	ba.Auctioneer = addr
	return nil
}

func (ba BaseAuction) GetSellingPoolAddress() string {
	return ba.SellingPoolAddress
}

func (ba *BaseAuction) SetSellingPoolAddress(addr string) error {
	ba.SellingPoolAddress = addr
	return nil
}

func (ba BaseAuction) GetPayingPoolAddress() string {
	return ba.PayingPoolAddress
}

func (ba *BaseAuction) SetPayingPoolAddress(addr string) error {
	ba.PayingPoolAddress = addr
	return nil
}

func (ba BaseAuction) GetStartPrice() sdk.Dec {
	return ba.StartPrice
}

func (ba *BaseAuction) SetStartPrice(price sdk.Dec) error {
	ba.StartPrice = price
	return nil
}

func (ba BaseAuction) GetSellingCoin() sdk.Coin {
	return ba.SellingCoin
}

func (ba *BaseAuction) SetSellingCoin(coin sdk.Coin) error {
	ba.SellingCoin = coin
	return nil
}

func (ba BaseAuction) GetPayingCoinDenom() string {
	return ba.PayingCoinDenom
}

func (ba *BaseAuction) SetPayingCoinDenom(denom string) error {
	ba.PayingCoinDenom = denom
	return nil
}

func (ba BaseAuction) GetVestingAddress() string {
	return ba.VestingAddress
}

func (ba *BaseAuction) SetVestingAddress(addr string) error {
	ba.VestingAddress = addr
	return nil
}

func (ba BaseAuction) GetVestingSchedules() []VestingSchedule {
	return ba.VestingSchedules
}

func (ba *BaseAuction) SetVestingSchedules(schedules []VestingSchedule) error {
	ba.VestingSchedules = schedules
	return nil
}

func (ba BaseAuction) GetWinningPrice() sdk.Dec {
	return ba.WinningPrice
}

func (ba *BaseAuction) SetWinningPrice(price sdk.Dec) error {
	ba.WinningPrice = price
	return nil
}

func (ba BaseAuction) GetTotalSellingCoin() sdk.Coin {
	return ba.TotalSellingCoin
}

func (ba *BaseAuction) SetTotalSellingCoin(coin sdk.Coin) error {
	ba.TotalSellingCoin = coin
	return nil
}

func (ba BaseAuction) GetStartTime() time.Time {
	return ba.StartTime
}

func (ba *BaseAuction) SetStartTime(t time.Time) error {
	ba.StartTime = t
	return nil
}

func (ba BaseAuction) GetEndTimes() []time.Time {
	return ba.EndTimes
}

func (ba *BaseAuction) SetEndTimes(t []time.Time) error {
	ba.EndTimes = t
	return nil
}

func (ba BaseAuction) GetStatus() AuctionStatus {
	return ba.Status
}

func (ba *BaseAuction) SetStatus(status AuctionStatus) error {
	ba.Status = status
	return nil
}

// Validate checks for errors on the Auction fields
func (ba BaseAuction) Validate() error {
	if ba.Type != AuctionTypeFixedPrice && ba.Type != AuctionTypeEnglish {
		return sdkerrors.Wrapf(ErrInvalidAuctionType, "unknown plan type: %s", ba.Type)
	}
	if _, err := sdk.AccAddressFromBech32(ba.Auctioneer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid auctioneer address %q: %v", ba.Auctioneer, err)
	}
	if _, err := sdk.AccAddressFromBech32(ba.SellingPoolAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid selling pool address %q: %v", ba.SellingPoolAddress, err)
	}
	if _, err := sdk.AccAddressFromBech32(ba.PayingPoolAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid paying pool address %q: %v", ba.PayingPoolAddress, err)
	}
	if _, err := sdk.AccAddressFromBech32(ba.VestingAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid vesting address %q: %v", ba.VestingAddress, err)
	}
	if !ba.StartPrice.IsPositive() {
		return sdkerrors.Wrapf(ErrInvalidStartPrice, "invalid start price: %f", ba.StartPrice)
	}
	if err := ba.SellingCoin.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid selling coin: %v", ba.SellingCoin)
	}
	// TODO: not implemented yet
	return nil
}

// NewFixedPriceAuction returns a new fixed price ba.
func NewFixedPriceAuction(baseAuction *BaseAuction) *FixedPriceAuction {
	return &FixedPriceAuction{
		BaseAuction: baseAuction,
	}
}

// NewEnglishAuction returns a new english ba.
func NewEnglishAuction(baseAuction *BaseAuction, maximumBidPrice sdk.Dec, extended uint32, extendRate sdk.Dec) *EnglishAuction {
	return &EnglishAuction{
		BaseAuction:     baseAuction,
		MaximumBidPrice: maximumBidPrice,
		Extended:        extended,
		ExtendRate:      extendRate,
	}
}

// AuctionI is an interface that inherits the BaseAuction and exposes common functions
// to get and set standard auction data.
type AuctionI interface {
	proto.Message

	GetId() uint64
	SetId(uint64) error

	GetType() AuctionType
	SetType(AuctionType) error

	GetAuctioneer() string
	SetAuctioneer(string) error

	GetSellingPoolAddress() string
	SetSellingPoolAddress(string) error

	GetPayingPoolAddress() string
	SetPayingPoolAddress(string) error

	GetStartPrice() sdk.Dec
	SetStartPrice(sdk.Dec) error

	GetSellingCoin() sdk.Coin
	SetSellingCoin(sdk.Coin) error

	GetPayingCoinDenom() string
	SetPayingCoinDenom(string) error

	GetVestingAddress() string
	SetVestingAddress(string) error

	GetVestingSchedules() []VestingSchedule
	SetVestingSchedules([]VestingSchedule) error

	GetWinningPrice() sdk.Dec
	SetWinningPrice(sdk.Dec) error

	GetTotalSellingCoin() sdk.Coin
	SetTotalSellingCoin(sdk.Coin) error

	GetStartTime() time.Time
	SetStartTime(time.Time) error

	GetEndTimes() []time.Time
	SetEndTimes([]time.Time) error

	GetStatus() AuctionStatus
	SetStatus(AuctionStatus) error
}

// SellingReserveAcc returns module account for the selling reserve pool account with the given selling coin denom.
func SellingReserveAcc(sellingCoinDenom string) sdk.AccAddress {
	return DeriveAddress(ReserveAddressType, ModuleName, SellingReserveAccPrefix+AccNameSplitter+sellingCoinDenom)
}

// PayingReserveAcc returns module account for the paying reserve pool account with the given selling coin denom.
func PayingReserveAcc(sellingCoinDenom string) sdk.AccAddress {
	return DeriveAddress(ReserveAddressType, ModuleName, PayingReserveAccPrefix+AccNameSplitter+sellingCoinDenom)
}

// VestingReserveAcc returns module account for the vesting reserve pool account with the given selling coin denom.
func VestingReserveAcc(sellingCoinDenom string) sdk.AccAddress {
	return DeriveAddress(ReserveAddressType, ModuleName, VestingReserveAccPrefix+AccNameSplitter+sellingCoinDenom)
}
