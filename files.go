package iface

import (
	"context"
	"os"

	"github.com/ipfs/go-cid"
)

type FileInfo interface {
	os.FileInfo
	CID() cid.Cid
}

type FilesAPI interface {
	Stat(ctx context.Context, path string) (FileInfo, error)
}

/* AFERO interface for reference

Chmod(name string, mode os.FileMode) : error
Chtimes(name string, atime time.Time, mtime time.Time) : error
Create(name string) : File, error
Mkdir(name string, perm os.FileMode) : error
MkdirAll(path string, perm os.FileMode) : error
Name() : string
Open(name string) : File, error
OpenFile(name string, flag int, perm os.FileMode) : File, error
Remove(name string) : error
RemoveAll(path string) : error
Rename(oldname, newname string) : error
Stat(name string) : os.FileInfo, error
*/
