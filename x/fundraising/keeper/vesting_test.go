package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestVestingQueue_RemainingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2022-01-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.3"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2022-04-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.3"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2022-08-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.4"),
			},
		},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidFixedPrice(auction.GetId(), s.addr(1), sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), sdk.OneDec(), parseCoin("15000000denom2"), true)

	err := s.keeper.ApplyVestingSchedules(s.ctx, auction)
	s.Require().NoError(err)

	vestingReserveAddr := auction.GetVestingReserveAddress()
	vestingReserveCoin := s.getBalance(vestingReserveAddr, auction.PayingCoinDenom)

	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId()) {
		vestingReserveCoin = vestingReserveCoin.Sub(vq.PayingCoin)
	}
	s.Require().True(vestingReserveCoin.IsZero())
}
