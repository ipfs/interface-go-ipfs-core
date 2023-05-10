package iface

import "errors"

var (
	// Deprecated: use github.com/ipfs/boxo/coreiface.ErrIsDir
	ErrIsDir = errors.New("this dag node is a directory")
	// Deprecated: use github.com/ipfs/boxo/coreiface.ErrNotFile
	ErrNotFile = errors.New("this dag node is not a regular file")
	// Deprecated: use github.com/ipfs/boxo/coreiface.ErrOffline
	ErrOffline = errors.New("this action must be run in online mode, try running 'ipfs daemon' first")
	// Deprecated: use github.com/ipfs/boxo/coreiface.ErrNotSupported
	ErrNotSupported = errors.New("operation not supported")
)
