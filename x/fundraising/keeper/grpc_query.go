package keeper

import (
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Keeper{}
