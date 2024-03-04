package subscription

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/BrobridgeOrg/gravity-sdk/v2/config_store"
	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
)

type ConsumerSetting struct {
	Name         string `json:"name"`
	Partitions   []int  `json:"partitions"`
	StartFromSeq uint64 `json:"startFromSeq"`
}

type SubscriptionSetting struct {
	Product   string             `json:"product"`
	Consumers []*ConsumerSetting `json:"consumers"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

func NewSubscriptionSetting() *SubscriptionSetting {
	return &SubscriptionSetting{
		Consumers: make([]*ConsumerSetting, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type SubscriptionClient struct {
	options     *Options
	client      *core.Client
	configStore *config_store.ConfigStore
}

func NewSubscriptionClient(client *core.Client, options *Options) *SubscriptionClient {

	sc := &SubscriptionClient{
		options: options,
		client:  client,
	}

	sc.configStore = config_store.NewConfigStore(client,
		config_store.WithDomain(options.Domain),
		config_store.WithCatalog("SUBSCRIPTIONS"),
	)

	err := sc.configStore.Init()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return sc
}

func (sc *SubscriptionClient) ListSubscriptions() ([]*SubscriptionSetting, error) {

	products := make([]*SubscriptionSetting, 0)

	// Preparing request
	req := &ListSubscriptionsRequest{}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(SubscriptionAPI+".LIST", sc.options.Domain)
	msg, err := sc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return products, err
	}

	// Parsing response
	resp := &ListSubscriptionsReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return products, err
	}

	if resp.Error != nil {
		return products, errors.New(resp.Error.Message)
	}

	return resp.Subscriptions, nil
}
