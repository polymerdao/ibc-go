package multihop

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

// Endpoint represents a Cosmos chain endpoint for queries.
// Endpoint is stateless from caller's perspective.
type Endpoint interface {
	ChainID() string
	Codec() codec.BinaryCodec
	ClientID() string
	// Get counterparty chain's height for its key/value proof query.
	GetKeyValueProofHeight() exported.Height
	// Get counterparty chain's consensus state height stored on this chain.
	GetConsensusHeight() exported.Height
	GetConsensusState(height exported.Height) (exported.ConsensusState, error)
	ConnectionID() string
	GetConnection() (*connectiontypes.ConnectionEnd, error)
	// Returns the proof of the `key`` at `height` within the ibc module store.
	QueryProofAtHeight(key []byte, height int64) ([]byte, clienttypes.Height, error)
	GetMerklePath(path string) (commitmenttypes.MerklePath, error)
	// UpdateClient updates the clientState of counterparty chain's header
	UpdateClient() error
	Counterparty() Endpoint
}

func EndpointToString(endpoint Endpoint) string {
	kvHeight := endpoint.GetKeyValueProofHeight()
	consHeight := endpoint.GetConsensusHeight()
	return fmt.Sprintf("{\"chain_id\": %q, \"client_id\": %q, \"connection_id\": %q, \"kv_proof_height\": %q, \"consensus_height\": %q}",
		endpoint.ChainID(), endpoint.ClientID(), endpoint.ConnectionID(), kvHeight, consHeight)
}

// Path contains two endpoints of chains that have a direct IBC connection, ie. a single-hop IBC path.
type Path struct {
	EndpointA Endpoint
	EndpointB Endpoint
}

func (p Path) String() string {
	return fmt.Sprintf("[%s, %s]", EndpointToString(p.EndpointA), EndpointToString(p.EndpointB))
}

// ChanPath represents a multihop channel path that spans 2 or more single-hop `Path`s.
type ChanPath struct {
	Paths        []*Path
	counterparty *ChanPath
}

// NewChanPath creates a new multi-hop ChanPath from a list of single-hop Paths.
func NewChanPath(paths []*Path) ChanPath {
	if len(paths) < 2 {
		panic(fmt.Sprintf("multihop channel path expects at least 2 single-hop paths, but got %d", len(paths)))
	}
	return ChanPath{
		Paths: paths,
	}
}

func (p ChanPath) Counterparty() *ChanPath {
	if p.counterparty != nil {
		return p.counterparty
	}
	p.counterparty = &ChanPath{
		counterparty: &p,
	}
	p.counterparty.Paths = make([]*Path, len(p.Paths))
	for i, hop := range p.Paths {
		reversedSinglePath := Path{EndpointA: hop.EndpointB, EndpointB: hop.EndpointA}
		p.counterparty.Paths[len(p.Paths)-i-1] = &reversedSinglePath
	}
	return p.counterparty
}

// UpdateClient updates the clientState{AB, BC, .. YZ} so chainA's consensusState is propogated to chainZ.
func (p ChanPath) UpdateClient() error {
	for _, path := range p.Paths {
		if err := path.EndpointB.UpdateClient(); err != nil {
			return err
		}
	}
	return nil
}

// GetConnectionHops returns the connection hops for the multihop channel.
func (p ChanPath) GetConnectionHops() []string {
	hops := make([]string, len(p.Paths))
	for i, path := range p.Paths {
		hops[i] = path.EndpointA.ConnectionID()
	}
	return hops
}

// GenerateProof generates a proof for the given key on the the source chain, which is to be verified on the dest
// chain.
func (p ChanPath) GenerateProof(
	key []byte,
	val []byte,
	doVerify bool,
) (result *channeltypes.MsgMultihopProofs, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = sdkerrors.Wrapf(channeltypes.ErrMultihopProofGeneration, "%v", r)
		}
	}()

	result = &channeltypes.MsgMultihopProofs{}

	// generate proof for key on source chain
	chainB := p.Source().Counterparty()
	keyProofHeight := chainB.GetKeyValueProofHeight()
	consStateAB := getTmConsensusState(chainB, chainB.GetConsensusHeight())
	result.KeyProof = queryProof(p.Source(), key, val, keyProofHeight, consStateAB.GetRoot(), doVerify)

	proofGenFuncs := []proofGenFunc{
		genConsensusStateProof,
		genConnProof,
	}
	linkedPathProofs, err := p.GenerateIntermediateStateProofs(proofGenFuncs)
	if err != nil {
		return nil, err
	}
	if len(linkedPathProofs) != len(proofGenFuncs) {
		return nil, sdkerrors.Wrapf(
			channeltypes.ErrMultihopProofGeneration,
			"expected %d linked path proofs for consensus, connections, and client states but got %d",
			len(proofGenFuncs), len(linkedPathProofs),
		)
	}
	result.ConsensusProofs = linkedPathProofs[0]
	result.ConnectionProofs = linkedPathProofs[1]

	return result, nil
}

// The Source chain
func (p ChanPath) Source() Endpoint {
	return p.Paths[0].EndpointA
}

// GenerateIntermediateStateProofs generates lists of connection, consensus, and client state proofs from the source to dest chains.
func (p ChanPath) GenerateIntermediateStateProofs(
	proofGenFuncs []proofGenFunc,
) (result [][]*channeltypes.MultihopProof, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = sdkerrors.Wrapf(channeltypes.ErrMultihopProofGeneration, "%v", r)
		}
	}()
	// initialize a 2-d slice of proofs, where 1st dim is the proof gen funcs, and 2nd dim is the path iter count
	result = make([][]*channeltypes.MultihopProof, len(proofGenFuncs))

	// iterate over all but last single-hop path
	iterCount := len(p.Paths) - 1
	for i := 0; i < iterCount; i++ {
		// Given 3 chains connected by 2 paths:
		// A -(path)-> B -(nextPath)-> C
		// , We need to generate proofs for chain A's key paths. The proof is verified with B's consensus state on C.
		// ie. proof to verify A's state on C.
		// The loop starts with the source chain as A, and ends with the dest chain as chain C.
		// NOTE: chain {A,B,C} are relatively referenced to the current iteration, not to be confused with the chainID
		// or endpointA/B.
		chainB, chainC := p.Paths[i].EndpointB, p.Paths[i+1].EndpointB
		heightAB := chainB.GetConsensusHeight()
		heightBC := chainC.GetConsensusHeight()
		consStateBC := getTmConsensusState(chainC, heightBC)

		for j, proofGenFunc := range proofGenFuncs {
			proof := proofGenFunc(chainB, heightAB, heightBC, consStateBC.GetRoot())
			result[j] = append([]*channeltypes.MultihopProof{proof}, result[j]...)
		}

	}

	return result, nil
}

func (p ChanPath) String() string {
	output := "["
	for i, path := range p.Paths {
		if i > 0 {
			output += ", "
		}
		output += path.String()
	}
	return output + "]"
}

type proofGenFunc func(Endpoint, exported.Height, exported.Height, exported.Root) *channeltypes.MultihopProof

// Generate a proof for A's consensusState stored on B using B's consensusState root stored on C.
func genConsensusStateProof(
	chainB Endpoint,
	heightAB, heightBC exported.Height,
	consStateBCRoot exported.Root,
) *channeltypes.MultihopProof {
	chainAID := chainB.Counterparty().ChainID()
	consStateAB, err := chainB.GetConsensusState(heightAB)
	panicIfErr(err, "chain [%s]'s consensus state on chain '%s' at height %s not found due to: %v",
		chainAID, chainB.ChainID(), heightAB, err,
	)
	bzConsStateAB, err := chainB.Codec().MarshalInterface(consStateAB)
	panicIfErr(err, "fail to marshal consensus state of chain '%s' on chain '%s' at height %s due to: %v",
		chainAID, chainB.ChainID(), heightAB, err,
	)
	return queryProof(
		chainB,
		host.FullConsensusStateKey(chainB.ClientID(), heightAB),
		bzConsStateAB,
		heightBC,
		consStateBCRoot,
		true,
	)
}

// Generate a proof for the connEnd denoting A stored on B using B's consensusState root stored on C.
func genConnProof(
	chainB Endpoint,
	heightAB, heightBC exported.Height,
	consStateBCRoot exported.Root,
) *channeltypes.MultihopProof {
	connAB, err := chainB.GetConnection()
	panicIfErr(err, "fail to get connection '%s' on chain '%s' due to: %v",
		chainB.ConnectionID(), chainB.ChainID(), err,
	)
	bzConnAB, err := chainB.Codec().Marshal(connAB)
	panicIfErr(err, "fail to marshal connection '%s' on chain '%s' due to: %v",
		chainB.ConnectionID(), chainB.ChainID(), err,
	)
	return queryProof(chainB, host.ConnectionKey(chainB.ConnectionID()), bzConnAB, heightBC, consStateBCRoot, true)
}

// queryProof queries the key-value pair or absence proof stored on A and optionally ensures the proof
// can be verified by A's consensus state root stored on B at heightAB. where A--B is connected by a
// single ibc connection.
//
// if doVerify is false, skip verification.
// If value is nil, do non-membership verification.
// If heightAB is nil, use the latest height of B's client state.
// consStateABRoot must be non-empty. It should be the root of the consensus state of clientAB at heightAB.
//
// Panic if proof generation or verification fails.
func queryProof(
	chainA Endpoint,
	key, value []byte,
	heightAB exported.Height,
	consStateABRoot exported.Root,
	doVerify bool,
) *channeltypes.MultihopProof {
	// provide context for error thrown
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("fail to generate proof on chain '%s' for key [%s] at height %s: %v",
				chainA.ChainID(), key, heightAB, r,
			))
		}
	}()

	if len(key) == 0 {
		panic("key must be non-empty")
	}

	chainB := chainA.Counterparty()
	// set optional params if not passed in
	if heightAB == nil {
		heightAB = chainB.GetConsensusHeight()
	}
	if consStateABRoot == nil || consStateABRoot.Empty() {
		panic("consensus state root must be non-empty")
	}

	keyMerklePath, err := chainB.GetMerklePath(string(key))
	panicIfErr(err, "fail to create merkle path [%s]: %v", key, err)

	bzProof, _, err := chainA.QueryProofAtHeight(key, int64(heightAB.GetRevisionHeight()))
	panicIfErr(err, "proof query error: %v", err)

	// only verify ke/value if value is not nil
	if doVerify {
		var proof commitmenttypes.MerkleProof
		err = chainA.Codec().Unmarshal(bzProof, &proof)
		panicIfErr(err, "fail to unmarshal chain [%s]'s proof on chain [%s] due to: %v",
			chainA.ChainID(), chainB.ChainID(), err,
		)
		if len(value) > 0 {
			// ensure key-value pair can be verified by consStateBC
			err = proof.VerifyMembership(
				commitmenttypes.GetSDKSpecs(), consStateABRoot,
				keyMerklePath, value,
			)
			panicIfErr(err, "verify membership failed: %v", err)
		} else {
			err = proof.VerifyNonMembership(
				commitmenttypes.GetSDKSpecs(), consStateABRoot,
				keyMerklePath,
			)
			panicIfErr(err, "verify non-membership failed: %v", err)
		}
	}

	return &channeltypes.MultihopProof{
		Proof:       bzProof,
		Value:       value,
		PrefixedKey: &keyMerklePath,
	}
}

func getTmConsensusState(end Endpoint, height exported.Height) *tmclient.ConsensusState {
	consState, err := end.GetConsensusState(height)
	panicIfErr(err, "fail to get consensus state of chain '%s' at height %s due to: %v",
		end.ChainID(), height, err,
	)
	tmConsState, ok := consState.(*tmclient.ConsensusState)
	if !ok {
		panic(fmt.Sprintf("expected consensus state to be tendermint consensus state, got: %T", consState))
	}
	return tmConsState
}

func panicIfErr(err error, format string, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprintf(format, args...))
	}
}
