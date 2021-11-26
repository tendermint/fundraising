package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
)

var (
	_ AuctionI = (*FixedPriceAuction)(nil)
	_ AuctionI = (*EnglishAuction)(nil)
)

// NewBaseAuction creates a new BaseAuction object
//nolint:interfacer
func NewBaseAuction(
	id uint64, name string, typ AuctionType, auctioneerAddr string,
	sellingPoolAddr string, payingPoolAddr string, startPrice sdk.Dec, sellingCoin sdk.Coin,
	payingCoinDenom string, vestingAddr string, vestingSchedules []VestingSchedule,
	startTime time.Time, endTimes []time.Time, status AuctionStatus,
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
		StartTime:          startTime,
		EndTimes:           endTimes,
		Status:             status,
	}
}

func (auction BaseAuction) GetId() uint64 { //nolint:golint
	return auction.Id
}

func (auction *BaseAuction) SetId(id uint64) error { //nolint:golint
	auction.Id = id
	return nil
}

func (auction BaseAuction) GetType() AuctionType {
	return auction.Type
}

func (auction *BaseAuction) SetType(typ AuctionType) error {
	auction.Type = typ
	return nil
}

func (auction BaseAuction) GetAuctioneer() string {
	return auction.Auctioneer
}

func (auction *BaseAuction) SetAuctioneer(addr string) error {
	auction.Auctioneer = addr
	return nil
}

func (auction BaseAuction) GetSellingPoolAddress() string {
	return auction.SellingPoolAddress
}

func (auction *BaseAuction) SetSellingPoolAddress(addr string) error {
	auction.SellingPoolAddress = addr
	return nil
}

func (auction BaseAuction) GetPayingPoolAddress() string {
	return auction.PayingPoolAddress
}

func (auction *BaseAuction) SetPayingPoolAddress(addr string) error {
	auction.PayingPoolAddress = addr
	return nil
}

func (auction BaseAuction) GetStartPrice() sdk.Dec {
	return auction.StartPrice
}

func (auction *BaseAuction) SetStartPrice(price sdk.Dec) error {
	auction.StartPrice = price
	return nil
}

func (auction BaseAuction) GetSellingCoin() sdk.Coin {
	return auction.SellingCoin
}

func (auction *BaseAuction) SetSellingCoin(coin sdk.Coin) error {
	auction.SellingCoin = coin
	return nil
}

func (auction BaseAuction) GetPayingCoinDenom() string {
	return auction.PayingCoinDenom
}

func (auction *BaseAuction) SetPayingCoinDenom(denom string) error {
	auction.PayingCoinDenom = denom
	return nil
}

func (auction BaseAuction) GetVestingAddress() string {
	return auction.VestingAddress
}

func (auction *BaseAuction) SetVestingAddress(addr string) error {
	auction.VestingAddress = addr
	return nil
}

func (auction BaseAuction) GetVestingSchedules() []VestingSchedule {
	return auction.VestingSchedules
}

func (auction *BaseAuction) SetVestingSchedules(schedules []VestingSchedule) error {
	auction.VestingSchedules = schedules
	return nil
}

func (auction BaseAuction) GetStartTime() time.Time {
	return auction.StartTime
}

func (auction *BaseAuction) SetStartTime(t time.Time) error {
	auction.StartTime = t
	return nil
}

func (auction BaseAuction) GetEndTimes() []time.Time {
	return auction.EndTimes
}

func (auction *BaseAuction) SetEndTimes(t []time.Time) error {
	auction.EndTimes = t
	return nil
}

func (auction BaseAuction) GetStatus() AuctionStatus {
	return auction.Status
}

func (auction *BaseAuction) SetStatus(status AuctionStatus) error {
	auction.Status = status
	return nil
}

// Validate checks for errors on the Auction fields
func (auction BaseAuction) Validate() error {
	// TODO: not implemented yet
	return nil
}

// NewFixedPriceAuction returns a new fixed price auction.
func NewFixedPriceAuction(baseAuction *BaseAuction) *FixedPriceAuction {
	return &FixedPriceAuction{
		BaseAuction: baseAuction,
	}
}

// NewEnglishAuction returns a new english auction.
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

	GetStartTime() time.Time
	SetStartTime(time.Time) error

	GetEndTimes() []time.Time
	SetEndTimes([]time.Time) error

	GetStatus() AuctionStatus
	SetStatus(AuctionStatus) error
}
