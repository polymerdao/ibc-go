package types

import (
	types "github.com/cosmos/ibc-go/v3/modules/apps/interchain-queries/types_host"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
)

// DefaultGenesis creates and returns the interchain accounts GenesisState
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		QueriesGenesisState: DefaultQueriesGenesis(),
	}
}

// NewGenesisState creates and returns a new GenesisState instance from the provided queries genesis state types
func NewGenesisState(queriesGenesisState QueriesGenesisState) *GenesisState {
	return &GenesisState{
		QueriesGenesisState: queriesGenesisState,
	}
}

// Validate performs basic validation of the interchain accounts GenesisState
func (gs GenesisState) Validate() error {
	if err := gs.QueriesGenesisState.Validate(); err != nil {
		return err
	}
	return nil
}

// DefaultControllerGenesis creates and returns the default interchain queires GenesisState
func DefaultQueriesGenesis() QueriesGenesisState {
	return QueriesGenesisState{
		Params: types.DefaultParams(),
	}
}

// NewControllerGenesisState creates a returns a new QueriesGenesisState instance
func NewQueriesGenesisState(channels []ActiveChannel, ports []string, queriesParams types.Params) QueriesGenesisState {
	return QueriesGenesisState{
		ActiveChannels: channels,
		Ports:          ports,
		Params:         queriesParams,
	}
}

// Validate performs basic validation of the QueriesGenesisState
func (gs QueriesGenesisState) Validate() error {
	for _, ch := range gs.ActiveChannels {
		if err := host.ChannelIdentifierValidator(ch.ChannelId); err != nil {
			return err
		}

		if err := host.PortIdentifierValidator(ch.PortId); err != nil {
			return err
		}
	}

	for _, port := range gs.Ports {
		if err := host.PortIdentifierValidator(port); err != nil {
			return err
		}
	}

	if err := gs.Params.Validate(); err != nil {
		return err
	}

	return nil
}
