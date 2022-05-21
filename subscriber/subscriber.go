package subscriber

import (
	"os"

	"github.com/BrobridgeOrg/gravity-sdk/core"
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

func (s *Subscriber) Subscribe(productName string, handler func(*nats.Msg), opts ...SubscriptionOpt) (*Subscription, error) {

	subscription := NewSubscription(s, s.options.Domain, productName, handler, opts...)
	err := subscription.Subscribe()
	if err != nil {
		return nil, err
	}

	return subscription, nil
}
