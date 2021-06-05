package iface

import (
	"context"
	"fmt"
)

func ExamplePinAPI_Ls() {
	var api CoreAPI

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pins, errC := api.Pin().Ls(ctx)
	for pin := range pins {
		fmt.Println(pin.Path())
	}
	if err := <-errC; err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func ExamplePinAPI_Verify() {
	var api CoreAPI

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	status, errC := api.Pin().Verify(ctx)
	for pin := range status {
		if !pin.Ok() {
			for _, missing := range pin.BadNodes() {
				fmt.Printf("missing node %s: %s\n", missing.Path().Cid(), missing.Err())
			}
		}
	}
	if err := <-errC; err != nil {
		fmt.Printf("failed to complete pin verification request: %s\n", err)
	}
}
