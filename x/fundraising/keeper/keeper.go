package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

var (
	// Set this to "true" using testing build flag to enable AddAllowedBidder msg in Makefile
	enableAddAllowedBidder = "false"

	// EnableAddAllowedBidder indicates whether msgServer accepts MsgAddAllowedBidder or not.
	// Never set this to true in production environment. Doing that will expose serious attack vector.
	EnableAddAllowedBidder = false
)

func init() {
	var err error
	EnableAddAllowedBidder, err = strconv.ParseBool(enableAddAllowedBidder)
	if err != nil {
		panic(err)
	}
}

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	memKey     sdk.StoreKey
	paramSpace paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper

	blockedAddrs map[string]bool
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	memKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	blockedAddrs map[string]bool,
) *Keeper {
	// Ensure fundraising module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// Set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      key,
		memKey:        memKey,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		blockedAddrs:  blockedAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams returns the parameters for the fundraising module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the parameters for the fundraising module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetCodec returns codec.Codec object used by the keeper.
func (k Keeper) GetCodec() codec.BinaryCodec { return k.cdc }
