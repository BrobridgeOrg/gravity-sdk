//go:build !ci
// +build !ci

package product

import (
	"fmt"

	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
)

// ExampleProductClient_CreateProduct demonstrates the creation of a new data product using the Gravity SDK.
// This example first establishes a connection to the Gravity service using the client and options.
// It then initializes a product client with a specified domain.
// After setting up the client, the example creates a new data product with specified settings, such as name and description.
// Finally, it prints the name of the newly created product.
// Note: Replace '0.0.0.0:32803' with your actual Gravity service address.
func ExampleProductClient_CreateProduct() {

	client := core.NewClient()

	// Connect to Gravity
	options := core.NewOptions()
	err := client.Connect("0.0.0.0:32803", options)
	if err != nil {
		panic(err)
	}

	// Initializing data product client
	pcOpts := NewOptions()
	pcOpts.Domain = "default"
	productClient := NewProductClient(client, pcOpts)

	// Create a new data product
	ps, err := productClient.CreateProduct(&ProductSetting{
		Name:        "sdk_example",
		Description: "SDK Example",
		Enabled:     false,
		Stream:      "",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(ps.Name)

	// Output:
	// sdk_example
}

// ExampleProductClient_GetProduct demonstrates how to retrieve information about a specific data product using the Gravity SDK.
// This example begins by establishing a connection to the Gravity service using a client and options.
// A product client is then initialized with a domain specification.
// The example proceeds to retrieve information for a data product named "sdk_example".
// Upon successful retrieval, it prints the name of the product.
// Note: Replace '0.0.0.0:32803' with your actual Gravity service address.
// Ensure that the product "sdk_example" exists before running this example.
func ExampleProductClient_GetProduct() {

	client := core.NewClient()

	// Connect to Gravity
	options := core.NewOptions()
	err := client.Connect("0.0.0.0:32803", options)
	if err != nil {
		panic(err)
	}

	// Initializing data product client
	pcOpts := NewOptions()
	pcOpts.Domain = "default"
	productClient := NewProductClient(client, pcOpts)

	// Get product information
	info, err := productClient.GetProduct("sdk_example")
	if err != nil {
		panic(err)
	}

	fmt.Println(info.Setting.Name)

	// Output:
	// sdk_example
}

// ExampleProductClient_ListProducts demonstrates how to list all products using the Gravity SDK.
// This example initializes a Gravity SDK client, creates a product client, and then lists all products.
// It finally prints the name and description of each product.
func ExampleProductClient_ListProducts() {

	client := core.NewClient()

	// Connect to Gravity
	options := core.NewOptions()
	err := client.Connect("0.0.0.0:32803", options)
	if err != nil {
		panic(err)
	}

	// Initializing data product client
	pcOpts := NewOptions()
	pcOpts.Domain = "default"
	productClient := NewProductClient(client, pcOpts)

	// Get product list
	products, err := productClient.ListProducts()
	if err != nil {
		panic(err)
	}

	for _, p := range products {

		if p.Setting.Name != "sdk_example" {
			continue
		}

		// Data product information
		fmt.Println("Name:", p.Setting.Name)
		fmt.Println("Is enabled:", p.Setting.Enabled)
	}

	// Output:
	// Name: sdk_example
	// Is enabled: false
}

// ExampleProductClient_DeleteProduct demonstrates how to delete a data product using the Gravity SDK.
// This example first establishes a connection to the Gravity service using the client and options.
// It then initializes a product client for the specified domain.
// The example proceeds to delete a data product named "sdk_example".
// Finally, it prints a confirmation message upon successful deletion.
// Note: Replace '0.0.0.0:32803' with your actual Gravity service address.
// Ensure that the product "sdk_example" exists before running this example.
func ExampleProductClient_DeleteProduct() {

	client := core.NewClient()

	// Connect to Gravity
	options := core.NewOptions()
	err := client.Connect("0.0.0.0:32803", options)
	if err != nil {
		panic(err)
	}

	// Initializing data product client
	pcOpts := NewOptions()
	pcOpts.Domain = "default"
	productClient := NewProductClient(client, pcOpts)

	err = productClient.DeleteProduct("sdk_example")
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted")

	// Output:
	// Deleted
}
