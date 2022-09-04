package subscriber

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshotPull(t *testing.T) {

	chunkSize := 100

	// Prepare dummy snapshot request to simulate RPC
	targetCount := 1000
	sr := NewSnapshotRequestDummyImpl()
	for i := 0; i < targetCount; i++ {
		sr.add([]byte(fmt.Sprintf("%d", i+1)))
	}

	s := NewSnapshot(
		WithSnapshotPipelineID(1),
		WithSnapshotSubscriberID("test_subscriber"),
		WithSnapshotChunkSize(chunkSize),
		WithSnapshotRequest(sr),
		WithSnapshotCollections([]string{"test_collection"}),
	)

	err := s.Create()
	assert.Equal(t, nil, err)
	defer s.Close()

	counter := uint64(0)
	for {
		col, records, err := s.Pull()
		assert.Equal(t, nil, err)

		if len(records) == 0 {
			break
		}

		assert.Equal(t, chunkSize, len(records))

		for _, r := range records {
			counter++
			assert.Equal(t, fmt.Sprintf("%d", counter), string(r))
		}

		assert.Equal(t, "test_collection", col)
		assert.Equal(t, fmt.Sprintf("%d", counter), string(s.getLastPosition().lastKey))
	}

	assert.Equal(t, uint64(targetCount), counter)
}

func TestSnapshotPullUnstable(t *testing.T) {

	chunkSize := 100

	// Prepare dummy snapshot request to simulate RPC
	targetCount := 1000
	sr := NewSnapshotRequestDummyImpl()
	unstableCount := 5
	sr.unstableCount = unstableCount
	for i := 0; i < targetCount; i++ {
		sr.add([]byte(fmt.Sprintf("%d", i+1)))
	}

	s := NewSnapshot(
		WithSnapshotPipelineID(1),
		WithSnapshotSubscriberID("test_subscriber"),
		WithSnapshotChunkSize(chunkSize),
		WithSnapshotRequest(sr),
		WithSnapshotCollections([]string{"test_collection"}),
	)

	err := s.Create()
	assert.Equal(t, nil, err)
	defer s.Close()

	unstableCounter := 0
	counter := uint64(0)
	for {
		col, records, err := s.Pull()
		if err != nil && err.Error() == "Unstable" {
			unstableCounter++
			continue
		}

		if len(records) == 0 {
			break
		}

		assert.Equal(t, chunkSize, len(records))

		for _, r := range records {
			counter++
			assert.Equal(t, fmt.Sprintf("%d", counter), string(r))
		}

		assert.Equal(t, "test_collection", col)
		assert.Equal(t, fmt.Sprintf("%d", counter), string(s.getLastPosition().lastKey))
	}

	assert.Equal(t, uint64(targetCount), counter)
	assert.Equal(t, unstableCount, unstableCounter)
}
