package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	zlog "github.com/rs/zerolog/log"
	v1 "github.com/vayzur/inferno/pkg/api/v1"
	"github.com/vayzur/inferno/pkg/errs"
)

func (s *Server) GetAllNodes(c fiber.Ctx) error {
	nodes, err := s.nodeService.ListNodes(c.RequestCtx())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.Status(http.StatusOK).JSON(nodes)
}

func (s *Server) GetActiveNodes(c fiber.Ctx) error {
	nodes, err := s.nodeService.ListActiveNodes(c.RequestCtx())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.Status(http.StatusOK).JSON(nodes)
}

func (s *Server) GetNode(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	node, err := s.nodeService.GetNode(c.RequestCtx(), nodeID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.Status(http.StatusOK).JSON(node)
}

func (s *Server) CreateNode(c fiber.Ctx) error {
	node := new(v1.Node)
	if err := c.Bind().JSON(node); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	node.Metadata.ID = uuid.NewString()
	node.Metadata.CreationTimestamp = time.Now()

	if err := s.nodeService.PutNode(c.RequestCtx(), node); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.Status(http.StatusCreated).JSON(node)
}

func (s *Server) DeleteNode(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	if err := s.nodeService.DelNode(context.Background(), nodeID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) UpdateNodeStatus(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	nodeStatus := new(v1.NodeStatus)
	if err := c.Bind().JSON(nodeStatus); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.nodeService.UpdateNodeStatus(c.RequestCtx(), nodeID, nodeStatus); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}
	zlog.Info().Str("component", "apiserver").Str("nodeID", nodeID).Msg("node status updated")
	return c.SendStatus(fiber.StatusOK)
}
