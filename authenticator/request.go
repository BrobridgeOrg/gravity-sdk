package authenticator

import (
	"errors"
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

	// Preparing payload
	payload := packet_pb.Payload{
		Data: data,
	}

	payloadData, _ := proto.Marshal(&payload)

	// Preparing packet
	packet := packet_pb.Packet{
		AppID:   auth.options.Key.GetAppID(),
		Payload: payloadData,
	}

	// Encrypt
	if encrypted {
		payload, err := auth.options.Key.Encryption().Encrypt(packet.Payload)
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

	// Parsing data
	err = proto.Unmarshal(resp.Data, &packet)
	if err != nil {
		log.Error(err)
		return []byte(""), errors.New("Forbidden")
	}

	// Decrypt
	if encrypted {
		data, err = auth.options.Key.Encryption().Decrypt(packet.Payload)
		if err != nil {
			return []byte(""), errors.New("Forbidden")
		}

		return data, nil
	}

	return packet.Payload, nil
}
