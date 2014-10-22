package kernel

import (
	"errors"

	"code.google.com/p/go.net/context"
)

type PluginInfo struct {
	Creater                   func() (interface{}, error)
	ServerInit      func(ctx context.Context, server *Server) error
	ServerShutdown  func(ctx context.Context, server *Server) error
	RequestInit     func(ctx context.Context, sapi *Sapi, obj interface{}) error
	RequestShutdown func(ctx context.Context, sapi *Sapi, obj interface{}) error
}

var pluginMap map[string]PluginInfo
var PluginJump = errors.New("kernel plugin jump  here")
var PluginStop = errors.New("kernel plugin return here")

func RegisterPlugin(name string, hookInfo PluginInfo) {
	pluginMap[name] = hookInfo
}

func serverInit(ctx context.Context, server *Server) {
	for _, info := range pluginMap {
		if nil != info.ServerInit {
			info.ServerInit(ctx, server)
		}
	}
}

func serverShutdown(ctx context.Context, server *Server) {
	for _, info := range pluginMap {
		if nil != info.ServerShutdown {
			info.ServerShutdown(ctx, server)
		}
	}
}

func requestInit(ctx context.Context, sapi *Sapi) {
	for name, info := range pluginMap {
		if nil == info.Creater {
			continue
		}

		obj, err := info.Creater()
		if nil != err {
			continue
		}
		sapi.plugins[name] = obj

		if nil != info.RequestInit {
			info.RequestInit(ctx, sapi, obj)
		}
	}
}

func requestShutdown(ctx context.Context, sapi *Sapi) {
	for name, info := range pluginMap {
		if nil == info.Creater || nil == info.RequestShutdown {
			continue
		}

		info.RequestShutdown(ctx, sapi, sapi.plugins[name])
	}
}

func init() {
	pluginMap = make(map[string]PluginInfo)
}

