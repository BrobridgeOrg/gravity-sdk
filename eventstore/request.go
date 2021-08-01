package eventstore

import (
	"errors"
	"fmt"
	"time"

	packet_pb "github.com/BrobridgeOrg/gravity-api/packet"
	"github.com/golang/protobuf/proto"
)

func (es *EventStore) Request(eventstoreID string, method string, data []byte) ([]byte, error) {

	conn := es.client.GetConnection()

	// Preparing payload
	payload := packet_pb.Payload{
		Data: data,
	}

	payloadData, _ := proto.Marshal(&payload)

	// Preparing packet
	packet := packet_pb.Packet{
		AppID:   es.options.Key.GetAppID(),
		Payload: payloadData,
	}

	keyInfo := es.options.Key

	// Encrypt
	encrypted, err := keyInfo.Encryption().Encrypt(packet.Payload)
	if err != nil {
		return []byte(""), err
	}

	packet.Payload = encrypted

	msg, _ := proto.Marshal(&packet)

	// Send request
	channel := fmt.Sprintf("%s.eventstore.%s.%s", es.options.Domain, eventstoreID, method)
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
	plain, err := keyInfo.Encryption().Decrypt(packet.Payload)
	if err != nil {
		return []byte(""), errors.New("Forbidden")
	}

	return plain, nil
}
