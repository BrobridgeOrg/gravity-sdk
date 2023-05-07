package config_store

import (
	"fmt"
	"time"

	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
	"github.com/nats-io/nats.go"
)

const (
	configBucket = "GVT_%s_%s"
)

type ConfigOp int32

const (
	ConfigCreate ConfigOp = iota
	ConfigUpdate
	ConfigDelete
)

var configOps = map[ConfigOp]string{
	ConfigCreate: "Create",
	ConfigUpdate: "Update",
	ConfigDelete: "Delete",
}

func (co ConfigOp) String() string {
	return configOps[co]
}

type ConfigEntry struct {
	Operation ConfigOp
	Key       string
	Value     []byte
	Revision  uint64
	Created   time.Time
	Delta     uint64
}

type ConfigStore struct {
	client       *core.Client
	domain       string
	catalog      string
	ttl          time.Duration
	watcher      nats.KeyWatcher
	eventHandler func(*ConfigEntry)
	kv           nats.KeyValue
}

func NewConfigStore(client *core.Client, opts ...func(*ConfigStore)) *ConfigStore {

	cs := &ConfigStore{
		client: client,
	}

	for _, opt := range opts {
		opt(cs)
	}

	return cs
}

func WithDomain(domain string) func(*ConfigStore) {
	return func(cs *ConfigStore) {
		cs.domain = domain
	}
}

func WithCatalog(catalog string) func(*ConfigStore) {
	return func(cs *ConfigStore) {
		cs.catalog = catalog
	}
}

func WithTTL(ttl time.Duration) func(*ConfigStore) {
	return func(cs *ConfigStore) {
		cs.ttl = ttl
	}
}

func WithEventHandler(fn func(*ConfigEntry)) func(*ConfigStore) {
	return func(cs *ConfigStore) {
		cs.eventHandler = fn
	}
}

func (cs *ConfigStore) Init() error {

	// Preparing JetStream
	js, err := cs.client.GetJetStream()
	if err != nil {
		return err
	}

	bucket := fmt.Sprintf(configBucket, cs.domain, cs.catalog)

	// Attempt to create KV store
	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:      bucket,
		Description: "Gravity configuration store",
		TTL:         cs.ttl,
	})
	if err != nil {
		return err
	}

	// Load configuration from KV store
	kv, err = js.KeyValue(bucket)
	if err != nil {
		return err
	}

	cs.kv = kv

	// Do not load and watch if not event handler has been set
	if cs.eventHandler == nil {
		return nil
	}
	/*
		// Getting all entries
		keys, err := kv.Keys()
		if len(keys) > 0 {

			entries := make([]nats.KeyValueEntry, len(keys))
			for i, key := range keys {

				entry, err := kv.Get(key)
				if err != nil {
					return err
				}

				entries[i] = entry
			}

			for _, entry := range entries {
				cs.eventHandler(ConfigCreate, entry.Key(), entry.Value())
			}
		}
	*/

	if cs.eventHandler == nil {
		return nil
	}

	// Watching event of KV Store for real-time updating
	watcher, err := kv.WatchAll()
	if err != nil {
		return err
	}

	cs.watcher = watcher

	go func() {

		for entry := range watcher.Updates() {

			if entry == nil {
				continue
			}

			var op ConfigOp
			switch entry.Operation() {
			case nats.KeyValuePut:
				op = ConfigUpdate
			case nats.KeyValueDelete:
				op = ConfigDelete
			case nats.KeyValuePurge:
				op = ConfigDelete
			}

			ce := &ConfigEntry{
				Operation: op,
				Key:       entry.Key(),
				Value:     entry.Value(),
				Revision:  entry.Revision(),
				Created:   entry.Created(),
				Delta:     entry.Delta(),
			}

			cs.eventHandler(ce)
		}
	}()

	return nil
}

func (cs *ConfigStore) Put(key string, value []byte) (uint64, error) {
	return cs.kv.Put(key, value)
}

func (cs *ConfigStore) Update(key string, value []byte, revision uint64) (uint64, error) {
	return cs.kv.Update(key, value, revision)
}

func (cs *ConfigStore) Get(key string) (nats.KeyValueEntry, error) {
	return cs.kv.Get(key)
}

func (cs *ConfigStore) Delete(key string) error {
	return cs.kv.Delete(key)
}

func (cs *ConfigStore) Keys() ([]string, error) {
	return cs.kv.Keys()
}
