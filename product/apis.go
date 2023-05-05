package product

import "github.com/BrobridgeOrg/gravity-sdk/core"

const (
	ProductAPI = "$GVT.%s.API.PRODUCT"
)

// Product APIs
type ListProductsRequest struct {
}

type ListProductsReply struct {
	core.ErrorReply

	Products []*ProductSetting `json:"products"`
}

type CreateProductRequest struct {
	Setting *ProductSetting `json:"setting"`
}

type CreateProductReply struct {
	core.ErrorReply
	Setting *ProductSetting `json:"setting"`
}

type UpdateProductRequest struct {
	Name    string          `json:"name"`
	Setting *ProductSetting `json:"setting"`
}

type UpdateProductReply struct {
	core.ErrorReply
	Setting *ProductSetting `json:"setting"`
}

type DeleteProductRequest struct {
	Name string `json:"name"`
}

type DeleteProductReply struct {
	core.ErrorReply
}

type InfoProductRequest struct {
	Name string `json:"name"`
}

type InfoProductReply struct {
	core.ErrorReply
	Setting *ProductSetting `json:"setting"`
}

type PurgeProductRequest struct {
	Name string `json:"name"`
}

type PurgeProductReply struct {
	core.ErrorReply
}

type PrepareSubscriptionRequest struct {
	Product string `json:"product"`
}

type PrepareSubscriptionReply struct {
	core.ErrorReply
}

// Snapshot APIs
type GetProductSnapshotInfoRequest struct {
	Product string `json:"product"`
}

type GetProductSnapshotInfoReply struct {
	core.ErrorReply
	Product    string           `json:"product"`
	Partitions map[int32]uint64 `json:"partitions"`
}

type FetchProductSnapshotRequest struct {
	Product         string `json:"product"`
	Partition       int32  `json:"partition"`
	LastPrimaryKey  string `json:"lastPrimaryKey"`
	DeliverySubject string `json:"deliverySubject"`
	Limit           int64  `json:"limit"`
}

type FetchProductSnapshotReply struct {
	core.ErrorReply
	Product        string `json:"product"`
	Partition      int32  `json:"partition"`
	LastPrimaryKey string `json:"lastPrimaryKey"`
}
