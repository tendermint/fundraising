package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestMsgCreateFixedPriceAuction() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	for _, tc := range []struct {
		name string
		msg  *types.MsgCreateFixedPriceAuction
		err  error
	}{
		{
			"valid message",
			types.NewMsgCreateFixedPriceAuction(
				suite.addrs[0].String(),
				suite.StartPrice("1.0"),
				suite.SellingCoin(1_000_000_000_000),
				suite.PayingCoinDenom(),
				suite.VestingSchedules(),
				types.ParseTime("2021-12-01T00:00:00Z"),
				types.ParseTime("2022-01-01T00:00:00Z"),
			),
			nil,
		},
	} {
		suite.Run(tc.name, func() {
			_, err := suite.srv.CreateFixedPriceAuction(ctx, tc.msg)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
				return
			}
			suite.Require().NoError(err)

			auction, found := suite.keeper.GetAuction(suite.ctx, 1)
			fmt.Println("found: ", found)
			fmt.Println("auction: ", auction)
		})
	}
}

func (suite *KeeperTestSuite) TestMsgCreateEnglishAuction() {
	// TODO: not implemented yet
}

func (suite *KeeperTestSuite) TestMsgCancelAuction() {
	// TODO: not implemented yet
}

func (suite *KeeperTestSuite) TestMsgPlaceBid() {
	// TODO: not implemented yet
}
