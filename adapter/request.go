package adapter

type Request struct {
	IsCompleted bool
	Key         []byte
	EventName   string
	Payload     []byte
}
