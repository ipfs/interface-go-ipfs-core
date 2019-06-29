package fsnode

import (
	"context"
	"strings"
	"sync"

	"github.com/billziss-gh/cgofuse/fuse"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	fs "github.com/ipfs/interface-go-ipfs-core/filesystem/interface"
)

type BaseNode struct {
	sync.RWMutex
	path    string
	version uint

	metadata Metadata
}

// implementation of Metadata interface
type Metadata struct {
	//apiPath string
	fStat     fuse.Stat_t
	ipfsFlags fs.Flags
}

func (md *Metadata) Type() fs.Kind {
	t := coreiface.TUnknown

	switch md.fStat.Mode & fuse.S_IFMT {
	case fuse.S_IFREG:
		t = coreiface.TFile
	case fuse.S_IFDIR:
		t = coreiface.TDirectory
	case fuse.S_IFLNK:
		t = coreiface.TSymlink
	}

	return fs.Kind(t)
}

func (md *Metadata) Flags() fs.Flags {
	return md.ipfsFlags
}

func (md *Metadata) Size() uint {
	return uint(md.fStat.Size)
}

// conversion between fStat and interface
func (rb *BaseNode) Metadata(_ context.Context) (fs.Metadata, error) {
	return &rb.metadata, nil
}

func NewBase(name string) BaseNode {
	return BaseNode{path: name}
}

func (rb *BaseNode) String() string {
	return rb.path
}

func (rb *BaseNode) Remove(_ context.Context) error {
	return fs.ErrNotImplemented
}

func (rb *BaseNode) Create(_ context.Context, _ fs.Kind) error {
	return fs.ErrNotImplemented
}

//Core API Path interface
func (rb *BaseNode) IsValid() error {
	return fs.ErrNotImplemented
}

func (rb *BaseNode) Mutable() bool {
	return false
}

func (rb *BaseNode) Namespace() string {
	i := strings.IndexRune(rb.path[1:], '/')
	if i == -1 {
		return "root"
	}
	return rb.path[1:i]
}

const (
	IRWXA = fuse.S_IRWXU | fuse.S_IRWXG | fuse.S_IRWXO
	IRXA  = IRWXA &^ (fuse.S_IWUSR | fuse.S_IWGRP | fuse.S_IWOTH)
)

func (rb *BaseNode) InitMetadata(_ context.Context) (*fuse.Stat_t, error) {
	now := fuse.Now()
	rb.metadata.fStat.Birthtim,
		rb.metadata.fStat.Atim,
		rb.metadata.fStat.Mtim,
		rb.metadata.fStat.Ctim = now, now, now, now //!!!!
	rb.metadata.fStat.Mode = IRXA
	return &rb.metadata.fStat, nil
}

func (rb *BaseNode) typeCheck(nodeType fs.Kind) (err error) {
	if !typeCheck(rb.metadata.fStat.Mode, nodeType) {
		return fs.ErrIOType
	}
	return nil
}

func typeCheck(pMode uint32, nodeType fs.Kind) bool {
	switch nodeType {
	case fs.UfsFile:
		return fuse.S_IFREG == (pMode & fuse.S_IFDIR)
	case fs.UfsDirectory:
		return fuse.S_IFDIR == (pMode & fuse.S_IFDIR)
	case fs.UfsSymlink:
		return fuse.S_IFLNK == (pMode & fuse.S_IFDIR)
	default:
		return false
	}
}

/*
func genQueryID(path string, md fs.Metadata) (cid.Cid) {
	mdBuf := make([]byte, 16)
	binary.LittleEndian.PutUint64(mdBuf, uint64(md.Size()))
	binary.LittleEndian.PutUint64(mdBuf[8:], uint64(md.Type()))

	prefix := cid.V1Builder{Codec: cid.DagCBOR, MhType: multihash.BLAKE2B_MIN}
	return prefix.Sum(append([]byte(path), mdBuf...))
}
*/
