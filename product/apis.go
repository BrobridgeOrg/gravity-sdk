package product

import "github.com/BrobridgeOrg/gravity-sdk/core"

const (
	ProductAPI = "$GVT.%s.API.PRODUCT"
)

type ListProductsRequest struct {
}

type ListProductsReply struct {
	core.ErrorReply

	Products []*ProductInfo `json:"products"`
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
	State   *ProductState   `json:"state"`
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
