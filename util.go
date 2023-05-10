package iface

import (
	"context"
	"io"
)

// Deprecated: use github.com/ipfs/boxo/coreiface.Reader
type Reader interface {
	ReadSeekCloser
	Size() uint64
	CtxReadFull(context.Context, []byte) (int, error)
}

// A ReadSeekCloser implements interfaces to read, copy, seek and close.
//
// Deprecated: use github.com/ipfs/boxo/coreiface.ReadSeekCloser
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
	io.WriterTo
}
