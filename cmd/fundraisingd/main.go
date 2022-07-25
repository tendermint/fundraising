package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/tendermint/fundraising/app"
	fundraisingtypes "github.com/tendermint/fundraising/pkg/types"
)

func main() {
	rootCmd, _ := fundraisingtypes.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.DefaultChainID,
		app.ModuleBasics,
		app.New,
	)
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
