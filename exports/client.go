package main

/*
#include <stdint.h>
#include "./error.h"
#include "./client.h"
*/
import "C"
import (
	"time"
	"unsafe"

	"github.com/BrobridgeOrg/gravity-sdk/core"
	pointer "github.com/mattn/go-pointer"
)

//export NewClientOptions
func NewClientOptions() *C.ClientOptions {

	opts := core.NewOptions()

	// Convert to C struct
	copts := (*C.ClientOptions)(C.malloc(C.size_t(unsafe.Sizeof(C.ClientOptions{}))))
	copts.appID = C.CString(opts.AppID)
	copts.appKey = C.CString(opts.AppKey)
	copts.pingInterval = C.longlong(int64(opts.PingInterval))
	copts.maxPingsOutstanding = C.int(opts.MaxPingsOutstanding)
	copts.maxReconnects = C.int(opts.MaxReconnects)

	return copts
}

//export NewClient
func NewClient() *C.Client {

	client := core.NewClient()

	c := (*C.Client)(C.malloc(C.size_t(unsafe.Sizeof(C.Client{}))))
	c.instance = pointer.Save(client)

	return c
}

//export ClientConnect
func ClientConnect(c *C.Client, host *C.char, options *C.ClientOptions) *C.GravityError {

	client := pointer.Restore(c.instance).(*core.Client)

	opts := core.NewOptions()
	opts.AppID = C.GoString(options.appID)
	opts.AppKey = C.GoString(options.appKey)
	opts.PingInterval = time.Duration(options.pingInterval)
	opts.MaxPingsOutstanding = int(options.maxPingsOutstanding)
	opts.MaxReconnects = int(options.maxReconnects)
	opts.DisconnectHandler = func() {
		C.callClientEventHandler(c.disconnectHandler)
	}
	opts.ReconnectHandler = func() {
		C.callClientEventHandler(c.reconnectHandler)
	}

	err := client.Connect(C.GoString(host), opts)
	if err != nil {
		return NewError(err.Error())
	}

	return nil
}

//export ClientSetDisconnectHandler
func ClientSetDisconnectHandler(c *C.Client, callback C.ClientEventHandler) {
	c.disconnectHandler = callback
}

//export ClientSetReconnectHandler
func ClientSetReconnectHandler(c *C.Client, callback C.ClientEventHandler) {
	c.reconnectHandler = callback
}
