package multihop

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// Endpoint represents a Cosmos chain endpoint for queries.
// Endpoint is stateless from caller's perspective.
type Endpoint interface {
	ChainID() string
	Codec() codec.BinaryCodec
	ClientID() string
	AppHash() []byte
	LastAppHash() []byte
	GetLatestHeight() exported.Height
	GetCurrentHeight() exported.Height
	GetClientState() exported.ClientState
	GetConsensusState(height exported.Height) (exported.ConsensusState, bool)
	ConnectionID() string
	GetConnection() (*connectiontypes.ConnectionEnd, error)
	// Returns the value of the `key`` at `height` within the ibc module store and optionally the proof
	QueryStateAtHeight(key []byte, height int64, doProof bool) ([]byte, []byte, error)

	// QueryMinimumConsensusHeight returns the minimum height within the provided range at which the consensusState exists (processedHeight)
	// and the corresponding consensus state height (consensusHeight).
	QueryMinimumConsensusHeight(minConsensusHeight exported.Height, maxConsensusHeight exported.Height) (exported.Height, exported.Height, error)
	GetMerklePath(path string) (commitmenttypes.MerklePath, error)
	Counterparty() Endpoint
	UpdateClient() error
}

// Path contains two endpoints of chains that have a direct IBC connection, ie. a single-hop IBC path.
type Path struct {
	EndpointA Endpoint
	EndpointB Endpoint
}

// ChanPath represents a multihop channel path that spans 2 or more single-hop `Path`s.
type ChanPath struct {
	Paths []*Path
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

// Counterparty returns the ChanPath as seen in the reverse direction.
func (p ChanPath) Counterparty() ChanPath {
	paths := make([]*Path, len(p.Paths))
	for i, hop := range p.Paths {
		reversedSinglePath := Path{EndpointA: hop.EndpointB, EndpointB: hop.EndpointA}
		paths[len(p.Paths)-i-1] = &reversedSinglePath
	}
	return NewChanPath(paths)
}

// GetConnectionHops returns the connection hops for the multihop channel.
func (p ChanPath) GetConnectionHops() []string {
	hops := make([]string, len(p.Paths))
	for i, path := range p.Paths {
		hops[i] = path.EndpointA.ConnectionID()
	}
	return hops
}

// Source returns the Source chain for the chan path
func (p ChanPath) Source() Endpoint {
	return p.Paths[0].EndpointA
}

// QueryMultihopProof returns a multi-hop proof for the given key on the the source chain along with the proofHeight, which is to be verified on the dest chain.
// The proofHeight is the consensus state height for the immediate counterparty on the receiving chain. This is the known/trusted consensus state which starts
// the multi-hop proof verification.
func (p ChanPath) QueryMultihopProof(key []byte, minProofHeight exported.Height) (multihopProof channeltypes.MsgMultihopProofs, proofHeight exported.Height, err error) {

	N := len(p.Paths) - 1

	if N < 1 {
		err = fmt.Errorf("multihop proof query requires channel path length >= 2")
		return
	}

	heights, err := p.calcProofPath(0, minProofHeight)
	if err != nil {
		return
	}
	fmt.Printf("proofHeights:\n")
	for i, pf := range heights {
		fmt.Printf("i=%d proofHeight=%v consensusHeight=%v\n", i, pf.proofHeight, pf.consensusHeight)
	}

	// update minProofHeight
	minProofHeight, _ = heights[len(heights)-1].consensusHeight.Decrement()

	// query the proof of the key/value on the source chain at a height provable on the next chain.
	fmt.Printf("PROOF: key proven on %s at height=%v\n", p.Source().ChainID(), minProofHeight)
	if multihopProof.KeyProof, err = queryProof(p.Source(), key, minProofHeight, false); err != nil {
		return
	}

	multihopProof.ConsensusProofs = make([]*channeltypes.MultihopProof, N)
	multihopProof.ConnectionProofs = make([]*channeltypes.MultihopProof, N)

	// query proofs of consensus/connection states on intermediate chains
	if err = p.queryIntermediateProofs(0, heights[1:], multihopProof.ConsensusProofs, multihopProof.ConnectionProofs); err != nil {
		return
	}

	proofHeight = heights[1].proofHeight
	fmt.Printf("[QueryMultihopProof] returning with proofHeight=%v\n", proofHeight)

	return
}

type ProofHeights struct {
	proofHeight     exported.Height
	consensusHeight exported.Height
}

// calcProofPath
func (p ChanPath) calcProofPath(pathIdx int, minProofHeight exported.Height) ([]*ProofHeights, error) {
	chain := p.Paths[pathIdx].EndpointB

	// find minimum consensus height provable on the next chain
	nextProofHeight, consensusHeight, _ := chain.QueryMinimumConsensusHeight(minProofHeight, nil)
	fmt.Printf("chain=%s minProofHeight=%v nextProofHeight=%v consensusState=%v\n", chain.ChainID(), minProofHeight, nextProofHeight, consensusHeight)

	if nextProofHeight == nil {
		if err := chain.UpdateClient(); err != nil {
			return nil, err
		}

		nextProofHeight = chain.GetLatestHeight()
		consensusHeight = chain.GetClientState().GetLatestHeight()
		fmt.Printf("Update: chain=%s nextProofHeight=%v consensusState=%v\n", chain.ChainID(), nextProofHeight, consensusHeight)
	}

	N := len(p.Paths) - 1

	if pathIdx == N {
		height := ProofHeights{
			proofHeight:     nextProofHeight,
			consensusHeight: consensusHeight,
		}
		return []*ProofHeights{&height}, nil
	}

	heights, err := p.calcProofPath(pathIdx+1, nextProofHeight)
	if err != nil {
		return nil, err
	}
	height := ProofHeights{
		proofHeight:     nextProofHeight,
		consensusHeight: consensusHeight,
	}
	heights = append(heights, &height)
	return heights, nil
}

// queryIntermediateProofs recursively queries intermediate chains in a multi-hop channel path for consensus state
// and connection proofs. It stops at the second to last path since the consensus and connection state on the
// final hop is already known on the destination.
func (p ChanPath) queryIntermediateProofs(
	pathIdx int,
	proofHeights []*ProofHeights,
	consensusProofs []*channeltypes.MultihopProof,
	connectionProofs []*channeltypes.MultihopProof,
) (err error) {

	chain := p.Paths[pathIdx].EndpointB
	N := len(p.Paths) - 2

	pf := proofHeights[N-pathIdx] // TODO: collect in reverse order?

	fmt.Printf("PROOF: consensusState/%v proven on %s at height=%v\n", pf.consensusHeight, chain.ChainID(), pf.proofHeight)
	if consensusProofs[N-pathIdx], err = queryConsensusStateProof(chain, pf.proofHeight, pf.consensusHeight); err != nil {
		return
	}

	if connectionProofs[N-pathIdx], err = queryConnectionProof(chain, pf.proofHeight); err != nil {
		return
	}

	// no need to query min consensus height on final chain
	if pathIdx == N {
		return
	}

	return p.queryIntermediateProofs(pathIdx+1, proofHeights, consensusProofs, connectionProofs)
}

// queryConsensusStateProof queries a chain for a proof at `proofHeight` for a consensus state at `consensusHeight`
func queryConsensusStateProof(
	chain Endpoint,
	proofHeight exported.Height,
	consensusHeight exported.Height,
) (*channeltypes.MultihopProof, error) {
	key := host.FullConsensusStateKey(chain.ClientID(), consensusHeight)
	return queryProof(chain, key, proofHeight, true)
}

// queryConnectionProof queries a chain for a proof at `proofHeight` for a connection
func queryConnectionProof(
	chain Endpoint,
	proofHeight exported.Height,
) (*channeltypes.MultihopProof, error) {
	key := host.ConnectionKey(chain.ConnectionID())
	return queryProof(chain, key, proofHeight, true)
}

// queryProof queries a (non-)membership proof for the key on the specified chain and
// returns the proof
//
// if doValue, the queried value is added to the proof this is required for
// intermediate consensus/connection multihop proofs
func queryProof(
	chain Endpoint,
	key []byte,
	height exported.Height,
	doValue bool,
) (*channeltypes.MultihopProof, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("key must be non-empty")
	}

	if height == nil {
		return nil, fmt.Errorf("height must be non-nil")
	}

	keyMerklePath, err := chain.GetMerklePath(string(key))
	if err != nil {
		return nil, fmt.Errorf("fail to create merkle path on chain '%s' with path '%s' due to: %v",
			chain.ChainID(), key, err)
	}

	valueBytes, proof, err := chain.QueryStateAtHeight(key, int64(height.GetRevisionHeight()), true)
	if err != nil {
		return nil, fmt.Errorf("fail to generate proof on chain '%s' for key '%s' at height %d due to: %v",
			chain.ChainID(), key, height, err,
		)
	}

	if !doValue {
		valueBytes = nil
	}

	return &channeltypes.MultihopProof{
		Proof:       proof,
		Value:       valueBytes,
		PrefixedKey: &keyMerklePath,
	}, nil
}
