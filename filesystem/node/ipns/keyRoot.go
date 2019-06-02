package name

import (
	"context"

	"github.com/billziss-gh/cgofuse/fuse"
)

type keyRoot struct {
	softDirRoot
	keyAPI coreiface.KeyAPI
}

func (kr *keyRoot) InitMetadata(ctx context.Context) (*fuse.Stat_t, error) {
	nodeStat, err := kr.softDirRoot.InitMetadata(ctx)
	if err != nil {
		return nodeStat, err
	}

	nodeStat.Mode |= fuse.S_IWUSR
	return nodeStat, nil
}

func (kr *keyRoot) YieldIo(ctx context.Context, nodeType FsType) (io interface{}, err error) {
	if err := in.recordBase.typeCheck(nodeType); err != nil {
		return nil, err
	}

	keys, err := kr.keyAPI.List(ctx)
	if err != nil {
		return nil, err
	}

	keyChan := make(chan directoryStringEntry)
	asyncContext := deriveTimerContext(ctx, entryTimeout)
	go func() {
		defer close(keyChan)
		for _, key := range keys {
			select {
			case <-asyncContext.Done():
				return
			case keyChan <- directoryStringEntry{string: key.Name()}:
				continue
			}
		}

	}()
	return backgroundDir(asyncContext, len(keys), keyChan)
}
