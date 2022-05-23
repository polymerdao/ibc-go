package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/ibc-go/v3/modules/apps/interchain-queries/types"
	icqtypes "github.com/cosmos/ibc-go/v3/modules/apps/interchain-queries/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
)

// GetCmdParams returns the command handler for the submodule parameter querying.
func GetCmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current interchain-queries submodule parameters",
		Long:    "Query the current interchain-queries submodule parameters",
		Args:    cobra.NoArgs,
		Example: fmt.Sprintf("%s query interchain-queries params", version.AppName),
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queriesClient := types.NewQueriesClient(clientCtx)

			res, err := queriesClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdPacketEvents returns the command handler for the packet events querying.
func GetCmdPacketEvents() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "packet-events [channel-id] [sequence]",
		Short:   "Query the interchain-queries submodule packet events",
		Long:    "Query the interchain-queries submodule packet events for a particular channel and sequence",
		Args:    cobra.ExactArgs(2),
		Example: fmt.Sprintf("%s query interchain-queries packet-events channel-0 100", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			channelID, portID := args[0], icqtypes.PortID
			if err := host.ChannelIdentifierValidator(channelID); err != nil {
				return err
			}

			seq, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			searchEvents := []string{
				fmt.Sprintf("%s.%s='%s'", channeltypes.EventTypeRecvPacket, channeltypes.AttributeKeyDstChannel, channelID),
				fmt.Sprintf("%s.%s='%s'", channeltypes.EventTypeRecvPacket, channeltypes.AttributeKeyDstPort, portID),
				fmt.Sprintf("%s.%s='%d'", channeltypes.EventTypeRecvPacket, channeltypes.AttributeKeySequence, seq),
			}

			result, err := tx.QueryTxsByEvents(clientCtx, searchEvents, 1, 1, "")
			if err != nil {
				return err
			}

			var resEvents []abci.Event
			for _, r := range result.Txs {
				resEvents = append(resEvents, r.Events...)
			}

			return clientCtx.PrintString(sdk.StringifyEvents(resEvents).String())
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
