module github.com/BrobridgeOrg/gravity-sdk

go 1.15

require (
	github.com/BrobridgeOrg/broton v0.0.7
	github.com/BrobridgeOrg/compton v0.0.0-20220315055918-5c0ee958b7d2
	github.com/BrobridgeOrg/gravity-api v0.2.25
	github.com/BrobridgeOrg/gravity-dispatcher v0.0.0-20220206181110-2ea65aa048be // indirect
	github.com/BrobridgeOrg/gravity-snapshot v0.0.0-20220213082459-fd38b9058d95 // indirect
	github.com/BrobridgeOrg/gravity-transmitter-postgres v0.0.0-20210824004507-cf5a2b3fb5e7 // indirect
	github.com/BrobridgeOrg/schemer v0.0.10 // indirect
	github.com/BrobridgeOrg/sequential-data-flow v0.0.1 // indirect
	github.com/bytedance/sonic v1.1.1 // indirect
	github.com/cfsghost/buffered-input v0.0.1
	github.com/cfsghost/parallel-chunked-flow v0.0.6 // indirect
	github.com/doug-martin/goqu/v9 v9.18.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/copier v0.3.2 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d
	github.com/prometheus/common v0.9.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api

//replace github.com/BrobridgeOrg/gravity-sdk => ../gravity-sdk

//replace github.com/BrobridgeOrg/broton => ../../broton
