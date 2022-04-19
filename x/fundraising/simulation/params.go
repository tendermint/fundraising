package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals.
// on the simulation.
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyAuctionCreationFee),
			func(r *rand.Rand) string {
				bz, err := GenAuctionCreationFee(r).MarshalJSON()
				if err != nil {
					panic(err)
				}
				return string(bz)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyExtendedPeriod),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", GenExtendedPeriod(r))
			},
		),
	}
}
