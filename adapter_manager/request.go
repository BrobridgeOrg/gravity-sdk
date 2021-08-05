package adapter_manager

import (
	"errors"
	"time"

	packet_pb "github.com/BrobridgeOrg/gravity-api/packet"
	"github.com/golang/protobuf/proto"
)

func (am *AdapterManager) request(method string, data []byte, encrypted bool) ([]byte, error) {

	// Getting endpoint from client object
	endpoint, err := am.GetEndpoint()
	if err != nil {
		return []byte(""), err
	}

	conn := endpoint.GetConnection()
	key := am.options.Key

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
