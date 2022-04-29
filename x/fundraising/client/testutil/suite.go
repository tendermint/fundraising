package testutil

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	utilcli "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tm-db"

	chain "github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/x/fundraising/client/cli"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

type TxCmdTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	denom1 string
	denom2 string
}

func NewAppConstructor(encodingCfg cosmoscmd.EncodingConfig) network.AppConstructor {
	return func(val network.Validator) servertypes.Application {
		return chain.New(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			encodingCfg,
			simapp.EmptyAppOptions{},
			baseapp.SetPruning(store.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}
}

// SetupTest creates a new network for _each_ integration test. We create a new
// network for each test because there are some state modifications that are
// needed to be made in order to make useful queries. However, we don't want
// these state changes to be present in other tests.
func (s *TxCmdTestSuite) SetupTest() {
	s.T().Log("setting up integration test suite")

	keeper.EnableAddAllowedBidder = true

	encodingCfg := cosmoscmd.MakeEncodingConfig(chain.ModuleBasics)

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	cfg.AppConstructor = NewAppConstructor(encodingCfg)
	cfg.GenesisState = chain.ModuleBasics.DefaultGenesis(cfg.Codec)
	cfg.AccountTokens = sdk.NewInt(100_000_000_000_000) // node0token denom
	cfg.StakingTokens = sdk.NewInt(100_000_000_000_000) // stake denom

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)
	s.denom1, s.denom2 = fmt.Sprintf("%stoken", s.network.Validators[0].Moniker), s.cfg.BondDenom

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

// TearDownTest cleans up the current test network after each test in the suite.
func (s *TxCmdTestSuite) TearDownTest() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *TxCmdTestSuite) TestNewCreateFixedAmountPlanCmd() {
	val := s.network.Validators[0]

	startTime := time.Now()
	endTime := startTime.AddDate(0, 1, 0)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"valid case",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
					StartPrice:      sdk.MustNewDecFromStr("1.0"),
					SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom: s.denom2,
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("1.0"),
						},
					},
					StartTime: startTime,
					EndTime:   endTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid case #1: invalid end time",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
					StartPrice:      sdk.MustNewDecFromStr("1.0"),
					SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom: s.denom2,
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("1.0"),
						},
					},
					StartTime: startTime,
					EndTime:   startTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
		{
			"invalid case #2: invalid vesting schedule",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
					StartPrice:      sdk.MustNewDecFromStr("1.0"),
					SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom: s.denom2,
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: startTime,
							Weight:      sdk.MustNewDecFromStr("1.0"),
						},
					},
					StartTime: startTime,
					EndTime:   endTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
		{
			"invalid case #3: invalid vesting schedule",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
					StartPrice:      sdk.MustNewDecFromStr("1.0"),
					SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom: s.denom2,
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("5.0"),
						},
					},
					StartTime: startTime,
					EndTime:   endTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
		{
			"invalid case #4: invalid vesting schedule",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
					StartPrice:      sdk.MustNewDecFromStr("1.0"),
					SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom: s.denom2,
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("0.5"),
						},
						{
							ReleaseTime: endTime.AddDate(0, 6, 0),
							Weight:      sdk.MustNewDecFromStr("0.51"),
						},
					},
					StartTime: startTime,
					EndTime:   endTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
		{
			"invalid case #5: invalid vesting schedule",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
					StartPrice:      sdk.MustNewDecFromStr("1.0"),
					SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom: s.denom2,
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("0.5"),
						},
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("0.5"),
						},
					},
					StartTime: startTime,
					EndTime:   endTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCreateFixedPriceAuctionCmd()
			clientCtx := val.ClientCtx

			out, err := utilcli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewCreateBatchAuctionCmd() {
	val := s.network.Validators[0]

	startTime := time.Now()
	endTime := startTime.AddDate(0, 1, 0)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"valid case",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
					StartPrice:        sdk.MustNewDecFromStr("0.5"),
					MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
					SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom:   s.denom2,
					MaxExtendedRound:  2,
					ExtendedRoundRate: sdk.MustNewDecFromStr("0.15"),
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("1.0"),
						},
					},
					StartTime: startTime,
					EndTime:   endTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid case #1: invalid end time",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
					StartPrice:        sdk.MustNewDecFromStr("0.5"),
					MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
					SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom:   s.denom2,
					MaxExtendedRound:  2,
					ExtendedRoundRate: sdk.MustNewDecFromStr("0.15"),
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("1.0"),
						},
					},
					StartTime: startTime,
					EndTime:   startTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
		{
			"invalid case #2: invalid vesting schedule",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
					StartPrice:        sdk.MustNewDecFromStr("0.5"),
					MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
					SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom:   s.denom2,
					MaxExtendedRound:  2,
					ExtendedRoundRate: sdk.MustNewDecFromStr("0.15"),
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: startTime,
							Weight:      sdk.MustNewDecFromStr("1.0"),
						},
					},
					StartTime: startTime,
					EndTime:   startTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
		{
			"invalid case #3: invalid vesting schedule",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
					StartPrice:        sdk.MustNewDecFromStr("0.5"),
					MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
					SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom:   s.denom2,
					MaxExtendedRound:  2,
					ExtendedRoundRate: sdk.MustNewDecFromStr("0.15"),
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("5.0"),
						},
					},
					StartTime: startTime,
					EndTime:   startTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
		{
			"invalid case #4: invalid vesting schedule",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
					StartPrice:        sdk.MustNewDecFromStr("0.5"),
					MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
					SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom:   s.denom2,
					MaxExtendedRound:  2,
					ExtendedRoundRate: sdk.MustNewDecFromStr("0.15"),
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("0.5"),
						},
						{
							ReleaseTime: endTime.AddDate(0, 6, 0),
							Weight:      sdk.MustNewDecFromStr("0.51"),
						},
					},
					StartTime: startTime,
					EndTime:   startTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
		{
			"invalid case #5: invalid vesting schedule",
			[]string{
				testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
					StartPrice:        sdk.MustNewDecFromStr("0.5"),
					MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
					SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
					PayingCoinDenom:   s.denom2,
					MaxExtendedRound:  2,
					ExtendedRoundRate: sdk.MustNewDecFromStr("0.15"),
					VestingSchedules: []types.VestingSchedule{
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("0.5"),
						},
						{
							ReleaseTime: endTime.AddDate(0, 3, 0),
							Weight:      sdk.MustNewDecFromStr("0.5"),
						},
					},
					StartTime: startTime,
					EndTime:   startTime,
				}.String()).Name(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			true, &sdk.TxResponse{}, 1,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCreateBatchAuctionCmd()
			clientCtx := val.ClientCtx

			out, err := utilcli.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewCancelAuctionCmd() {
	val := s.network.Validators[0]

	// Create a fixed price auction
	_, err := MsgCreateFixedPriceAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
			StartPrice:      sdk.MustNewDecFromStr("1.0"),
			SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom: s.denom2,
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(1, 0, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now().AddDate(0, 1, 0),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"valid case",
			[]string{
				fmt.Sprint(1),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid case #1: auction not found",
			[]string{
				fmt.Sprint(5),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 38,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewCancelAuctionCmd()
			clientCtx := val.ClientCtx

			out, err := utilcli.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewPlaceBidCmd() {
	val := s.network.Validators[0]

	// Create a fixed price auction
	_, err := MsgCreateFixedPriceAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
			StartPrice:      sdk.MustNewDecFromStr("1.0"),
			SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom: s.denom2,
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(0, 6, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	// Add allowed bidder
	_, err = MsgAddAllowedBidderExec(
		val.ClientCtx,
		val.Address.String(),
		1,
		sdk.NewInt(100_000_000),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"valid case",
			[]string{
				fmt.Sprint(1),
				"fixed-price",
				sdk.MustNewDecFromStr("1.0").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid case #1: incorrect auction type",
			[]string{
				fmt.Sprint(1),
				"batch-worth",
				sdk.MustNewDecFromStr("1.0").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 7,
		},
		{
			"invalid case #2: incorrect auction type",
			[]string{
				fmt.Sprint(1),
				"batch-many",
				sdk.MustNewDecFromStr("1.0").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 7,
		},
		{
			"invalid case #3: incorrect start price",
			[]string{
				fmt.Sprint(1),
				"fixed-price",
				sdk.MustNewDecFromStr("0.1").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 3,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewPlaceBidCmd()
			clientCtx := val.ClientCtx

			out, err := utilcli.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewModifyBidCmd() {
	val := s.network.Validators[0]

	// Create a batch auction
	_, err := MsgCreateBatchAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
			StartPrice:        sdk.MustNewDecFromStr("0.5"),
			MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
			SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom:   s.denom2,
			MaxExtendedRound:  2,
			ExtendedRoundRate: sdk.MustNewDecFromStr("0.2"),
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(0, 6, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	// Add allowed bidder
	_, err = MsgAddAllowedBidderExec(
		val.ClientCtx,
		val.Address.String(),
		1,
		sdk.NewInt(100_000_000),
	)
	s.Require().NoError(err)

	// Place a bid
	_, err = MsgPlaceBidExec(
		val.ClientCtx,
		val.Address.String(),
		1,
		"batch-worth",
		sdk.MustNewDecFromStr("0.55"),
		sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"valid case",
			[]string{
				fmt.Sprint(1),
				fmt.Sprint(1),
				sdk.MustNewDecFromStr("0.6").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
		{
			"invalid case #1: bid price and coin amount must be changed",
			[]string{
				fmt.Sprint(1),
				fmt.Sprint(1),
				sdk.MustNewDecFromStr("0.5").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 18,
		},
		{
			"invalid case #2: bid not found",
			[]string{
				fmt.Sprint(1),
				fmt.Sprint(5),
				sdk.MustNewDecFromStr("0.5").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 38,
		},
		{
			"invalid case #3: min bid price",
			[]string{
				fmt.Sprint(1),
				fmt.Sprint(1),
				sdk.MustNewDecFromStr("0.05").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 13,
		},
		{
			"invalid case #4: incorrect denom",
			[]string{
				fmt.Sprint(1),
				fmt.Sprint(1),
				sdk.MustNewDecFromStr("0.6").String(),
				sdk.NewCoin(s.denom1, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 8,
		},
		{
			"invalid case #5: bid price or coin amount cannot be lower",
			[]string{
				fmt.Sprint(1),
				fmt.Sprint(1),
				sdk.MustNewDecFromStr("0.2").String(),
				sdk.NewCoin(s.denom2, sdk.NewInt(50_000_000)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 10)).String()),
			},
			false, &sdk.TxResponse{}, 18,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewModifyBidCmd()
			clientCtx := val.ClientCtx

			out, err := utilcli.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

type QueryCmdTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	denom1 string
	denom2 string
}

func (s *QueryCmdTestSuite) SetupTest() {
	s.T().Log("setting up integration test suite")

	keeper.EnableAddAllowedBidder = true

	encodingCfg := cosmoscmd.MakeEncodingConfig(chain.ModuleBasics)

	cfg := network.DefaultConfig()
	cfg.NumValidators = 2
	cfg.AppConstructor = NewAppConstructor(encodingCfg)
	cfg.GenesisState = chain.ModuleBasics.DefaultGenesis(cfg.Codec)
	cfg.AccountTokens = sdk.NewInt(100_000_000_000_000) // node0token denom
	cfg.StakingTokens = sdk.NewInt(100_000_000_000_000) // stake denom

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)
	s.denom1, s.denom2 = fmt.Sprintf("%stoken", s.network.Validators[0].Moniker), s.cfg.BondDenom

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *QueryCmdTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *QueryCmdTestSuite) TestNewQueryParamsCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			"happy case",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
		},
		{
			"with specific height",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.NewQueryParamsCmd()

			out, err := utilcli.ExecTestCLICmd(clientCtx, cmd, tc.args)

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().NotEqual("internal", err.Error())
			} else {
				s.Require().NoError(err)

				var params types.Params
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &params))
				s.Require().NotEmpty(params.AuctionCreationFee)
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewQueryAuctionsCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	types.RegisterInterfaces(clientCtx.InterfaceRegistry)

	// Create a fixed price auction
	_, err := MsgCreateFixedPriceAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
			StartPrice:      sdk.MustNewDecFromStr("1.0"),
			SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom: s.denom2,
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(1, 0, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now().AddDate(0, 1, 0),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	// Create a batch auction
	_, err = MsgCreateBatchAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
			StartPrice:        sdk.MustNewDecFromStr("0.5"),
			MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
			SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom:   s.denom2,
			MaxExtendedRound:  2,
			ExtendedRoundRate: sdk.MustNewDecFromStr("0.2"),
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(0, 6, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	for _, tc := range []struct {
		name        string
		args        []string
		expectedErr string
		postRun     func(resp types.QueryAuctionsResponse)
	}{
		{
			"happy case",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			"",
			func(resp types.QueryAuctionsResponse) {
				s.Require().Len(resp.Auctions, 2)
			},
		},
	} {
		s.Run(tc.name, func() {
			cmd := cli.NewQueryAuctionsCmd()

			out, err := utilcli.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)

			if tc.expectedErr == "" {
				s.Require().NoError(err)
				var resp types.QueryAuctionsResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewQueryAuctionCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	types.RegisterInterfaces(clientCtx.InterfaceRegistry)

	// Create a fixed price auction
	_, err := MsgCreateFixedPriceAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
			StartPrice:      sdk.MustNewDecFromStr("1.0"),
			SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom: s.denom2,
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(1, 0, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now().AddDate(0, 1, 0),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	for _, tc := range []struct {
		name        string
		args        []string
		expectedErr string
		postRun     func(resp types.QueryAuctionResponse)
	}{
		{
			"happy case",
			[]string{
				strconv.Itoa(1),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			"",
			func(resp types.QueryAuctionResponse) {
				auction, err := types.UnpackAuction(resp.Auction)
				s.Require().NoError(err)
				s.Require().Equal(types.AuctionTypeFixedPrice, auction.GetType())
			},
		},
	} {
		s.Run(tc.name, func() {
			cmd := cli.NewQueryAuctionCmd()

			out, err := utilcli.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)

			if tc.expectedErr == "" {
				s.Require().NoError(err)
				var resp types.QueryAuctionResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewQueryAllowedBiddersCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	types.RegisterInterfaces(clientCtx.InterfaceRegistry)

	// Create a fixed price auction
	_, err := MsgCreateFixedPriceAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
			StartPrice:      sdk.MustNewDecFromStr("1.0"),
			SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom: s.denom2,
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(1, 0, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now().AddDate(0, 1, 0),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	// Add allowed bidder
	_, err = MsgAddAllowedBidderExec(
		val.ClientCtx,
		val.Address.String(),
		1,
		sdk.NewInt(100_000_000),
	)
	s.Require().NoError(err)

	for _, tc := range []struct {
		name        string
		args        []string
		expectedErr string
		postRun     func(resp types.QueryAllowedBiddersResponse)
	}{
		{
			"happy case",
			[]string{
				strconv.Itoa(1),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			"",
			func(resp types.QueryAllowedBiddersResponse) {
				s.Require().Len(resp.AllowedBidders, 1)
			},
		},
	} {
		s.Run(tc.name, func() {
			cmd := cli.NewQueryAllowedBiddersCmd()

			out, err := utilcli.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)

			if tc.expectedErr == "" {
				s.Require().NoError(err)
				var resp types.QueryAllowedBiddersResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewQueryAllowedBidderCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	types.RegisterInterfaces(clientCtx.InterfaceRegistry)

	// Create a fixed price auction
	_, err := MsgCreateFixedPriceAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.FixedPriceAuctionRequest{
			StartPrice:      sdk.MustNewDecFromStr("1.0"),
			SellingCoin:     sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom: s.denom2,
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(1, 0, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now().AddDate(0, 1, 0),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	// Add allowed bidder
	maxBidAmt := sdk.NewInt(100_000_000)
	_, err = MsgAddAllowedBidderExec(
		val.ClientCtx,
		val.Address.String(),
		1,
		maxBidAmt,
	)
	s.Require().NoError(err)

	for _, tc := range []struct {
		name        string
		args        []string
		expectedErr string
		postRun     func(resp types.QueryAllowedBidderResponse)
	}{
		{
			"happy case",
			[]string{
				strconv.Itoa(1),
				val.Address.String(),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			"",
			func(resp types.QueryAllowedBidderResponse) {
				s.Require().Equal(val.Address.String(), resp.AllowedBidder.Bidder)
				s.Require().Equal(maxBidAmt, resp.AllowedBidder.MaxBidAmount)
			},
		},
	} {
		s.Run(tc.name, func() {
			cmd := cli.NewQueryAllowedBidderCmd()

			out, err := utilcli.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)

			if tc.expectedErr == "" {
				s.Require().NoError(err)
				var resp types.QueryAllowedBidderResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *TxCmdTestSuite) TestNewQueryBidsCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	types.RegisterInterfaces(clientCtx.InterfaceRegistry)

	// Create a batch auction
	_, err := MsgCreateBatchAuctionExec(
		val.ClientCtx,
		val.Address.String(),
		testutil.WriteToNewTempFile(s.T(), cli.BatchAuctionRequest{
			StartPrice:        sdk.MustNewDecFromStr("0.5"),
			MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
			SellingCoin:       sdk.NewInt64Coin(s.denom1, 100_000_000_000),
			PayingCoinDenom:   s.denom2,
			MaxExtendedRound:  2,
			ExtendedRoundRate: sdk.MustNewDecFromStr("0.2"),
			VestingSchedules: []types.VestingSchedule{
				{
					ReleaseTime: time.Now().AddDate(0, 6, 0),
					Weight:      sdk.MustNewDecFromStr("1.0"),
				},
			},
			StartTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 3, 0),
		}.String()).Name(),
	)
	s.Require().NoError(err)

	// Add allowed bidder
	_, err = MsgAddAllowedBidderExec(
		val.ClientCtx,
		val.Address.String(),
		1,
		sdk.NewInt(100_000_000),
	)
	s.Require().NoError(err)

	// Place a bid #1
	_, err = MsgPlaceBidExec(
		val.ClientCtx,
		val.Address.String(),
		1,
		"batch-worth",
		sdk.MustNewDecFromStr("0.55"),
		sdk.NewCoin(s.denom2, sdk.NewInt(20_000_000)),
	)
	s.Require().NoError(err)

	// Place a bid #2
	_, err = MsgPlaceBidExec(
		val.ClientCtx,
		val.Address.String(),
		1,
		"batch-many",
		sdk.MustNewDecFromStr("0.6"),
		sdk.NewCoin(s.denom1, sdk.NewInt(5_000_000)),
	)
	s.Require().NoError(err)

	for _, tc := range []struct {
		name        string
		args        []string
		expectedErr string
		postRun     func(resp types.QueryBidsResponse)
	}{
		{
			"happy case",
			[]string{
				strconv.Itoa(1),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			"",
			func(resp types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 2)
			},
		},
	} {
		s.Run(tc.name, func() {
			cmd := cli.NewQueryBidsCmd()

			out, err := utilcli.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)

			if tc.expectedErr == "" {
				s.Require().NoError(err)
				var resp types.QueryBidsResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}
