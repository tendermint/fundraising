package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestGenesisState_Validate(t *testing.T) {
	validAddr := sdk.AccAddress(crypto.AddressHash([]byte("validAddr")))
	validAuction := types.NewFixedPriceAuction(
		&types.BaseAuction{
			Id:                    1,
			Type:                  types.AuctionTypeFixedPrice,
			Auctioneer:            validAddr.String(),
			SellingReserveAddress: types.SellingReserveAddress(1).String(),
			PayingReserveAddress:  types.PayingReserveAddress(1).String(),
			StartPrice:            math.LegacyMustNewDecFromStr("0.5"),
			SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
			PayingCoinDenom:       "denom2",
			VestingReserveAddress: types.VestingReserveAddress(1).String(),
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
					Weight:      math.LegacyMustNewDecFromStr("0.5"),
				},
				{
					ReleaseTime: types.MustParseRFC3339("2023-12-01T00:00:00Z"),
					Weight:      math.LegacyMustNewDecFromStr("0.5"),
				},
			},
			StartTime: types.MustParseRFC3339("2022-01-01T00:00:00Z"),
			EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-01T00:00:00Z")},
			Status:    types.AuctionStatusStarted,
		},
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
	)

	validAllowedBidder := types.AllowedBidder{
		AuctionId:    1,
		Bidder:       validAddr.String(),
		MaxBidAmount: math.NewInt(10_000_000),
	}

	validBid := types.Bid{
		AuctionId: 1,
		Id:        1,
		Bidder:    validAddr.String(),
		Price:     math.LegacyMustNewDecFromStr("0.5"),
		Coin:      sdk.NewInt64Coin("denom2", 50_000_000),
	}

	validVestingQueue := types.VestingQueue{
		AuctionId:   1,
		Auctioneer:  validAddr.String(),
		PayingCoin:  sdk.NewInt64Coin("denom2", 100_000_000),
		ReleaseTime: types.MustParseRFC3339("2022-12-20T00:00:00Z"),
		Released:    false,
	}

	tests := []struct {
		desc      string
		configure func(*types.GenesisState)
		valid     bool
	}{
		{
			desc: "default is valid",
			configure: func(genState *types.GenesisState) {
				params := types.DefaultParams()
				genState.Params = params
			},
			valid: true,
		},
		{
			desc: "valid genesis state",
			configure: func(genState *types.GenesisState) {
				params := types.DefaultParams()
				auctionAny, _ := types.PackAuction(validAuction)

				genState.Params = params
				genState.AuctionList = []*codectypes.Any{auctionAny}
				genState.AllowedBidderList = []types.AllowedBidder{validAllowedBidder}
				genState.BidList = []types.Bid{validBid}
				genState.VestingQueueList = []types.VestingQueue{validVestingQueue}

				// TODO fix when add a new field
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "invalid auction - unsupported auction type",
			configure: func(genState *types.GenesisState) {
				auctionAny, _ := types.PackAuction(types.NewFixedPriceAuction(
					&types.BaseAuction{
						Id:                    1,
						Type:                  types.AuctionTypeNil,
						Auctioneer:            validAddr.String(),
						SellingReserveAddress: types.SellingReserveAddress(1).String(),
						PayingReserveAddress:  types.PayingReserveAddress(1).String(),
						StartPrice:            math.LegacyMustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom2",
						VestingReserveAddress: types.VestingReserveAddress(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
								Weight:      math.LegacyMustNewDecFromStr("0.5"),
							},
							{
								ReleaseTime: types.MustParseRFC3339("2023-06-01T00:00:00Z"),
								Weight:      math.LegacyMustNewDecFromStr("0.5"),
							},
						},
						StartTime: types.MustParseRFC3339("2021-12-10T00:00:00Z"),
						EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
						Status:    types.AuctionStatusStarted,
					},
					sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				))

				genState.AuctionList = []*codectypes.Any{auctionAny}
			},
			valid: false,
		},
		{
			desc: "invalid auction - duplicate denom for selling and paying",
			configure: func(genState *types.GenesisState) {
				auctionAny, _ := types.PackAuction(types.NewFixedPriceAuction(
					&types.BaseAuction{
						Id:                    1,
						Type:                  types.AuctionTypeFixedPrice,
						Auctioneer:            validAddr.String(),
						SellingReserveAddress: types.SellingReserveAddress(1).String(),
						PayingReserveAddress:  types.PayingReserveAddress(1).String(),
						StartPrice:            math.LegacyMustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAddress(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"),
								Weight:      math.LegacyMustNewDecFromStr("0.5"),
							},
							{
								ReleaseTime: types.MustParseRFC3339("2022-12-01T00:00:00Z"),
								Weight:      math.LegacyMustNewDecFromStr("0.5"),
							},
						},
						StartTime: types.MustParseRFC3339("2021-12-10T00:00:00Z"),
						EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
						Status:    types.AuctionStatusStarted,
					},
					sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				))

				genState.AuctionList = []*codectypes.Any{auctionAny}
			},
			valid: false,
		},
		{
			desc: "invalid auction - invalid sum of vesting schedule weights",
			configure: func(genState *types.GenesisState) {
				auctionAny, _ := types.PackAuction(types.NewFixedPriceAuction(
					&types.BaseAuction{
						Id:                    1,
						Type:                  types.AuctionTypeFixedPrice,
						Auctioneer:            validAddr.String(),
						SellingReserveAddress: types.SellingReserveAddress(1).String(),
						PayingReserveAddress:  types.PayingReserveAddress(1).String(),
						StartPrice:            math.LegacyMustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAddress(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"),
								Weight:      math.LegacyMustNewDecFromStr("0.9"),
							},
							{
								ReleaseTime: types.MustParseRFC3339("2022-12-01T00:00:00Z"),
								Weight:      math.LegacyMustNewDecFromStr("0.5"),
							},
						},
						StartTime: types.MustParseRFC3339("2021-12-10T00:00:00Z"),
						EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
						Status:    types.AuctionStatusStarted,
					},
					sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				))

				genState.AuctionList = []*codectypes.Any{auctionAny}
			},
			valid: false,
		},
		{
			desc: "invalid auction - invalid auctioneer address",
			configure: func(genState *types.GenesisState) {
				auctionAny, _ := types.PackAuction(types.NewFixedPriceAuction(
					&types.BaseAuction{
						Id:                    1,
						Type:                  types.AuctionTypeFixedPrice,
						Auctioneer:            "invalid",
						SellingReserveAddress: types.SellingReserveAddress(1).String(),
						PayingReserveAddress:  types.PayingReserveAddress(1).String(),
						StartPrice:            math.LegacyMustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAddress(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"),
								Weight:      math.LegacyMustNewDecFromStr("0.9"),
							},
							{
								ReleaseTime: types.MustParseRFC3339("2022-12-01T00:00:00Z"),
								Weight:      math.LegacyMustNewDecFromStr("0.5"),
							},
						},
						StartTime: types.MustParseRFC3339("2021-12-10T00:00:00Z"),
						EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
						Status:    types.AuctionStatusStarted,
					},
					sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				))

				genState.AuctionList = []*codectypes.Any{auctionAny}
			},
			valid: false,
		},
		{
			desc: "invalid bid - invalid bidder address",
			configure: func(genState *types.GenesisState) {
				genState.BidList = []types.Bid{
					{
						AuctionId: 1,
						Id:        1,
						Bidder:    "invalid",
						Price:     math.LegacyMustNewDecFromStr("0.5"),
						Coin:      sdk.NewInt64Coin("denom2", 50_000_000),
					},
				}
			},
			valid: false,
		},
		{
			desc: "invalid bid - invalid coin amount",
			configure: func(genState *types.GenesisState) {
				genState.BidList = []types.Bid{
					{
						AuctionId: 1,
						Id:        1,
						Bidder:    validAddr.String(),
						Price:     math.LegacyMustNewDecFromStr("0.5"),
						Coin:      sdk.NewInt64Coin("denom2", 0),
					},
				}
			},
			valid: false,
		},
		{
			desc: "invalid bid - invalid price",
			configure: func(genState *types.GenesisState) {
				genState.BidList = []types.Bid{
					{
						AuctionId: 1,
						Id:        1,
						Bidder:    validAddr.String(),
						Price:     math.LegacyMustNewDecFromStr("0"),
						Coin:      sdk.NewInt64Coin("denom2", 100_000),
					},
				}
			},
			valid: false,
		},
		{
			desc: "invalid allowed bidder - invalid max bid amount",
			configure: func(genState *types.GenesisState) {
				genState.AllowedBidderList = []types.AllowedBidder{
					{
						AuctionId:    1,
						Bidder:       validAddr.String(),
						MaxBidAmount: math.NewInt(0),
					},
				}
			},
			valid: false,
		},
		{
			desc: "invalid vesting queue - invalid auctioneer address",
			configure: func(genState *types.GenesisState) {
				params := types.DefaultParams()
				genState.Params = params
				genState.VestingQueueList = []types.VestingQueue{
					{
						AuctionId:   2,
						Auctioneer:  "",
						PayingCoin:  sdk.NewInt64Coin("denom2", 100_000_000),
						ReleaseTime: types.MustParseRFC3339("2022-12-20T00:00:00Z"),
						Released:    false,
					},
				}
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			genState := types.DefaultGenesis()
			tc.configure(genState)
			err := genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
