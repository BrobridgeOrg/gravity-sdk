module github.com/BrobridgeOrg/gravity-sdk

go 1.15

require (
	github.com/BrobridgeOrg/broton v0.0.3
	github.com/BrobridgeOrg/gravity-api v0.2.17
	github.com/cfsghost/parallel-chunked-flow v0.0.6
	github.com/golang/protobuf v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/kr/pretty v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nats-io/nats-server/v2 v2.1.2 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api
//replace github.com/BrobridgeOrg/broton => ../../broton
