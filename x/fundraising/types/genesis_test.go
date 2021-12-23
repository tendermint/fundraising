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
	validAcc := sdk.AccAddress(crypto.AddressHash([]byte("validAcc")))
	validAuction := types.NewFixedPriceAuction(
		&types.BaseAuction{
			Id:                    1,
			Type:                  types.AuctionTypeFixedPrice,
			Auctioneer:            validAcc.String(),
			SellingReserveAddress: types.SellingReserveAcc(1).String(),
			PayingReserveAddress:  types.PayingReserveAcc(1).String(),
			StartPrice:            sdk.MustNewDecFromStr("0.5"),
			SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
			PayingCoinDenom:       "denom2",
			VestingReserveAddress: types.VestingReserveAcc(1).String(),
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: types.ParseTime("2022-06-01T00:00:00Z"),
					Weight:      sdk.MustNewDecFromStr("0.5"),
				},
				{
					ReleaseTime: types.ParseTime("2022-12-01T00:00:00Z"),
					Weight:      sdk.MustNewDecFromStr("0.5"),
				},
			},
			WinningPrice:  sdk.ZeroDec(),
			RemainingCoin: sdk.NewInt64Coin("denom1", 1_000_000_000_000),
			StartTime:     types.ParseTime("2021-12-10T00:00:00Z"),
			EndTimes:      []time.Time{types.ParseTime("2022-12-20T00:00:00Z")},
			Status:        types.AuctionStatusStarted,
		},
	)
	validBid := types.Bid{
		AuctionId: 1,
		Sequence:  1,
		Bidder:    validAcc.String(),
		Price:     sdk.MustNewDecFromStr("0.5"),
		Coin:      sdk.NewInt64Coin("denom2", 50_000_000),
	}
	validVestingQueue := types.VestingQueue{
		AuctionId:   1,
		Auctioneer:  validAcc.String(),
		PayingCoin:  sdk.NewInt64Coin("denom2", 100_000_000),
		ReleaseTime: types.ParseTime("2022-12-20T00:00:00Z"),
		Vested:      false,
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
						Auctioneer:            validAcc.String(),
						SellingReserveAddress: types.SellingReserveAcc(1).String(),
						PayingReserveAddress:  types.PayingReserveAcc(1).String(),
						StartPrice:            sdk.MustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom2",
						VestingReserveAddress: types.VestingReserveAcc(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.ParseTime("2022-06-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
							{
								ReleaseTime: types.ParseTime("2022-12-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
						},
						WinningPrice:  sdk.ZeroDec(),
						RemainingCoin: sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						StartTime:     types.ParseTime("2021-12-10T00:00:00Z"),
						EndTimes:      []time.Time{types.ParseTime("2022-12-20T00:00:00Z")},
						Status:        types.AuctionStatusStarted,
					},
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
						Auctioneer:            validAcc.String(),
						SellingReserveAddress: types.SellingReserveAcc(1).String(),
						PayingReserveAddress:  types.PayingReserveAcc(1).String(),
						StartPrice:            sdk.MustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAcc(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.ParseTime("2022-06-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
							{
								ReleaseTime: types.ParseTime("2022-12-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
						},
						WinningPrice:  sdk.ZeroDec(),
						RemainingCoin: sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						StartTime:     types.ParseTime("2021-12-10T00:00:00Z"),
						EndTimes:      []time.Time{types.ParseTime("2022-12-20T00:00:00Z")},
						Status:        types.AuctionStatusStarted,
					},
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
						Auctioneer:            validAcc.String(),
						SellingReserveAddress: types.SellingReserveAcc(1).String(),
						PayingReserveAddress:  types.PayingReserveAcc(1).String(),
						StartPrice:            sdk.MustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAcc(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.ParseTime("2022-06-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.9"),
							},
							{
								ReleaseTime: types.ParseTime("2022-12-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
						},
						WinningPrice:  sdk.ZeroDec(),
						RemainingCoin: sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						StartTime:     types.ParseTime("2021-12-10T00:00:00Z"),
						EndTimes:      []time.Time{types.ParseTime("2022-12-20T00:00:00Z")},
						Status:        types.AuctionStatusStarted,
					},
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
						SellingReserveAddress: types.SellingReserveAcc(1).String(),
						PayingReserveAddress:  types.PayingReserveAcc(1).String(),
						StartPrice:            sdk.MustNewDecFromStr("0.5"),
						SellingCoin:           sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						PayingCoinDenom:       "denom1",
						VestingReserveAddress: types.VestingReserveAcc(1).String(),
						VestingSchedules: []types.VestingSchedule{
							{
								ReleaseTime: types.ParseTime("2022-06-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.9"),
							},
							{
								ReleaseTime: types.ParseTime("2022-12-01T00:00:00Z"),
								Weight:      sdk.MustNewDecFromStr("0.5"),
							},
						},
						WinningPrice:  sdk.ZeroDec(),
						RemainingCoin: sdk.NewInt64Coin("denom1", 1_000_000_000_000),
						StartTime:     types.ParseTime("2021-12-10T00:00:00Z"),
						EndTimes:      []time.Time{types.ParseTime("2022-12-20T00:00:00Z")},
						Status:        types.AuctionStatusStarted,
					},
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
						Sequence:  1,
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
						Sequence:  1,
						Bidder:    validAcc.String(),
						Price:     sdk.MustNewDecFromStr("0.5"),
						Coin:      sdk.NewInt64Coin("denom2", 0),
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
