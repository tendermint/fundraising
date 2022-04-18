package simulation

import (
	"bytes"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding fundraising type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.AuctionKeyPrefix):
			var aA, aB types.BaseAuction
			cdc.MustUnmarshal(kvA.Value, &aA)
			cdc.MustUnmarshal(kvB.Value, &aB)
			return fmt.Sprintf("%v\n%v", aA, aB)

		case bytes.Equal(kvA.Key[:1], types.LastBidIdKeyPrefix):
			fmt.Printf("> kvA.Key %X\n: ", kvA.Key[:1])
			fmt.Println("> kvA.Value: ", kvA.Value)
			fmt.Printf("> kvB.Key %X\n: ", kvB.Key[:1])
			fmt.Println("> kvB.Value: ", kvB.Value)

			var bA, bB gogotypes.UInt64Value
			cdc.MustUnmarshal(kvA.Value, &bA)
			cdc.MustUnmarshal(kvB.Value, &bB)
			return fmt.Sprintf("%d\n%d", bA, bB)

		case bytes.Equal(kvA.Key[:1], types.BidKeyPrefix):
			var bA, bB types.Bid
			cdc.MustUnmarshal(kvA.Value, &bA)
			cdc.MustUnmarshal(kvB.Value, &bB)
			return fmt.Sprintf("%v\n%v", bA, bB)

		case bytes.Equal(kvA.Key[:1], types.VestingQueueKeyPrefix):
			var vA, vB types.VestingQueue
			cdc.MustUnmarshal(kvA.Value, &vA)
			cdc.MustUnmarshal(kvB.Value, &vB)
			return fmt.Sprintf("%v\n%v", vA, vB)

		default:
			panic(fmt.Sprintf("invalid fundraising key prefix %X", kvA.Key[:1]))
		}
	}
}
