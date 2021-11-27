/* Code generated by cmd/cgo; DO NOT EDIT. */

/* package command-line-arguments */


#line 1 "cgo-builtin-export-prolog"

#include <stddef.h> /* for ptrdiff_t below */

#ifndef GO_CGO_EXPORT_PROLOGUE_H
#define GO_CGO_EXPORT_PROLOGUE_H

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
#endif

#endif

/* Start of preamble from import "C" comments.  */


#line 3 "adapter.go"

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

#line 1 "cgo-generated-wrapper"

#line 3 "client.go"

#include <stdint.h>
#include "./error.h"
#include "./client.h"

#line 1 "cgo-generated-wrapper"

#line 3 "subscriber.go"

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

#line 1 "cgo-generated-wrapper"


/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef _GoString_ GoString;
#endif
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif

extern AdapterOptions* NewAdapterOptions();
extern Adapter* NewAdapterWithClient(Client* client, AdapterOptions* options);
extern GravityError* AdapterRegister(Adapter* s, char* componentName, char* subscriberID, char* subscriberName);
extern GravityError* AdapterPublish(Adapter* ca, char* eventName, char* payload);
extern ClientOptions* NewClientOptions();
extern Client* NewClient();
extern GravityError* ClientConnect(Client* c, char* host, ClientOptions* options);
extern void ClientSetDisconnectHandler(Client* c, ClientEventHandler callback);
extern void ClientSetReconnectHandler(Client* c, ClientEventHandler callback);
extern SubscriberOptions* NewSubscriberOptions();
extern Subscriber* NewSubscriber(SubscriberOptions* options);
extern Subscriber* NewSubscriberWithClient(Client* client, SubscriberOptions* options);
extern GravityError* SubscriberConnect(Subscriber* s, char* host);
extern void SubscriberSetEventHandler(Subscriber* s, EventHandler callback);
extern void SubscriberSetSnapshotHandler(Subscriber* s, EventHandler callback);
extern GravityError* SubscriberRegister(Subscriber* s, char* subscriberType, char* componentName, char* subscriberID, char* subscriberName);
extern int SubscriberGetPipelineCount(Subscriber* s);
extern GravityError* SubscriberAddAllPipelines(Subscriber* s);
extern GravityError* SubscriberSubscribeToPipelines(Subscriber* s, GoSlice pipelines);
extern GravityError* SubscriberSubscribeToCollection(Subscriber* s, char* collection, char** tables, int length);
extern void SubscriberStart(Subscriber* s);
extern void SubscriberMessageAck(SubscriberMessage* m);

#ifdef __cplusplus
}
#endif