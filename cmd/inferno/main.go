package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/vayzur/inferno/internal/config"
	"github.com/vayzur/inferno/internal/server"
	"github.com/vayzur/inferno/pkg/client/spark"
	"github.com/vayzur/inferno/pkg/controller"
	"github.com/vayzur/inferno/pkg/flock"
	"github.com/vayzur/inferno/pkg/httputil"
	"github.com/vayzur/inferno/pkg/service"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/vayzur/inferno/pkg/storage/etcd"
	"github.com/vayzur/inferno/pkg/storage/resources"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		zlog.Fatal().Err(err).Msg("config load failed")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.AppConfig.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		zlog.Fatal().Err(err).Msg("etcd connect failed")
	}

	defer etcdClient.Close()

	etcdStorege := etcd.NewEtcdStorage(etcdClient)

	inboundStore := resources.NewInboundStore(etcdStorege)
	nodeStore := resources.NewNodeStore(etcdStorege)

	httpClient := httputil.New(time.Second * 5)
	sparkClient := spark.NewSparkClient(httpClient)

	inboundService := service.NewInboundService(inboundStore, sparkClient)
	nodeService := service.NewNodeSerivce(nodeStore)

	controllerManager := controller.NewControllerManager(nodeService, inboundService)

	// TODO: using etcd distributed lock instead of file lock

	nodeControllerLock := flock.NewFlock("/tmp/inferno-node-controller.lock")
	if err := nodeControllerLock.TryLock(); err == nil {
		controllerManager.StartNodeMonitor(
			ctx,
			config.AppConfig.NodeMonitorPeriod,
			config.AppConfig.NodeMonitorGracePeriod,
		)
		defer nodeControllerLock.Unlock()
	}

	inboundControllerLock := flock.NewFlock("/tmp/inferno-inbound-controller.lock")
	if err := inboundControllerLock.TryLock(); err == nil {
		controllerManager.StartInboundMonitor(
			ctx,
			config.AppConfig.InboundMonitorPeriod,
		)
		defer inboundControllerLock.Unlock()
	}

	serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.Address, config.AppConfig.Port)

	apiserver := server.NewServer(serverAddr, inboundService, nodeService)

	go func() {
		if config.AppConfig.TLS.Enabled {
			zlog.Fatal().Err(apiserver.StartTLS())
		} else {
			zlog.Fatal().Err(apiserver.Start())
		}
	}()

	defer apiserver.Stop()

	zlog.Info().Str("component", "apiserver").Msg("server started")
	<-ctx.Done()
	zlog.Info().Str("component", "apiserver").Msg("server stopped")
}
