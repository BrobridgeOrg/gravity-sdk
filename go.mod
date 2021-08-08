module github.com/BrobridgeOrg/gravity-sdk

go 1.15

require (
	github.com/BrobridgeOrg/broton v0.0.7
	github.com/BrobridgeOrg/gravity-api v0.2.22
	github.com/cfsghost/buffered-input v0.0.1
	github.com/cfsghost/parallel-chunked-flow v0.0.6
	github.com/golang/protobuf v1.5.2
	github.com/json-iterator/go v1.1.10
	github.com/nats-io/nats.go v1.11.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/protobuf v1.26.0
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api

//replace github.com/BrobridgeOrg/gravity-sdk => ../gravity-sdk

//replace github.com/BrobridgeOrg/broton => ../../broton
