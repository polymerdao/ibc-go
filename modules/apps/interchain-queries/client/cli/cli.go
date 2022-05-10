package cli

import (
	"github.com/spf13/cobra"

	controllercli "github.com/cosmos/ibc-go/v3/modules/apps/interchain-queries/controller/client/cli"
	hostcli "github.com/cosmos/ibc-go/v3/modules/apps/interchain-queries/host/client/cli"
)

// GetQueryCmd returns the query commands for the interchain-queries submodule
func GetQueryCmd() *cobra.Command {
	icqQueryCmd := &cobra.Command{
		Use:                        "interchain-queries",
		Aliases:                    []string{"icq"},
		Short:                      "interchain-queries subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	icqQueryCmd.AddCommand(
		controllercli.GetQueryCmd(),
		hostcli.GetQueryCmd(),
	)

	return icqQueryCmd
}
