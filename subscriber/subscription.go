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

func DeliverNew() func(*Subscription) {
	return func(s *Subscription) {
		s.subOpts = append(s.subOpts, nats.DeliverNew())
	}
}

func StartSequence(seq uint64) func(*Subscription) {
	return func(s *Subscription) {
		s.startSequence = seq
		s.subOpts = append(s.subOpts, nats.StartSequence(seq))
	}
}

func InitialLoad(enabled bool) func(*Subscription) {
	return func(s *Subscription) {
		s.enabledInitialLoad = enabled
	}
}

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
