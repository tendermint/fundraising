package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestBid() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)

	_, err := suite.srv.PlaceBid(ctx, suite.sampleFixedPriceBids[0])
	suite.Require().NoError(err)

	_, err = suite.srv.PlaceBid(ctx, suite.sampleFixedPriceBids[1])
	suite.Require().NoError(err)

	fmt.Println("msgs: ", suite.sampleFixedPriceBids)

	bids := []types.Bid{}
	suite.keeper.IterateBidsByAuctionId(suite.ctx, auction.GetId(), func(bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})

	fmt.Println("len: ", len(bids))
	for _, b := range bids {
		fmt.Println("b: ", b)
	}

	// bidderAcc := suite.addrs[0]
	// price := sdk.MustNewDecFromStr("")
	// coin := sdk.NewInt64Coin("", 1_000_000)
	// nextSeq := suite.keeper.GetNextSequenceWithUpdate(suite.ctx)

	// bid := types.Bid{
	// 	AuctionId: auction.GetId(),
	// 	Sequence:  nextSeq,
	// 	Bidder:    bidderAcc.String(),
	// 	Price:     price,
	// 	Coin:      coin,
	// }

	// suite.keeper.SetBid(
	// 	suite.ctx,
	// 	auction.GetId(),
	// 	nextSeq,
	// 	bidderAcc,
	// 	bid,
	// )

	// Iterate 테스트 필요
}
