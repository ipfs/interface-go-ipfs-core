package fs

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/ipfs/go-unixfs"
	corepath "github.com/ipfs/interface-go-ipfs-core/path"
)

type (
	ParseFn func(ctx context.Context, name string) (Node, error)
	Kind    int
	Flag    string
	Flags   []struct {
		Flag
		bool
	}
)

const (
	FRead  Flag = "File Read"
	FWrite      = "File Write"
	FSync       = "File Sync"
)

const (
	UfsFile      = Kind(unixfs.TFile)
	UfsDirectory = Kind(unixfs.TDirectory)
	UfsHAMT      = Kind(unixfs.THAMTShard)
	UfsSymlink   = Kind(unixfs.TSymlink)
)

type FsError interface {
	error
}

var (
	ErrNoLink         = FsError(errors.New("not a symlink"))
	ErrInvalidHandle  = FsError(errors.New("invalid handle"))
	ErrNoKey          = FsError(errors.New("key not found"))
	ErrInvalidPath    = FsError(errors.New("invalid path"))
	ErrInvalidArg     = FsError(errors.New("invalid argument"))
	ErrReadOnly       = FsError(errors.New("read only section"))
	ErrIOType         = FsError(errors.New("node does not impliment requested interface"))
	ErrUnexpected     = FsError(errors.New("unexpected node type"))
	ErrNotInitialized = FsError(errors.New("node metadata is not initialized"))
	ErrRoot           = FsError(errors.New("root initialization exception"))
	ErrRecurse        = FsError(errors.New("hit recursion limit"))
	//
	ErrNoHandler      = FsError(errors.New("no handler registered for this request"))
	ErrNotImplemented = FsError(errors.New("operation not implemented"))
)

// Provides standard convention wrappers around an index
type FileSystem interface {
	Index // a namespace registry that is used to parse `name` strings into `Node`s

	// Create node using respective API for its namespace
	Create(ctx context.Context, name string, nodeType Kind) error

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

	OpenFile(name string, flags Flags, perm os.FileMode) (FsFile, error)
	OpenDirectory(path string) (FsDirectory, error)
	OpenReference(string) (FsReference, error)
}

type Index interface {
	sync.Locker
	// mounts the parser function on the filesystem at "namespace"
	// returns "unmount" closer

	// rename consider: Bind(Plan9, Spring)
	Mount(namespace string, nodeParser ParseFn) (unmount io.Closer, err error)
	// NOTE: it may be better to return the equivalent of a "KeepAlive" function rather than a closer
	// if the caller disappears, we'll never release
	//Mount(namespace string, nodeParser ParseFn) (keepmounted KeepaliveFn), err error)

	// returns a list of mounted namespaces
	Mounts() []string

	//TODO: consider
	// Bind() P9; "duplicates some piece of existing name space at another point in the name space"

	/*
		Lookup is implicitly shallow
		if the named string resolves to a reference, a reference object shall be returned

		If the named string exists within its namespace
		a node with its metadata initialized shall be returned

		If named string does not exist within its namespace:
		return a node interface, with metadata uninitialized, and os.ErrNotExist for error
		(expectation is for node.Create() to be called immediately after)

		In any other case:
		If the named string is not valid for any reason, nil and an error shall be returned
	*/
	Lookup(ctx context.Context, name string) (Node, error)
}

type Node interface {
	corepath.Path
	Version() uint
	// Cid shall return a cid for the Node itself
	// the purpose is to allow for cache coherency between API boundaries
	// simillar to Plan9's `qid` and `version` fields
	// if CID is the same between requests, assume cache safe
	//Cid() cid.Cid

	Metadata(context.Context) (Metadata, error)
	//RWLocker Maybe?

	// YieldIo yields an Fs IO type (FsFile, FsDirectory, ...)
	// that cooresponds to the nodes Metadata.Type()
	// If the implimentation does not support IO for this type,
	// then `ErrIOType` shall be returned
	YieldIo(ctx context.Context) (io interface{}, err error)

	Remove(context.Context) error
	Create(context.Context, Kind) error // should this return an Node?
}

/* TODO: Move to fuse or delete
type TimeSpec interface {
	// Update if...
	// any operation...
	Accessed() *time.Time
	// Data itself was...
	Modified() *time.Time
	// Metadata itself...
	Changed() *time.Time
}

type IdSpec struct {
	Owner, // Unix: `uid`; could be libp2p.Id ?
	Group, // `gid`
	Special uint // `rdev`
}
*/

type FsMode uint

type Metadata interface {
	Size() uint
	// Spec history: TODO: move to notes or somewhere else
	// UFS (Unix actual) implimentations:
	// v4 (1973) C-int (2 byte minimum) (PDP ASM -> C)
	// v7 (1979) C-signed-long (4 byte minimum) (
	// POSIX (1988) "shall be signed integer types" RATIONALE: None.
	// Plan9 (1992-current) 9c-vlong (C-long-long; 8byte minimum)
	// IPFS UFSv1 Go-uint64
	// Go-uint (4 byte minimum)
	// There's no reason for this to be signed outside of C and POSIX as far as I can tell

	//TODO: rename? IPLD uses "kind" for these names; Kind() FsKind
	Type() Kind
	Flags() Flags

	//TODO: consider renaming Cid() -> Version(); Type() -> Format()
	//lines up better with both Plan9 and IPLD

	//Cid() cid.Cid // can return nil
}

type FsFile interface {
	io.Reader
	io.Seeker
	io.Closer
	Size() (int64, error)
	Write(buff []byte, ofst int64) (int, error)
	Sync() error
	Truncate(size uint) error
	Record() Node
}

type FsDirectory interface {
	Entries() int // count
	List(ctx context.Context, offset int64) <-chan string
}

type FsReference interface {
	Target() string
}

/* Theoretical -
type FsPipe interface {
}
creation side: store a real file containing your peerid
client side: FS reads file, interprets peerid, establishes p2p<->sockets that act like a named pipe?
*/
