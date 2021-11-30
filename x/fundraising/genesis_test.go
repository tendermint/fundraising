package fundraising_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/tendermint/fundraising/testutil/keeper"
	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		// this line is used by starport scaffolding # genesis/test/state
	}

	ctx, k := keepertest.Fundraising(t)
	fundraising.InitGenesis(ctx, *k, genesisState)
	got := fundraising.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}
