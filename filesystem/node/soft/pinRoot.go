package fsnode

import (
	"context"
	"time"

	//TODO: consider /ipfs/interface-go-ipfs-filesystem/node ?
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	fs "github.com/ipfs/interface-go-ipfs-core/filesystem/interface"
	coreoptions "github.com/ipfs/interface-go-ipfs-core/options"
)

type pinRoot struct {
	SoftDirRoot
	pinAPI coreiface.PinAPI
}

func PinParser(pinAPI coreiface.PinAPI, epoch time.Time) fs.ParseFn {
	return func(_ context.Context, path string) (fs.Node, error) {
		if path != "" {
			return nil, fs.ErrInvalidPath
		}
		return &pinRoot{pinAPI: pinAPI, SoftDirRoot: csd(path, epoch)}, nil
	}
}

func (pr *pinRoot) Version() uint {
	pr.version++ //TODO: we need a "has changed" signal from the pinservice to be cache friendly
	return pr.version
}

func (pr *pinRoot) YieldIo(ctx context.Context) (io interface{}, err error) {
	//TODO: some way to ping the pinapi or coreapi here
	return pr, nil
}

func (pr *pinRoot) Read(ctx context.Context, offset int64) <-chan string {
	pins, err := pr.pinAPI.Ls(ctx, coreoptions.Pin.Type.Recursive())
	if err != nil {
		return nil
	}
	return pinMux(ctx, pins...)
}

func (pr *pinRoot) Entries() int {
	pins, err := pr.pinAPI.Ls(ctx, coreoptions.Pin.Type.Recursive())
	if err != nil {
		return 0
	}
	return len(pins)
}
