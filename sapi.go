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
	
	Server *Server
	Ext interface{}

	plugins map[string]interface{}
}

func (p *Sapi) Plugin(name string) interface{} {
	return p.plugins[name]
}

func (p *Sapi) Print(param ...interface{}) {
	fmt.Fprint(p.Stdout, param...)
}

func (p *Sapi) Println(param ...interface{}) {
	fmt.Fprintln(p.Stdout, param...)
}

var FireFunc func(ctx context.Context, sapi *Sapi)

func FireAction(ctx context.Context, sapi *Sapi, do func(ctx context.Context, sapi *Sapi)) {
	
	/*
	defer func() {
		r := recover()
		if nil!=r {
			sapi.Server.Logger.Info(r)
		}
	}()*/

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

func NewSapi() *Sapi {
	ret := &Sapi{}
	ret.Server = serverObj
	ret.plugins = make(map[string]interface{})
	return ret
}
