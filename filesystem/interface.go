package fs

import (
	"context"
	"io"
	"os"

	coreiface "github.com/ipfs/interface-go-ipfs-core"
	corepath "github.com/ipfs/interface-go-ipfs-core/path"
)

type (
	ParseFn func(ctx context.Context, name string) (FsNode, error)
	FsType  = coreiface.FileType
	// TODO: os.FileInfo-like interface
	//FileMetadata
)

// Provides standard convention wrappers around an index
type FileSystem interface {
	Index // a namespace registry that is used to parse `name` strings into `FsNode`s

	// Create node using respective API for its namespace
	Create(ctx context.Context, name string, nodeType FsType) error

	/* XXX: see PR discussion; we should avoid context.Value if we can; different requests require different arity, which may be a problem
	Request specific parameters are passed via context key's and values, at an API level
	If a type's creation method relies on values, the FS must provide them,
	else the create method shall return `ErrPartialContext`

	e.g. Creating a symlink wrapper at a package level would look as such
	`CreateLink(ctx, name, target ) {
	    fs.Create(context.WithValue(ctx, fs.TargetKey, target), name, unixfs.TSymlink)
	}`
	*/

	// Remove node using respective API for its namespace
	// caller provides values necessary to complete the request (if any), inside the context
	Remove(ctx context.Context, name string) error

	OpenFile(name string, flags flags, perm os.FileMode) (FsFile, error)
	OpenDirectory(path string) (FsDirectory, error)
	OpenReference(string) (FsReference, error)
}

type Index interface {
	// register namespace parser with index, returns release/teardown function
	Register(namespace string, nodeParser ParseFn) (io.Closer, error)
	//NOTE: it may be better to return the equivalent of a "KeepAlive" function rather than a closer
	// if the namespace "host" disappears, we'll never release

	// returns true if namespace is registered with the index
	// i.e. if "/ipfs/" was registered, Provides("/ipfs/") == true
	Provides(namespace string) bool

	/*
		Lookup is implicitly shallow
		if the named string resolves to a reference, a reference object shall be returned

		If the named string exists within its namespace
		a node with its metadata initialized shall be returned

		If named string does not exist within its namespace:
		returns a node interface, with metadata uninitialized, and os.ErrNotExist for error
		(expectation is for node.Create() to be called immediately after)

		In any other case:
		If the named string is not valid for any reason, nil and an error shall be returned
	*/
	Lookup(ctx context.Context, name string) (FsNode, error)
}

type FsNode interface {
	corepath.Path
	InitMetadata(context.Context) (FileMetadata, error)
	Metadata() FileMetadata
	//RWLocker Maybe?

	/* XXX: see PR discussion, this is not ideal; YieldIo should not require type
	YieldIo must ensure that the requested nodeType, and the interface returned, are compatible.
	e.g. `YieldIo(ctx, unixfs.TFile)`, is expected to return an object that satisfies the `FsFile` interface
	or `errIOType` if the implementation
	*/
	YieldIo(ctx context.Context, nodeType FsType) (io interface{}, err error)

	Remove(context.Context) error
	Create(context.Context, FsType) error // should this return something? FsNode, io interface of type FsType?
}

type FsFile interface {
	io.Reader
	io.Seeker
	io.Closer
	Size() (int64, error)
	Write(buff []byte, ofst int64) (int, error)
	Sync() error
	Truncate(size uint) error
	Record() FsNode
}

type FsDirectory interface {
	Entries() int // count
	Read(ctx context.Context, offset int64) <-chan FsNode
	Record() FsNode
}

type FsReference interface {
	Target() string
}
