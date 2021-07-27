package subscriber_manager

import (
	"time"

	packet_pb "github.com/BrobridgeOrg/gravity-api/packet"
	"github.com/golang/protobuf/proto"
)

func (sm *SubscriberManager) request(method string, data []byte, encrypted bool) ([]byte, error) {

	// Getting endpoint from client object
	endpoint, err := sm.GetEndpoint()
	if err != nil {
		return []byte(""), err
	}

	conn := endpoint.GetConnection()
	key := sm.options.Key

	// Preparing packet
	packet := packet_pb.Packet{
		AppID:   key.GetAppID(),
		Payload: data,
	}

	// Encrypt
	if encrypted {
		payload, err := key.Encryption().Encrypt(data)
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

	// Decrypt
	if encrypted {
		data, err = key.Encryption().Decrypt(resp.Data)
		if err != nil {
			return []byte(""), err
		}

		return data, nil
	}

	return resp.Data, nil
}
