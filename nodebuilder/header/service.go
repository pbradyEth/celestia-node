package header

import (
	"context"

	"github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/celestia-node/header/p2p"
	"github.com/celestiaorg/celestia-node/header/sync"
)

// Module exposes the functionality needed for querying headers from the network.
// Any method signature changed here needs to also be changed in the API struct.
//
//go:generate mockgen -destination=mocks/api.go -package=mocks . Module
type Module interface {
	// GetByHeight returns the ExtendedHeader at the given height, blocking
	// until header has been processed by the store or context deadline is exceeded.
	GetByHeight(context.Context, uint64) (*header.ExtendedHeader, error)
	// Head returns the ExtendedHeader of the chain head.
	Head(context.Context) (*header.ExtendedHeader, error)
	// IsSyncing returns the status of sync
	IsSyncing() bool
}

// API is a wrapper around Module for the RPC.
// TODO(@distractedm1nd): These structs need to be autogenerated.
type API struct {
	GetByHeight func(context.Context, uint64) (*header.ExtendedHeader, error)
	Head        func(context.Context) (*header.ExtendedHeader, error)
	IsSyncing   func() bool
}

// service represents the header service that can be started / stopped on a node.
// service's main function is to manage its sub-services. service can contain several
// sub-services, such as Exchange, ExchangeServer, Syncer, and so forth.
type service struct {
	ex header.Exchange

	syncer    *sync.Syncer
	sub       header.Subscriber
	p2pServer *p2p.ExchangeServer
	store     header.Store
}

// NewHeaderService creates a new instance of header service.
func NewHeaderService(
	syncer *sync.Syncer,
	sub header.Subscriber,
	p2pServer *p2p.ExchangeServer,
	ex header.Exchange,
	store header.Store) Module {
	return &service{
		syncer:    syncer,
		sub:       sub,
		p2pServer: p2pServer,
		ex:        ex,
		store:     store,
	}
}

func (s *service) GetByHeight(ctx context.Context, height uint64) (*header.ExtendedHeader, error) {
	return s.store.GetByHeight(ctx, height)
}

func (s *service) Head(ctx context.Context) (*header.ExtendedHeader, error) {
	return s.store.Head(ctx)
}

func (s *service) IsSyncing() bool {
	return !s.syncer.State().Finished()
}
