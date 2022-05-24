package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for the interchain-queries module
func GetQueryCmd() *cobra.Command {
	icqQueryCmd := &cobra.Command{
		Use:                        "interchain-queires",
		Aliases:                    []string{"icq"},
		Short:                      "interchain-queires query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	icqQueryCmd.AddCommand(
		GetCmdParams(),
	)

	return icqQueryCmd
}

// NewTxCmd returns the transaction commands for the interchain-queries moduel
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "interchain-queires",
		Short:                      "interchain-queires transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewQueriesTxCmd(),
	)

	return txCmd
}
