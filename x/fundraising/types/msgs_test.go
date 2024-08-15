package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestMsgCreateFixedPriceAuction(t *testing.T) {
	testCases := []struct {
		expectedErr string
		msg         *types.MsgCreateFixedPriceAuction
	}{
		{
			"", // empty means no error expected
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"start price must be positive: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"selling coin amount must be positive: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 0),
				"denom1",
				[]types.VestingSchedule{},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"selling coin denom must not be the same as paying coin denom: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom2",
				[]types.VestingSchedule{},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"end time must be set after start time: invalid request",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				time.Now(),
				time.Now().AddDate(-1, 0, 0),
			),
		},
		{
			"vesting weight must be positive: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
						math.LegacyZeroDec(),
					},
				},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"vesting weight must not be greater than 1: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
						math.LegacyMustNewDecFromStr("1.1"),
					},
				},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"release time must be set after the end time: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						types.MustParseRFC3339("2022-06-01T22:08:41+00:00"),
						math.LegacyMustNewDecFromStr("1.0"),
					},
				},
				time.Now(),
				time.Now().AddDate(1, 0, 0),
			),
		},
		{
			"release time must be chronological: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
						math.LegacyMustNewDecFromStr("0.5"),
					},
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 3, 0),
						math.LegacyMustNewDecFromStr("0.5"),
					},
				},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"total vesting weight must be equal to 1: invalid vesting schedules",
			types.NewMsgCreateFixedPriceAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
						math.LegacyMustNewDecFromStr("0.5"),
					},
					{
						time.Now().AddDate(0, 1, 0).AddDate(1, 0, 0),
						math.LegacyMustNewDecFromStr("0.3"),
					},
				},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"invalid auctioneer address: empty address string is not allowed: invalid address",
			types.NewMsgCreateFixedPriceAuction(
				"",
				math.LegacyMustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCreateFixedPriceAuction{}, tc.msg)

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgCreateBatchAuction(t *testing.T) {
	testCases := []struct {
		expectedErr string
		msg         *types.MsgCreateBatchAuction
	}{
		{
			"", // empty means no error expected
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"start price must be positive: invalid request",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"minimum price must be positive: invalid request",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.1"),
				math.LegacyMustNewDecFromStr("0"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"selling coin amount must be positive: invalid request",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 0),
				"denom1",
				[]types.VestingSchedule{},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"selling coin denom must not be the same as paying coin denom: invalid request",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom2",
				[]types.VestingSchedule{},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"end time must be set after start time: invalid request",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(-1, 0, 0),
			),
		},
		{
			"vesting weight must be positive: invalid vesting schedules",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
						math.LegacyZeroDec(),
					},
				},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"vesting weight must not be greater than 1: invalid vesting schedules",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
						math.LegacyMustNewDecFromStr("1.1"),
					},
				},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"release time must be set after the end time: invalid vesting schedules",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now(),
						math.LegacyMustNewDecFromStr("1.0"),
					},
				},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(1, 0, 0),
			),
		},
		{
			"release time must be chronological: invalid vesting schedules",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
						math.LegacyMustNewDecFromStr("0.5"),
					},
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 3, 0),
						math.LegacyMustNewDecFromStr("0.5"),
					},
				},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"total vesting weight must be equal to 1: invalid vesting schedules",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{
					{
						time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
						math.LegacyMustNewDecFromStr("0.5"),
					},
					{
						time.Now().AddDate(0, 1, 0).AddDate(1, 0, 0),
						math.LegacyMustNewDecFromStr("0.3"),
					},
				},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"extend rate must be positive: invalid request",
			types.NewMsgCreateBatchAuction(
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				uint32(2),
				math.LegacyMustNewDecFromStr("-0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
		{
			"invalid auctioneer address: empty address string is not allowed: invalid address",
			types.NewMsgCreateBatchAuction(
				"",
				math.LegacyMustNewDecFromStr("0.5"),
				math.LegacyMustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom2", 10_000_000_000_000),
				"denom1",
				[]types.VestingSchedule{},
				uint32(2),
				math.LegacyMustNewDecFromStr("0.05"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCreateBatchAuction{}, tc.msg)

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
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
		{
			"invalid auctioneer address: empty address string is not allowed: invalid address",
			types.NewMsgCancelAuction(
				"",
				uint64(1),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgCancelAuction{}, tc.msg)

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgPlaceBid(t *testing.T) {
	testCases := []struct {
		expectedErr string
		msg         *types.MsgPlaceBid
	}{
		{
			"", // empty means no error expected
			types.NewMsgPlaceBid(
				uint64(1),
				sdk.AccAddress(crypto.AddressHash([]byte("Bidder"))).String(),
				types.BidTypeBatchWorth,
				math.LegacyOneDec(),
				sdk.NewInt64Coin("denom2", 1000000),
			),
		},
		{
			"bid price must be positive value: invalid request",
			types.NewMsgPlaceBid(
				uint64(1),
				sdk.AccAddress(crypto.AddressHash([]byte("Bidder"))).String(),
				types.BidTypeBatchWorth,
				math.LegacyZeroDec(),
				sdk.NewInt64Coin("denom2", 1000000),
			),
		},
		{
			"invalid coin amount: 0: invalid request",
			types.NewMsgPlaceBid(
				uint64(1),
				sdk.AccAddress(crypto.AddressHash([]byte("Bidder"))).String(),
				types.BidTypeBatchWorth,
				math.LegacyOneDec(),
				sdk.NewInt64Coin("denom2", 0),
			),
		},
		{
			"invalid bidder address: empty address string is not allowed: invalid address",
			types.NewMsgPlaceBid(
				uint64(1),
				"",
				types.BidTypeBatchWorth,
				math.LegacyOneDec(),
				sdk.NewInt64Coin("denom2", 1000000),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgPlaceBid{}, tc.msg)

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgModifyBid(t *testing.T) {
	testCases := []struct {
		expectedErr string
		msg         *types.MsgModifyBid
	}{
		{
			"", // empty means no error expected
			types.NewMsgModifyBid(
				uint64(1),
				sdk.AccAddress(crypto.AddressHash([]byte("Bidder"))).String(),
				uint64(0),
				math.LegacyOneDec(),
				sdk.NewInt64Coin("denom2", 1000000),
			),
		},
		{
			"bid price must be positive value: invalid request",
			types.NewMsgModifyBid(
				uint64(1),
				sdk.AccAddress(crypto.AddressHash([]byte("Bidder"))).String(),
				uint64(0),
				math.LegacyZeroDec(),
				sdk.NewInt64Coin("denom2", 1000000),
			),
		},
		{
			"invalid coin amount: 0: invalid request",
			types.NewMsgModifyBid(
				uint64(1),
				sdk.AccAddress(crypto.AddressHash([]byte("Bidder"))).String(),
				uint64(0),
				math.LegacyOneDec(),
				sdk.NewInt64Coin("denom2", 0),
			),
		},
		{
			"invalid bidder address: empty address string is not allowed: invalid address",
			types.NewMsgModifyBid(
				uint64(1),
				"",
				uint64(0),
				math.LegacyOneDec(),
				sdk.NewInt64Coin("denom2", 1000000),
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgModifyBid{}, tc.msg)

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestAddAllowedBidder(t *testing.T) {
	testCases := []struct {
		expectedErr string
		msg         *types.MsgAddAllowedBidder
	}{
		{
			"", // empty means no error expected
			types.NewMsgAddAllowedBidder(
				1,
				types.AllowedBidder{
					1,
					sdk.AccAddress(crypto.AddressHash([]byte("Bidder"))).String(),
					math.NewInt(100_000_000),
				},
			),
		},
		{
			"invalid bidder address: empty address string is not allowed: invalid address",
			types.NewMsgAddAllowedBidder(
				1,
				types.AllowedBidder{
					1,
					"",
					math.NewInt(100_000_000),
				},
			),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgAddAllowedBidder{}, tc.msg)
		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}
