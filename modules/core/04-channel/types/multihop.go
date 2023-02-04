package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	connectiontypes "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
)

type ConnectionEnd = connectiontypes.ConnectionEnd

// GetMultihopConnectionEnd returns the final connectionEnd from the counterparty perspective
func (m *MsgMultihopProofs) GetMultihopConnectionEnd(cdc codec.BinaryCodec) (*ConnectionEnd, error) {
	var connectionEnd ConnectionEnd
	if err := cdc.Unmarshal(m.ConnectionProofs[0].Value, &connectionEnd); err != nil {
		return nil, err
	}
	return &connectionEnd, nil
}

// GetMultihopCounterpartyConsensus returns the final consensusState from the counterparty perspective (e.g. the source chain state).
func (m *MsgMultihopProofs) GetMultihopCounterpartyConsensus(cdc codec.BinaryCodec) (consensusState exported.ConsensusState, err error) {
	err = cdc.UnmarshalInterface(m.ConsensusProofs[0].Value, &consensusState)
	return
}

// GetCounterpartyHops returns the counter party connectionHops
func (m *MsgMultihopProofs) GetCounterpartyHops(
	cdc codec.BinaryCodec,
	lastConnection *ConnectionEnd,
) ([]string, error) {
	var counterpartyHops []string
	for _, connData := range m.ConnectionProofs {
		var connectionEnd ConnectionEnd
		if err := cdc.Unmarshal(connData.Value, &connectionEnd); err != nil {
			return nil, err
		}
		counterpartyHops = append(counterpartyHops, connectionEnd.GetCounterparty().GetConnectionID())
	}

	counterpartyHops = append(counterpartyHops, lastConnection.GetCounterparty().GetConnectionID())

	return counterpartyHops, nil
}
