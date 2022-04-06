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

		auctioneerAcc := account.GetAddress()
		ranDenom := shuffleDenoms(denoms)
		startPrice := sdk.NewDecFromInt(simtypes.RandomAmount(r, sdk.NewInt(10)))
		sellingCoin := sdk.NewInt64Coin(ranDenom, 1)
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
			auctioneerAcc.String(),
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

		return simtypes.OperationMsg{}, []simtypes.FutureOperation{}, nil
	}
}

// SimulateMsgPlaceBid generates a MsgStake with random values
// nolint: interfacer
func SimulateMsgPlaceBid(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO:
		// Auction must exist
		// Bidder must be in allowed bidder list
		// Depending on auction type, place a bid accordingly

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
	return denoms[0]
}
