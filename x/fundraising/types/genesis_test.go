package types_test

import (
	"testing"
	time "time"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

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
			StartPrice:            sdk.MustNewDecFromStr("0.5"),
			SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
			PayingCoinDenom:       "denom2",
			VestingReserveAddress: types.VestingReserveAddress(1).String(),
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
					Weight:      sdk.MustNewDecFromStr("0.5"),
				},
				{
					ReleaseTime: types.MustParseRFC3339("2023-12-01T00:00:00Z"),
					Weight:      sdk.MustNewDecFromStr("0.5"),
				},
			},
			StartTime: types.MustParseRFC3339("2022-01-01T00:00:00Z"),
			EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-01T00:00:00Z")},
			Status:    types.AuctionStatusStarted,
		},
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
	)

	validAllowedBidder := types.AllowedBidder{
		Bidder:       validAddr.String(),
		MaxBidAmount: sdk.NewInt(10_000_000),
	}

	validBid := types.Bid{
		AuctionId: 1,
		Id:        1,
		Bidder:    validAddr.String(),
		Price:     sdk.MustNewDecFromStr("0.5"),
		Coin:      sdk.NewInt64Coin("denom2", 50_000_000),
	}

	validVestingQueue := types.VestingQueue{
		AuctionId:   1,
		Auctioneer:  validAddr.String(),
		PayingCoin:  sdk.NewInt64Coin("denom2", 100_000_000),
		ReleaseTime: types.MustParseRFC3339("2022-12-20T00:00:00Z"),
		Released:    false,
	}

	for _, tc := range []struct {
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
				genState.Auctions = []*codectypes.Any{auctionAny}
				genState.AllowedBidderRecords = []types.AllowedBidderRecord{
					{
						AuctionId:     1,
						AllowedBidder: validAllowedBidder,
					},
				}
				genState.Bids = []types.Bid{validBid}
				genState.VestingQueues = []types.VestingQueue{validVestingQueue}
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
						StartPrice:            sdk.MustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom2",
						VestingReserveAddress: types.VestingReserveAddress(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
							{
								ReleaseTime: types.MustParseRFC3339("2023-06-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
						},
						StartTime: types.MustParseRFC3339("2021-12-10T00:00:00Z"),
						EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
						Status:    types.AuctionStatusStarted,
					},
					sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				))

				genState.Auctions = []*codectypes.Any{auctionAny}
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
						StartPrice:            sdk.MustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAddress(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
							{
								ReleaseTime: types.MustParseRFC3339("2022-12-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
						},
						StartTime: types.MustParseRFC3339("2021-12-10T00:00:00Z"),
						EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
						Status:    types.AuctionStatusStarted,
					},
					sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				))

				genState.Auctions = []*codectypes.Any{auctionAny}
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
						StartPrice:            sdk.MustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAddress(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.9"),
							},
							{
								ReleaseTime: types.MustParseRFC3339("2022-12-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
						},
						StartTime: types.MustParseRFC3339("2021-12-10T00:00:00Z"),
						EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
						Status:    types.AuctionStatusStarted,
					},
					sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				))

				genState.Auctions = []*codectypes.Any{auctionAny}
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
						StartPrice:            sdk.MustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAddress(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.9"),
							},
							{
								ReleaseTime: types.MustParseRFC3339("2022-12-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
						},
						StartTime: types.MustParseRFC3339("2021-12-10T00:00:00Z"),
						EndTimes:  []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
						Status:    types.AuctionStatusStarted,
					},
					sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				))

				genState.Auctions = []*codectypes.Any{auctionAny}
			},
			valid: false,
		},
		{
			desc: "invalid bid - invalid bidder address",
			configure: func(genState *types.GenesisState) {
				genState.Bids = []types.Bid{
					{
						AuctionId: 1,
						Id:        1,
						Bidder:    "invalid",
						Price:     sdk.MustNewDecFromStr("0.5"),
						Coin:      sdk.NewInt64Coin("denom2", 50_000_000),
					},
				}
			},
			valid: false,
		},
		{
			desc: "invalid bid - invalid coin amount",
			configure: func(genState *types.GenesisState) {
				genState.Bids = []types.Bid{
					{
						AuctionId: 1,
						Id:        1,
						Bidder:    validAddr.String(),
						Price:     sdk.MustNewDecFromStr("0.5"),
						Coin:      sdk.NewInt64Coin("denom2", 0),
					},
				}
			},
			valid: false,
		},
		{
			desc: "invalid bid - invalid price",
			configure: func(genState *types.GenesisState) {
				genState.Bids = []types.Bid{
					{
						AuctionId: 1,
						Id:        1,
						Bidder:    validAddr.String(),
						Price:     sdk.MustNewDecFromStr("0"),
						Coin:      sdk.NewInt64Coin("denom2", 100_000),
					},
				}
			},
			valid: false,
		},
		{
			desc: "invalid allowed bidder - invalid max bid amount",
			configure: func(genState *types.GenesisState) {
				genState.AllowedBidderRecords = []types.AllowedBidderRecord{
					{
						AuctionId: 1,
						AllowedBidder: types.AllowedBidder{
							Bidder:       validAddr.String(),
							MaxBidAmount: sdk.NewInt(0),
						},
					},
				}
			},
			valid: false,
		},
		{
			desc: "invalid allowed bidder - auction id cannot be 0",
			configure: func(genState *types.GenesisState) {
				genState.AllowedBidderRecords = []types.AllowedBidderRecord{
					{
						AuctionId: 0,
						AllowedBidder: types.AllowedBidder{
							Bidder:       validAddr.String(),
							MaxBidAmount: sdk.NewInt(100_000),
						},
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
				genState.VestingQueues = []types.VestingQueue{
					{
						AuctionId:   1,
						Auctioneer:  "",
						PayingCoin:  sdk.NewInt64Coin("denom2", 100_000_000),
						ReleaseTime: types.MustParseRFC3339("2022-12-20T00:00:00Z"),
						Released:    false,
					},
				}
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			genState := types.DefaultGenesisState()
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
