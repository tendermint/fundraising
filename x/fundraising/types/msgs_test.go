package types_test

import (
	"testing"
	time "time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestMsgCreateFixedPriceAuction(t *testing.T) {
	auctioneerAcc := sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer")))
	startTime, _ := time.Parse(time.RFC3339, "2021-11-01T22:00:00+00:00")
	endTime := startTime.AddDate(0, 1, 0) // add 1 month

	testCases := []struct {
		expectedErr string
		msg         *types.MsgCreateFixedPriceAuction
	}{
		{
			"", // empty means no error expected
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				startTime,
				endTime,
			),
		},
		{
			"start price must be positve: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				startTime,
				endTime,
			),
		},
		{
			"selling coin amount must be positive: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 0),
				"denom1",
				[]types.VestingSchedule{},
				startTime,
				endTime,
			),
		},
		{
			"vesting weight must be positive: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						types.ParseTime("2022-06-01T22:08:41+00:00"),
						sdk.ZeroDec(),
					},
				},
				startTime,
				endTime,
			),
		},
		{
			"each vesting weight must not be greater than 1: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						types.ParseTime("2022-06-01T22:08:41+00:00"),
						sdk.MustNewDecFromStr("1.1"),
					},
				},
				startTime,
				endTime,
			),
		},
		{
			"release time must be after the end time: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						types.ParseTime("2022-06-01T22:08:41+00:00"),
						sdk.MustNewDecFromStr("1.0"),
					},
				},
				startTime,
				types.ParseTime("2022-06-05T22:08:41+00:00"),
			),
		},
		{
			"total vesting weight must be equal to 1: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						types.ParseTime("2022-06-01T22:00:00+00:00"),
						sdk.MustNewDecFromStr("0.5"),
					},
					{
						types.ParseTime("2022-12-01T22:00:00+00:00"),
						sdk.MustNewDecFromStr("0.3"),
					},
				},
				startTime,
				endTime,
			),
		},
		{
			"end time must be greater than start time: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				startTime,
				startTime.AddDate(-1, 0, 0),
			),
		},
		{
			"release time must be chronological: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						types.ParseTime("2022-12-01T22:00:00+00:00"),
						sdk.MustNewDecFromStr("0.5"),
					},
					{
						types.ParseTime("2022-06-01T22:00:00+00:00"),
						sdk.MustNewDecFromStr("0.5"),
					},
				},
				startTime,
				endTime,
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCreateFixedPriceAuction{}, tc.msg)
		require.Equal(t, types.TypeMsgCreateFixedPriceAuction, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetAuctioneer(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgCreateEnglishAuction(t *testing.T) {
	// TODO: not implemented yet
}

func TestMsgCancelAuction(t *testing.T) {
	testCases := []struct {
		expectedErr string
		msg         *types.MsgCancelAuction
	}{
		{
			"", // empty means no error expected
			types.NewMsgCancelAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				uint64(1),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCancelAuction{}, tc.msg)
		require.Equal(t, types.TypeMsgCancelAuction, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetAuctioneer(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgPlaceBid(t *testing.T) {
	bidderAcc := sdk.AccAddress(crypto.AddressHash([]byte("Bidder")))

	testCases := []struct {
		expectedErr string
		msg         *types.MsgPlaceBid
	}{
		{
			"", // empty means no error expected
			types.NewMsgPlaceBid(
				uint64(1),
				bidderAcc.String(),
				sdk.OneDec(),
				sdk.NewInt64Coin("denom2", 1000000),
			),
		},
		{
			"bid price must be positve value: invalid request",
			types.NewMsgPlaceBid(
				uint64(1),
				bidderAcc.String(),
				sdk.ZeroDec(),
				sdk.NewInt64Coin("denom2", 1000000),
			),
		},
		{
			"bid price must be positve value: invalid request",
			types.NewMsgPlaceBid(
				uint64(1),
				bidderAcc.String(),
				sdk.ZeroDec(),
				sdk.NewInt64Coin("denom2", 0),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgPlaceBid{}, tc.msg)
		require.Equal(t, types.TypeMsgPlaceBid, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(legacy.Cdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetBidder(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}
