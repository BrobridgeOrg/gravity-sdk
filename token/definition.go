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

func (ts *TokenSetting) CheckPermission(perm string) bool {

	if _, ok := ts.Permissions["ADMIN"]; ok {
		return true
	}

	return false
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

func (ts *TokenSetting) CheckSubscription(subscriptionID string) bool {

	if ts.CheckPermission("ADMIN") {
		return true
	}

	if _, ok := ts.Subscription.Subscriptions[subscriptionID]; ok {
		return true
	}

	return false
}
