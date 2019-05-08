package fs

import (
	"context"
	"os"
	//pretend we're not the same package
	//fs "github.com/ipfs/interface-go-ipfs-core/filesystem"
)

// This file represents a 3rd party user of the API

// Client side using pkg defaults
func client() {
	f, err := fs.OpenFile("/ipfs/Qm.../file.ext", flags, perm)
	if err != nil {
		//...
	}
	defer f.Close()
	x, y := f.Read(someByteSlice)
}

// done

// In addition, clients can also implement their own...
// nodes,
type myNode struct {
	myData bool
}

func (mn *myNode) YieldIo(ctx context.Context, nodeType fs.FsType) {
	if !mn.myData {
		return nil, errIOType // we don't have the data needed to perform this request
	}
	return constructFileIo(mn.MyData), nil //we have what we need to construct a FsFile interface, so return that
}

// namespace parsers
func myParser(ctx context.Context, name string) (FsNode, error) {
	if condition(name) {
		return myNode{false}, nil
	}
	return myNode{true}, nil
}

// and filing systems
func aDifferentFileSystem(ctx context.Context) {
	// something like this
	fsCtx := deriveFrom(ctx)
	return &colonDelimitedRegistry{fsCtx}, nil
}

func clientImp() {
	myFs := aDifferentFileSystem(ctx)

	_, err = myFs.OpenFile(":myName", flags, perm)
	err == os.ErrNotExist // or some other standard error

	namespaceCloser, err := myFs.Register(":myName", myParser)

	myF, err := myFs.OpenFile(":myName", flags, perm)
	myF.Read()

	namespaceCloser() // requests for `/myName` are no longer valid according to the filesystem interface (`myFs`)

	_, err = myFs.OpenFile(":myName", flags, perm)
	err == os.ErrNotExist

	_, err = myF.Read()
	err == os.ErrInvalid // or some other standard error
}
