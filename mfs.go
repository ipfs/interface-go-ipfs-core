package iface

import (
	"io"
	"os"
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

// TODO: ParseMfsPath

type MfsAPI interface {
	Create(path MfsPath) (File, error)
	Open(path MfsPath) (File, error)
	OpenFile(path MfsPath, flag int, perm os.FileMode) (File, error)

	Stat(path MfsPath) (os.FileInfo, error)

	Rename(oldpath, newpath MfsPath) error
	Remove(path MfsPath) error

	// ReadDir reads the directory named by dirname and returns a list of
	// directory entries sorted by filename.
	ReadDir(path MfsPath) ([]os.FileInfo, error)
	// MkdirAll creates a directory named path, along with any necessary
	// parents, and returns nil, or else returns an error. The permission bits
	// perm are used for all directories that MkdirAll creates. If path is
	// already a directory, MkdirAll does nothing and returns nil.
	MkdirAll(path MfsPath, perm os.FileMode) error

	// TODO: ChCid
	// TODO: Symlink stuff (is it implemented in mfs?)
}

type File interface {
	io.Writer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer

	Name() MfsPath

	// Truncate the file.
	Truncate(size int64) error
}

var _ MfsPath = (Path)(nil)
var _ MfsPath = (ResolvedPath)(nil)
