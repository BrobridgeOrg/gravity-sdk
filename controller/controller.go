package controller

import (
	core "github.com/BrobridgeOrg/gravity-sdk/core"
)

type Controller struct {
	client  *core.Client
	host    string
	options *Options
}

func NewController(options *Options) *Controller {
	return &Controller{
		options: options,
	}
}

func NewControllerWithClient(client *core.Client, options *Options) *Controller {
	return &Controller{
		client:  client,
		options: options,
	}
}

func (sub *Controller) Connect(host string, options *core.Options) error {
	sub.client = core.NewClient()
	return sub.client.Connect(host, options)
}

func (sub *Controller) Disconnect() {
	sub.client.Disconnect()
}
