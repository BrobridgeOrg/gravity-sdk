package adapter

type Request struct {
	IsCompleted bool
	EventName   string
	Payload     []byte
}
