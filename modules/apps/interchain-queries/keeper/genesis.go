package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icqtypes "github.com/cosmos/ibc-go/v3/modules/apps/interchain-queries/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
)

// InitGenesis controller initializes the interchain accounts controller application state from a provided genesis state
func InitGenesis(ctx sdk.Context, keeper Keeper, state icqtypes.GenesisState) {
	for _, portID := range state.Ports {
		if !keeper.IsBound(ctx, portID) {
			cap := keeper.BindPort(ctx, portID)
			if err := keeper.ClaimCapability(ctx, cap, host.PortPath(portID)); err != nil {
				panic(fmt.Sprintf("could not claim port capability: %v", err))
			}
		}
	}

	for _, ch := range state.ActiveChannels {
		keeper.SetActiveChannelID(ctx, ch.ConnectionId, ch.PortId, ch.ChannelId)
	}

	for _, acc := range state.InterchainAccounts {
		keeper.SetInterchainAccountAddress(ctx, acc.ConnectionId, acc.PortId, acc.AccountAddress)
	}

	keeper.SetParams(ctx, state.Params)
}

// InitGenesis initializes the interchain accounts host application state from a provided genesis state
func InitGenesis(ctx sdk.Context, keeper Keeper, state icqtypes.GenesisState) {
	if !keeper.IsBound(ctx, state.Port) {
		cap := keeper.BindPort(ctx, state.Port)
		if err := keeper.ClaimCapability(ctx, cap, host.PortPath(state.Port)); err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}

	for _, ch := range state.ActiveChannels {
		keeper.SetActiveChannelID(ctx, ch.ConnectionId, ch.PortId, ch.ChannelId)
	}

	for _, acc := range state.InterchainAccounts {
		keeper.SetInterchainAccountAddress(ctx, acc.ConnectionId, acc.PortId, acc.AccountAddress)
	}

	keeper.SetParams(ctx, state.Params)
}

// ExportGenesis returns the interchain accounts host exported genesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) icqtypes.GenesisState {
	return icqtypes.GenesisState(
		keeper.GetAllActiveChannels(ctx),
		keeper.GetAllPorts(ctx),
		keeper.GetParams(ctx),
	)
}
