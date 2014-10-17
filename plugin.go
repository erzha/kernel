package kernel

import (
	"errors"

	"code.google.com/p/go.net/context"
)

type PluginInfo struct {
	Creater                   func() (interface{}, error)
	HookPluginServerInit      func(ctx context.Context, server *Server) error
	HookPluginServerShutdown  func(ctx context.Context, server *Server) error
	HookPluginRequestInit     func(ctx context.Context, sapi *Sapi) error
	HookPluginRequestShutdown func(ctx context.Context, sapi *Sapi) error
}

var pluginMap map[string]PluginInfo
var PluginJump = errors.New("kernel plugin jump  here")
var PluginStop = errors.New("kernel plugin return here")

func RegisterPlugin(name string, hookInfo PluginInfo) {
	pluginMap[name] = hookInfo
}

func serverInit(ctx context.Context, server *Server) {

}

func serverShutdown(ctx context.Context, server *Server) {

}

func requestInit(ctx context.Context, sapi *Sapi) {

}

func requestShutdown(ctx context.Context, sapi *Sapi) {

}

func init() {
	pluginMap = make(map[string]PluginInfo)
}

