package ipfs

import "github.com/ipfs/interface-go-ipfs-core/filesystem/interface"

//TODO: store these in the daemon/ipfs-node scope, or elsewhere
// have something  extract the FS from the daemon (`core.fsFrom(IpfsNode)``)
var 	pkgRoot fs.FileSystem

func init() {
	pkgRoot, err = newDefaultFs()
	if err != nil {
		panic(err)
	}
}