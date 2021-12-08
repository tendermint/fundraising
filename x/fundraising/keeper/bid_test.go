package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (suite *KeeperTestSuite) TestBid() {
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)

	bidderAcc := suite.addrs[0]
	price := sdk.MustNewDecFromStr("")
	coin := sdk.NewInt64Coin("", 1_000_000)
	nextSeq := suite.keeper.GetNextSequenceWithUpdate(suite.ctx)

	bid := types.Bid{
		AuctionId: auction.GetId(),
		Sequence:  nextSeq,
		Bidder:    bidderAcc.String(),
		Price:     price,
		Coin:      coin,
	}

	suite.keeper.SetBid(
		suite.ctx,
		auction.GetId(),
		nextSeq,
		bidderAcc,
		bid,
	)
}
