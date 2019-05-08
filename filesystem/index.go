package fs

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ipfs/go-unixfs"
)

// This file represents a partial and theoretical
// pkg-level implementation of a filesystem{}

var (
	pkgRoot FileSystem
	closers []io.Closer
)

func init() {

	// we depend on data from the coreapi to initalize our API nodes
	// so fetch it or something and store it on the FS
	daemon := fallbackApi()
	ctx := deriveCtx(daemon.ctx) // if the daemon cancels, so should we

	pkgRoot = NewFileSystem(ctx)
	for _, pair := range [...]struct {
		string
		ParseFn
	}{
		{"/", rootParser},
		{"/ipfs", pinRootParser},  // requests for "/ipfs" are directed at pinRootParser(ctx, requestString)
		{"/ipfs/", coreAPIParser}, // all requests beneath "/ipfs/" are directed at coreAPIParser(ctx, requestString)
		{"/ipns", keyRootParser},
		{"/ipns/", nameAPIParser},
		{filesRootPrefix, filesAPIParser},
	} {
		closer, err := pkgRoot.Register(pair.string, pair.ParseFn)
		if err != nil {
			if err == errRegistered {
				panic(fmt.Sprtinf("%q is already registered", pair.string))
			}
			panic(fmt.Sprtinf("cannot register %q: %s", pair.string, err))
		}

		// store these somewhere important and call them before you exit
		// this will release the namespace at the FS level
		closers = append(closers, closer) // in our example we do nothing with them :^)
	}
}

func NewDefaultFileSystem(parentCtx context.Context) (FileSystem, error) {
	// something like this
	fsCtx := deriveFrom(parentCtx)
	// go onCancel(fsCtx) { callClosers() } ()
	return &PathParserRegistry{fsCtx}, nil
}

type PathParserRegistry struct {
	sync.Mutex
	ctx         context.Context
	nodeParsers map[string]ParseFn
}

func (rr *PathParserRegistry) Register(subrootPath string, nodeParser ParseFn) (io.Closer, error) {
	rr.Lock()
	val, ok := rr.nodeParsers[subrootPath]
	if ok || val != nil {
		return nil, errRegistered
	}

	rr.nodeParsers[subrootPath] = nodeParsers
	return func() {
		ri.Lock()
		delete(ri.subSystem, subrootPath)
		ri.Unlock()
	}, nil
}

func (PathParserRegistry) Lookup(ctx context.Context, name string) (FsPath, error) {
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

func Lookup(ctx context.Context, name string) (FsPath, error) {
	return pkgRoot.Lookup(ctx, name)
}

// ...

func OpenFile(name string, flags flags, perm os.FileMode) (FsFile, error) {
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
