#ifndef CLIENT_H
#define CLIENT_H

typedef void (*ClientEventHandler)();

inline void callClientEventHandler(ClientEventHandler handler) {
	handler();
}

typedef struct {
	char *appID;
	char *appKey;
	long long int pingInterval;
	int maxPingsOutstanding;
	int maxReconnects;
} ClientOptions;

typedef struct {
	void *instance;
	ClientEventHandler disconnectHandler;
	ClientEventHandler reconnectHandler;
} Client;

#endif
