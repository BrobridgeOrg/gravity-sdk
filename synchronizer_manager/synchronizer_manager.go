package synchronizer_manager

import (
	"errors"
	"os"
	"time"

	synchronizer_manager_pb "github.com/BrobridgeOrg/gravity-api/service/synchronizer_manager"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type SynchronizerManager struct {
	client  *core.Client
	options *Options
}

func NewSynchronizerManager(options *Options) *SynchronizerManager {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	return &SynchronizerManager{
		options: options,
	}
}

func NewSynchronizerManagerWithClient(client *core.Client, options *Options) *SynchronizerManager {

	synchronizer := NewSynchronizerManager(options)
	synchronizer.client = client

	return synchronizer
}

func (sm *SynchronizerManager) Connect(host string, options *core.Options) error {
	sm.client = core.NewClient()
	return sm.client.Connect(host, options)
}

func (sm *SynchronizerManager) Disconnect() {
	sm.client.Disconnect()
}

func (sm *SynchronizerManager) Register(synchronizerID string) error {

	conn := sm.client.GetConnection()

	request := synchronizer_manager_pb.RegisterSynchronizerRequest{
		SynchronizerID: synchronizerID,
	}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request("gravity.synchronizer_manager.register", msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply synchronizer_manager_pb.RegisterSynchronizerReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (sm *SynchronizerManager) GetSynchronizers() ([]*Synchronizer, error) {

	conn := sm.client.GetConnection()

	request := synchronizer_manager_pb.GetSynchronizersRequest{}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request("gravity.synchronizer_manager.getSynchronizers", msg, time.Second*10)
	if err != nil {
		return nil, err
	}

	var reply synchronizer_manager_pb.GetSynchronizersReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	synchronizers := make([]*Synchronizer, 0, len(reply.Synchronizers))
	for _, sub := range reply.Synchronizers {
		synchronizers = append(synchronizers, &Synchronizer{
			ID: sub.SynchronizerID,
		})
	}

	return synchronizers, nil
}

func (sm *SynchronizerManager) ReleasePipelines(synchronizerID string, pipelines []uint64) error {

	conn := sm.client.GetConnection()

	request := synchronizer_manager_pb.ReleasePipelinesRequest{
		SynchronizerID: synchronizerID,
		Pipelines:      pipelines,
	}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request("gravity.synchronizer_manager.releasePipelines", msg, time.Second*10)
	if err != nil {
		return err
	}

	var reply synchronizer_manager_pb.ReleasePipelinesReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (sm *SynchronizerManager) GetPipelines(synchronizerID string) ([]uint64, error) {

	conn := sm.client.GetConnection()

	request := synchronizer_manager_pb.GetPipelinesRequest{
		SynchronizerID: synchronizerID,
	}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request("gravity.synchronizer_manager.getPipelines", msg, time.Second*10)
	if err != nil {
		return nil, err
	}

	var reply synchronizer_manager_pb.GetPipelinesReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	return reply.Pipelines, nil
}
