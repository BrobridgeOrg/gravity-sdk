package adapter

import (
	"errors"
	"time"

	packet_pb "github.com/BrobridgeOrg/gravity-api/packet"
	"github.com/golang/protobuf/proto"
)

var (
	ErrInvalidPacket = errors.New("adapter: invalid packet")
)

type Request struct {
	IsCompleted bool
	Key         []byte
	EventName   string
	Payload     []byte
}

func (ac *AdapterConnector) request(method string, data []byte, encrypted bool) ([]byte, error) {

	// Getting endpoint from client object
	endpoint, err := ac.GetEndpoint()
	if err != nil {
		return []byte(""), err
	}

	conn := endpoint.GetConnection()
	key := ac.options.Key

	// Preparing payload
	payload := packet_pb.Payload{
		Data: data,
	}

	payloadData, _ := proto.Marshal(&payload)

	// Preparing packet
	packet := packet_pb.Packet{
		AppID:   key.GetAppID(),
		Payload: payloadData,
	}

	// Encrypt
	if encrypted {
		payload, err := key.Encryption().Encrypt(packet.Payload)
		if err != nil {
			return []byte(""), err
		}

		packet.Payload = payload
	}

	msg, _ := proto.Marshal(&packet)

	// Send request
	resp, err := conn.Request(endpoint.Channel(method), msg, time.Second*10)
	if err != nil {
		return []byte(""), err
	}

	// Parsing data
	err = proto.Unmarshal(resp.Data, &packet)
	if err != nil {
		log.Error(err)
		return []byte(""), errors.New("Forbidden")
	}

	if packet.Error {
		return []byte(""), errors.New(packet.Reason)
	}

	// Decrypt
	if encrypted {
		data, err = key.Encryption().Decrypt(packet.Payload)
		if err != nil {
			return []byte(""), errors.New("Forbidden")
		}

		return data, nil
	}

	return packet.Payload, nil
}
