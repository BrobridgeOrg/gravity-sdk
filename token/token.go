package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/BrobridgeOrg/gravity-sdk/v2/config_store"
	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
)

var (
	ErrTokenNotFound      = errors.New("product not found")
	ErrTokenExistsAlready = errors.New("product exists already")
	ErrInvalidTokenToken  = errors.New("invalid product token")
)

type TokenSetting struct {
	ID          string                 `json:"id"`
	Description string                 `json:"desc"`
	Enabled     bool                   `json:"enabled"`
	Permissions map[string]*Permission `json:"permissions"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

type TokenClient struct {
	options     *Options
	client      *core.Client
	configStore *config_store.ConfigStore
}

func NewTokenClient(client *core.Client, options *Options) *TokenClient {

	pc := &TokenClient{
		options: options,
		client:  client,
	}

	pc.configStore = config_store.NewConfigStore(client,
		config_store.WithDomain(options.Domain),
		config_store.WithCatalog("TOKEN"),
	)

	err := pc.configStore.Init()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return pc
}

func (pc *TokenClient) ListAvailablePermissions() (map[string]string, error) {

	// Preparing request
	req := &ListAvailablePermissionsRequest{}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(TokenAPI+".LIST_AVAILABLE_PERMISSIONS", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return nil, err
	}

	// Parsing response
	resp := &ListAvailablePermissionsReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	return resp.Permissions, nil
}

func (pc *TokenClient) CreateToken(tokenID string, productSetting *TokenSetting) (string, *TokenSetting, error) {

	// Preparing request
	req := &CreateTokenRequest{
		TokenID: tokenID,
		Setting: productSetting,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(TokenAPI+".CREATE", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return "", nil, err
	}

	// Parsing response
	resp := &CreateTokenReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return "", nil, err
	}

	if resp.Error != nil {
		return "", nil, errors.New(resp.Error.Message)
	}

	return resp.Token, resp.Setting, nil
}

func (pc *TokenClient) DeleteToken(tokenID string) error {

	// Preparing request
	req := &DeleteTokenRequest{
		TokenID: tokenID,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(TokenAPI+".DELETE", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return err
	}

	// Parsing response
	resp := &DeleteTokenReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	return nil
}

func (pc *TokenClient) UpdateToken(tokenID string, productSetting *TokenSetting) (*TokenSetting, error) {

	// Preparing request
	req := &UpdateTokenRequest{
		TokenID: tokenID,
		Setting: productSetting,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(TokenAPI+".UPDATE", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return nil, err
	}

	// Parsing response
	resp := &UpdateTokenReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	return resp.Setting, nil
}

func (pc *TokenClient) GetToken(tokenID string) (string, *TokenSetting, error) {

	// Preparing request
	req := &InfoTokenRequest{
		TokenID: tokenID,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(TokenAPI+".INFO", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return "", nil, err
	}

	// Parsing response
	resp := &InfoTokenReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return "", nil, err
	}

	if resp.Error != nil {
		return "", nil, errors.New(resp.Error.Message)
	}

	return resp.Token, resp.Setting, nil
}

func (pc *TokenClient) ListTokens() ([]*TokenSetting, error) {

	products := make([]*TokenSetting, 0)

	// Preparing request
	req := &ListTokensRequest{}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(TokenAPI+".LIST", pc.options.Domain)
	msg, err := pc.client.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return products, err
	}

	// Parsing response
	resp := &ListTokensReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return products, err
	}

	if resp.Error != nil {
		return products, errors.New(resp.Error.Message)
	}

	return resp.Tokens, nil
}
