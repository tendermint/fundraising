package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.DefaultChainID,
		app.ModuleBasics,
		app.New,
	)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
