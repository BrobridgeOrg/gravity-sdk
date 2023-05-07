package token

import "github.com/BrobridgeOrg/gravity-sdk/v2/core"

const (
	TokenAPI = "$GVT.%s.API.TOKEN"
)

type ListAvailablePermissionsRequest struct {
}

type ListAvailablePermissionsReply struct {
	core.ErrorReply

	Permissions map[string]string `json:"permissions"`
}

type ListTokensRequest struct {
}

type ListTokensReply struct {
	core.ErrorReply

	Tokens []*TokenSetting `json:"products"`
}

type CreateTokenRequest struct {
	TokenID string        `json:"tokenID"`
	Setting *TokenSetting `json:"setting"`
}

type CreateTokenReply struct {
	core.ErrorReply
	Token   string        `json:"token"`
	Setting *TokenSetting `json:"setting"`
}

type UpdateTokenRequest struct {
	TokenID string        `json:"token"`
	Setting *TokenSetting `json:"setting"`
}

type UpdateTokenReply struct {
	core.ErrorReply
	Setting *TokenSetting `json:"setting"`
}

type DeleteTokenRequest struct {
	TokenID string `json:"tokenID"`
}

type DeleteTokenReply struct {
	core.ErrorReply
}

type InfoTokenRequest struct {
	TokenID string `json:"token"`
}

type InfoTokenReply struct {
	core.ErrorReply
	Token   string        `json:"token"`
	Setting *TokenSetting `json:"setting"`
}
