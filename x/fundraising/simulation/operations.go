package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
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
	OpWeightMsgCancelAuction           = "op_weight_msg_cancel_auction"
	OpWeightMsgPlaceBid                = "op_weight_msg_place_bid"
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

var (
	testCoinDenoms = []string{
		"denoma",
		"denomb",
		"denomc",
		"denomd",
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

	var weightMsgCancelAuction int
	appParams.GetOrGenerate(cdc, OpWeightMsgCancelAuction, &weightMsgCancelAuction, nil,
		func(_ *rand.Rand) {
			weightMsgCancelAuction = appparams.DefaultWeightMsgCancelAuction
		},
	)

	var weightMsgPlaceBid int
	appParams.GetOrGenerate(cdc, OpWeightMsgPlaceBid, &weightMsgPlaceBid, nil,
		func(_ *rand.Rand) {
			weightMsgPlaceBid = appparams.DefaultWeightMsgPlaceBid
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
			weightMsgCancelAuction,
			SimulateMsgCancelAuction(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgPlaceBid,
			SimulateMsgPlaceBid(ak, bk, k),
		),
	}
}

// SimulateMsgCreateFixedPriceAuction generates a MsgCreateFixedAmountPlan with random values
// nolint: interfacer
func SimulateMsgCreateFixedPriceAuction(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)
		_, hasNeg := spendable.SafeSub(params.AuctionCreationFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedPriceAuction, "insufficient balance for auction creation fee"), nil, nil
		}

		auctioneer := account.GetAddress()
		startPrice := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 10)), 1) // 0.1 ~ 1.0
		sellingCoin := sdk.NewInt64Coin(testCoinDenoms[r.Intn(len(testCoinDenoms))], int64(simtypes.RandIntBetween(r, 100000000000, 100000000000000)))
		payingCoinDenom := sdk.DefaultBondDenom
		vestingSchedules := []types.VestingSchedule{}
		startTime := ctx.BlockTime()
		endTime := ctx.BlockTime().AddDate(0, 0, 1+r.Intn(5))

		if _, err := fundBalances(ctx, r, bk, auctioneer, testCoinDenoms); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedPriceAuction, "failed to fund auctioneer"), nil, err
		}

		_, hasNeg = spendable.SafeSub(sdk.NewCoins(sellingCoin))
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFixedPriceAuction, "insufficient balance for auction creation"), nil, nil
		}

		msg := types.NewMsgCreateFixedPriceAuction(
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

		return genAndDeliverTxWithFees(txCtx, Gas, Fees)
	}
}

// SimulateMsgCreateBatchAuction generates a MsgCreateRatioPlan with random values
// nolint: interfacer
func SimulateMsgCreateBatchAuction(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		params := k.GetParams(ctx)
		_, hasNeg := spendable.SafeSub(params.AuctionCreationFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateBatchAuction, "insufficient balance for auction creation fee"), nil, nil
		}

		auctioneer := account.GetAddress()
		startPrice := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 10)), 1) // 0.1 ~ 1.0
		minBidPrice := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 10)), 2)
		sellingCoin := sdk.NewInt64Coin(testCoinDenoms[r.Intn(len(testCoinDenoms))], int64(simtypes.RandIntBetween(r, 100000000000, 100000000000000)))
		payingCoinDenom := sdk.DefaultBondDenom
		vestingSchedules := []types.VestingSchedule{}
		maxExtendedRound := uint32(simtypes.RandIntBetween(r, 1, 5))
		extendedRoundRate := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 3)), 1) // 0.1 ~ 0.3
		startTime := ctx.BlockTime()
		endTime := ctx.BlockTime().AddDate(0, simtypes.RandIntBetween(r, 1, 24), 0)

		if _, err := fundBalances(ctx, r, bk, auctioneer, testCoinDenoms); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateBatchAuction, "failed to fund auctioneer"), nil, err
		}

		_, hasNeg = spendable.SafeSub(sdk.NewCoins(sellingCoin))
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateBatchAuction, "insufficient balance for auction creation"), nil, nil
		}

		msg := types.NewMsgCreateBatchAuction(
			auctioneer.String(),
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

		return genAndDeliverTxWithFees(txCtx, Gas, Fees)
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

		return genAndDeliverTxWithFees(txCtx, Gas, Fees)
	}
}

// SimulateMsgPlaceBid generates a MsgStake with random values
// nolint: interfacer
func SimulateMsgPlaceBid(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		auctions := k.GetAuctions(ctx)
		if len(auctions) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPlaceBid, "no auction to place a bid"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		// Select a random auction
		auction := auctions[simtypes.RandIntBetween(r, 0, len(auctions))]

		if auction.GetStatus() != types.AuctionStatusStarted {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPlaceBid, "auction must be started status"), nil, nil
		}

		bidder := account.GetAddress().String()

		// Random number to be used to choose bid type
		n := r.Int()

		var bid types.Bid
		switch auction.GetType() {
		case types.AuctionTypeFixedPrice:
			bid.Type = types.BidTypeFixedPrice
			bid.Price = auction.GetStartPrice()
			if n%2 == 0 {
				bid.Coin = sdk.NewInt64Coin(auction.GetPayingCoinDenom(), int64(simtypes.RandIntBetween(r, 100000, 1000000000)))
			} else {
				bid.Coin = sdk.NewInt64Coin(auction.GetSellingCoin().Denom, int64(simtypes.RandIntBetween(r, 100000, 1000000000)))
			}
		case types.AuctionTypeBatch:
			bid.Price = auction.GetStartPrice().Add(sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 10)), 1)) // 0.1 ~ 1.0
			if n%2 == 0 {
				bid.Type = types.BidTypeBatchWorth
				bid.Coin = sdk.NewInt64Coin(auction.GetPayingCoinDenom(), int64(simtypes.RandIntBetween(r, 100000, 1000000000)))
			} else {
				bid.Type = types.BidTypeBatchMany
				bid.Coin = sdk.NewInt64Coin(auction.GetSellingCoin().Denom, int64(simtypes.RandIntBetween(r, 100000, 1000000000)))
			}
		}

		if !sdk.NewCoins(bid.Coin).IsAllLTE(spendable) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPlaceBid, "insufficient funds"), nil, nil
		}

		// Increase the bidder's maximum bid amount if they had previous bids
		maxBidAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())
		_, found := k.GetAuction(ctx, auction.GetId())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPlaceBid, "auction not found"), nil, nil
		}

		prevMaxBidAmt, found := auction.GetAllowedBiddersMap()[bidder]
		if found {
			maxBidAmt = maxBidAmt.Add(prevMaxBidAmt)
		}

		k.AddAllowedBidders(ctx, auction.GetId(), []types.AllowedBidder{
			{Bidder: bidder, MaxBidAmount: maxBidAmt},
		})

		msg := types.NewMsgPlaceBid(auction.GetId(), bidder, bid.Type, bid.Price, bid.Coin)

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

		fmt.Println("genAndDeliverTxWithFees")

		return genAndDeliverTxWithFees(txCtx, Gas, Fees)
	}
}

// fundBalances mints random amount of coins with the provided coin denoms and
// send them to the simulated account.
func fundBalances(ctx sdk.Context, r *rand.Rand, bk types.BankKeeper, acc sdk.AccAddress, denoms []string) (mintCoins sdk.Coins, err error) {
	for _, denom := range denoms {
		mintCoins = mintCoins.Add(sdk.NewInt64Coin(denom, int64(simtypes.RandIntBetween(r, 1e14, 1e15))))
	}

	if err := bk.MintCoins(ctx, minttypes.ModuleName, mintCoins); err != nil {
		return nil, err
	}

	if err := bk.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, acc, mintCoins); err != nil {
		return nil, err
	}
	return mintCoins, nil
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

// ranDenom returns randomly shuffled denom.
func ranDenom(denoms []string) string {
	denoms2 := make([]string, len(denoms))
	copy(denoms2, denoms)
	rand.Shuffle(len(denoms), func(i, j int) {
		denoms2[i], denoms2[j] = denoms2[j], denoms2[i]
	})
	return denoms2[0]
}

// genAndDeliverTxWithFees generates a transaction with given fee and delivers it.
func genAndDeliverTxWithFees(txCtx simulation.OperationInput, gas uint64, fees sdk.Coins) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
	account := txCtx.AccountKeeper.GetAccount(txCtx.Context, txCtx.SimAccount.Address)
	spendable := txCtx.Bankkeeper.SpendableCoins(txCtx.Context, account.GetAddress())

	var err error

	_, hasNeg := spendable.SafeSub(txCtx.CoinsSpentInMsg)
	if hasNeg {
		return simtypes.NoOpMsg(txCtx.ModuleName, txCtx.MsgType, "message doesn't leave room for fees"), nil, err
	}

	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, txCtx.MsgType, "unable to generate fees"), nil, err
	}
	return genAndDeliverTx(txCtx, fees, gas)
}

// genAndDeliverTx generates a transactions and delivers it.
func genAndDeliverTx(txCtx simulation.OperationInput, fees sdk.Coins, gas uint64) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
	account := txCtx.AccountKeeper.GetAccount(txCtx.Context, txCtx.SimAccount.Address)
	tx, err := helpers.GenTx(
		txCtx.TxGen,
		[]sdk.Msg{txCtx.Msg},
		fees,
		gas,
		txCtx.Context.ChainID(),
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		txCtx.SimAccount.PrivKey,
	)

	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, txCtx.MsgType, "unable to generate mock tx"), nil, err
	}

	_, _, err = txCtx.App.Deliver(txCtx.TxGen.TxEncoder(), tx)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, txCtx.MsgType, "unable to deliver tx"), nil, err
	}

	return simtypes.NewOperationMsg(txCtx.Msg, true, "", txCtx.Cdc), nil, nil

}
