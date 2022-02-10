package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateFixedPriceAuction{},
		&MsgCreateEnglishAuction{},
		&MsgCancelAuction{},
		&MsgPlaceBid{},
		&MsgAddAllowedBidder{},
	)

	registry.RegisterInterface(
		"cosmos.fundraising.v1beta1.AuctionI",
		(*AuctionI)(nil),
		&FixedPriceAuction{},
		&EnglishAuction{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)
