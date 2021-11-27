#ifndef CLIENT_H
#define CLIENT_H

typedef void (*ClientEventHandler)();

static inline void callClientEventHandler(ClientEventHandler handler) {
	handler();
}

typedef struct {
	char *appID;
	char *appKey;
	int64_t pingInterval;
	int maxPingsOutstanding;
	int maxReconnects;
} ClientOptions;

typedef struct {
	void *instance;
	ClientEventHandler disconnectHandler;
	ClientEventHandler reconnectHandler;
} Client;

#endif
