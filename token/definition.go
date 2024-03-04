package token

import (
	"errors"
	"time"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type SubscriptionInfo struct {
	Subscriptions map[string]string `json:"subscriptions"` // key: subscription id, value: product id
}

type TokenSetting struct {
	ID           string                 `json:"id"`
	Description  string                 `json:"desc"`
	Enabled      bool                   `json:"enabled"`
	Permissions  map[string]*Permission `json:"permissions"`
	Subscription *SubscriptionInfo      `json:"subscription"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

func (ts *TokenSetting) GetSubscriptionByProduct(product string) (string, error) {

	if ts.Subscription == nil {
		return "", ErrSubscriptionNotFound
	}

	for id, p := range ts.Subscription.Subscriptions {
		if p == product {
			return id, nil
		}
	}

	return "", ErrSubscriptionNotFound
}
