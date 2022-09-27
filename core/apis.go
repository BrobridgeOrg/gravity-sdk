package core

const (
	CoreAPI = "$GVT.%s.API.CORE"
)

type Permissions map[string]string

var AvailablePermissions = Permissions{

	// Administrator
	"ADMIN": "Administrator",

	// Product
	"PRODUCT.LIST":          "List available products",
	"PRODUCT.CREATE":        "Create product",
	"PRODUCT.DELETE":        "Delete specific product",
	"PRODUCT.UPDATE":        "Update specific product",
	"PRODUCT.PURGE":         "Purge specific product",
	"PRODUCT.INFO":          "Get specific product information",
	"PRODUCT.SUBSCRIPTION":  "Subscribe to specific product",
	"PRODUCT.SNAPSHOT.READ": "Read snapshot of specific product",
	"PRODUCT.ACL":           "Update ACL of specific product",

	// Token
	"TOKEN.LIST":   "List available tokens",
	"TOKEN.CREATE": "Create token",
	"TOKEN.DELETE": "Delete specific token",
	"TOKEN.UPDATE": "Update specific token",
	"TOKEN.INFO":   "Get specific token information",
}

type AuthenticateRequest struct {
	Token string `json:"token"`
}

type AuthenticateReply struct {
	ErrorReply

	Durable     string   `json:"durable"`
	Permissions []string `json:"permissions"`
}
