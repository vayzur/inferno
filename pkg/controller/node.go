package controller

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"
)

func (c *ControllerManager) StartNodeMonitor(ctx context.Context, nodeMonitorPeriod, nodeMonitorGracePeriod time.Duration) {
	go c.runNodeMonitor(ctx, nodeMonitorPeriod, nodeMonitorGracePeriod)
}

func (c *ControllerManager) runNodeMonitor(ctx context.Context, nodeMonitorPeriod, nodeMonitorGracePeriod time.Duration) {
	ticker := time.NewTicker(nodeMonitorPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "controller").Msg("node monitor started")
	for {
		select {
		case <-ctx.Done():
			zlog.Info().Str("component", "controller").Msg("node monitor stopped")
			return
		case <-ticker.C:
			nodes, err := c.nodeService.ListNodes(ctx)
			if err != nil {
				zlog.Error().Err(err).Str("component", "controller").Msg("failed to get nodes")
				continue
			}

			now := time.Now()
			for _, node := range nodes {
				if now.Sub(node.Status.LastHeartbeatTime) >= nodeMonitorGracePeriod {
					node.Status.Status = false
					if err := c.nodeService.PutNode(ctx, node); err != nil {
						zlog.Error().Err(err).Str("component", "controller").Msg("failed to update node status")
						continue
					}
				}
			}
		}
	}
}
