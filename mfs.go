package iface

import (
	"context"
	"io"
	"os"
)

type MfsAPI interface {
	Create(ctx context.Context, path Path) (File, error)
	Open(ctx context.Context, path Path) (File, error)
	OpenFile(ctx context.Context, path Path, flag int, perm os.FileMode) (File, error)

	Stat(ctx context.Context, path Path) (os.FileInfo, error)

	Rename(ctx context.Context, oldpath, newpath Path) error
	Copy(ctx context.Context, oldpath, newpath Path) error
	Remove(ctx context.Context, path Path) error

	// ReadDir reads the directory named by dirname and returns a list of
	// directory entries sorted by filename.
	ReadDir(ctx context.Context, path Path) ([]os.FileInfo, error)
	// MkdirAll creates a directory named path, along with any necessary
	// parents, and returns nil, or else returns an error. The permission bits
	// perm are used for all directories that MkdirAll creates. If path is
	// already a directory, MkdirAll does nothing and returns nil.
	MkdirAll(ctx context.Context, path Path, perm os.FileMode) error

	Flush(ctx context.Context, path Path) error
	// TODO: ChCid
	// TODO: Symlink stuff (is it implemented in mfs?)
}

type File interface {
	io.Writer
	io.WriterAt
	io.Reader
	io.Seeker
	io.Closer

	// Truncate the file.
	Truncate(size int64) error
}
