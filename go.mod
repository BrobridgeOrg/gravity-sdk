module github.com/BrobridgeOrg/gravity-sdk

go 1.15

require (
	github.com/BrobridgeOrg/broton v0.0.3
	github.com/BrobridgeOrg/gravity-api v0.2.15
	github.com/cfsghost/parallel-chunked-flow v0.0.6
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
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api
//replace github.com/BrobridgeOrg/broton => ../../broton
