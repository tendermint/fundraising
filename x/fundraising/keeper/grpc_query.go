package keeper

import (
	"github.com/tendermint/fundraising/x/fundraising/types"
)

var _ types.QueryServer = Keeper{}
