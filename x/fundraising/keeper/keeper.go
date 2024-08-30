package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		addressCodec address.Codec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message.
		// Typically, this should be the x/gov module account.
		authority string

		Schema         collections.Schema
		Params         collections.Item[types.Params]
		MatchedBidsLen collections.Map[uint64, int64]
		AllowedBidder  collections.Map[collections.Pair[uint64, sdk.AccAddress], types.AllowedBidder]
		VestingQueue   collections.Map[collections.Pair[uint64, time.Time], types.VestingQueue]
		BidSeq         collections.Map[uint64, uint64]
		Bid            collections.Map[collections.Pair[uint64, uint64], types.Bid]
		AuctionSeq     collections.Sequence
		Auction        collections.Map[uint64, types.AuctionI]
		// this line is used by starport scaffolding # collection/type

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		distrKeeper   types.DistrKeeper

		hooks types.FundraisingHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistrKeeper,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:            cdc,
		addressCodec:   addressCodec,
		storeService:   storeService,
		authority:      authority,
		logger:         logger,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		distrKeeper:    distrKeeper,
		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		MatchedBidsLen: collections.NewMap(sb, types.MatchedBidsLenKey, "matchedBidsLen", collections.Uint64Key, collections.Int64Value),
		AllowedBidder:  collections.NewMap(sb, types.AllowedBidderKey, "allowedBidder", collections.PairKeyCodec(collections.Uint64Key, sdk.LengthPrefixedAddressKey(sdk.AccAddressKey)), codec.CollValue[types.AllowedBidder](cdc)),
		VestingQueue:   collections.NewMap(sb, types.VestingQueueKey, "vestingQueue", collections.PairKeyCodec(collections.Uint64Key, sdk.TimeKey), codec.CollValue[types.VestingQueue](cdc)),
		BidSeq:         collections.NewMap(sb, types.BidCountKey, "bid_seq", collections.Uint64Key, collections.Uint64Value),
		Bid:            collections.NewMap(sb, types.BidKey, "bid", collections.PairKeyCodec(collections.Uint64Key, collections.Uint64Key), codec.CollValue[types.Bid](cdc)),
		AuctionSeq:     collections.NewSequence(sb, types.AuctionCountKey, "auction_seq"),
		Auction:        collections.NewMap(sb, types.AuctionKey, "auction", collections.Uint64Key, codec.CollInterfaceValue[types.AuctionI](cdc)),
		// this line is used by starport scaffolding # collection/instantiate
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// PayCreationFee sends the auction creation fee to the fee collector account.
func (k Keeper) PayCreationFee(ctx context.Context, auctioneerAddr sdk.AccAddress) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	if err := k.distrKeeper.FundCommunityPool(ctx, params.AuctionCreationFee, auctioneerAddr); err != nil {
		return err
	}
	return nil
}

// PayPlaceBidFee sends the fee when placing a bid for an auction to the fee collector account.
func (k Keeper) PayPlaceBidFee(ctx context.Context, bidderAddr sdk.AccAddress) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	if err := k.distrKeeper.FundCommunityPool(ctx, params.PlaceBidFee, bidderAddr); err != nil {
		return err
	}
	return nil
}

// ReserveSellingCoin reserves the selling coin to the selling reserve account.
func (k Keeper) ReserveSellingCoin(ctx context.Context, auctionId uint64, auctioneerAddr sdk.AccAddress, sellingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, auctioneerAddr, types.SellingReserveAddress(auctionId), sdk.NewCoins(sellingCoin)); err != nil {
		return err
	}
	return nil
}

// ReservePayingCoin reserves paying coin to the paying reserve account.
func (k Keeper) ReservePayingCoin(ctx context.Context, auctionId uint64, bidderAddr sdk.AccAddress, payingCoin sdk.Coin) error {
	if err := k.bankKeeper.SendCoins(ctx, bidderAddr, types.PayingReserveAddress(auctionId), sdk.NewCoins(payingCoin)); err != nil {
		return err
	}
	return nil
}
