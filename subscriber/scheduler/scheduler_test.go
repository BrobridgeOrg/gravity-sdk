package scheduler

import (
	"sync"
	"testing"
)

var testScheduler *Scheduler
var testTaskCount uint64 = 100
var testEmptyTaskFunc = func(task *Task, data interface{}) TaskState {
	return TASK_STATE_IDLE
}

func TestScheduler(t *testing.T) {

	opts := NewOptions()
	testScheduler = NewScheduler(opts)
	go testScheduler.Start()
}

func TestSchedulerAddTasks(t *testing.T) {

	// Added tasks to scheduler
	for i := uint64(0); i < testTaskCount; i++ {
		testScheduler.AddTask(i, i, testEmptyTaskFunc)
	}

	// Check tasks
	for i := uint64(0); i < testTaskCount; i++ {
		task, ok := testScheduler.tasks[i]
		if !ok {
			t.Fail()
		}

		if task.ID != i {
			t.Fail()
		}
	}
}

func TestSchedulerDeleteTasks(t *testing.T) {

	for i := uint64(0); i < testTaskCount; i++ {
		testScheduler.DeleteTask(i)
	}

	for i := uint64(0); i < testTaskCount; i++ {
		_, ok := testScheduler.tasks[i]
		if ok {
			t.Fail()
		}
	}
}

func TestSchedulerTaskExecute(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(100)

	// Added tasks to scheduler
	for i := uint64(0); i < testTaskCount; i++ {
		task := testScheduler.AddTask(i, i, func(task *Task, data interface{}) TaskState {
			wg.Done()
			return TASK_STATE_IDLE
		})

		// Set task state to be prepared
		task.SetState(TASK_STATE_PREPARED)
	}

	testScheduler.awake()

	wg.Wait()

	// Check tasks
	for i := uint64(0); i < testTaskCount; i++ {
		task, ok := testScheduler.tasks[i]
		if !ok {
			t.Fail()
		}

		state := task.GetState()
		if state != TASK_STATE_IDLE {
			t.Fail()
		}
	}

	TestSchedulerDeleteTasks(t)
}

func TestSchedulerTrigger(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(100)

	// Added tasks to scheduler
	for i := uint64(0); i < testTaskCount; i++ {
		task := testScheduler.AddTask(i, i, func(task *Task, data interface{}) TaskState {
			wg.Done()
			return TASK_STATE_IDLE
		})

		// Set task state to be prepared
		task.SetState(TASK_STATE_PREPARED)
	}

	testScheduler.Trigger(SIGNAL_AWAKE)

	wg.Wait()

	// Check tasks
	for i := uint64(0); i < testTaskCount; i++ {
		task, ok := testScheduler.tasks[i]
		if !ok {
			t.Fail()
		}

		state := task.GetState()
		if state != TASK_STATE_IDLE {
			t.Fail()
		}
	}

	TestSchedulerDeleteTasks(t)
}

func TestSchedulerPrivateData(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(100)

	// Added tasks to scheduler
	for i := uint64(0); i < testTaskCount; i++ {

		func(i uint64) {
			task := testScheduler.AddTask(i, i, func(task *Task, data interface{}) TaskState {
				if data.(uint64) == i {
					wg.Done()
				}

				return TASK_STATE_IDLE
			})

			// Set task state to be prepared
			task.SetState(TASK_STATE_PREPARED)
		}(i)
	}

	testScheduler.Trigger(SIGNAL_AWAKE)

	wg.Wait()

	// Check tasks
	for i := uint64(0); i < testTaskCount; i++ {
		task, ok := testScheduler.tasks[i]
		if !ok {
			t.Fail()
		}

		state := task.GetState()
		if state != TASK_STATE_IDLE {
			t.Fail()
		}
	}

	TestSchedulerDeleteTasks(t)
}
