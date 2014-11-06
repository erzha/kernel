// Copyright 2014 The erzha Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package kernel

import (
	"fmt"
	"io"

	"golang.org/x/net/context"
)

/*
Sapi represent the request in network application model.

this Sapi is the kernel Sapi, the diy server may have self sapi struct, but it must have a member named Kernel pointed to the *Kernel.Sapi.

for ex:

	//github.com/erzha/http/server
	package server
	import "kernel"

	type Sapi struct {
		Kernel *kernel.Sapi
	}



*/
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

//The func type for Kernel.Sapi to call the user defined logic
var FireFunc func(ctx context.Context, sapi *Sapi)

func FireAction(ctx context.Context, sapi *Sapi, do func(ctx context.Context, sapi *Sapi)) {

	requestDone := make(chan bool)
	go func() {
		defer func() {
			close(requestDone)
			r := recover()
			if nil!=r {
				sapi.Server.Logger.Info("panic occured", r)
			}
		}()
		
		if PluginStop == requestInit(ctx, sapi) {
			return
		}

		do(ctx, sapi)
		requestShutdown(ctx, sapi)
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
