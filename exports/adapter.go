package main

/*
#include <stdint.h>
#include <stdbool.h>
#include "./error.h"
#include "./client.h"

typedef struct {
	char *endpoint;
	char *domain;
	int batchSize;
	bool verbose;
} AdapterOptions;

typedef struct {
	void *instance;
} Adapter;
*/
import "C"
import (
	"unsafe"

	"github.com/BrobridgeOrg/gravity-sdk/adapter"
	"github.com/BrobridgeOrg/gravity-sdk/core"
	pointer "github.com/mattn/go-pointer"
)

func getAdapterNativeOptions(options *C.AdapterOptions) *adapter.Options {

	opts := adapter.NewOptions()
	opts.Endpoint = C.GoString(options.endpoint)
	opts.Domain = C.GoString(options.domain)
	opts.BatchSize = int(options.batchSize)
	opts.Verbose = bool(options.verbose)

	return opts
}

//export NewAdapterOptions
func NewAdapterOptions() *C.AdapterOptions {

	opts := adapter.NewOptions()

	// Convert to C struct
	aopts := (*C.AdapterOptions)(C.malloc(C.size_t(unsafe.Sizeof(C.AdapterOptions{}))))
	aopts.endpoint = C.CString(opts.Endpoint)
	aopts.domain = C.CString(opts.Domain)
	aopts.batchSize = C.int(opts.BatchSize)
	aopts.verbose = C.bool(opts.Verbose)

	return aopts
}

//export NewAdapterWithClient
func NewAdapterWithClient(client *C.Client, options *C.AdapterOptions) *C.Adapter {

	c := pointer.Restore(client.instance).(*core.Client)
	opts := getAdapterNativeOptions(options)

	s := adapter.NewAdapterConnectorWithClient(c, opts)

	// Prepare subscriber struct
	sub := (*C.Adapter)(C.malloc(C.size_t(unsafe.Sizeof(C.Adapter{}))))
	sub.instance = pointer.Save(s)

	return sub
}

//export AdapterRegister
func AdapterRegister(s *C.Adapter, componentName *C.char, subscriberID *C.char, subscriberName *C.char) *C.GravityError {

	a := pointer.Restore(s.instance).(*adapter.AdapterConnector)

	compName := C.GoString(componentName)
	adapterID := C.GoString(subscriberID)
	adapterName := C.GoString(subscriberName)

	err := a.Register(compName, adapterID, adapterName)
	if err != nil {
		return NewError(err.Error())
	}

	return nil
}

//export AdapterPublish
func AdapterPublish(ca *C.Adapter, eventName *C.char, payload *C.char) *C.GravityError {

	a := pointer.Restore(ca.instance).(*adapter.AdapterConnector)

	err := a.Publish(C.GoString(eventName), StrToBytes(C.GoString(payload)), nil)
	if err != nil {
		return NewError(err.Error())
	}

	return nil
}
