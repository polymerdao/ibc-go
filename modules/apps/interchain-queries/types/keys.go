package types

const (
	// ModuleName defines the interchain accounts module name
	ModuleName = "interchainquery"

	// PortID is the default port id that the interchain query submodules binds to
	PortID = "icq"

	// Version defines the current version for interchain accounts
	Version = "icq-1"

	// StoreKey is the store key string for interchain accounts
	StoreKey = ModuleName

	// RouterKey is the message route for interchain accounts
	RouterKey = ModuleName

	// QuerierRoute is the querier route for interchain accounts
	QuerierRoute = ModuleName
)
