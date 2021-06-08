package subscriber

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	workerCount int
	idle        chan *Pipeline
	states      sync.Map
}

func NewScheduler(taskCount int, workerCount int) *Scheduler {
	return &Scheduler{
		workerCount: workerCount,
		idle:        make(chan *Pipeline, taskCount),
	}
}

func (scheduler *Scheduler) Initialize() {

	for i := 0; i < scheduler.workerCount; i++ {
		go scheduler.startWorker(i)
	}
}

func (scheduler *Scheduler) startWorker(workerID int) {

	log.WithFields(logrus.Fields{
		"worker": workerID,
	}).Info("Started wroker")

	for {
		select {
		case pipeline := <-scheduler.idle:

			taskState := scheduler.GetTaskState(pipeline.id)

			// Not idle
			if taskState.state != SchedulerTaskState_Idle {
				continue
			}

			taskState.SetState(SchedulerTaskState_Busy)

			// Fetching data from specfic pipeline
			err := pipeline.Pull()
			if err != nil {
				log.WithFields(logrus.Fields{
					"worker":   workerID,
					"pipeline": pipeline.id,
				}).Error(err)

				taskState.SetState(SchedulerTaskState_Idle)
			}

			if pipeline.isSuspended {
				taskState.SetState(SchedulerTaskState_Suspend)
				continue
			}
		}
	}
}

func (scheduler *Scheduler) AddPipeline(pipeline *Pipeline) {

	taskState := NewTaskState()
	taskState.pipeline = pipeline
	taskState.state = SchedulerTaskState_Idle
	scheduler.states.Store(pipeline.id, taskState)

	scheduler.idle <- pipeline
}

func (scheduler *Scheduler) GetTaskState(pipelineID uint64) *TaskState {

	state, ok := scheduler.states.Load(pipelineID)
	if !ok {
		return nil
	}

	return state.(*TaskState)
}

func (scheduler *Scheduler) Awake(pipelineID uint64) {

	log.WithFields(logrus.Fields{
		"pipeline": pipelineID,
	}).Info("Awake pipeline")

	taskState := scheduler.GetTaskState(pipelineID)
	if taskState == nil {
		return
	}

	if taskState.state != SchedulerTaskState_Suspend {
		return
	}

	taskState.SetState(SchedulerTaskState_Idle)
	pipeline := taskState.pipeline
	pipeline.Awake()

	scheduler.idle <- pipeline
}

func (scheduler *Scheduler) Idle(pipelineID uint64) {

	taskState := scheduler.GetTaskState(pipelineID)
	if taskState == nil {
		return
	}

	if taskState.state == SchedulerTaskState_Idle {
		return
	}

	taskState.SetState(SchedulerTaskState_Idle)
	pipeline := taskState.pipeline

	scheduler.idle <- pipeline
}
