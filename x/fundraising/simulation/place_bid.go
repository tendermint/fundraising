package simulation

import (
	"math/rand"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// SimulateMsgPlaceBid generates a MsgPlaceBid with random values
// nolint: interfacer
func SimulateMsgPlaceBid(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgPlaceBid{}
		auctions, err := k.Auctions(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "failed to get auctions"), nil, nil
		}
		if len(auctions) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "no auction to place a bid"), nil, nil
		}

		// Select a random auction
		auction := auctions[r.Intn(len(auctions))]

		if auction.GetStatus() != types.AuctionStatusStarted {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "auction must be started status"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		bidder := account.GetAddress()
		sellingCoinDenom := auction.GetSellingCoin().Denom
		payingCoinDenom := auction.GetPayingCoinDenom()

		if _, err := fundBalances(ctx, r, bk, bidder, testCoinDenoms); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "failed to fund bidder"), nil, err
		}

		var bid types.Bid
		switch auction.GetType() {
		case types.AuctionTypeFixedPrice:
			bid.Type = types.BidTypeFixedPrice
			bid.Price = auction.GetStartPrice()
			if r.Int()%2 == 0 {
				bid.Coin = sdk.NewInt64Coin(payingCoinDenom, int64(simtypes.RandIntBetween(r, 100000, 1000000000)))
			} else {
				bid.Coin = sdk.NewInt64Coin(sellingCoinDenom, int64(simtypes.RandIntBetween(r, 100000, 1000000000)))
			}
		case types.AuctionTypeBatch:
			bid.Price = auction.GetStartPrice().Add(math.LegacyNewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 5)), 1)) // StartPrice + 0.1 ~ 0.5
			if r.Int()%2 == 0 {
				bid.Type = types.BidTypeBatchWorth
				bid.Coin = sdk.NewInt64Coin(payingCoinDenom, int64(simtypes.RandIntBetween(r, 100000, 1000000000)))
			} else {
				bid.Type = types.BidTypeBatchMany
				bid.Coin = sdk.NewInt64Coin(sellingCoinDenom, int64(simtypes.RandIntBetween(r, 100000, 1000000000)))
			}
		}

		bidReserveAmt := bid.ConvertToPayingAmount(payingCoinDenom)
		maxBidAmt := bid.ConvertToSellingAmount(payingCoinDenom)

		if !bk.SpendableCoins(ctx, bidder).AmountOf(payingCoinDenom).GT(bidReserveAmt) {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "insufficient balance to place a bid"), nil, nil
		}

		allowedBidder, err := k.AllowedBidder.Get(ctx, collections.Join(auction.GetId(), bidder))
		if err == nil {
			maxBidAmt = maxBidAmt.Add(allowedBidder.MaxBidAmount)
		}

		newAllowedBidder := types.NewAllowedBidder(auction.GetId(), bidder, maxBidAmt)
		if err := k.AddAllowedBidders(ctx, auction.GetId(), []types.AllowedBidder{newAllowedBidder}); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "failed to add allowed bidders"), nil, err
		}

		msg = types.NewMsgPlaceBid(
			auction.GetId(),
			bidder.String(),
			bid.Type,
			bid.Price,
			bid.Coin,
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
