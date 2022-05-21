package product

import "github.com/BrobridgeOrg/gravity-sdk/core"

func NotFoundSnapshotViewErr() *core.Error {
	return &core.Error{
		Code:    44404,
		Message: "Not found snapshot view",
	}
}
