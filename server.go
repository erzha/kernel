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
	sigIntC		chan bool //接收SIG_INT信号，用于强制结束程序
	sigIntCount int	//SIG_INT信号次数
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
	runtime.GOMAXPROCS(runtime.NumCPU())

	go p.handleControlSignal()
	serverInit(serverCtx, p)

	chanHandlerDone := make(chan bool)
	go func() {
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
