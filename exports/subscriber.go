package main

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>
#include "./error.h"
#include "./client.h"

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
	char *appID;
	char *accessKey;
	bool verbose;
	SubscriberInitialLoadOptions initialLoad;
} SubscriberOptions;

typedef struct {
	void *instance;
	uint64_t pipelineId;
	char *eventName;
	char *collection;
	char *payload;
} SubscriberMessage;

typedef void (*EventHandler)(SubscriberMessage *message);

inline void callEventHandler(EventHandler handler, SubscriberMessage *message) {
	handler(message);
}

typedef struct {
	void *instance;
	EventHandler eventHandler;
	EventHandler snapshotHandler;
} Subscriber;
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/BrobridgeOrg/gravity-sdk/core/keyring"
	subscriber "github.com/BrobridgeOrg/gravity-sdk/subscriber"
	gravity_sdk_types_record "github.com/BrobridgeOrg/gravity-sdk/types/record"
	jsoniter "github.com/json-iterator/go"
	pointer "github.com/mattn/go-pointer"
)

func getSubscriberNativeOptions(options *C.SubscriberOptions) *subscriber.Options {

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

	var appID string
	if options.appID == nil {
		appID = "anonymous"
	} else {
		appID = C.GoString(options.appID)
	}

	var accessKey string = ""
	if options.accessKey != nil {
		accessKey = C.GoString(options.accessKey)
	}

	opts.Key = keyring.NewKey(appID, accessKey)

	return opts
}

func initalizeHandlers(s *C.Subscriber, sub *subscriber.Subscriber) {

	// Initializing event handlers
	sub.SetEventHandler(func(msg *subscriber.Message) {

		if s.eventHandler == nil {
			return
		}

		message := (*C.SubscriberMessage)(C.malloc(C.size_t(unsafe.Sizeof(C.SubscriberMessage{}))))
		message.instance = pointer.Save(msg)
		message.pipelineId = C.ulonglong(msg.Pipeline.GetID())

		// Getting record information
		event := msg.Payload.(*subscriber.DataEvent)
		record := event.Payload

		message.collection = C.CString(record.Table)
		message.eventName = C.CString(record.EventName)

		// Getting payload
		fields := record.GetFields()
		mapData := gravity_sdk_types_record.ConvertFieldsToMap(fields)
		payload, _ := jsoniter.Marshal(mapData)
		message.payload = C.CString(string(payload))

		C.callEventHandler(s.eventHandler, message)
		/*
			// Release
			C.free(unsafe.Pointer(message.collection))
			C.free(unsafe.Pointer(message.eventName))
			C.free(unsafe.Pointer(message.payload))
			C.free(unsafe.Pointer(message))
		*/
	})

	sub.SetSnapshotHandler(func(msg *subscriber.Message) {

		if s.snapshotHandler == nil {
			return
		}

		message := (*C.SubscriberMessage)(C.malloc(C.size_t(unsafe.Sizeof(C.SubscriberMessage{}))))
		message.instance = pointer.Save(msg)
		message.pipelineId = C.ulonglong(msg.Pipeline.GetID())

		// Getting record information
		event := msg.Payload.(*subscriber.SnapshotEvent)
		record := event.Payload

		message.collection = C.CString(event.Collection)
		message.eventName = C.CString("snapshot")

		// Getting payload
		fields := record.Payload.Map.Fields
		mapData := gravity_sdk_types_record.ConvertFieldsToMap(fields)
		payload, _ := jsoniter.Marshal(mapData)
		message.payload = C.CString(string(payload))

		C.callEventHandler(s.snapshotHandler, message)

		/*
			// Release
			C.free(unsafe.Pointer(message.collection))
			C.free(unsafe.Pointer(message.payload))
			C.free(unsafe.Pointer(message))
		*/
	})
}

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
	sopts.appID = C.CString("anonymous")
	sopts.accessKey = nil
	sopts.verbose = C.bool(opts.Verbose)
	sopts.initialLoad.enabled = C.bool(opts.InitialLoad.Enabled)
	sopts.initialLoad.mode = C.CString(opts.InitialLoad.Mode)
	sopts.initialLoad.omittedCount = C.ulonglong(opts.InitialLoad.OmittedCount)

	return sopts
}

//export NewSubscriber
func NewSubscriber(options *C.SubscriberOptions) *C.Subscriber {

	opts := getSubscriberNativeOptions(options)

	s := subscriber.NewSubscriber(opts)

	// Prepare subscriber struct
	sub := (*C.Subscriber)(C.malloc(C.size_t(unsafe.Sizeof(C.Subscriber{}))))
	sub.instance = pointer.Save(s)
	sub.eventHandler = nil
	sub.snapshotHandler = nil

	initalizeHandlers(sub, s)

	return sub
}

//export NewSubscriberWithClient
func NewSubscriberWithClient(client *C.Client, options *C.SubscriberOptions) *C.Subscriber {

	c := pointer.Restore(client.instance).(*core.Client)
	opts := getSubscriberNativeOptions(options)

	s := subscriber.NewSubscriberWithClient(c, opts)

	// Prepare subscriber struct
	sub := (*C.Subscriber)(C.malloc(C.size_t(unsafe.Sizeof(C.Subscriber{}))))
	sub.instance = pointer.Save(s)
	sub.eventHandler = nil
	sub.snapshotHandler = nil

	initalizeHandlers(sub, s)

	return sub
}

//export SubscriberConnect
func SubscriberConnect(s *C.Subscriber, host *C.char) *C.GravityError {

	sub := pointer.Restore(s.instance).(*subscriber.Subscriber)

	opts := core.NewOptions()

	err := sub.Connect(C.GoString(host), opts)
	if err != nil {
		return NewError(err.Error())
	}

	return nil
}

//export SubscriberSetEventHandler
func SubscriberSetEventHandler(s *C.Subscriber, callback C.EventHandler) {
	s.eventHandler = callback
}

//export SubscriberSetSnapshotHandler
func SubscriberSetSnapshotHandler(s *C.Subscriber, callback C.EventHandler) {
	s.snapshotHandler = callback
}

//export SubscriberRegister
func SubscriberRegister(s *C.Subscriber, subscriberType *C.char, componentName *C.char, subscriberID *C.char, subscriberName *C.char) *C.GravityError {

	sub := pointer.Restore(s.instance).(*subscriber.Subscriber)

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

//export SubscriberGetPipelineCount
func SubscriberGetPipelineCount(s *C.Subscriber) C.int {

	sub := pointer.Restore(s.instance).(*subscriber.Subscriber)

	count, err := sub.GetPipelineCount()
	if err != nil {
		return C.int(-1)
	}

	return C.int(count)
}

//export SubscriberAddAllPipelines
func SubscriberAddAllPipelines(s *C.Subscriber) *C.GravityError {

	sub := pointer.Restore(s.instance).(*subscriber.Subscriber)

	err := sub.AddAllPipelines()
	if err != nil {
		return NewError(err.Error())
	}

	return nil
}

//export SubscriberSubscribeToPipelines
func SubscriberSubscribeToPipelines(s *C.Subscriber, pipelines []C.ulonglong) *C.GravityError {

	sub := pointer.Restore(s.instance).(*subscriber.Subscriber)

	for len(pipelines) == 0 {
		return nil
	}

	list := make([]uint64, 0, len(pipelines))
	for i, p := range pipelines {
		list[i] = uint64(p)
	}

	err := sub.SubscribeToPipelines(list)
	if err != nil {
		return NewError(err.Error())
	}

	return nil
}

//export SubscriberSubscribeToCollection
func SubscriberSubscribeToCollection(s *C.Subscriber, collection *C.char, tables **C.char, length C.int) *C.GravityError {

	sub := pointer.Restore(s.instance).(*subscriber.Subscriber)

	// We need to do some pointer arithmetic.
	start := unsafe.Pointer(tables)
	pointerSize := unsafe.Sizeof(tables)

	inSlice := make([]string, int(length))
	for i := 0; i < int(length); i++ {
		// Copy each input string into a Go string and add it to the slice.
		pointer := (**C.char)(unsafe.Pointer(uintptr(start) + uintptr(i)*pointerSize))
		inSlice[i] = C.GoString(*pointer)
	}

	err := sub.SubscribeToCollections(map[string][]string{
		C.GoString(collection): inSlice,
	})
	if err != nil {
		return NewError(err.Error())
	}

	return nil

}

//export SubscriberStart
func SubscriberStart(s *C.Subscriber) {

	sub := pointer.Restore(s.instance).(*subscriber.Subscriber)

	sub.Start()
}

//export SubscriberStop
func SubscriberStop(s *C.Subscriber) {

	sub := pointer.Restore(s.instance).(*subscriber.Subscriber)

	sub.Stop()
}

//export SubscriberMessageAck
func SubscriberMessageAck(m *C.SubscriberMessage) {

	msg := pointer.Restore(m.instance).(*subscriber.Message)
	msg.Ack()

	// Release
	C.free(unsafe.Pointer(m.collection))
	C.free(unsafe.Pointer(m.eventName))
	C.free(unsafe.Pointer(m.payload))
	C.free(unsafe.Pointer(m))
}
