package cmd

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectype "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

type (
	// CosmosApp implements the common methods for a Cosmos SDK-based application
	// specific blockchain.
	CosmosApp interface {
		// Name is the assigned name of the app.
		Name() string

		// The application types codec.
		// NOTE: This should be sealed before being returned.
		LegacyAmino() *codec.LegacyAmino

		// BeginBlocker updates every begin block.
		BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock

		// EndBlocker updates every end block.
		EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock

		// InitChainer updates at chain (i.e app) initialization.
		InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain

		// LoadHeight loads the app at a given height.
		LoadHeight(height int64) error

		// ExportAppStateAndValidators exports the state of the application for a genesis file.
		ExportAppStateAndValidators(
			forZeroHeight bool,
			jailAllowedAddrs []string,
			modulesToExport []string,
		) (servertypes.ExportedApp, error)

		// ModuleAccountAddrs are registered module account addreses.
		ModuleAccountAddrs() map[string]bool
	}

	// EncodingConfig specifies the concrete encoding types to use for a given app.
	// This is provided for compatibility between protobuf and amino implementations.
	EncodingConfig struct {
		InterfaceRegistry codectype.InterfaceRegistry
		Marshaler         codec.Codec
		TxConfig          client.TxConfig
		Amino             *codec.LegacyAmino
	}
)

// makeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func makeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectype.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig(moduleBasics module.BasicManager) EncodingConfig {
	encodingConfig := makeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	moduleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	moduleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
