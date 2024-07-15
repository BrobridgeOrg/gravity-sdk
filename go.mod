module github.com/BrobridgeOrg/gravity-sdk/v2

go 1.21

require (
	github.com/json-iterator/go v1.1.12
	github.com/nats-io/nats.go v1.33.1
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api

//replace github.com/BrobridgeOrg/gravity-sdk => ../gravity-sdk

//replace github.com/BrobridgeOrg/broton => ../../broton
