package types_test

import (
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestSanitizeBids(t *testing.T) {
	sampleBids := []types.Bid{
		{
			AuctionId: 1,
			Bidder:    sdk.AccAddress(crypto.AddressHash([]byte("Bidder1"))).String(),
			Id:        1,
			Price:     sdk.MustNewDecFromStr("0.10"),
			Coin:      sdk.NewInt64Coin("denom1", 1),
		},
		{
			AuctionId: 1,
			Bidder:    sdk.AccAddress(crypto.AddressHash([]byte("Bidder2"))).String(),
			Id:        2,
			Price:     sdk.MustNewDecFromStr("1.10"),
			Coin:      sdk.NewInt64Coin("denom1", 1),
		},
		{
			AuctionId: 1,
			Bidder:    sdk.AccAddress(crypto.AddressHash([]byte("Bidder3"))).String(),
			Id:        3,
			Price:     sdk.MustNewDecFromStr("0.35"),
			Coin:      sdk.NewInt64Coin("denom1", 1),
		},
		{
			AuctionId: 1,
			Bidder:    sdk.AccAddress(crypto.AddressHash([]byte("Bidder4"))).String(),
			Id:        4,
			Price:     sdk.MustNewDecFromStr("0.77"),
			Coin:      sdk.NewInt64Coin("denom1", 1),
		},
	}

	bids := types.SanitizeReverseBids(sampleBids)

	require.True(t, sort.SliceIsSorted(bids, func(i, j int) bool {
		return bids[i].Price.GT(bids[j].Price)
	}))
}
