package types

import (
	"fmt"
	time "time"

	proto "github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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

	GetRemainingCoin() sdk.Coin
	SetRemainingCoin(sdk.Coin) error

	GetStartTime() time.Time
	SetStartTime(time.Time) error

	GetEndTimes() []time.Time
	SetEndTimes([]time.Time) error

	GetStatus() AuctionStatus
	SetStatus(AuctionStatus) error
}

// NewBaseAuction creates a new BaseAuction object
//nolint:interfacer
func NewBaseAuction(
	id uint64, typ AuctionType, auctioneerAddr string, sellingPoolAddr string,
	payingPoolAddr string, startPrice sdk.Dec, sellingCoin sdk.Coin,
	payingCoinDenom string, vestingAddr string, vestingSchedules []VestingSchedule,
	winningPrice sdk.Dec, remainingCoin sdk.Coin, startTime time.Time,
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
		RemainingCoin:      remainingCoin,
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

func (ba BaseAuction) GetRemainingCoin() sdk.Coin {
	return ba.RemainingCoin
}

func (ba *BaseAuction) SetRemainingCoin(coin sdk.Coin) error {
	ba.RemainingCoin = coin
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

// UnmarshalBid unmarshals bid from a store value.
func UnmarshalBid(cdc codec.BinaryCodec, value []byte) (b Bid, err error) {
	err = cdc.Unmarshal(value, &b)
	return b, err
}

// PackAuction converts AuctionI to Any.
func PackAuction(auction AuctionI) (*codectypes.Any, error) {
	any, err := codectypes.NewAnyWithValue(auction)
	if err != nil {
		return nil, err
	}
	return any, nil
}

// UnpackAuction converts Any to AuctionI.
func UnpackAuction(any *codectypes.Any) (AuctionI, error) {
	if any == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot unpack nil")
	}

	if any.TypeUrl == "" {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "empty type url")
	}

	var auction AuctionI
	v := any.GetCachedValue()
	if v == nil {
		registry := codectypes.NewInterfaceRegistry()
		RegisterInterfaces(registry)
		if err := registry.UnpackAny(any, &auction); err != nil {
			return nil, err
		}
		return auction, nil
	}

	auction, ok := v.(AuctionI)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "cannot unpack auction from %T", v)
	}

	return auction, nil
}

// UnpackAuctions converts Any slice to AuctionIs.
func UnpackAuctions(auctionsAny []*codectypes.Any) ([]AuctionI, error) {
	auctions := make([]AuctionI, len(auctionsAny))
	for i, any := range auctionsAny {
		p, err := UnpackAuction(any)
		if err != nil {
			return nil, err
		}
		auctions[i] = p
	}
	return auctions, nil
}

// SellingReserveAcc returns module account for the selling reserve pool account with the given auction id.
func SellingReserveAcc(auctionId uint64) sdk.AccAddress {
	return DeriveAddress(ReserveAddressType, ModuleName, SellingReserveAccPrefix+AccNameSplitter+fmt.Sprint(auctionId))
}

// PayingReserveAcc returns module account for the paying reserve pool account with the given auction id.
func PayingReserveAcc(auctionId uint64) sdk.AccAddress {
	return DeriveAddress(ReserveAddressType, ModuleName, PayingReserveAccPrefix+AccNameSplitter+fmt.Sprint(auctionId))
}

// VestingReserveAcc returns module account for the vesting reserve pool account with the given auction id.
func VestingReserveAcc(auctionId uint64) sdk.AccAddress {
	return DeriveAddress(ReserveAddressType, ModuleName, VestingReserveAccPrefix+AccNameSplitter+fmt.Sprint(auctionId))
}

// IsAuctionStarted returns true if the start time of the auction is passed over the given time t.
func IsAuctionStarted(startTime time.Time, t time.Time) bool {
	return !startTime.After(t)
}

// IsAuctionFinished returns true if the end time of the auction is passed over the given time t.
func IsAuctionFinished(endTime time.Time, t time.Time) bool {
	return !endTime.Before(t)
}
