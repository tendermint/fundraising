package simulation_test

import (
	"fmt"
	"testing"

	simappparams "cosmossdk.io/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/fundraising/x/fundraising/simulation"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestDecodeFarmingStore(t *testing.T) {
	cdc := simappparams.MakeTestEncodingConfig().Codec
	dec := simulation.NewDecodeStore(cdc)

	baseAuction := types.BaseAuction{StartPrice: sdk.NewDec(0)}
	bid := types.Bid{Price: sdk.NewDec(0)}
	vestingQueue := types.VestingQueue{}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.AuctionKeyPrefix, Value: cdc.MustMarshal(&baseAuction)},
			{Key: types.BidKeyPrefix, Value: cdc.MustMarshal(&bid)},
			{Key: types.VestingQueueKeyPrefix, Value: cdc.MustMarshal(&vestingQueue)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Auction", fmt.Sprintf("%v\n%v", baseAuction, baseAuction)},
		{"Bid", fmt.Sprintf("%v\n%v", bid, bid)},
		{"VestingQueue", fmt.Sprintf("%v\n%v", vestingQueue, vestingQueue)},
		{"other", ""},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				got := dec(kvPairs.Pairs[i], kvPairs.Pairs[i])
				require.EqualValues(t, tt.expectedLog, got, tt.name)
			}
		})
	}
}
