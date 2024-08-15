package fundraising

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/tendermint/fundraising/api/fundraising/fundraising/v1"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	moduloOpts := &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListAllowedBidder",
					Use:       "list-allowed-bidder",
					Short:     "List all AllowedBidder",
				},
				{
					RpcMethod:      "GetAllowedBidder",
					Use:            "get-allowed-bidder [id]",
					Short:          "Gets a AllowedBidder",
					Alias:          []string{"show-allowed-bidder"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auctionId"}},
				},
				{
					RpcMethod: "ListVestingQueue",
					Use:       "list-vesting-queue",
					Short:     "List all VestingQueue",
				},
				{
					RpcMethod: "ListBid",
					Use:       "list-bid",
					Short:     "List all Bid",
				},
				{
					RpcMethod:      "GetBid",
					Use:            "get-bid [id]",
					Short:          "Gets a Bid by id",
					Alias:          []string{"show-bid"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "ListAuction",
					Use:       "list-auction",
					Short:     "List all auction",
				},
				{
					RpcMethod:      "GetAuction",
					Use:            "get-auction [id]",
					Short:          "Gets a auction by id",
					Alias:          []string{"show-auction"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateFixedPriceAuction",
					Use:            "create-fixed-price-auction [start-price] [selling-coin] [paying-coin-denom] [vesting-schedules] [start-time] [end-time]",
					Short:          "Send a CreateFixedPriceAuction tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "startPrice"}, {ProtoField: "sellingCoin"}, {ProtoField: "payingCoinDenom"}, {ProtoField: "vestingSchedules"}, {ProtoField: "startTime"}, {ProtoField: "endTime"}},
				},
				{
					RpcMethod:      "CreateBatchAuction",
					Use:            "create-batch-auction [start-price] [min-bid-price] [selling-coin] [paying-coin-denom] [vesting-schedules] [max-extended-round] [extended-round-rate] [start-time] [end-time]",
					Short:          "Send a CreateBatchAuction tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "startPrice"}, {ProtoField: "minBidPrice"}, {ProtoField: "sellingCoin"}, {ProtoField: "payingCoinDenom"}, {ProtoField: "vestingSchedules"}, {ProtoField: "maxExtendedRound"}, {ProtoField: "extendedRoundRate"}, {ProtoField: "startTime"}, {ProtoField: "endTime"}},
				},
				{
					RpcMethod:      "CancelAuction",
					Use:            "cancel-auction [auction-id]",
					Short:          "Send a CancelAuction tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auctionId"}},
				},
				{
					RpcMethod:      "PlaceBid",
					Use:            "place-bid [auction-id] [bid-type] [price] [coin]",
					Short:          "Send a PlaceBid tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auctionId"}, {ProtoField: "bidType"}, {ProtoField: "price"}, {ProtoField: "coin"}},
				},
				{
					RpcMethod:      "ModifyBid",
					Use:            "modify-bid [auction-id] [bid-id] [price] [coin]",
					Short:          "Send a ModifyBid tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auctionId"}, {ProtoField: "bidId"}, {ProtoField: "price"}, {ProtoField: "coin"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
	if keeper.EnableAddAllowedBidder {
		moduloOpts.Tx.RpcCommandOptions = append(moduloOpts.Tx.RpcCommandOptions, &autocliv1.RpcCommandOptions{
			RpcMethod:      "AddAllowedBidder",
			Use:            "add-allowed-bidder [auction-id] [allowed-bidder]",
			Short:          "Send a AddAllowedBidder tx",
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "auctionId"}, {ProtoField: "allowedBidder"}},
		})
	}

	return moduloOpts
}
