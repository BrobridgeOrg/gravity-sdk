package synchronizer_manager

import (
	"errors"
	"os"

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

	sm := &SynchronizerManager{
		options: options,
	}

	return sm
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

func (sm *SynchronizerManager) GetEndpoint() (*core.Endpoint, error) {
	return sm.client.ConnectToEndpoint(sm.options.Endpoint, sm.options.Domain, nil)
}

func (sm *SynchronizerManager) Register(synchronizerID string) error {

	request := synchronizer_manager_pb.RegisterSynchronizerRequest{
		SynchronizerID: synchronizerID,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := sm.request("synchronizer_manager.register", msg, true)
	if err != nil {
		return err
	}

	var reply synchronizer_manager_pb.RegisterSynchronizerReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (sm *SynchronizerManager) GetSynchronizers() ([]*Synchronizer, error) {

	request := synchronizer_manager_pb.GetSynchronizersRequest{}
	msg, _ := proto.Marshal(&request)

	respData, err := sm.request("synchronizer_manager.getSynchronizers", msg, true)
	if err != nil {
		return nil, err
	}

	var reply synchronizer_manager_pb.GetSynchronizersReply
	err = proto.Unmarshal(respData, &reply)
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

	request := synchronizer_manager_pb.ReleasePipelinesRequest{
		SynchronizerID: synchronizerID,
		Pipelines:      pipelines,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := sm.request("synchronizer_manager.releasePipelines", msg, true)
	if err != nil {
		return err
	}

	var reply synchronizer_manager_pb.ReleasePipelinesReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.New(reply.Reason)
	}

	return nil
}

func (sm *SynchronizerManager) GetPipelines(synchronizerID string) ([]uint64, error) {

	request := synchronizer_manager_pb.GetPipelinesRequest{
		SynchronizerID: synchronizerID,
	}
	msg, _ := proto.Marshal(&request)

	respData, err := sm.request("synchronizer_manager.getPipelines", msg, true)
	if err != nil {
		return nil, err
	}

	var reply synchronizer_manager_pb.GetPipelinesReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, errors.New(reply.Reason)
	}

	return reply.Pipelines, nil
}
