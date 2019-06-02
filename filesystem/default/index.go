package ipfs

import (
	"context"
	"encoding/binary"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-unixfs"
	fs "github.com/ipfs/interface-go-ipfs-core/filesystem/interface"
	"github.com/multiformats/go-multihash"
)

func Lookup(ctx context.Context, name string) (fs.FsNode, error) {
	return pkgRoot.Lookup(ctx, name)
}

// api's may define Metadata.Cid() however they like
// for the default, we use this QID-like generator
func GenQueryID(path string, md fs.Metadata) (cid.Cid, error) {
	mdBuf := make([]byte, 16)
	binary.LittleEndian.PutUint64(mdBuf, uint64(md.Size()))
	binary.LittleEndian.PutUint64(mdBuf[8:], uint64(md.Type()))

	prefix := cid.V1Builder{Codec: cid.DagCBOR, MhType: multihash.BLAKE2B_MIN}
	return prefix.Sum(append([]byte(path), mdBuf...))
}

type PathParserRegistry struct {
	sync.Mutex
	nodeParsers map[string]fs.ParseFn
}

func (rr *PathParserRegistry) Mount(subrootPath string, nodeParser fs.ParseFn) (io.Closer, error) {
	rr.Lock()
	val, ok := rr.nodeParsers[subrootPath]
	if ok || val != nil {
		return nil, errRegistered
	}

	rr.nodeParsers[subrootPath] = nodeParsers
	return func() {
		ri.Lock()
		//TODO: somehow trigger cascading close for open subsystem handles
		// or note that open handles are still valid, but new handles will not be made
		delete(ri.subSystem, subrootPath)
		ri.Unlock()
	}, nil
}

func (PathParserRegistry) Lookup(ctx context.Context, name string) (FsNode, error) {
	//NOTE: we can use a pkg level cache here, and fallback to the parser only when necessary

	/* very simple map lookup
	   path is compared against registered subsystems
	   the function associated with the most specific matching prefix wins
	   (see registration in init() for additional context)
	*/
	var (
		subLookup        ParseFn
		highestPrecision int
	)
	for subSystem, parser := range ri.nodeParsers {
		if strings.HasPrefix(name, subSystem) {
			if precision := len(subSystem); precision > highestPrecision {
				highestPrecision = precision
				subLookup = parser
			}
		}

		if subLookup == nil {
			return nil, errNoHandler
		}

		//TODO: [important] should we pass in relative (to namespace) paths here instead of absolute?
		return subLookup(ctx, name)
	}
}

// `pkg/log`-like wrappers for pkg level FS
// simply so we can say `pkg.lock` instead of `pkg.default.lock`
// this is more justified in client.go rather than here
func Lock() {
	pkgRoot.Lock()
}

func Unlock() {
	pkgRoot.Unlock()
}

func OpenFile(name string, flags OFlags, perm os.FileMode) (FsFile, error) {
	fs.Lock()
	defer fs.Unlock()

	fsNode, err := fs.Lookup(ctx, name)
	if err != nil {
		return nil, err
	}

	/*
		insert logic pertinent to your own file system implementation here
		i.e. flag parsing, permission checking, link-resolution, etc.
		e.g. If you are implementing a POSIX compliant FS,
		you may want to return an error if write permissions where requested, for a name which resides in a read-only API

		Descriptor management (if any) is also client defined
		For our Go example, we simple use the interface directly, and Close() it when we're done
		If we were implementing POSIX, we could wrap that interface inside of PosixOpen(),
		assign it an integer value, hold on to the interface, and Close it later when PosixClose(ourInt) is called
		Something like this
		`wrappedIo, err := posixRoot.YieldOrGetHandle(fsNode, perm, unixfs.TFile)
		if err != nil {
		    //...
		}
		return wrappedIo, nil`
	*/

	// YieldIo is expected to handle type/capability checking and return conforming errors
	// XXX: see PR discussion; YieldIo should not require type
	return fsNode.YieldIo(ctx, unixfs.TFile)
}
