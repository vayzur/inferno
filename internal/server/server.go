package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/vayzur/inferno/internal/auth"
	"github.com/vayzur/inferno/internal/config"
	"github.com/vayzur/inferno/pkg/service"
)

type Server struct {
	addr           string
	app            *fiber.App
	inboundService *service.InboundService
	nodeService    *service.NodeSerivce
}

func NewServer(addr string, inboundService *service.InboundService, nodeSerivce *service.NodeSerivce) *Server {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
	})
	s := &Server{
		addr:           addr,
		app:            app,
		inboundService: inboundService,
		nodeService:    nodeSerivce,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// s.app.Use(authMiddleware)

	s.app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	s.app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	nodes := v1.Group("/nodes")
	nodes.Get("", s.GetAllNodes)
	nodes.Get("/active", s.GetActiveNodes)
	nodes.Get("/:nodeID", s.GetNode)
	nodes.Post("", s.CreateNode)
	nodes.Delete("/:nodeID", s.DeleteNode)
	nodes.Patch("/:nodeID/status", s.UpdateNodeStatus)

	inbounds := nodes.Group("/:nodeID/inbounds")
	inbounds.Get("", s.GetAllInbounds)
	inbounds.Get("/:tag", s.GetInbound)
	inbounds.Post("", s.CreateInbound)
	inbounds.Delete("/:tag", s.DeleteInbound)
}

func (s *Server) StartTLS() error {
	return s.app.Listen(s.addr, fiber.ListenConfig{
		DisableStartupMessage: true,
		CertFile:              config.AppConfig.TLS.CertFile,
		CertKeyFile:           config.AppConfig.TLS.KeyFile,
		EnablePrefork:         config.AppConfig.Prefork,
	})
}

func (s *Server) Start() error {
	return s.app.Listen(s.addr, fiber.ListenConfig{
		DisableStartupMessage: true,
		EnablePrefork:         config.AppConfig.Prefork,
	})
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}

func authMiddleware(c fiber.Ctx) error {
	h := c.Get("Authorization")
	if h == "" {
		return fiber.ErrUnauthorized
	}

	if err := auth.VerifyRollingHash(h); err != nil {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
