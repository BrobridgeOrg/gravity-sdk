package adapter

import (
	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
)

// ExampleAdapterConnector_Publish demonstrates how to publish an event to the Gravity service using the AdapterConnector.
// This example starts by establishing a connection to the Gravity service using the core client.
// It then creates an AdapterConnector with the client and options.
// The example proceeds to publish an event with a payload and metadata to a specific event name.
// Note: Replace '0.0.0.0:32803' with your actual Gravity service address.
// No output is produced as the result of the publishing operation is not printed.
func ExampleAdapterConnector_Publish() {

	client := core.NewClient()

	// Connect to Gravity
	options := core.NewOptions()
	err := client.Connect("0.0.0.0:32803", options)
	if err != nil {
		panic(err)
	}

	// Create adapter connector
	acOpts := NewOptions()
	ac := NewAdapterConnectorWithClient(client, acOpts)

	// Publish event
	meta := map[string]string{
		"example": "example",
	}
	_, err = ac.Publish("sdk_adapter_event", []byte("sdk_payload"), meta)
	if err != nil {
		panic(err)
	}

	// Output:
}
