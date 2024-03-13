// Package wshooks provides registering and calling hooks
package wshooks

import (
	"github.com/gorilla/websocket"
)

// Hooks to run when a new client connects
var onConnectHooks []func(conn *websocket.Conn)

// RegisterOnConnectHandler registers a handler that gets invoked when a new client connects.
func RegisterOnConnectHandler(handler func(conn *websocket.Conn)) {
	onConnectHooks = append(onConnectHooks, handler)
}

// CallOnConnectHandlers is the hook that gets called when a client connects
func CallOnConnectHandlers(conn *websocket.Conn) {
	for _, v := range onConnectHooks {
		v(conn)
	}
}
