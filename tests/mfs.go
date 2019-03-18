package tests

import (
	"context"
	"testing"

	"github.com/ipfs/interface-go-ipfs-core"
)

func (tp *provider) TestMfs(t *testing.T) {
	t.Run("TestFilesDirs", tp.TestFilesDirs)
}

func (tp *provider) TestFilesDirs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api, err := tp.makeAPI(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}

	l, err := api.Files().ReadDir(ctx, iface.FilePath("/test"))
	if len(l) != 0 {
		return
	}
}
