package fsnode

func csd(path string, now fuse.Timespec) softDirRoot {
	sd := softDirRoot{recordBase: crb(path)}
	meta := &sd.recordBase.metadata
	meta.Birthtim, meta.Mtim, meta.Ctim = now, now, now // !!!
	meta.Atim = fuse.Now()
	return sd
}

func crb(path string) fsnode.BaseNode {
	return fsnode.BaseNode{path: path, ioHandles: make(nodeHandles)}
}