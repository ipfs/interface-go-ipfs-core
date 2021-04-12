package iface

import (
	"context"
)

type DNSAPI interface {
	LookupTXT(ctx context.Context, name string) (txt []string, err error)
}
