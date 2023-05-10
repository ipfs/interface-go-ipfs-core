package iface

import (
	ipld "github.com/ipfs/go-ipld-format"
)

// APIDagService extends ipld.DAGService
//
// Deprecated: use github.com/ipfs/boxo/coreiface.APIDagService
type APIDagService interface {
	ipld.DAGService

	// Pinning returns special NodeAdder which recursively pins added nodes
	Pinning() ipld.NodeAdder
}
