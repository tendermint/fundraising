package simulation_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/tendermint/fundraising/x/fundraising/simulation"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestDecodeFarmingStore(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	baseAuction := types.BaseAuction{}
	lastBidId := gogotypes.UInt64Value{}
	bid := types.Bid{}
	vestingQueue := types.VestingQueue{}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.AuctionKeyPrefix, Value: cdc.MustMarshal(&baseAuction)},
			{Key: types.LastBidIdKeyPrefix, Value: cdc.MustMarshal(&lastBidId)},
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
		{"LastBidId", fmt.Sprintf("%v\n%v", lastBidId, lastBidId)},
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
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
