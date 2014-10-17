package kernel

import (
	"fmt"
	"io"

	"code.google.com/p/go.net/context"
)

type Sapi struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	Ext interface{}

	plugins map[string]interface{}
}

func (p *Sapi) Plugin(name string) interface{} {
	return p.plugins[name]
}

var FireFunc func(ctx context.Context, sapi *Sapi)

func FireAction(ctx context.Context, sapi *Sapi, do func(ctx context.Context, sapi *Sapi)) {
	requestDone := make(chan bool)

	go func() {
		requestInit(ctx, sapi)
		do(ctx, sapi)
		requestShutdown(ctx, sapi)
		close(requestDone)
	}()

	select {
		case <-ctx.Done():
		case <-requestDone:
	}
}
