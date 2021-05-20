package pipeline_manager

import (
	"errors"
	"os"
	"time"

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

	return &PipelineManager{
		options: options,
	}
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

func (pm *PipelineManager) GetPipelineCount() (uint64, error) {

	conn := pm.client.GetConnection()

	request := pb.GetPipelineCountRequest{}
	msg, _ := proto.Marshal(&request)

	resp, err := conn.Request("gravity.pipeline_manager.getCount", msg, time.Second*10)
	if err != nil {
		return 0, err
	}

	var reply pb.GetPipelineCountReply
	err = proto.Unmarshal(resp.Data, &reply)
	if err != nil {
		return 0, err
	}

	if !reply.Success {
		return 0, errors.New(reply.Reason)
	}

	return reply.Count, nil
}
