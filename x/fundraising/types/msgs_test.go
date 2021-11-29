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
	distributedTime1, _ := time.Parse(time.RFC3339, "2022-06-01T22:08:41+00:00")
	distributedTime2, _ := time.Parse(time.RFC3339, "2022-12-01T22:08:41+00:00")

	testCases := []struct {
		expectedErr string
		msg         *types.MsgCreateFixedPriceAuction
	}{
		{
			"", // empty means no error expected
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("ugdex", 10_000_000_000_000),
				"uatom",
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
				sdk.NewInt64Coin("ugdex", 10_000_000_000_000),
				"uatom",
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
				sdk.NewInt64Coin("ugdex", 0),
				"uatom",
				[]types.VestingSchedule{},
				startTime,
				endTime,
			),
		},
		{
			"vesting weight must be positive: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("ugdex", 10_000_000_000_000),
				"uatom",
				[]types.VestingSchedule{
					types.NewVestingSchedule(distributedTime1, sdk.ZeroDec()),
				},
				startTime,
				endTime,
			),
		},
		{
			"total vesting weight must not greater than 1: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("ugdex", 10_000_000_000_000),
				"uatom",
				[]types.VestingSchedule{
					types.NewVestingSchedule(distributedTime1, sdk.MustNewDecFromStr("1.1")),
				},
				startTime,
				endTime,
			),
		},
		{
			"total vesting weight must be equal to 1: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("ugdex", 10_000_000_000_000),
				"uatom",
				[]types.VestingSchedule{
					types.NewVestingSchedule(distributedTime1, sdk.MustNewDecFromStr("0.5")),
					types.NewVestingSchedule(distributedTime2, sdk.MustNewDecFromStr("0.3")),
				},
				startTime,
				endTime,
			),
		},
		{
			"end time must be greater than start time: invalid auction end time",
			types.NewMsgCreateFixedPriceAuction(
				auctioneerAcc.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("ugdex", 10_000_000_000_000),
				"uatom",
				[]types.VestingSchedule{},
				startTime,
				startTime.AddDate(-1, 0, 0),
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

func TestMsgCancelFundraising(t *testing.T) {
	testCases := []struct {
		expectedErr string
		msg         *types.MsgCancelFundraising
	}{
		{
			"", // empty means no error expected
			types.NewMsgCancelFundraising(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				uint64(1),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCancelFundraising{}, tc.msg)
		require.Equal(t, types.TypeMsgCancelFundraising, tc.msg.Type())
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
				sdk.NewInt64Coin("ugdex", 1000000),
			),
		},
		{
			"bid price must be positve value: invalid request",
			types.NewMsgPlaceBid(
				uint64(1),
				bidderAcc.String(),
				sdk.ZeroDec(),
				sdk.NewInt64Coin("ugdex", 1000000),
			),
		},
		{
			"bid price must be positve value: invalid request",
			types.NewMsgPlaceBid(
				uint64(1),
				bidderAcc.String(),
				sdk.ZeroDec(),
				sdk.NewInt64Coin("ugdex", 0),
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
