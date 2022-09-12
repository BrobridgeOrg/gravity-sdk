module github.com/BrobridgeOrg/gravity-sdk

go 1.15

require (
	github.com/BrobridgeOrg/broton v0.0.7
	github.com/BrobridgeOrg/gravity-api v0.2.28
	github.com/BrobridgeOrg/gravity-transmitter-postgres v0.0.0-20210824004507-cf5a2b3fb5e7 // indirect
	github.com/BrobridgeOrg/sequential-data-flow v0.0.2
	github.com/cfsghost/buffered-input v0.0.2
	github.com/cfsghost/parallel-chunked-flow v0.0.6
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/copier v0.3.2 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/nats-io/nats.go v1.16.0
	github.com/prometheus/common v0.4.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.8.0 // indirect
	google.golang.org/protobuf v1.26.0
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api

//replace github.com/BrobridgeOrg/gravity-sdk => ../gravity-sdk

//replace github.com/BrobridgeOrg/broton => ../../broton
