package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// FixedPriceAuctionRequest defines CLI request for a fixed price auction.
type FixedPriceAuctionRequest struct {
	StartPrice       sdk.Dec                 `json:"start_price"`
	SellingCoin      sdk.Coin                `json:"selling_coin"`
	PayingCoinDenom  string                  `json:"paying_coin_denom"`
	VestingSchedules []types.VestingSchedule `json:"vesting_schedules"`
	StartTime        time.Time               `json:"start_time"`
	EndTime          time.Time               `json:"end_time"`
}

// ParseFixedPriceAuctionRequest reads the file and parses FixedPriceAuctionRequest.
func ParseFixedPriceAuctionRequest(fileName string) (req FixedPriceAuctionRequest, err error) {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return req, err
	}

	if err = json.Unmarshal(contents, &req); err != nil {
		return req, err
	}

	return req, nil
}

// String returns a human readable string representation of the request.
func (req FixedPriceAuctionRequest) String() string {
	result, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}
	return string(result)
}

// BatchAuctionRequest defines CLI request for an batch auction.
type BatchAuctionRequest struct {
	StartPrice        sdk.Dec                 `json:"start_price"`
	MinBidPrice       sdk.Dec                 `json:"min_bid_price"`
	SellingCoin       sdk.Coin                `json:"selling_coin"`
	PayingCoinDenom   string                  `json:"paying_coin_denom"`
	MaxExtendedRound  uint32                  `json:"max_extended_round"`
	ExtendedRoundRate sdk.Dec                 `json:"extended_round_rate"`
	VestingSchedules  []types.VestingSchedule `json:"vesting_schedules"`
	StartTime         time.Time               `json:"start_time"`
	EndTime           time.Time               `json:"end_time"`
}

// ParseBatchAuctionRequest reads the file and parses BatchAuctionRequest.
func ParseBatchAuctionRequest(fileName string) (req BatchAuctionRequest, err error) {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return req, err
	}

	if err = json.Unmarshal(contents, &req); err != nil {
		return req, err
	}

	return req, nil
}

// String returns a human readable string representation of the request.
func (req BatchAuctionRequest) String() string {
	result, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}
	return string(result)
}

// parseBidType parses bid type string and returns types.BidType.
func parseBidType(s string) (types.BidType, error) {
	switch strings.ToLower(s) {
	case "fixed-price", "fp", "f":
		return types.BidTypeFixedPrice, nil
	case "batch-worth", "bw", "w":
		return types.BidTypeBatchWorth, nil
	case "batch-many", "bm", "m":
		return types.BidTypeBatchMany, nil
	}
	return 0, fmt.Errorf("invalid bid type: %s", s)
}
