package iface

import (
	"context"

	"github.com/ipfs/go-dagwriter"
	"github.com/ipfs/go-fetcher"
)

// NodeAPI provides an interface for working directly with IPLD nodes
// through go-ipld-prime. It is a combination of DagWritingService (for writes)
// and Fetcher (for reads)
type NodeAPI interface {
	// NewSession returns an instance of the Fetcher
	NewSession(ctx context.Context) fetcher.Fetcher

	// DagWritingService implements methods to write dags
	dagwriter.DagWritingService
}
