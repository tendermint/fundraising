package simulation

import (
	"math/rand"
	"sync"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
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
	OpWeightMsgModifyBid               = "op_weight_msg_modify_bid"
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

	var weightMsgModifyBid int
	appParams.GetOrGenerate(cdc, OpWeightMsgModifyBid, &weightMsgModifyBid, nil,
		func(_ *rand.Rand) {
			weightMsgModifyBid = appparams.DefaultWeightMsgModifyBid
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
		simulation.NewWeightedOperation(
			weightMsgModifyBid,
			SimulateMsgModifyBid(ak, bk, k),
		),
	}
}

// SimulateMsgCreateFixedPriceAuction generates a MsgCreateFixedAmountPlan with random values
// nolint: interfacer
func SimulateMsgCreateFixedPriceAuction(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		fundAccountsOnce(r, ctx, bk, accs)

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)
		_, hasNeg := spendable.SafeSub(params.AuctionCreationFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedPriceAuction, "insufficient balance for auction creation fee"), nil, nil
		}

		auctioneerAcc := account.GetAddress()

		// mint pool coins to simulate the real-world cases
		// funds, err := fundBalances(ctx, r, bk, creatorAcc, testCoinDenoms)
		// if err != nil {
		// 	return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedAmountPlan, "unable to mint pool coins"), nil, nil
		// }

		startPrice := sdk.ZeroDec()
		sellingCoin := sdk.NewInt64Coin("", 1)
		payingCoinDenom := ""
		vestingSchedules := []types.VestingSchedule{}
		startTime := ctx.BlockTime()
		endTime := startTime.AddDate(0, simtypes.RandIntBetween(r, 1, 24), 0)

		msg := types.NewMsgCreateFixedPriceAuction(
			auctioneerAcc.String(),
			startPrice,
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			startTime,
			endTime,
		)

		encoding := cosmoscmd.MakeEncodingConfig(simapp.ModuleBasics)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           encoding.TxConfig,
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

		return simtypes.OperationMsg{}, []simtypes.FutureOperation{}, nil
	}
}

// SimulateMsgPlaceBid generates a MsgStake with random values
// nolint: interfacer
func SimulateMsgPlaceBid(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		return simtypes.OperationMsg{}, []simtypes.FutureOperation{}, nil
	}
}

// SimulateMsgCancelAuction generates a SimulateMsgCancelAuction with random values
// nolint: interfacer
func SimulateMsgCancelAuction(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		return simtypes.OperationMsg{}, []simtypes.FutureOperation{}, nil
	}
}

// SimulateMsgModifyBid generates a MsgHarvest with random values
// nolint: interfacer
func SimulateMsgModifyBid(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		return simtypes.OperationMsg{}, []simtypes.FutureOperation{}, nil
	}
}

var once sync.Once

func fundAccountsOnce(r *rand.Rand, ctx sdk.Context, bk types.BankKeeper, accs []simtypes.Account) {
	once.Do(func() {
		denoms := []string{"denom1", "denom2", "denom3", "denom4", "denom5", "denom6"}
		maxAmt := sdk.NewInt(1_000_000_000_000_000)
		for _, acc := range accs {
			var coins sdk.Coins
			for _, denom := range denoms {
				coins = coins.Add(sdk.NewCoin(denom, simtypes.RandomAmount(r, maxAmt)))
			}
			if err := bk.MintCoins(ctx, types.ModuleName, coins); err != nil {
				panic(err)
			}
			if err := bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acc.Address, coins); err != nil {
				panic(err)
			}
		}
	})
}
