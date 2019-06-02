package ipfs

import (
	"context"
	"time"
	coreiface "github.com/ipfs/interface-go-ipfs-core"


 "github.com/ipfs/interface-go-ipfs-core/filesystem/node/soft"
"github.com/ipfs/go-ipfs-http-client"
)

type defaultFs struct {
	// index
	PathParserRegistry

//defaultFS specific
ctx context.Context
		epoch time.Time

//IPFS parser specifics
//ipfsNode core.IpfsNode
core coreiface.CoreAPI
}
/* From DRAFT
func NewDefaultFileSystem(parentCtx context.Context) (FileSystem, error) {
	// something like this
	fsCtx := deriveFrom(parentCtx)
	// go onCancel(fsCtx) { callClosers() } ()
	return &PathParserRegistry{fsCtx}, nil
}
*/

func newDefaultFs() (FileSystem, error) {
	root :=  defaultFs {
		ctx:context.TODO(), // cancelable something or other
	epoch:time.Now(),
	}

	// we depend on data from the coreapi to initalize our API nodes
	// so fetch it or something and store it on the FS
	//daemon := fallbackApi()
	//ctx := deriveCtx(daemon.ctx) // if the daemon cancels, so should we

	//root.ipfsNode = daemon

	// mount base subsystems
	epoch := time.Now()

// TODO: connect to daemon or fallback [new constructor]
core, err := httpapi.NewLocalApi()
if err != nil {
	return nil, err
}
	for _, pair := range [...]struct {
		string
		ParseFn
	}{
		{"/", fsnode.RootParser(epoch)},
		{"/ipfs", inode.PinParser(core.Pin(), epoch)},  // requests for "/ipfs" are directed at pinRootParser(ctx, requestString)
		{"/ipfs/", coreAPIParser}, // all requests beneath "/ipfs/" are directed at coreAPIParser(ctx, requestString)
		{"/ipns", keyRootParser},
		{"/ipns/", nameAPIParser},
		{filesRootPrefix, filesAPIParser},
	} {
		closer, err := root.Register(pair.string, pair.ParseFn)
		if err != nil {
			if err == errRegistered {
				//TODO: [discuss] consider having Plan9 style unions; Mount() would require flags (union contents go to: front, back, replace)
				// doing this complicates our io.Closer consturction, but may be worth having
				return nil, fmt.Errorf("%q is already registered", pair.string)
			}
			return nil, fmt.Errorf("cannot register %q: %s", pair.string, err)
		}

		// store these somewhere important and call them before you exit
		// this will release the namespace at the FS level
		root.closers = append(root.closers, closer) // in our example we do nothing with them :^)
	}
}

type DescriptorTable interface {
}

func newDescriptorTable() {
}

