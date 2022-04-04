package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

var (
	// Set this to "true" using testing build flag to enable AddAllowedBidder msg in Makefile
	enableAddAllowedBidder = "false"

	// EnableAddAllowedBidder indicates whether msgServer accepts MsgAddAllowedBidder or not.
	// Never set this to true in production environment. Doing that will expose serious attack vector.
	// Default is false, which means AddAllowedBidder can't be executed through message level.
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
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	memKey        sdk.StoreKey
	paramSpace    paramtypes.Subspace
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper
	hooks         types.FundraisingHooks
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	memKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistrKeeper,
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
		distrKeeper:   distrKeeper,
	}
}

// SetHooks sets the fundraising hooks.
func (k *Keeper) SetHooks(fk types.FundraisingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set fundraising hooks twice")
	}

	k.hooks = fk

	return k
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

// ReserveCreationFee reserves the auction creation fee to the fee collector account.
func (k Keeper) ReserveCreationFee(ctx sdk.Context, auctioneerAddr sdk.AccAddress) error {
	params := k.GetParams(ctx)

	if err := k.distrKeeper.FundCommunityPool(ctx, params.AuctionCreationFee, auctioneerAddr); err != nil {
		return sdkerrors.Wrap(err, "failed to reserve auction creation fee to the community pool")
	}
	return nil
}

// ReserveSellingCoin reserves the selling coin to the selling reserve account.
func (k Keeper) ReserveSellingCoin(ctx sdk.Context, auctionId uint64, auctioneerAddr sdk.AccAddress, sellingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, auctioneerAddr, types.SellingReserveAddress(auctionId), sdk.NewCoins(sellingCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to reserve selling coin")
	}
	return nil
}

// ReservePayingCoin reserves paying coin to the paying reserve account.
func (k Keeper) ReservePayingCoin(ctx sdk.Context, auctionId uint64, bidderAddr sdk.AccAddress, payingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, bidderAddr, types.PayingReserveAddress(auctionId), sdk.NewCoins(payingCoin)); err != nil {
		return sdkerrors.Wrap(err, "failed to reserve paying coin")
	}
	return nil
}
