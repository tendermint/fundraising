package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		NewCreateFixedPriceAuction(),
		NewCreateEnglishAuction(),
		NewCancelAuction(),
		NewPlaceBid(),
	)
	return cmd
}

func NewCreateFixedPriceAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-fixed-price-auction [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a fixed price auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a fixed price auction.
The auction details must be provided through a JSON file. 
		
Example:
$ %s tx %s create-fixed-price-auction <path/to/auction.json> --from mykey 

Where auction.json contains:

{}
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Println(clientCtx)

			// TODO: not implemented yet

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCreateEnglishAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-english-auction [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a english auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a english auction.
The auction details must be provided through a JSON file. 
		
Example:
$ %s tx %s create-english-auction <path/to/auction.json> --from mykey 

Where auction.json contains:

{}
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Println(clientCtx)

			// TODO: not implemented yet

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel [auction-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Cancel the auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel the auction with the id. 
		
Example:
$ %s tx %s cancel 1 --from mykey 
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Println(clientCtx)

			// TODO: not implemented yet

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewPlaceBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bid [auction-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Bid for the auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Bid for the auction with the price and coin. 
		
Example:
$ %s tx %s cancel 1 --from mykey 
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Println(clientCtx)

			// TODO: not implemented yet (consider auction type)

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
