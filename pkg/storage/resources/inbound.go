package resources

import (
	"context"
	"encoding/json"
	"fmt"

	zlog "github.com/rs/zerolog/log"
	v1 "github.com/vayzur/inferno/pkg/api/v1"
	"github.com/vayzur/inferno/pkg/storage"
)

type InboundStore struct {
	store storage.Storage
}

func NewInboundStore(store storage.Storage) *InboundStore {
	return &InboundStore{store: store}
}

func (s *InboundStore) GetInbound(ctx context.Context, nodeID, tag string) (*v1.Inbound, error) {
	key := fmt.Sprintf("/inbounds/%s/%s", nodeID, tag)
	resp, err := s.store.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get inbound %s/%s: %w", nodeID, tag, err)
	}
	var inbound v1.Inbound
	if err := json.Unmarshal(resp, &inbound); err != nil {
		return nil, fmt.Errorf("unmarshal inbound %s/%s: %w", nodeID, tag, err)
	}

	return &inbound, nil
}

func (s *InboundStore) PutInbound(ctx context.Context, nodeID string, inbound *v1.Inbound) error {
	val, err := json.Marshal(inbound)
	if err != nil {
		return fmt.Errorf("marshal inbound %s/%s: %w", nodeID, inbound.Config.Tag, err)
	}

	key := fmt.Sprintf("/inbounds/%s/%s", nodeID, inbound.Config.Tag)
	if err := s.store.Put(ctx, key, string(val)); err != nil {
		return fmt.Errorf("put inbound %s/%s: %w", nodeID, inbound.Config.Tag, err)
	}

	return nil
}

func (s *InboundStore) DelInbound(ctx context.Context, nodeID, tag string) error {
	key := fmt.Sprintf("/inbounds/%s/%s", nodeID, tag)
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete inbound %s/%s: %w", nodeID, tag, err)
	}
	return nil
}

func (s *InboundStore) ListInbounds(ctx context.Context, nodeID string) ([]*v1.Inbound, error) {
	key := fmt.Sprintf("/inbounds/%s/", nodeID)
	resp, err := s.store.List(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("list inbounds %s: %w", nodeID, err)
	}

	var inbounds []*v1.Inbound
	for k, v := range resp {
		var inbound v1.Inbound
		if err := json.Unmarshal(v, &inbound); err != nil {
			zlog.Error().Err(err).Str("component", "inbound").Str("nodeID", nodeID).Str("tag", k).Msg("unmarshal failed")
			continue
		}
		inbounds = append(inbounds, &inbound)
	}

	return inbounds, nil
}
