package authenticator

import (
	"fmt"
	"time"

	packet_pb "github.com/BrobridgeOrg/gravity-api/packet"
	"github.com/golang/protobuf/proto"
)

func (auth *Authenticator) request(method string, data []byte, encrypted bool) ([]byte, error) {

	// Getting endpoint from client object
	endpoint, err := auth.GetEndpoint()
	if err != nil {
		return []byte(""), err
	}

	// Prepare channel string
	baseChannel, err := auth.GetChannel()
	if err != nil {
		return nil, err
	}

	channel := fmt.Sprintf("%s.%s", baseChannel, method)

	conn := endpoint.GetConnection()

	// Preparing packet
	packet := packet_pb.Packet{
		AppID:   "auth",
		Payload: data,
	}

	// Encrypt
	if encrypted {
		payload, err := auth.encryption.Encrypt(data)
		if err != nil {
			return []byte(""), err
		}

		packet.Payload = payload
	}

	msg, _ := proto.Marshal(&packet)

	// Send request
	resp, err := conn.Request(channel, msg, time.Second*10)
	if err != nil {
		return []byte(""), err
	}

	// Decrypt
	if encrypted {
		data, err = auth.encryption.Decrypt(resp.Data)
		if err != nil {
			return []byte(""), err
		}

		return data, nil
	}

	return resp.Data, nil
}
