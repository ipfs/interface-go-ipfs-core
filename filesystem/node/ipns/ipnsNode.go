package name

import (
	"context"
	"runtime"

	"github.com/billziss-gh/cgofuse/fuse"
)

func (nn *ipnsNode) String() string {
	path := nn.fusePath.String()
	if nn.key != nil {
		return gopath.Join(nn.key.Path().String(), path)
	}
	return path
}

func (_ *ipnsNode) NameSpace() string {
	return "ipns"
}

func (nn *ipnsNode) YieldIo(ctx context.Context, nodeType FsType) (interface{}, error) {
	if err := nn.recordBase.typeCheck(nodeType); err != nil {
		return nil, err
	}

	if nn.key == nil { // use core for non-keyed paths
		return nn.ipfsNode.YieldIo(ctx, nodeType)
	}

	// check that our key is still valid
	if err = checkAPIKeystore(ctx, nn.core.Key(), nn.key); err != nil {
		nn.key = nil
		return nil, err
	}

	if nn.path == "/" {
		switch nn.metadata.Mode & fuse.S_IFMT {
		case fuse.S_IFREG:
			return keyYieldFileIO(ctx, nn.key, nn.core)
		case fuse.S_IFLNK:
			var (
				ipldNode ipld.Node
				target   string
				lnk      *link
			)
			if ipldNode, err = nn.core.ResolveNode(ctx, nn.key.Path()); err != nil {
				goto linkEnd
			}

			if target, err = ipldReadLink(ipldNode); err != nil {
				goto linkEnd
			}
			lnk = &link{target: target}
			_, err = lnk.InitMetadata(ctx)

		linkEnd:
			return lnk, err

		case fuse.S_IFDIR:
			// fallback to MFS IO handler to list out root contents
			break
		default:
			return nil, errUnexpected
		}
	}

	//handle other nodes via MFS
	nn.fsRootIndex.Lock()
	defer nn.fsRootIndex.Unlock()
	nn.root, err = nn.fsRootIndex.Request(keyName)
	switch err {
	case nil:
		break
	default:
		return nil, err
	case errNotInitialized:
		nn.root, err = ipnsToMFSRoot(ctx, key.Path(), nn.core)
		if err != nil {
			return nil, err
		}

		nn.fsRootIndex.Register(keyName, mroot)
		if nn.root, err = nn.fsRootIndex.Request(keyName); err != nil {
			return nil, err
		}
	}

	nnIO, err := nn.mfsNode.YieldIo(ctx)
	if err != nil {
		runtime.SetFinalizer(nnIO, ipnsKeyRootFree)
	}
	return nnIO, err
}
