package fsnode

import (
	"context"
	"time"

	"github.com/billziss-gh/cgofuse/fuse"
)

func csd(path string, metaTimes time.Time) softDirRoot {
	sd := softDirRoot{recordBase: crb(path)}
	now := fuse.NewTimespec(metaTimes)
	meta := &sd.recordBase.metadata
	meta.Birthtim, meta.Mtim, meta.Ctim = now, now, now // !!!
	meta.Atim = fuse.Now()
	return sd
}

func crb(path string) fsnode.BaseNode {
	return fsnode.BaseNode{path: path, ioHandles: make(nodeHandles)}
}

func stringStream(ctx context.Context, strings ...string) <-chan string {
	stringChan := make(chan string)
	go func() {
		for _, s := range strings {
			select {
			case ctx.Done():
				return
			case stringChan <- s:
			}
		}
		close(stringChan)
	}()
	return stringChan
}

func pinMux(ctx context.Context, pins ...coreiface.Pin) <-chan string {
	pinChan := make(chan string)
	go func() {
		for _, pin := range pins {
			select {
			case ctx.Done():
				return
			case pinChan <- pin.String():
			}
		}
		close(pinChan)
	}()
	return pinChan
}
