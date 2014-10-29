// Copyright 2014 The erzha Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package kernel

import (
	"os"
	"runtime"
	"os/signal"
	"syscall"

	"code.google.com/p/go.net/context"
	
	"github.com/erzha/econf"
	"github.com/erzha/elog"
	
)

var serverObj *Server

var serverCtx context.Context
var serverCancel context.CancelFunc

type ServerHandler interface {
	Serve(ctx context.Context, p *Server)
}


type Server struct {
	Handler ServerHandler
	PluginOrder []string
	Conf *econf.Conf
	Logger *elog.Logger
	sigIntC		chan bool //recv sigint sigkill to kill process
	sigIntCount int
}

func newServer() *Server {
	p := &Server{}
	p.sigIntC = make(chan bool)
	return p
}

func Boot(handler ServerHandler) {
	
	parseArgs()
	
	serverObj = newServer()
	serverObj.Conf = initConf()
	serverObj.Logger = initLogger(serverObj.Conf)
	serverObj.Handler = handler
	serverObj.Boot()
}

func (p *Server) Basedir() string {
	return *flagBasedir
}

func (p *Server) Boot() {
	p.Logger.Debug("kernel_server_boot")
	runtime.GOMAXPROCS(runtime.NumCPU())

	go p.handleControlSignal()
	serverInit(serverCtx, p)

	chanHandlerDone := make(chan bool)
	go func() {
		p.Logger.Debug("kernel_server_boot_handler")
		p.Handler.Serve(serverCtx, p)
		close(chanHandlerDone)
	}()

	select {
		case <-chanHandlerDone:
		case <-serverCtx.Done():
	}

	serverShutdown(serverCtx, p)
}

func (p *Server) handleControlSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGUSR1)

	for true {
		s := <-c
		switch s {
		case syscall.SIGINT:
			serverCancel()
		}
	}
}

func init() {
	serverCtx, serverCancel = context.WithCancel(context.Background())
}
