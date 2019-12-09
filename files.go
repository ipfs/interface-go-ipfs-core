package iface

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-cidutil/cidenc"
	"github.com/ipfs/go-mfs"
)

type FilesCopyOptions struct {
	Flush bool
}

type FilesMoveOptions struct {
	Flush bool
}

type FilesListOptions struct {
	CidEncoder cidenc.Encoder
	Long       bool
}

type FilesRemoveOptions struct {
	Force     bool
	Recursive bool
}

type FilesStatOptions struct {
	CidEncoder   cidenc.Encoder
	WithLocality bool
}

type FilesWriteOptions struct {
	Create            bool
	MakeParents       bool
	Truncate          bool
	Flush             bool
	RawLeaves         bool
	RawLeavesOverride bool
	Offset            int64
	// Count is the number of bytes to write. 0 (default) writes everything.
	Count      int64
	CidBuilder cid.Builder
}

type FilesReadOptions struct {
	Offset int64
	// Count is the number of bytes to read. 0 (default) reads everything.
	Count int64
}

type FilesMkdirOptions struct {
	MakeParents bool
	Flush       bool
	CidBuilder  cid.Builder
}

type FilesChangeCidOptions struct {
	Flush      bool
	CidBuilder cid.Builder
}

type FileInfo struct {
	Cid            cid.Cid
	Size           uint64
	CumulativeSize uint64
	Blocks         int
	Type           string
	WithLocality   bool   `json:",omitempty"`
	Local          bool   `json:",omitempty"`
	SizeLocal      uint64 `json:",omitempty"`
}

// FilesAPI specifies an interface to interact with the Mutable File System
// layer.
type FilesAPI interface {
	Copy(ctx context.Context, src, dst string, opts *FilesCopyOptions) error
	Move(ctx context.Context, src, dst string, opts *FilesMoveOptions) error
	List(ctx context.Context, path string, opts *FilesListOptions) ([]mfs.NodeListing, error)
	Remove(ctx context.Context, path string, opts *FilesRemoveOptions) error
	Stat(ctx context.Context, path string, opts *FilesStatOptions) (*FileInfo, error)
	Read(ctx context.Context, path string, opts *FilesReadOptions) (io.ReadCloser, error)
	Write(ctx context.Context, path string, r io.Reader, opts *FilesWriteOptions) error
	Mkdir(ctx context.Context, path string, opts *FilesMkdirOptions) error
	Flush(ctx context.Context, path string) (cid.Cid, error)
}
