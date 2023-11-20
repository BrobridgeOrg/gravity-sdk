package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/BrobridgeOrg/gravity-sdk/v2/config_store"
	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
	"github.com/nats-io/nats.go"
)

const (
	ProductEventStream = "GVT_%s_DP_%s"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductExistsAlready = errors.New("product exists already")
	ErrInvalidProductName   = errors.New("invalid product name")
)

// ProductSetting defines the structure of a Data Product setting.
// It includes information such as name, description, Schema, and event handling rules.
type ProductSetting struct {
	Name            string                 `json:"name"`            // A unique identifier for the product.
	Description     string                 `json:"desc"`            // A brief description of the product.
	Enabled         bool                   `json:"enabled"`         // Flag indicating whether the product is active or not.
	Rules           map[string]*Rule       `json:"rules"`           // A map of event handling rules associated with the product.
	Schema          map[string]interface{} `json:"schema"`          // The data schema defining the structure of the product's data.
	EnabledSnapshot bool                   `json:"enabledSnapshot"` // Flag indicating whether snapshots are enabled for the product.
	Snapshot        *SnapshotSetting       `json:"snapshot"`        // Configuration settings for product snapshots.
	Stream          string                 `json:"stream"`          // The name of the data stream associated with the product.
	CreatedAt       time.Time              `json:"createdAt"`       // Timestamp indicating when the product was created.
	UpdatedAt       time.Time              `json:"updatedAt"`       // Timestamp indicating the last update to the product.
}

// ProductState represents the current state of a data product.
type ProductState struct {
	EventCount uint64    `json:"eventCount"` // The total number of events processed by the product.
	Bytes      uint64    `json:"bytes"`      // The total number of bytes processed by the product.
	FirstTime  time.Time `json:"firstTime"`  // Timestamp of the first event processed by the product.
	LastTime   time.Time `json:"lastTime"`   // Timestamp of the most recent event processed by the product.
}

// ProductInfo encapsulates both the settings and the current state of a data product.
type ProductInfo struct {
	Setting *ProductSetting `json:"setting"` // The configuration settings of the product.
	State   *ProductState   `json:"state"`   // The current operational state of the product.
}

type ProductClient struct {
	options     *Options
	client      *core.Client
	configStore *config_store.ConfigStore
}

func NewProductClient(client *core.Client, options *Options) *ProductClient {

	pc := &ProductClient{
		options: options,
		client:  client,
	}

	pc.configStore = config_store.NewConfigStore(client,
		config_store.WithDomain(options.Domain),
		config_store.WithCatalog("PRODUCT"),
	)

	err := pc.configStore.Init()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return pc
}

// CreateProduct creates a new product with the specified settings.
// It returns a pointer to the newly created ProductSetting or an error if the creation fails.
func (pc *ProductClient) CreateProduct(productSetting *ProductSetting) (*ProductSetting, error) {

	if len(productSetting.Name) == 0 {
		return nil, ErrInvalidProductName
	}

	if len(productSetting.Stream) == 0 {
		productSetting.Stream = fmt.Sprintf(ProductEventStream, pc.options.Domain, productSetting.Name)
	}

	// Preparing request
	req := &CreateProductRequest{
		Setting: productSetting,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(ProductAPI+".CREATE", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return nil, err
	}

	// Parsing response
	resp := &CreateProductReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	return resp.Setting, nil
}

// DeleteProduct removes the product identified by the given name.
// It returns an error if the deletion fails.
func (pc *ProductClient) DeleteProduct(name string) error {

	// Preparing request
	req := &DeleteProductRequest{
		Name: name,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(ProductAPI+".DELETE", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return err
	}

	// Parsing response
	resp := &DeleteProductReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	return nil
}

// UpdateProduct updates the settings of an existing product identified by name.
// It returns the updated ProductSetting or an error if the update fails.
func (pc *ProductClient) UpdateProduct(name string, productSetting *ProductSetting) (*ProductSetting, error) {

	// Preparing request
	req := &UpdateProductRequest{
		Name:    name,
		Setting: productSetting,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(ProductAPI+".UPDATE", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return nil, err
	}

	// Parsing response
	resp := &UpdateProductReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	return resp.Setting, nil
}

// PurgeProduct purge a data product without deleting it.
// It returns an error if the purge operation fails.
func (pc *ProductClient) PurgeProduct(name string) error {

	// Preparing request
	req := &PurgeProductRequest{
		Name: name,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(ProductAPI+".PURGE", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return err
	}

	// Parsing response
	resp := &PurgeProductReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	return nil
}

// GetProduct retrieves the information of a product identified by the given name.
// It returns a pointer to ProductInfo or an error if the retrieval fails.
func (pc *ProductClient) GetProduct(name string) (*ProductInfo, error) {

	// Preparing request
	req := &InfoProductRequest{
		Name: name,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(ProductAPI+".INFO", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return nil, err
	}

	// Parsing response
	resp := &InfoProductReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	return &ProductInfo{
		Setting: resp.Setting,
		State:   resp.State,
	}, nil
}

// ListProducts returns a list of all available products.
// It returns a slice of ProductInfo pointers or an error if the listing fails.
func (pc *ProductClient) ListProducts() ([]*ProductInfo, error) {

	products := make([]*ProductInfo, 0)

	// Preparing request
	req := &ListProductsRequest{}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(ProductAPI+".LIST", pc.options.Domain)
	reqMsg := nats.NewMsg(apiPath)
	reqMsg.Data = reqData

	msg, err := pc.client.RequestMsg(reqMsg, time.Second*30)
	if err != nil {
		return products, err
	}

	// Parsing response
	resp := &ListProductsReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return products, err
	}

	if resp.Error != nil {
		return products, errors.New(resp.Error.Message)
	}

	return resp.Products, nil
}

// CreateSnapshot creates a snapshot of a product identified by productName.
// Additional options can be passed using SnapshotOpt variadic parameters.
// It returns a pointer to the created Snapshot or an error if the creation fails.
func (pc *ProductClient) CreateSnapshot(productName string, opts ...SnapshotOpt) (*Snapshot, error) {

	// Check whether product exists or not
	_, err := pc.configStore.Get(productName)
	if err == nats.ErrKeyNotFound {
		return nil, ErrProductNotFound
	}

	s, err := NewSnapshot(pc, productName, opts...)
	if err != nil {
		return nil, err
	}

	return s, nil

}
