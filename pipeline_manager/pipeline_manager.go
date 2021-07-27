package pipeline_manager

import (
	"errors"
	"os"

	pb "github.com/BrobridgeOrg/gravity-api/service/pipeline_manager"
	core "github.com/BrobridgeOrg/gravity-sdk/core"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type PipelineManager struct {
	client  *core.Client
	options *Options
}

func NewPipelineManager(options *Options) *PipelineManager {

	log.Out = os.Stdout
	log.SetLevel(logrus.ErrorLevel)

	if options.Verbose {
		log.SetLevel(logrus.InfoLevel)
	}

	pm := &PipelineManager{
		options: options,
	}

	return pm
}

func NewPipelineManagerWithClient(client *core.Client, options *Options) *PipelineManager {

	pipeline := NewPipelineManager(options)
	pipeline.client = client

	return pipeline
}

func (pm *PipelineManager) Connect(host string, options *core.Options) error {
	pm.client = core.NewClient()
	return pm.client.Connect(host, options)
}

func (pm *PipelineManager) Disconnect() {
	pm.client.Disconnect()
}

func (pm *PipelineManager) GetEndpoint() (*core.Endpoint, error) {
	return pm.client.ConnectToEndpoint(pm.options.Endpoint, pm.options.Domain, nil)
}

func (pm *PipelineManager) GetPipelineCount() (uint64, error) {

	request := pb.GetPipelineCountRequest{}
	msg, _ := proto.Marshal(&request)

	respData, err := pm.request("pipeline_manager.getCount", msg, true)
	if err != nil {
		return 0, err
	}

	var reply pb.GetPipelineCountReply
	err = proto.Unmarshal(respData, &reply)
	if err != nil {
		return 0, err
	}

	if !reply.Success {
		return 0, errors.New(reply.Reason)
	}

	return reply.Count, nil
}
