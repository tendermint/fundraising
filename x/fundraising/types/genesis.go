package types

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AllowedBidderList: []AllowedBidder{},
		VestingQueueList:  []VestingQueue{},
		BidList:           []Bid{},
		AuctionList:       []*codectypes.Any{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in allowedBidder
	allowedBidderIndexMap := make(map[string]struct{})

	for _, elem := range gs.AllowedBidderList {
		index := fmt.Sprint(elem.AuctionId)
		if _, ok := allowedBidderIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for allowedBidder")
		}
		allowedBidderIndexMap[index] = struct{}{}

		if err := elem.Validate(); err != nil {
			return err
		}
	}
	// Check for duplicated index in vestingQueue
	vestingQueueIndexMap := make(map[string]struct{})

	for _, elem := range gs.VestingQueueList {
		index := fmt.Sprint(elem.AuctionId)
		if _, ok := vestingQueueIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for vestingQueue")
		}
		vestingQueueIndexMap[index] = struct{}{}

		if err := elem.Validate(); err != nil {
			return err
		}
	}
	// Check for duplicated ID in bid
	bidIdMap := make(map[uint64]bool)
	for _, elem := range gs.BidList {
		if _, ok := bidIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for bid")
		}
		bidIdMap[elem.Id] = true

		if err := elem.Validate(); err != nil {
			return err
		}
	}

	// Check for duplicated ID in auction
	auctionIdMap := make(map[uint64]bool)
	for _, elem := range gs.AuctionList {
		auction, err := UnpackAuction(elem)
		if err != nil {
			return err
		}
		if _, ok := auctionIdMap[auction.GetId()]; ok {
			return fmt.Errorf("duplicated id for auction")
		}
		auctionIdMap[auction.GetId()] = true

		if err := auction.Validate(); err != nil {
			return err
		}
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}

// Validate validates Bid.
func (b Bid) Validate() error {
	if _, err := sdk.AccAddressFromBech32(b.Bidder); err != nil {
		return err
	}
	if !b.Price.IsPositive() {
		return fmt.Errorf("bid price must be positive value: %s", b.Price.String())
	}
	if err := b.Coin.Validate(); err != nil {
		return err
	}
	if !b.Coin.Amount.IsPositive() {
		return fmt.Errorf("coin amount must be positive: %s", b.Coin.Amount.String())
	}
	return nil
}

// Validate validates VestingQueue.
func (q VestingQueue) Validate() error {
	if _, err := sdk.AccAddressFromBech32(q.Auctioneer); err != nil {
		return err
	}
	if err := q.PayingCoin.Validate(); err != nil {
		return fmt.Errorf("paying coin is invalid: %v", err)
	}
	return nil
}
