package gravity_sdk_types_pipeline_event

import (
	"github.com/golang/protobuf/proto"
)

func Marshal(event *PipelineEvent) ([]byte, error) {
	return proto.Marshal(event)
}

func Unmarshal(data []byte, event *PipelineEvent) error {
	return proto.Unmarshal(data, event)
}
