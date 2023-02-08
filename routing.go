package iface

import (
	"context"

	path "github.com/ipfs/interface-go-ipfs-core/path"
)

// RoutingAPI specifies the interface to the routing layer.
type RoutingAPI interface {
	// Get retrieves the best value for a given key
	Get(context.Context, path.Path) ([]byte, error)

	// Put sets a value for a given key
	Put(context.Context, path.Path, []byte) error
}
