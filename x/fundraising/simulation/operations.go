package simulation

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/tendermint/spm/cosmoscmd"

	appparams "github.com/tendermint/fundraising/app/params"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Simulation operation weights constants.
const (
	OpWeightMsgCreateFixedPriceAuction = "op_weight_msg_create_fixed_price_auction"
	OpWeightMsgCreateBatchAuction      = "op_weight_msg_create_batch_auction"
	OpWeightMsgPlaceBid                = "op_weight_msg_place_bid"
	OpWeightMsgCancelAuction           = "op_weight_msg_cancel_auction"
)

var (
	Gas  = uint64(20000000)
	Fees = sdk.Coins{
		{
			Denom:  "stake",
			Amount: sdk.NewInt(0),
		},
	}
)

func init() {
	keeper.EnableAddAllowedBidder = true
}

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgCreateFixedPriceAuction int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateFixedPriceAuction, &weightMsgCreateFixedPriceAuction, nil,
		func(_ *rand.Rand) {
			weightMsgCreateFixedPriceAuction = appparams.DefaultWeightMsgCreateFixedPriceAuction
		},
	)

	var weightMsgCreateBatchAuction int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateBatchAuction, &weightMsgCreateBatchAuction, nil,
		func(_ *rand.Rand) {
			weightMsgCreateBatchAuction = appparams.DefaultWeightMsgCreateBatchAuction
		},
	)

	var weightMsgPlaceBid int
	appParams.GetOrGenerate(cdc, OpWeightMsgPlaceBid, &weightMsgPlaceBid, nil,
		func(_ *rand.Rand) {
			weightMsgPlaceBid = appparams.DefaultWeightMsgPlaceBid
		},
	)

	var weightMsgCancelAuction int
	appParams.GetOrGenerate(cdc, OpWeightMsgCancelAuction, &weightMsgCancelAuction, nil,
		func(_ *rand.Rand) {
			weightMsgCancelAuction = appparams.DefaultWeightMsgCancelAuction
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateFixedPriceAuction,
			SimulateMsgCreateFixedPriceAuction(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCreateBatchAuction,
			SimulateMsgCreateBatchAuction(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgPlaceBid,
			SimulateMsgPlaceBid(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgCancelAuction,
			SimulateMsgCancelAuction(ak, bk, k),
		),
	}
}

// SimulateMsgCreateFixedPriceAuction generates a MsgCreateFixedAmountPlan with random values
// nolint: interfacer
func SimulateMsgCreateFixedPriceAuction(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		denoms := genDenoms(r)
		fundAccountsOnce(r, ctx, bk, accs, denoms)

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)

		_, hasNeg := spendable.SafeSub(params.AuctionCreationFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedPriceAuction, "insufficient balance for auction creation fee"), nil, nil
		}

		auctioneer := account.GetAddress().String()
		ranDenom := shuffleDenoms(denoms)
		startPrice := sdk.NewDecFromInt(simtypes.RandomAmount(r, sdk.NewInt(10)))
		sellingCoin := sdk.NewCoin(ranDenom, simtypes.RandomAmount(r, sdk.NewInt(100_000_000_000_000)))
		payingCoinDenom := sdk.DefaultBondDenom
		vestingSchedules := []types.VestingSchedule{}
		startTime := ctx.BlockTime()
		endTime := ctx.BlockTime().AddDate(0, simtypes.RandIntBetween(r, 1, 24), 0)

		// TODO: is this logic reasonable to have?
		for _, auction := range k.GetAuctions(ctx) {
			if auction.GetSellingCoin().Denom == sellingCoin.Denom {
				return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedPriceAuction, "auction already exists with the same selling coin denom"), nil, nil
			}
		}

		msg := types.NewMsgCreateFixedPriceAuction(
			auctioneer,
			startPrice,
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			startTime,
			endTime,
		)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           cosmoscmd.MakeEncodingConfig(simapp.ModuleBasics).TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgCreateBatchAuction generates a MsgCreateRatioPlan with random values
// nolint: interfacer
func SimulateMsgCreateBatchAuction(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		denoms := genDenoms(r)
		fundAccountsOnce(r, ctx, bk, accs, denoms)

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)

		_, hasNeg := spendable.SafeSub(params.AuctionCreationFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedPriceAuction, "insufficient balance for auction creation fee"), nil, nil
		}

		auctioneer := account.GetAddress().String()
		ranDenom := shuffleDenoms(denoms)
		startPrice := sdk.NewDecFromInt(simtypes.RandomAmount(r, sdk.NewInt(10)))
		minBidPrice := simtypes.RandomDecAmount(r, startPrice)
		sellingCoin := sdk.NewCoin(ranDenom, simtypes.RandomAmount(r, sdk.NewInt(100_000_000_000_000)))
		payingCoinDenom := sdk.DefaultBondDenom
		vestingSchedules := []types.VestingSchedule{}
		maxExtendedRound := uint32(simtypes.RandIntBetween(r, 1, 5))
		extendedRoundRate := sdk.NewDecFromIntWithPrec(simtypes.RandomAmount(r, sdk.NewInt(30)), 2)
		startTime := ctx.BlockTime()
		endTime := ctx.BlockTime().AddDate(0, simtypes.RandIntBetween(r, 1, 24), 0)

		// TODO: is this logic reasonable to have?
		for _, auction := range k.GetAuctions(ctx) {
			if auction.GetSellingCoin().Denom == sellingCoin.Denom {
				return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedPriceAuction, "auction already exists with the same selling coin denom"), nil, nil
			}
		}

		msg := types.NewMsgCreateBatchAuction(
			auctioneer,
			startPrice,
			minBidPrice,
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			maxExtendedRound,
			extendedRoundRate,
			startTime,
			endTime,
		)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           cosmoscmd.MakeEncodingConfig(simapp.ModuleBasics).TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgPlaceBid generates a MsgStake with random values
// nolint: interfacer
func SimulateMsgPlaceBid(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO: Bidder must be in allowed bidder list
		// Depending on auction type, place a bid accordingly

		auctions := k.GetAuctions(ctx)
		if len(auctions) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPlaceBid, "no auction to place a bid"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		bidder := account.GetAddress().String()
		auction := auctions[simtypes.RandIntBetween(r, 0, len(auctions))]

		var bidType types.BidType
		var price sdk.Dec
		var coin sdk.Coin

		switch auction.GetType() {
		case types.AuctionTypeFixedPrice:
			bidType = types.BidTypeFixedPrice
			price = auction.GetStartPrice()
			coin = sdk.NewCoin(auction.GetPayingCoinDenom(), simtypes.RandomAmount(r, sdk.NewInt(500_000_000)))
		case types.AuctionTypeBatch:
			// TODO: randomize worth or many
			bidType = types.BidTypeBatchWorth
		default:
		}

		// TODO: check if bidder has sufficient amount of coin to place a bid

		msg := types.NewMsgPlaceBid(
			auction.GetId(),
			bidder,
			bidType,
			price,
			coin,
		)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           cosmoscmd.MakeEncodingConfig(simapp.ModuleBasics).TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgCancelAuction generates a SimulateMsgCancelAuction with random values
// nolint: interfacer
func SimulateMsgCancelAuction(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		auctions := k.GetAuctions(ctx)
		r.Shuffle(len(auctions), func(i, j int) {
			auctions[i], auctions[j] = auctions[j], auctions[i]
		})

		var simAccount simtypes.Account
		var auction types.AuctionI
		skip := true

		// Find an auction that is not started yet
		for _, a := range auctions {
			if !a.ShouldAuctionStarted(ctx.BlockTime()) {
				auction = a
				skip = false
				break
			}
		}
		if skip {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCancelAuction, "no auction to cancel"), nil, nil
		}

		accs = shuffleSimAccounts(r, accs)

		// Only the auctioneer can cancel the auction
		for _, acc := range accs {
			if acc.Address.Equals(auction.GetAuctioneer()) {
				simAccount = acc
			}
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		auctioneer := account.GetAddress().String()

		msg := types.NewMsgCancelAuction(auctioneer, auction.GetId())

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           cosmoscmd.MakeEncodingConfig(simapp.ModuleBasics).TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

var once sync.Once

func genDenoms(r *rand.Rand) []string {
	denoms := []string{}
	for i := 1; i <= simtypes.RandIntBetween(r, 10, 20); i++ {
		denoms = append(denoms, "denom"+fmt.Sprint(i))
	}
	return denoms
}

func fundAccountsOnce(r *rand.Rand, ctx sdk.Context, bk types.BankKeeper, accs []simtypes.Account, denoms []string) {
	once.Do(func() {
		maxAmt := sdk.NewInt(1_000_000_000_000_000)
		for _, acc := range accs {
			var coins sdk.Coins
			for _, denom := range denoms {
				coins = coins.Add(sdk.NewCoin(denom, simtypes.RandomAmount(r, maxAmt)))
			}
			if err := bk.MintCoins(ctx, minttypes.ModuleName, coins); err != nil {
				panic(err)
			}
			if err := bk.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, acc.Address, coins); err != nil {
				panic(err)
			}
		}
	})
}

// shuffleDenoms returns randomly shuffled denom.
func shuffleDenoms(denoms []string) string {
	denoms2 := make([]string, len(denoms))
	copy(denoms2, denoms)
	rand.Shuffle(len(denoms), func(i, j int) {
		denoms2[i], denoms2[j] = denoms2[j], denoms2[i]
	})
	return denoms2[0]
}

// shuffleSimAccounts returns randomly shuffled simulation accounts.
func shuffleSimAccounts(r *rand.Rand, accs []simtypes.Account) []simtypes.Account {
	accs2 := make([]simtypes.Account, len(accs))
	copy(accs2, accs)
	r.Shuffle(len(accs2), func(i, j int) {
		accs2[i], accs2[j] = accs2[j], accs2[i]
	})
	return accs2
}
