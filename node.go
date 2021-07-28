package iface

import (
	"github.com/ipfs/go-dagwriter"
	"github.com/ipfs/go-fetcher"
)

// NodeAPI provides an interface for working directly with IPLD nodes
// through go-ipld-prime. It is a combination of DagWritingService (for writes)
// and Fetcher (for reads)
type NodeAPI interface {
	// fetcher.Factory provides the interface to get new dag fetchers
	fetcher.Factory

	// DagWritingService implements methods to write dags
	dagwriter.DagWritingService
}
