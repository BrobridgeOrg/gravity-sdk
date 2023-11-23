package subscriber

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const productSubject = "$GVT.%s.DP.%s.>"
const productEventSubject = "$GVT.%s.DP.%s.%s.EVENT.>"
const productEventConsumer = "GVT_%s_SUBSCRIBER_%s"

type Subscription struct {
	subscriber          *Subscriber
	handler             func(*nats.Msg)
	domain              string
	productName         string
	startSequence       uint64
	enabledInitialLoad  bool
	partitions          []int
	nativeSubscriptions map[string]*nats.Subscription
	subOpts             []nats.SubOpt
}

// SubscriptionOpt is a type of function that modifies a Subscription.
type SubscriptionOpt func(*Subscription)

func NewSubscription(s *Subscriber, domain string, productName string, handler func(*nats.Msg), opts ...SubscriptionOpt) *Subscription {

	sub := &Subscription{
		subscriber:          s,
		handler:             handler,
		domain:              domain,
		productName:         productName,
		startSequence:       1,
		partitions:          []int{-1},
		nativeSubscriptions: make(map[string]*nats.Subscription),
		subOpts:             make([]nats.SubOpt, 0),
	}

	for _, opt := range opts {
		opt(sub)
	}

	return sub
}

// DeliverNew configures the Subscriber to begin receiving events from the most recently produced events.
// This option is used when the Subscriber needs to start receiving the latest data, ignoring previous events.
func DeliverNew() func(*Subscription) {
	return func(s *Subscription) {
		s.subOpts = append(s.subOpts, nats.DeliverNew())
	}
}

// StartSequence sets the starting sequence number from which the Subscriber begins to receive events.
// This option is used to specify a particular point in the event sequence to start receiving messages from.
// seq: The sequence number from which to start receiving events.
func StartSequence(seq uint64) func(*Subscription) {
	return func(s *Subscription) {
		s.startSequence = seq
		s.subOpts = append(s.subOpts, nats.StartSequence(seq))
	}
}

// InitialLoad determines whether to receive an initial copy of all existing data when first subscribing and interfacing with a Data Product.
// This option allows the subscriber to get a snapshot of all existing data before continuing to receive real-time data change events.
// enabled: Set to true to enable receiving the initial data load.
func InitialLoad(enabled bool) func(*Subscription) {
	return func(s *Subscription) {
		s.enabledInitialLoad = enabled
	}
}

// Partition specifies the particular partitions of a data product to subscribe to.
// This option is used to increase parallel computing capabilities by subscribing only to specific partitions of the data divided into 256 parts by Gravity.
// partitions: A list of partition indices to subscribe to.
func Partition(partitions ...int) func(*Subscription) {
	return func(s *Subscription) {

		if len(partitions) == 0 {
			return
		}

		for _, p := range partitions {
			if p == -1 {
				s.partitions = []int{-1}
				return
			}
		}

		s.partitions = partitions
	}
}

func (sub *Subscription) subscribe(subject string) error {

	js, err := sub.subscriber.client.GetJetStream()
	if err != nil {
		return err
	}

	var opts []nats.SubOpt

	opts = append(sub.subOpts, nats.AckAll(), nats.MaxAckPending(20480))

	// Specific consumer if subscriber name has been set
	if len(sub.subscriber.name) > 0 {
		consumerName := fmt.Sprintf(productEventConsumer, sub.domain, sub.subscriber.name)
		opts = append(opts, nats.Durable(consumerName))
	}

	s, err := js.Subscribe(subject, func(msg *nats.Msg) {
		sub.handler(msg)
	}, opts...)
	if err != nil {
		return err
	}

	s.SetPendingLimits(-1, -1)
	sub.subscriber.client.GetConnection().Flush()

	sub.nativeSubscriptions[subject] = s

	return nil
}

// Subscribe initiates the subscription process for the specified data product in the Subscription object.
// This function starts the reception of messages based on the configuration set in the Subscription.
// It connects to the Gravity service and begins handling incoming messages using the handler function defined earlier.
// Returns an error if the subscription process fails or if the Subscription is not properly configured.
func (sub *Subscription) Subscribe() error {

	// Subscribe to multiple partitions
	for _, p := range sub.partitions {

		var partition string
		if p == -1 {
			// All partitions
			partition = "*"
		} else {
			// Specific parition
			partition = fmt.Sprintf("%d", p)
		}

		subject := fmt.Sprintf(productEventSubject, sub.domain, sub.productName, partition)

		log.WithFields(logrus.Fields{
			"subject": subject,
		}).Info("Subscribing to subject")

		err := sub.subscribe(subject)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close terminates the subscription and closes the connection associated with it.
// This function should be called to cleanly shutdown the subscription, ensuring that all resources are released properly.
// Returns an error if the closing process encounters any issues.
func (sub *Subscription) Close() error {

	// Unsubscribe all partitions
	for _, s := range sub.nativeSubscriptions {
		err := s.Unsubscribe()
		if err != nil {
			return err
		}
	}

	return nil
}
