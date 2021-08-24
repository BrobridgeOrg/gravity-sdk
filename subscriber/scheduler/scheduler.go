package scheduler

import (
	"sync"
)

type SignalType int32

const (
	SIGNAL_AWAKE SignalType = iota
)

type SchedulerState int32

const (
	STATE_IDLE SchedulerState = iota
	STATE_AWAKE
)

type Scheduler struct {
	options   *Options
	state     SchedulerState
	prepared  chan *Task
	signal    chan SignalType
	tasks     map[uint64]*Task
	mutex     sync.RWMutex
	completed chan struct{}
}

func NewScheduler(options *Options) *Scheduler {
	return &Scheduler{
		options:   options,
		state:     STATE_IDLE,
		signal:    make(chan SignalType, 1000),
		prepared:  make(chan *Task),
		tasks:     make(map[uint64]*Task),
		completed: make(chan struct{}, 10000),
	}
}

func (scheduler *Scheduler) Start() {

	// Initializing workers
	for i := 0; i < scheduler.options.WorkerCount; i++ {
		go scheduler.startWorker(i)
	}

	scheduler.startLoop()
}

func (scheduler *Scheduler) Stop() {
	close(scheduler.prepared)
	close(scheduler.signal)
}

func (scheduler *Scheduler) startLoop() {

	// Initializing signals
	for signal := range scheduler.signal {
		switch signal {
		case SIGNAL_AWAKE:
			scheduler.awake()
		}
	}
}

func (scheduler *Scheduler) startWorker(workerID int) {

	for task := range scheduler.prepared {
		task.Execute()
		scheduler.completed <- struct{}{}
	}
}

func (scheduler *Scheduler) awake() {

	scheduler.mutex.Lock()
	defer scheduler.mutex.Unlock()

	scheduler.state = STATE_AWAKE

	// Scanning all tasks
	count := 0
	for _, task := range scheduler.tasks {

		state := task.GetState()
		if state == TASK_STATE_PREPARED {
			task.SetState(TASK_STATE_RUNNING)
			scheduler.prepared <- task
			count++
		}
	}

	if count > 0 {

		// Waiting for all tasks
		completed := 0
		for range scheduler.completed {
			completed++
			if completed == count {
				break
			}
		}
	}

	scheduler.state = STATE_IDLE

	// One more time to make sure everything's done
	if count > 0 {
		scheduler.Trigger(SIGNAL_AWAKE)
	}
}

func (scheduler *Scheduler) AddTask(taskID uint64, privData interface{}, fn func(*Task, interface{}) TaskState) *Task {

	scheduler.mutex.Lock()
	defer scheduler.mutex.Unlock()

	task := NewTask(taskID, privData, fn)
	scheduler.tasks[taskID] = task
	return task
}

func (scheduler *Scheduler) DeleteTask(taskID uint64) {
	scheduler.mutex.Lock()
	defer scheduler.mutex.Unlock()
	delete(scheduler.tasks, taskID)
}

func (scheduler *Scheduler) GetTask(taskID uint64) *Task {

	scheduler.mutex.RLock()
	defer scheduler.mutex.RUnlock()

	if task, ok := scheduler.tasks[taskID]; ok {
		return task
	}

	return nil
}

func (scheduler *Scheduler) GetAllTasks() []*Task {

	scheduler.mutex.RLock()
	defer scheduler.mutex.RUnlock()

	tasks := make([]*Task, 0, len(scheduler.tasks))
	for _, task := range scheduler.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (scheduler *Scheduler) Trigger(signal SignalType) {
	select {
	case scheduler.signal <- signal:
	default:
	}
}

func (scheduler *Scheduler) GetState() SchedulerState {
	scheduler.mutex.RLock()
	defer scheduler.mutex.RUnlock()
	return scheduler.state
}
