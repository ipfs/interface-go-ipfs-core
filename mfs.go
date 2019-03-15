package iface

import (
	"context"
	"io"
	"os"
	"strings"
)

// MfsPath
//
// Note that Path and ResolvedPath interfaces satisfy this interface
//
// This interface is inspired by https://godoc.org/github.com/src-d/go-billy,
// but doesn't implement it
type MfsPath interface {
	// String returns the path as a string.
	String() string

	// Mutable returns false if the data pointed to by this path in guaranteed
	// to not change.
	//
	// Note that resolved mutable path can be immutable.
	Mutable() bool
}

type filesPath struct {
	p string
}

func (p *filesPath) String() string {
	return p.p
}

func (p *filesPath) Mutable() bool {
	return !(strings.HasPrefix(p.p, "/ipfs") || strings.HasPrefix(p.p, "/ipld"))
}

func FilePath(p string) MfsPath {
	// TODO: more validation
	return &filesPath{
		p: p,
	}
}

type MfsAPI interface {
	Create(ctx context.Context, path MfsPath) (File, error)
	Open(ctx context.Context, path MfsPath) (File, error)
	OpenFile(ctx context.Context, path MfsPath, flag int, perm os.FileMode) (File, error)

	Stat(ctx context.Context, path MfsPath) (os.FileInfo, error)

	Rename(ctx context.Context, oldpath, newpath MfsPath) error
	Remove(ctx context.Context, path MfsPath) error

	// ReadDir reads the directory named by dirname and returns a list of
	// directory entries sorted by filename.
	ReadDir(ctx context.Context, path MfsPath) ([]os.FileInfo, error)
	// MkdirAll creates a directory named path, along with any necessary
	// parents, and returns nil, or else returns an error. The permission bits
	// perm are used for all directories that MkdirAll creates. If path is
	// already a directory, MkdirAll does nothing and returns nil.
	MkdirAll(ctx context.Context, path MfsPath, perm os.FileMode) error

	// TODO: ChCid
	// TODO: Symlink stuff (is it implemented in mfs?)
}

type File interface {
	io.Writer
	io.WriterAt
	io.Reader
	io.Seeker
	io.Closer

	Name() MfsPath

	// Truncate the file.
	Truncate(size int64) error
}

var _ MfsPath = (Path)(nil)
var _ MfsPath = (ResolvedPath)(nil)
