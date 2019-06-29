package fsnode

import (
	"context"
	"time"

	"github.com/billziss-gh/cgofuse/fuse"
	fs "github.com/ipfs/interface-go-ipfs-core/filesystem/interface"
)

type SoftDirRoot struct {
	BaseNode
	//fs *FUSEIPFS
}

type mountRoot struct {
	SoftDirRoot
	subroots []string
}

func Lookup(ctx context.Context, name string) (FsPath, error) {
}

const (
	filesNamespace  = "files"
	filesRootPath   = "/" + filesNamespace
	filesRootPrefix = filesRootPath + "/"
)

//TODO: add arg "roots ...string"?
func RootParser(epoch time.Time) fs.ParseFn {
	return func(_ context.Context, _ string) (FsPath, error) {
		return &mountRoot{subroots: []string{"/ipfs", "/ipns", filesRootPrefix},
			csd(path, epoch)}, nil
	}
}

func (sd *SoftDirRoot) InitMetadata(ctx context.Context) (*fuse.Stat_t, error) {
	nodeStat, err := sd.BaseNode.InitMetadata(ctx)
	if err != nil {
		return nodeStat, err
	}
	nodeStat.Mode |= fuse.S_IFDIR
	return nodeStat, nil
}

func (mr *mountRoot) YieldIo(ctx context.Context) (io interface{}, err error) {
	return mr, nil
}

func (mr *mountRoot) Read(ctx context.Context, offset int64) <-chan string {
	return stringStream(ctx, mr.subroots)
}

func (mr *mountRoot) Entries() int {
	return len(mr.subroots)
}
