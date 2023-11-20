package main

import (
	"fmt"

	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
	"github.com/BrobridgeOrg/gravity-sdk/v2/product"
)

func main() {

	// Preparing options
	options := core.NewOptions()

	// Default configurations
	//options.Token = ""
	//options.PingInterval = 10 * time.Second
	//options.MaxPingsOutstanding = 3
	//options.MaxReconnects = -1

	client := core.NewClient()

	// Connect to Gravity
	err := client.Connect("0.0.0.0:32803", options)
	if err != nil {
		panic(err)
	}

	// Initializing data product client
	pcOpts := product.NewOptions()
	pcOpts.Domain = "default"

	productClient := product.NewProductClient(
		client,
		pcOpts,
	)

	// Get product list
	products, err := productClient.ListProducts()
	if err != nil {
		panic(err)
	}

	for i, p := range products {

		fmt.Println("Data Product", i+1)

		// Data product information
		fmt.Println("Name:", p.Setting.Name)
		fmt.Println("Description:", p.Setting.Description)
		fmt.Println("Is enabled:", p.Setting.Enabled)
		fmt.Println("Updated at:", p.Setting.UpdatedAt)
		fmt.Println("Created at:", p.Setting.CreatedAt)

		// Event information
		fmt.Println("Size:", p.State.Bytes)
		fmt.Println("Event count:", p.State.EventCount)
		fmt.Println("Last time:", p.State.LastTime)
	}
}
