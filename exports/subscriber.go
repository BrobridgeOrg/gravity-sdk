package main

/*
#include <stdint.h>
#include <stdbool.h>

typedef struct {
	bool enabled;
	char *mode;
	uint64_t omittedCount;
} SubscriberInitialLoadOptions;

typedef struct {
	char *endpoint;
	char *domain;
	int workerCount;
	int bufferSize;
	int chunkSize;
	bool verbose;
	SubscriberInitialLoadOptions initialLoad;
} SubscriberOptions;

typedef struct {
	char *message;
} GravityError;

typedef struct {
	int x;
} SubscriberMessage;

typedef void (*EventHandler)(SubscriberMessage *message);

static inline void callEventHandler(EventHandler handler, SubscriberMessage *message) {
	handler(message);
}
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/BrobridgeOrg/gravity-sdk/core"
	subscriber "github.com/BrobridgeOrg/gravity-sdk/subscriber"
	pointer "github.com/mattn/go-pointer"
)

//export NewSubscriberOptions
func NewSubscriberOptions() *C.SubscriberOptions {

	opts := subscriber.NewOptions()

	// Convert to C struct
	sopts := (*C.SubscriberOptions)(C.malloc(C.size_t(unsafe.Sizeof(C.SubscriberOptions{}))))
	sopts.endpoint = C.CString(opts.Endpoint)
	sopts.domain = C.CString(opts.Domain)
	sopts.workerCount = C.int(opts.WorkerCount)
	sopts.bufferSize = C.int(opts.BufferSize)
	sopts.chunkSize = C.int(opts.ChunkSize)
	sopts.verbose = C.bool(opts.Verbose)
	sopts.initialLoad.enabled = C.bool(opts.InitialLoad.Enabled)
	sopts.initialLoad.mode = C.CString(opts.InitialLoad.Mode)
	sopts.initialLoad.omittedCount = C.ulonglong(opts.InitialLoad.OmittedCount)

	return sopts
}

//export NewSubscriber
func NewSubscriber(options *C.SubscriberOptions) unsafe.Pointer {

	opts := subscriber.NewOptions()
	opts.Endpoint = C.GoString(options.endpoint)
	opts.Domain = C.GoString(options.domain)
	opts.WorkerCount = int(options.workerCount)
	opts.BufferSize = int(options.bufferSize)
	opts.ChunkSize = int(options.chunkSize)
	opts.Verbose = bool(options.verbose)
	opts.InitialLoad.Enabled = bool(options.initialLoad.enabled)
	opts.InitialLoad.Mode = C.GoString(options.initialLoad.mode)
	opts.InitialLoad.OmittedCount = uint64(options.initialLoad.omittedCount)

	s := subscriber.NewSubscriber(opts)

	return pointer.Save(s)
}

//export SubscriberConnect
func SubscriberConnect(s unsafe.Pointer, host *C.char) *C.GravityError {

	sub := pointer.Restore(s).(*subscriber.Subscriber)

	opts := core.NewOptions()

	err := sub.Connect(C.GoString(host), opts)
	if err != nil {
		return NewError(err.Error())
	}

	return nil
}

//export SubscriberSetEventHandler
func SubscriberSetEventHandler(s unsafe.Pointer, callback C.EventHandler) {

	sub := pointer.Restore(s).(*subscriber.Subscriber)

	sub.SetEventHandler(func(msg *subscriber.Message) {
		message := (*C.SubscriberMessage)(C.malloc(C.size_t(unsafe.Sizeof(C.SubscriberMessage{}))))
		//		C.makeEventCallback(message, callback)
		C.callEventHandler(callback, message)
	})
}

//export SubscriberSetSnapshotHandler
func SubscriberSetSnapshotHandler(s unsafe.Pointer, callback C.EventHandler) {

	sub := pointer.Restore(s).(*subscriber.Subscriber)

	sub.SetSnapshotHandler(func(msg *subscriber.Message) {
		message := (*C.SubscriberMessage)(C.malloc(C.size_t(unsafe.Sizeof(C.SubscriberMessage{}))))
		//		C.makeSnapshotCallback(message, callback)
		//callback(message)
		C.callEventHandler(callback, message)
	})
}

//export SubscriberRegister
func SubscriberRegister(s unsafe.Pointer, subscriberType *C.char, componentName *C.char, subscriberID *C.char, subscriberName *C.char) *C.GravityError {

	sub := pointer.Restore(s).(*subscriber.Subscriber)

	subType := C.GoString(subscriberType)
	compName := C.GoString(componentName)
	subID := C.GoString(subscriberID)
	subName := C.GoString(subscriberName)

	t, ok := subscriber.SubscriberTypes[subType]
	if !ok {
		return NewError(fmt.Sprintf("No such subscriber type \"%s\"", subType))
	}

	err := sub.Register(t, compName, subID, subName)
	if err != nil {
		return NewError(err.Error())
	}

	return nil
}
