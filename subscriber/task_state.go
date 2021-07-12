package subscriber

import "github.com/sirupsen/logrus"

type SchedulerTaskState int32

const (
	SchedulerTaskState_Idle SchedulerTaskState = iota
	SchedulerTaskState_Busy
	SchedulerTaskState_Suspend
)

var SchedulerTaskState_name = map[SchedulerTaskState]string{
	0: "Idle",
	1: "Busy",
	2: "Suspend",
}

type TaskState struct {
	pipeline *Pipeline
	state    SchedulerTaskState
}

func NewTaskState() *TaskState {
	return &TaskState{}
}

func (ts *TaskState) SetState(state SchedulerTaskState) {

	log.WithFields(logrus.Fields{
		"pipeline": ts.pipeline.id,
		"from":     SchedulerTaskState_name[ts.state],
	}).Printf("Switch task state to %s", SchedulerTaskState_name[state])

	ts.state = state
}
