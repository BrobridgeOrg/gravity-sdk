package subscription

import "github.com/BrobridgeOrg/gravity-sdk/v2/core"

const (
	SubscriptionAPI = "$GVT.%s.API.SUBSCRIPTION"
)

type ListSubscriptionsRequest struct {
	Product string `json:"product"`
}

type ListSubscriptionsReply struct {
	core.ErrorReply

	Subscriptions []*SubscriptionSetting `json:"subscriptions"`
}

type CreateSubscriptionRequest struct {
	Setting *SubscriptionSetting `json:"setting"`
}

type CreateSubscriptionReply struct {
	core.ErrorReply
	SubscriptionID string               `json:"subscriptionID"`
	Setting        *SubscriptionSetting `json:"setting"`
}

type UpdateSubscriptionRequest struct {
	SubscriptionID string               `json:"subscriptionID"`
	Setting        *SubscriptionSetting `json:"setting"`
}

type UpdateSubscriptionReply struct {
	core.ErrorReply
	Setting *SubscriptionSetting `json:"setting"`
}

type DeleteSubscriptionRequest struct {
	SubscriptionID string `json:"subscriptionID"`
}

type DeleteSubscriptionReply struct {
	core.ErrorReply
}

type GetSubscriptionRequest struct {
	SubscriptionID string `json:"subscriptionID"`
}

type GetSubscriptionReply struct {
	core.ErrorReply
	SubscriptionID string               `json:"subscriptionID"`
	Setting        *SubscriptionSetting `json:"setting"`
}

type PrepareSubscriptionRequest struct {
	SubscriptionID string `json:"subscriptionID"`
}

type PrepareSubscriptionReply struct {
	core.ErrorReply
}
