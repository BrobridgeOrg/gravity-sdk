module github.com/BrobridgeOrg/gravity-sdk

go 1.15

require (
	github.com/BrobridgeOrg/EventStore v0.0.5
	github.com/BrobridgeOrg/broton v0.0.2
	github.com/BrobridgeOrg/gravity-adapter-nats v0.0.0-20201117192323-d62ed567fa57
	github.com/BrobridgeOrg/gravity-api v0.2.1
	github.com/cfsghost/parallel-chunked-flow v0.0.3
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/copier v0.3.0
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/lib/pq v1.10.1 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/prometheus/common v0.15.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/protobuf v1.23.0
)

replace github.com/BrobridgeOrg/gravity-api => ../gravity-api
