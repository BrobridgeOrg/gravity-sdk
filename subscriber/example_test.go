//go:build !ci
// +build !ci

package subscriber

import (
	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
	nats "github.com/nats-io/nats.go"
)

func ExampleSubscriber_Subscribe() {

	client := core.NewClient()

	// Connect to Gravity
	options := core.NewOptions()
	err := client.Connect("0.0.0.0:32803", options)
	if err != nil {
		panic(err)
	}

	// Create adapter connector
	acOpts := NewOptions()
	s := NewSubscriberWithClient("sdk_example_subscriber", client, acOpts)

	// Subscribe to specific data product
	sub, err := s.Subscribe("sdk_example", func(msg *nats.Msg) {
		// TODO: event handling
	}, Partition(-1), StartSequence(0))
	if err != nil {
		panic(err)
	}

	defer sub.Close()

	// Output:
}
