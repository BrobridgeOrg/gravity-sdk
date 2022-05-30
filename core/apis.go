package core

const (
	CoreAPI = "$GVT.%s.API.CORE"
)

type AuthenticateRequest struct {
	Token string `json:"token"`
}

type AuthenticateReply struct {
	ErrorReply

	Durable     string   `json:"durable"`
	Permissions []string `json:"permissions"`
}
