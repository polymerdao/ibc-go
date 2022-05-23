package cli

import (
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for the interchain-queries submodule
func GetQueryCmd() *cobra.Command {
	QueryCmd := &cobra.Command{
		Use:                        "interchain-queries",
		Aliases:                    []string{"icq"},
		Short:                      "interchain-queries subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	QueryCmd.AddCommand(
		GetCmdParams(),
		GetCmdPacketEvents(),
	)

	return QueryCmd
}
