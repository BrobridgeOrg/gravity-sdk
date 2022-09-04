package subscriber

import (
	"fmt"
	"sync"
	"testing"

	gravity_sdk_types_pipeline_event "github.com/BrobridgeOrg/gravity-sdk/types/pipeline_event"
	gravity_sdk_types_snapshot_record "github.com/BrobridgeOrg/gravity-sdk/types/snapshot_record"
	"github.com/stretchr/testify/assert"
)

func TestPipelineInitialization(t *testing.T) {

	p := NewPipeline(
		0,
		0,
		WithPipelineInitialLoad(false, "snapshot", 1000),
		WithPipelineRemoteLastSeq(0),
		WithPipelineSnapshotLastSeq(0),
		WithPipelineStateStore(NewStateStoreDummy()),
		WithPipelineChunkSize(2048),
		WithPipelineSubscriberID("subscriberTest"),
		WithPipelineRequest(NewPipelineRequestDummyImpl()),
		WithPipelineSnapshotRequest(NewSnapshotRequestDummyImpl()),
		//WithPipelineSubscription(sub.subscription),
		WithPipelineCollections([]string{"test"}),
	)

	err := p.Initialize()
	assert.Equal(t, nil, err)
}

func TestPipelinePull(t *testing.T) {

	chunkSize := 100

	// Prepare dummy pipeline request to simulate RPC
	targetCount := 1000
	pr := NewPipelineRequestDummyImpl()
	for i := 0; i < targetCount; i++ {
		pr.add([]byte(fmt.Sprintf("%d", i+1)))
	}

	p := NewPipeline(
		0,
		0,
		WithPipelineInitialLoad(false, "snapshot", 1000),
		WithPipelineRemoteLastSeq(0),
		WithPipelineSnapshotLastSeq(0),
		WithPipelineStateStore(NewStateStoreDummy()),
		WithPipelineChunkSize(chunkSize),
		WithPipelineSubscriberID("subscriberTest"),
		WithPipelineRequest(pr),
		WithPipelineSnapshotRequest(NewSnapshotRequestDummyImpl()),
		//WithPipelineSubscription(sub.subscription),
		WithPipelineCollections([]string{"test"}),
	)

	err := p.Initialize()
	assert.Equal(t, nil, err)

	counter := uint64(0)
	for {
		events, err := p.pull()
		assert.Equal(t, nil, err)

		if events == nil {
			break
		}

		assert.Equal(t, chunkSize, len(events))

		for _, event := range events {
			counter++
			assert.Equal(t, fmt.Sprintf("%d", counter), string(event))
		}

		assert.Equal(t, counter, p.lastSeq)
	}

	assert.Equal(t, uint64(targetCount), counter)
}

func TestPipelineAwake(t *testing.T) {

	pipelineID := uint64(1)
	msgs := make(chan *Message)

	// Prepare subscription
	var s Subscription
	s = NewSubscriptionImpl(
		WithSubscriptionEventHandler(func(msg *Message) {
			msgs <- msg
		}),
		WithSubscriptionSnapshotHandler(func(msg *Message) {
		}),
	)

	s.Start()
	defer s.Unsubscribe()

	// Prepare dummy pipeline request to simulate RPC
	targetCount := 1000
	pr := NewPipelineRequestDummyImpl()
	for i := 0; i < targetCount; i++ {

		// Preparing original pipeline event
		event := gravity_sdk_types_pipeline_event.PipelineEvent{}
		event.PipelineID = pipelineID
		event.Sequence = uint64(i + 1)

		data, err := gravity_sdk_types_pipeline_event.Marshal(&event)
		assert.Equal(t, nil, err)

		pr.add(data)
	}

	p := NewPipeline(
		pipelineID,
		0,
		WithPipelineInitialLoad(false, "snapshot", 1000),
		WithPipelineRemoteLastSeq(0),
		WithPipelineStateStore(NewStateStoreDummy()),
		WithPipelineChunkSize(100),
		WithPipelineSubscriberID("subscriberTest"),
		WithPipelineRequest(pr),
		WithPipelineSnapshotRequest(NewSnapshotRequestDummyImpl()),
		WithPipelineSubscription(s),
		WithPipelineCollections([]string{"test"}),
	)

	err := p.Initialize()
	assert.Equal(t, nil, err)

	// Simulate internal scheduler to awake pipeline until no events
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			// No more events
			if p.status == PIPELINE_STATUS_CLOSE {
				wg.Done()
				break
			}

			err = p.Awake()
			assert.Equal(t, nil, err)
		}
	}()

	// Message handling
	counter := 0
	for msg := range msgs {

		event := msg.Payload.(*DataEvent)
		assert.Equal(t, pipelineID, event.PipelineID)
		assert.Equal(t, uint64(counter+1), event.Sequence)
		msg.Ack()

		counter++
		if counter == targetCount {
			break
		}
	}

	select {
	case <-msgs:
		t.Error("message is much more than expected")
	default:
	}

	wg.Wait()
}

func TestPipelineAwakebyUnstableRequests(t *testing.T) {

	counter := 0
	pipelineID := uint64(1)
	msgs := make(chan *Message)

	// Prepare subscription
	var s Subscription
	s = NewSubscriptionImpl(
		WithSubscriptionEventHandler(func(msg *Message) {
			msgs <- msg
		}),
		WithSubscriptionSnapshotHandler(func(msg *Message) {
		}),
	)

	s.Start()
	defer s.Unsubscribe()

	// Prepare dummy pipeline request to simulate RPC
	targetCount := 1000
	pr := NewPipelineRequestDummyImpl()

	// Encounter connection problem 5 times
	unstableCount := 5
	pr.unstableCount = unstableCount

	for i := 0; i < targetCount; i++ {

		// Preparing original pipeline event
		event := gravity_sdk_types_pipeline_event.PipelineEvent{}
		event.PipelineID = pipelineID
		event.Sequence = uint64(i + 1)

		data, err := gravity_sdk_types_pipeline_event.Marshal(&event)
		assert.Equal(t, nil, err)

		pr.add(data)
	}

	p := NewPipeline(
		pipelineID,
		0,
		WithPipelineInitialLoad(false, "snapshot", 1000),
		WithPipelineRemoteLastSeq(0),
		WithPipelineStateStore(NewStateStoreDummy()),
		WithPipelineChunkSize(100),
		WithPipelineSubscriberID("subscriberTest"),
		WithPipelineRequest(pr),
		WithPipelineSnapshotRequest(NewSnapshotRequestDummyImpl()),
		WithPipelineSubscription(s),
		WithPipelineCollections([]string{"test"}),
	)

	err := p.Initialize()
	assert.Equal(t, nil, err)

	// Simulate internal scheduler to awake pipeline until no events
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		unstableCounter := 0
		for {
			// No more events
			if p.status == PIPELINE_STATUS_CLOSE && p.lastSeq == uint64(targetCount) {
				assert.Equal(t, unstableCount, unstableCounter)
				wg.Done()
				break
			}

			err := p.Awake()
			if err != nil && err.Error() == "Unstable" {
				unstableCounter++
			}
		}
	}()

	// Message handling
	for msg := range msgs {

		event := msg.Payload.(*DataEvent)
		assert.Equal(t, pipelineID, event.PipelineID)
		assert.Equal(t, uint64(counter+1), event.Sequence)
		msg.Ack()

		counter++
		if counter == targetCount {
			break
		}
	}

	select {
	case <-msgs:
		t.Error("message is much more than expected")
	default:
	}

	wg.Wait()
}

func TestPipelineAwake_Snapshot(t *testing.T) {

	pipelineID := uint64(1)
	msgs := make(chan *Message)

	// Prepare subscription
	var s Subscription
	s = NewSubscriptionImpl(
		WithSubscriptionEventHandler(func(msg *Message) {
		}),
		WithSubscriptionSnapshotHandler(func(msg *Message) {
			msgs <- msg
		}),
	)

	s.Start()
	defer s.Unsubscribe()

	// Prepare dummy snapshot request to simulate RPC
	targetCount := 1000
	sr := NewSnapshotRequestDummyImpl()
	for i := 0; i < targetCount; i++ {

		// Preparing original pipeline event
		r := gravity_sdk_types_snapshot_record.SnapshotRecord{}

		data, err := gravity_sdk_types_snapshot_record.Marshal(&r)
		assert.Equal(t, nil, err)

		sr.add(data)
	}

	// Prepare dummy pipeline request to simulate RPC
	pr := NewPipelineRequestDummyImpl()

	p := NewPipeline(
		pipelineID,
		0,
		WithPipelineInitialLoad(true, "snapshot", 1000),
		WithPipelineRemoteLastSeq(0),
		WithPipelineSnapshotLastSeq(uint64(targetCount)),
		WithPipelineStateStore(NewStateStoreDummy()),
		WithPipelineChunkSize(100),
		WithPipelineSubscriberID("subscriberTest"),
		WithPipelineRequest(pr),
		WithPipelineSnapshotRequest(sr),
		WithPipelineSubscription(s),
		WithPipelineCollections([]string{"test_collection"}),
	)

	err := p.Initialize()
	assert.Equal(t, nil, err)

	// Simulate internal scheduler to awake pipeline until no events
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err = p.Awake()
			assert.Equal(t, nil, err)

			if p.snapshot.isCompleted {
				break
			}
		}
	}()

	// Message handling
	counter := 0
	for msg := range msgs {

		event := msg.Payload.(*SnapshotEvent)
		assert.Equal(t, pipelineID, event.PipelineID)
		assert.Equal(t, "test_collection", event.Collection)
		msg.Ack()

		counter++
		if counter == targetCount {
			break
		}
	}

	assert.Equal(t, uint64(counter), p.lastSeq)

	select {
	case <-msgs:
		t.Error("message is much more than expected")
	default:
	}

	wg.Wait()
}
