package subscriber_manager

import (
	subscriber_manager_pb "github.com/BrobridgeOrg/gravity-api/service/subscriber_manager"
)

type Subscriber struct {
	ID        string
	Name      string
	Component string
	Type      subscriber_manager_pb.SubscriberType
}
