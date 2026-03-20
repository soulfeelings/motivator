package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-backend/internal/middleware"
	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/service"
	"github.com/hustlers/motivator-backend/pkg/response"
)

type QuestHandler struct {
	svc *service.QuestService
}

func NewQuestHandler(svc *service.QuestService) *QuestHandler {
	return &QuestHandler{svc: svc}
}

func (h *QuestHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	var req model.CreateQuestRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	log.Printf("trace=%s | creating quest company=%s", traceID, companyID)
	q, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Created(c, q)
}

func (h *QuestHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	log.Printf("trace=%s | listing quests company=%s", traceID, companyID)
	qs, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, qs)
}

func (h *QuestHandler) GetByID(c fiber.Ctx) error {
	questID := c.Params("questId")
	q, err := h.svc.GetByID(c.Context(), questID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "quest not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, q)
}

func (h *QuestHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	questID := c.Params("questId")
	log.Printf("trace=%s | deleting quest=%s", traceID, questID)
	if err := h.svc.Delete(c.Context(), questID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "quest not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

func (h *QuestHandler) Start(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	questID := c.Params("questId")
	companyID := c.Params("id")
	log.Printf("trace=%s | starting quest=%s", traceID, questID)
	if err := h.svc.Start(c.Context(), questID, companyID); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, fiber.Map{"started": true})
}

func (h *QuestHandler) GetMyTarget(c fiber.Ctx) error {
	questID := c.Params("questId")
	m := middleware.GetMembership(c)
	pair, err := h.svc.GetMyTarget(c.Context(), questID, m.ID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "you are not in this quest")
		}
		return response.InternalError(c)
	}
	// Don't reveal receiver's identity directly — just the pair ID and receiver name
	return response.Success(c, fiber.Map{
		"pair_id":     pair.ID,
		"receiver_id": pair.ReceiverID,
		"sent":        pair.SentAt != nil,
		"message":     pair.Message,
	})
}

func (h *QuestHandler) SendMessage(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	questID := c.Params("questId")
	m := middleware.GetMembership(c)
	var req model.SendMessageRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Message == "" {
		return response.BadRequest(c, "message is required")
	}
	log.Printf("trace=%s | member=%s sending quest message quest=%s", traceID, m.ID, questID)
	if err := h.svc.SendMessage(c.Context(), questID, m.ID, req.Message); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, fiber.Map{"sent": true})
}

func (h *QuestHandler) GetReceivedMessages(c fiber.Ctx) error {
	questID := c.Params("questId")
	m := middleware.GetMembership(c)
	messages, err := h.svc.GetReceivedMessages(c.Context(), questID, m.ID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, messages)
}

func (h *QuestHandler) StartVoting(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	questID := c.Params("questId")
	log.Printf("trace=%s | starting voting quest=%s", traceID, questID)
	if err := h.svc.StartVoting(c.Context(), questID); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"voting": true})
}

func (h *QuestHandler) Vote(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	questID := c.Params("questId")
	m := middleware.GetMembership(c)
	var req model.VoteRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	log.Printf("trace=%s | member=%s voting quest=%s", traceID, m.ID, questID)
	if err := h.svc.Vote(c.Context(), questID, m.ID, req.PairID); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, fiber.Map{"voted": true})
}

func (h *QuestHandler) Reveal(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	questID := c.Params("questId")
	log.Printf("trace=%s | revealing quest=%s", traceID, questID)
	if err := h.svc.Reveal(c.Context(), questID); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"revealed": true})
}

func (h *QuestHandler) Complete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	questID := c.Params("questId")
	log.Printf("trace=%s | completing quest=%s", traceID, questID)
	if err := h.svc.Complete(c.Context(), questID); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"completed": true})
}

func (h *QuestHandler) ListPairs(c fiber.Ctx) error {
	questID := c.Params("questId")
	pairs, err := h.svc.ListPairs(c.Context(), questID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, pairs)
}
