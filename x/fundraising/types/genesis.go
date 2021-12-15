package types

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		Auctions:      []codectypes.Any{},
		Bids:          []Bid{},
		VestingQueues: []VestingQueue{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	id := uint64(0)

	var auctions []AuctionI
	for _, a := range gs.Auctions {
		auction, err := UnpackAuction(&a)
		if err != nil {
			return err
		}
		if err := auction.Validate(); err != nil {
			return err
		}
		if auction.GetId() < id {
			return fmt.Errorf("auctions must be sorted")
		}
		auctions = append(auctions, auction)
		id = auction.GetId() + 1
	}

	fmt.Println("auctions: ", auctions)

	for _, b := range gs.Bids {
		if err := b.Validate(); err != nil {
			return err
		}
	}

	for _, q := range gs.VestingQueues {
		if err := q.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates Bid.
func (b Bid) Validate() error {
	// TODO: not implemented yet
	return nil
}

// Validate validates VestingQueue.
func (q VestingQueue) Validate() error {
	// TODO: not implemented yet
	return nil
}
