package simulation

import (
	"math/rand"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// SimulateMsgCreateFixedPriceAuction generates a MsgCreateFixedAmountPlan with random values
// nolint: interfacer
func SimulateMsgCreateFixedPriceAuction(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgCreateFixedPriceAuction{}
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params, err := k.Params.Get(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "failed to get params"), nil, err
		}

		_, hasNeg := spendable.SafeSub(params.AuctionCreationFee...)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "insufficient balance for auction creation fee"), nil, nil
		}

		auctioneer := account.GetAddress()
		startPrice := math.LegacyNewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 10)), 1) // 0.1 ~ 1.0
		sellingCoin := sdk.NewInt64Coin(testCoinDenoms[r.Intn(len(testCoinDenoms))], int64(simtypes.RandIntBetween(r, 10000000000, 1000000000000)))
		payingCoinDenom := sdk.DefaultBondDenom
		vestingSchedules := make([]types.VestingSchedule, 0)
		startTime := ctx.BlockTime().AddDate(0, 0, simtypes.RandIntBetween(r, 0, 2))
		endTime := startTime.AddDate(0, simtypes.RandIntBetween(r, 1, 12), 0)

		if _, err := fundBalances(ctx, r, bk, auctioneer, testCoinDenoms); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "failed to fund auctioneer"), nil, err
		}

		// Call spendable coins here again to get the funded balances
		_, hasNeg = bk.SpendableCoins(ctx, account.GetAddress()).SafeSub(sdk.NewCoins(sellingCoin)...)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "insufficient balance to reserve selling coin"), nil, nil
		}

		msg = types.NewMsgCreateFixedPriceAuction(
			auctioneer.String(),
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
			TxGen:           txGen,
			Cdc:             nil,
			Msg:             msg,
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
