package subscriber_manager

import (
	"time"

	subscriber_manager_pb "github.com/BrobridgeOrg/gravity-api/service/subscriber_manager"
)

type Subscriber struct {
	ID          string
	Name        string
	Component   string
	Type        subscriber_manager_pb.SubscriberType
	LastCheck   time.Time
	AppID       string
	AccessKey   string
	Permissions []string
	Collections []string
}
