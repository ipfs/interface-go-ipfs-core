package iface

import (
	"context"
	"fmt"

	"github.com/ipfs/interface-go-ipfs-core/path"
)

func ExampleNameAPI_Search() {
	var api CoreAPI

	_ = func(ctx context.Context) (result path.Path, err error) {
		results, errC := api.Name().Search(ctx, "foobar")
		for result = range results {
			fmt.Println(result)
		}
		return result, <-errC
	}
}
