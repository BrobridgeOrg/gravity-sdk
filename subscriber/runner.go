package subscriber

import (
	"time"

	"github.com/BrobridgeOrg/gravity-sdk/subscriber/scheduler"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	scheduler *scheduler.Scheduler
	closed    chan struct{}
}

func NewRunner() *Runner {

	opts := scheduler.NewOptions()
	s := scheduler.NewScheduler(opts)

	return &Runner{
		scheduler: s,
		closed:    make(chan struct{}),
	}
}

func (runner *Runner) execute(task *scheduler.Task, privData interface{}) scheduler.TaskState {

	pipeline := privData.(*Pipeline)

	// Awake pipeline to work
	err := pipeline.Awake()
	if err != nil {
		log.WithFields(logrus.Fields{
			"pipeline": pipeline.id,
		}).Error(err)
	}

	// still open so set state to re-schedule
	if pipeline.status == PIPELINE_STATUS_OPEN || pipeline.status == PIPELINE_STATUS_HALF_OPEN {
		return scheduler.TASK_STATE_PREPARED
	}

	return scheduler.TASK_STATE_IDLE
}

func (runner *Runner) awakeAllTasks() {

	tasks := runner.scheduler.GetAllTasks()
	if len(tasks) == 0 {
		return
	}

	// Preparing all tasks to be ready to go
	for _, task := range tasks {
		task.SetState(scheduler.TASK_STATE_PREPARED)
	}

	runner.scheduler.Trigger(scheduler.SIGNAL_AWAKE)
}

func (runner *Runner) AddPipeline(pipeline *Pipeline) {
	task := runner.scheduler.AddTask(pipeline.id, pipeline, runner.execute)
	task.SetState(scheduler.TASK_STATE_PREPARED)

	// TODO: for testing
	pipeline.suspend()
}

func (runner *Runner) Awake(pipelineID uint64) {

	task := runner.scheduler.GetTask(pipelineID)
	if task == nil {
		return
	}

	task.SetState(scheduler.TASK_STATE_PREPARED)
	runner.scheduler.Trigger(scheduler.SIGNAL_AWAKE)
}

func (runner *Runner) Start() {

	log.Info("Starting runner")

	go runner.scheduler.Start()

	// First time to awake all tasks for initializing
	runner.awakeAllTasks()

	// Initializing delay timer
	idleDelay := time.NewTimer(time.Second * 10)
	defer idleDelay.Stop()

	// Try to awake pipeline every 10 seconds
	for {
		idleDelay.Reset(time.Second * 10)

		select {
		case <-runner.closed:
			log.Info("Closing runner")
			return
		case <-idleDelay.C:

			if runner.scheduler.GetState() == scheduler.STATE_AWAKE {
				continue
			}

			log.Info("Trying to awake all pipeline")

			runner.awakeAllTasks()
		}
	}
}

func (runner *Runner) Stop() {
	runner.closed <- struct{}{}
	runner.scheduler.Stop()
}
