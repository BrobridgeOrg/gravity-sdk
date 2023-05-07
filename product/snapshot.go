package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/BrobridgeOrg/gravity-sdk/v2/core"
	"github.com/nats-io/nats.go"
)

type SnapshotOpt func(*Snapshot)

const (
	SnapshotAPI          = "$GVT.%s.API.SNAPSHOT"
	SnapshotEventSubject = "$GVT.%s.SS.%s.>"
)

var (
	ErrSnapshotNotReady = errors.New("snapshot view is not ready")
)

type CreateSnapshotViewRequest struct {
	Subscriber string  `json:"subscriber"`
	Product    string  `json:"product"`
	Partitions []int32 `json:"partitions"`
}

type CreateSnapshotViewReply struct {
	core.ErrorReply

	ID         string    `json:"id"`
	Subscriber string    `json:"subscriber"`
	Product    string    `json:"product"`
	Stream     string    `json:"stream"`
	CreatedAt  time.Time `json:"createAt"`
}

type DeleteSnapshotViewRequest struct {
	ID string `json:"id"`
}

type DeleteSnapshotViewReply struct {
	core.ErrorReply

	ID string `json:"id"`
}

type PullSnapshotViewRequest struct {
	ID           string `json:"id"`
	Partition    int32  `json:"partition"`
	LastKey      string `json:"lastKey"`
	AfterLastKey bool   `json:"afterLastKey"`
}

type PullSnapshotViewReply struct {
	core.ErrorReply

	ID      string `json:"id"`
	LastKey string `json:"lastKey"`
	Count   int    `json:"count"`
}

type SnapshotSetting struct {

	// TODO
}

type Snapshot struct {
	pc          *ProductClient
	id          string
	product     string
	stream      string
	partitions  []int32
	isReady     bool
	sub         *nats.Subscription
	pending     chan *nats.Msg
	cancelFetch context.CancelFunc
	lastKey     []byte
}

func WithSnapshotPartitions(partitions []int32) SnapshotOpt {
	return func(s *Snapshot) {
		s.partitions = partitions
	}
}

func WithSnapshotPendingLimits(limit int) SnapshotOpt {
	return func(s *Snapshot) {
		s.pending = make(chan *nats.Msg, limit)
	}
}

func NewSnapshot(pc *ProductClient, productName string, opts ...SnapshotOpt) (*Snapshot, error) {

	s := &Snapshot{
		pc:         pc,
		product:    productName,
		partitions: make([]int32, 0),
		pending:    make(chan *nats.Msg, 4096),
		lastKey:    []byte(""),
	}

	// Apply options
	for _, o := range opts {
		o(s)
	}

	err := s.initialize()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Snapshot) initialize() error {

	nc := s.pc.client.GetConnection()

	// Create snapshot view
	req := &CreateSnapshotViewRequest{
		Product:    s.product,
		Partitions: s.partitions,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(SnapshotAPI+".VIEW.CREATE", s.pc.options.Domain)
	msg, err := nc.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return err
	}

	// Parsing response
	resp := &CreateSnapshotViewReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	s.id = resp.ID

	// Subscribe to stream for snapshot view
	js, err := s.pc.client.GetJetStream()
	if err != nil {
		return err
	}

	// Preparing subject
	subject := fmt.Sprintf(SnapshotEventSubject, s.pc.options.Domain, s.id)

	// Subscribe
	sub, err := js.Subscribe(subject, s.handleMessage)
	if err != nil {
		return err
	}

	sub.SetPendingLimits(-1, -1)

	s.sub = sub
	s.isReady = true

	return nil
}

func (s *Snapshot) release() error {

	if !s.isReady {
		return nil
	}

	s.isReady = false

	if s.sub != nil {

		// Unsubscribe stream
		err := s.sub.Unsubscribe()
		if err != nil {
			return err
		}
	}

	if s.cancelFetch != nil {
		s.cancelFetch()
		s.cancelFetch = nil
	}

	close(s.pending)

	if len(s.id) == 0 {
		return nil
	}

	nc := s.pc.client.GetConnection()

	// delete snapshot view
	req := &DeleteSnapshotViewRequest{
		ID: s.id,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(SnapshotAPI+".VIEW.DELETE", s.pc.options.Domain)
	msg, err := nc.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return err
	}

	// Parsing response
	resp := &DeleteSnapshotViewReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	return nil
}

func (s *Snapshot) handleMessage(msg *nats.Msg) {
	s.pending <- msg
}

func (s *Snapshot) pull() (int, error) {

	if s.sub == nil {
		return 0, ErrSnapshotNotReady
	}

	if len(s.id) == 0 {
		return 0, nil
	}

	nc := s.pc.client.GetConnection()

	//	lastKey := base64.StdEncoding.EncodeToString(s.lastKey)

	// Fetch data from snapshot view
	req := &PullSnapshotViewRequest{
		ID: s.id,
		//		LastKey:      lastKey,
		//		AfterLastKey: true,
	}

	reqData, _ := json.Marshal(req)

	// Send request
	apiPath := fmt.Sprintf(SnapshotAPI+".VIEW.PULL", s.pc.options.Domain)
	msg, err := nc.Request(apiPath, reqData, time.Second*30)
	if err != nil {
		return 0, err
	}

	// Parsing response
	resp := &PullSnapshotViewReply{}
	err = json.Unmarshal(msg.Data, resp)
	if err != nil {
		return 0, err
	}

	if resp.Error != nil {
		return 0, errors.New(resp.Error.Message)
	}

	if resp.Count == 0 {
		// No more records
		fmt.Println("EOF")
		return 0, nil
	}
	/*
		// Getting last key
		s.lastKey, err = base64.StdEncoding.DecodeString(resp.LastKey)
		if err != nil {
			return 0, err
		}
	*/
	return resp.Count, nil
}

func (s *Snapshot) Close() error {
	return s.release()
}

func (s *Snapshot) Fetch() (chan *nats.Msg, error) {

	if s.sub == nil {
		return nil, ErrSnapshotNotReady
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFetch = cancel

	go func(ctx context.Context) {
		for {

			select {
			case <-ctx.Done():
				return
			default:
			}

			// Check pending buffer
			if len(s.pending) > int(float32(cap(s.pending))*0.5) {
				time.Sleep(time.Millisecond * 300)
				continue
			}

			// Pull data from server
			count, err := s.pull()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// No more records in partition
			if count == 0 {
				s.release()
				return
			}
		}
	}(ctx)

	return s.pending, nil
}
