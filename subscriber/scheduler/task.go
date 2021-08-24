package scheduler

import (
	"sync"
)

type TaskState int32

const (
	TASK_STATE_IDLE TaskState = iota
	TASK_STATE_PREPARED
	TASK_STATE_RUNNING
)

type Task struct {
	ID       uint64
	State    TaskState
	PrivData interface{}
	Handler  func(*Task, interface{}) TaskState

	mutex sync.RWMutex
}

func NewTask(id uint64, privData interface{}, fn func(*Task, interface{}) TaskState) *Task {
	return &Task{
		ID:       id,
		State:    TASK_STATE_IDLE,
		PrivData: privData,
		Handler:  fn,
	}
}

func (task *Task) Execute() {
	state := task.Handler(task, task.PrivData)
	task.SetState(state)
}

func (task *Task) SetState(state TaskState) {
	task.mutex.Lock()
	task.State = state
	task.mutex.Unlock()

	//	fmt.Printf("SET TASK %d ------ STATE: %d\n", task.ID, state)
}

func (task *Task) GetState() TaskState {
	task.mutex.RLock()
	state := task.State
	task.mutex.RUnlock()

	//	fmt.Printf("GET TASK %d ------ STATE: %d\n", task.ID, state)
	return state
}
