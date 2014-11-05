// Copyright 2014 The erzha Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package kernel

import (
	"errors"

	"code.google.com/p/go.net/context"
)


/*
Plugin is used to make erzha's function extended handily


Creater will be called when a request arrive and before the user defined action.
the kernel will get a instance of the plugin and calls it's RequestInit method before the user defined action and calls it's
RequestShutdown method when user defined action ended.

in action, we can use the kernel.Sapi instance with the request to get the plugin's instance, for ex:

	func Execute(ctx context.Context, sapi *http.Sapi) {
		pluginSessionInstance := sapi.Kernel.Plugin("session").(session.*Session)
	}
*/
type PluginInfo struct {
	Creater			func() (interface{}, error)
	ServerInit      func(ctx context.Context, server *Server) error
	ServerShutdown  func(ctx context.Context, server *Server) error
	RequestInit     func(ctx context.Context, sapi *Sapi, obj interface{}) error
	RequestShutdown func(ctx context.Context, sapi *Sapi, obj interface{}) error
}

var pluginMap map[string]PluginInfo

//If it's returned after a plugin's RequestInit/Shutdown method, the request will be ended directly
var PluginStop = errors.New("kernel plugin stop here")

func RegisterPlugin(name string, hookInfo PluginInfo) {
	pluginMap[name] = hookInfo
}

func serverInit(ctx context.Context, server *Server) {
	for key, info := range pluginMap {
		if nil != info.ServerInit {
			info.ServerInit(ctx, server)
			serverObj.Logger.Sysf("server_init_done_%s", key)
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

func requestInit(ctx context.Context, sapi *Sapi) error {
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
			err := info.RequestInit(ctx, sapi, obj)
			switch err {
				case PluginStop:
					return PluginStop
				case nil:
				default:
					sapi.Server.Logger.Warning("request_init_error plugin:%s err:%s", name, err.Error())
					return err //if err != PluginStop, the request will alive to continue rather than dead
			}
		}
	}
	return nil
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

