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

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductExistsAlready = errors.New("product exists already")
	ErrInvalidProductName   = errors.New("invalid product name")
)

type ProductSetting struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"desc"`
	Enabled         bool                   `json:"enabled"`
	Rules           map[string]*Rule       `json:"rules"`
	Schema          map[string]interface{} `json:"schema"`
	EnabledSnapshot bool                   `json:"enabledSnapshot"`
	Snapshot        *SnapshotSetting       `json:"snapshot"`
	Stream          string                 `json:"stream"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
}

type ProductState struct {
	EventCount uint64    `json:"eventCount"`
	Bytes      uint64    `json:"bytes"`
	FirstTime  time.Time `json:"firstTime"`
	LastTime   time.Time `json:"lastTime"`
}

type ProductInfo struct {
	Setting *ProductSetting `json:"setting"`
	State   *ProductState   `json:"state"`
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

func (pc *ProductClient) CreateProduct(productSetting *ProductSetting) (*ProductSetting, error) {

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
