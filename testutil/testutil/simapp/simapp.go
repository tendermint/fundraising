package simapp

// DONTCOVER

import (
	"context"
	"encoding/json"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/testutil/sample"
)

// New creates application instance with in-memory database and disabled logging.
func New(chainID string) (*app.App, error) {
	var (
		db     = dbm.NewMemDB()
		logger = log.NewNopLogger()
	)

	appOptions := simtestutil.AppOptionsMap{
		flags.FlagHome:            app.DefaultNodeHome,
		server.FlagInvCheckPeriod: simcli.FlagPeriodValue,
	}

	fundraisingApp, err := app.New(
		logger,
		db,
		nil,
		true,
		appOptions,
		fauxMerkleModeOpt,
		baseapp.SetChainID(chainID),
	)
	if err != nil {
		return nil, err
	}

	var (
		genesisState = fundraisingApp.DefaultGenesis()
		privVal      = mock.NewPV()
		pubKey, _    = privVal.GetPubKey()

		// create validator set with single validator
		validator = tmtypes.NewValidator(pubKey, 1)
		valSet    = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

		// generate genesis account
		senderPubKey = sample.PubKey()
		acc          = authtypes.NewBaseAccount(sdk.AccAddress(senderPubKey.Address()), senderPubKey, 0, 0)
		balances     = banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100000000000000))),
		}
	)
	// init chain must be called to stop deliverState from being nil
	genesisState = genesisStateWithValSet(fundraisingApp.AppCodec(), genesisState, valSet, []authtypes.GenesisAccount{acc}, balances)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		return nil, err
	}

	if _, err := fundraisingApp.InitChain(
		&abci.RequestInitChain{
			ChainId:       chainID,
			AppStateBytes: stateBytes,
		},
	); err != nil {
		return nil, err
	}

	return fundraisingApp, nil
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *app.App, ctx context.Context, accNum int, accAmt math.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
}

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func addTestAddrs(app *app.App, ctx context.Context, accNum int, accAmt math.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, accAmt))

	for _, addr := range testAddrs {
		initAccountWithCoins(app, ctx, addr, initCoins)
	}

	return testAddrs
}

func initAccountWithCoins(app *app.App, ctx context.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

// FundAccount is a utility function that funds an account by minting and
// sending the coins to the address. This should be used for testing purposes
// only!
func FundAccount(bankKeeper bankkeeper.Keeper, ctx context.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func genesisStateWithValSet(
	cdc codec.Codec,
	genesisState app.GenesisState,
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) app.GenesisState {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = cdc.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, _ := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		pkAny, _ := codectypes.NewAnyWithValue(pk)
		valAddr := sdk.ValAddress(val.Address).String()
		validator := stakingtypes.Validator{
			OperatorAddress:   valAddr,
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), valAddr, math.LegacyOneDec()))

	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	)
	genesisState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesis)

	return genesisState
}
