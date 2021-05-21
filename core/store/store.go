package store

import (
	broton "github.com/BrobridgeOrg/broton"
)

type Store struct {
	options     *Options
	storeEngine *broton.Broton
}

func NewStore(options *Options) (*Store, error) {

	es, err := broton.NewBroton(options.StoreOptions)
	if err != nil {
		return nil, err
	}

	return &Store{
		storeEngine: es,
		options:     options,
	}, nil
}

func (ss *Store) GetEngine() *broton.Broton {
	return ss.storeEngine
}
