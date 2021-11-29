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
	startTime, _ := time.Parse(time.RFC3339, "2021-11-01T22:08:41+00:00") // needs to be deterministic for test
	endTime := startTime.AddDate(1, 0, 0)

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
			"start price must be positve 0.000000000000000000: invalid request",
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
			"selling coin amount must be positive 0ugdex: invalid request",
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
		// TODO: vesting schedules not covered
		{
			"end time 2020-11-01T22:08:41Z must be greater than start time 2021-11-01T22:08:41Z: invalid auction end time",
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
	// TODO: not implemented yet
}

func TestMsgPlaceBid(t *testing.T) {
	// TODO: not implemented yet
}
