package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/tendermint/fundraising/testutil/sample"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Fundraising returns a keeper of the fundraising module for testing purpose.
func Fundraising(t testing.TB) (*keeper.Keeper, sdk.Context) {
	cdc := sample.Codec()
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)

	keeper := initFundraising(cdc, db, stateStore)
	require.NoError(t, stateStore.LoadLatestVersion())

	return keeper, sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
}

func initFundraising(cdc codec.Codec, db *tmdb.MemDB, stateStore store.CommitMultiStore) *keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)

	return keeper.NewKeeper(cdc, storeKey, memStoreKey)
}
