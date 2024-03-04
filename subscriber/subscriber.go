package subscriber

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
	"github.com/BrobridgeOrg/gravity-sdk/v2/product"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Subscriber struct {
	client  *core.Client
	options *Options
	name    string
}

func NewSubscriber(name string, options *Options) *Subscriber {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	s := &Subscriber{
		options: options,
		name:    name,
	}

	// Auto-generate name for anonymous if name has not been set
	//if len(name) == 0 {
	//	s.name = generateAnonymousName()
	//}

	return s
}

func NewSubscriberWithClient(name string, client *core.Client, options *Options) *Subscriber {

	subscriber := NewSubscriber(name, options)
	subscriber.client = client

	return subscriber
}

func (s *Subscriber) GetName() string {
	return s.name
}

func (s *Subscriber) Connect(host string, options *core.Options) error {

	s.client = core.NewClient()
	return s.client.Connect(host, options)
}

func (s *Subscriber) Disconnect() error {
	return s.Disconnect()
}

// Subscribe sets up a subscription to a specified product on Gravity.
// This function allows for message consumption by registering a handler function to process incoming messages.
// productName: The name of the product to subscribe to.
// handler: A function that will be called to handle each incoming message. It takes a *nats.Msg as a parameter.
// opts: A variadic list of options to customize the subscription behavior.
// Returns a pointer to the Subscription object and an error if the subscription fails.
func (s *Subscriber) Subscribe(productName string, handler func(*nats.Msg), opts ...SubscriptionOpt) (*Subscription, error) {

	subscription := NewSubscription(s, s.options.Domain, productName, handler, opts...)
	err := subscription.Subscribe()
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// Reset removes the subscription to a specified product on Gravity.
// productName: The name of the product to reset the subscription.
// Returns an error if the reset fails.
func (s *Subscriber) Reset(productName string) error {

	// Preparing request
	req := &product.DeleteSubscriptionRequest{
		Product: productName,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(product.ProductAPI+".DELETE_SUBSCRIPTION", s.options.Domain)
	msg, err := s.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return err
	}

	// Parsing response
	resp := &product.DeleteSubscriptionReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	return nil
}
