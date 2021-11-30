package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/tendermint/fundraising/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

var (
	moduleAccountPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
	}
)

// Fundraising returns a keeper of the fundraising module for testing purpose.
func Fundraising(t testing.TB) (sdk.Context, *keeper.Keeper) {
	cdc := simapp.Codec()
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)

	paramKeeper := initParam(cdc, db, stateStore)
	authKeeper := initAuth(cdc, db, stateStore, paramKeeper)
	bankKeeper := initBank(cdc, db, stateStore, paramKeeper, authKeeper)
	fundraisingKeeper := initFundraising(cdc, db, stateStore, paramKeeper, authKeeper, bankKeeper)
	require.NoError(t, stateStore.LoadLatestVersion())

	return sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger()), fundraisingKeeper
}

func initParam(
	cdc codec.Codec,
	db *tmdb.MemDB,
	stateStore store.CommitMultiStore,
) paramskeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkeys := sdk.NewTransientStoreKey(paramstypes.TStoreKey)

	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeys, sdk.StoreTypeTransient, db)

	return paramskeeper.NewKeeper(cdc, types.Amino, storeKey, tkeys)
}

func initAuth(
	cdc codec.Codec,
	db *tmdb.MemDB,
	stateStore store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
) authkeeper.AccountKeeper {
	storeKey := sdk.NewKVStoreKey(authtypes.StoreKey)

	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)

	paramKeeper.Subspace(authtypes.ModuleName)
	authSubspace, _ := paramKeeper.GetSubspace(authtypes.ModuleName)

	return authkeeper.NewAccountKeeper(cdc, storeKey, authSubspace, authtypes.ProtoBaseAccount, moduleAccountPerms)
}

func initBank(
	cdc codec.Codec,
	db *tmdb.MemDB,
	stateStore store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
	authKeeper authkeeper.AccountKeeper,
) bankkeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(banktypes.StoreKey)

	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)

	paramKeeper.Subspace(banktypes.ModuleName)
	bankSubspace, _ := paramKeeper.GetSubspace(banktypes.ModuleName)

	// module account addresses
	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return bankkeeper.NewBaseKeeper(cdc, storeKey, authKeeper, bankSubspace, modAccAddrs)
}

func initFundraising(
	cdc codec.Codec,
	db *tmdb.MemDB,
	stateStore store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) *keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)

	paramSubspace, _ := paramKeeper.GetSubspace(types.ModuleName)

	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return keeper.NewKeeper(cdc, storeKey, memStoreKey, paramSubspace, accountKeeper, bankKeeper, modAccAddrs)
}
