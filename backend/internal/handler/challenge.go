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

type ChallengeHandler struct {
	svc *service.ChallengeService
}

func NewChallengeHandler(svc *service.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{svc: svc}
}

// Create godoc
// @Summary Create a 1v1 challenge
// @Description Challenge another member to a metric competition
// @Tags challenges
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.CreateChallengeRequest true "Challenge data"
// @Success 201 {object} model.Challenge
// @Router /companies/{id}/challenges [post]
func (h *ChallengeHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	m := middleware.GetMembership(c)

	var req model.CreateChallengeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.OpponentID == "" || req.Metric == "" || req.Target <= 0 {
		return response.BadRequest(c, "opponent_id, metric, and target are required")
	}

	log.Printf("trace=%s | creating challenge from=%s to=%s metric=%s", traceID, m.ID, req.OpponentID, req.Metric)

	challenge, err := h.svc.Create(c.Context(), companyID, m.ID, req)
	if err != nil {
		log.Printf("trace=%s | error creating challenge: %v", traceID, err)
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, challenge)
}

// Accept godoc
// @Summary Accept a challenge
// @Description Accept a pending challenge
// @Tags challenges
// @Param id path string true "Company ID"
// @Param challengeId path string true "Challenge ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/challenges/{challengeId}/accept [post]
func (h *ChallengeHandler) Accept(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	challengeID := c.Params("challengeId")
	m := middleware.GetMembership(c)

	log.Printf("trace=%s | accepting challenge=%s by=%s", traceID, challengeID, m.ID)

	if err := h.svc.Accept(c.Context(), challengeID, m.ID); err != nil {
		if errors.Is(err, model.ErrForbidden) {
			return response.Forbidden(c, "only the opponent can accept")
		}
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, fiber.Map{"accepted": true})
}

// Decline godoc
// @Summary Decline a challenge
// @Description Decline a pending challenge
// @Tags challenges
// @Param id path string true "Company ID"
// @Param challengeId path string true "Challenge ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/challenges/{challengeId}/decline [post]
func (h *ChallengeHandler) Decline(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	challengeID := c.Params("challengeId")
	m := middleware.GetMembership(c)

	log.Printf("trace=%s | declining challenge=%s by=%s", traceID, challengeID, m.ID)

	if err := h.svc.Decline(c.Context(), challengeID, m.ID); err != nil {
		if errors.Is(err, model.ErrForbidden) {
			return response.Forbidden(c, "only the opponent can decline")
		}
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, fiber.Map{"declined": true})
}

// ReportScore godoc
// @Summary Report challenge score
// @Description Report your score in an active challenge
// @Tags challenges
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param challengeId path string true "Challenge ID"
// @Param body body model.ReportScoreRequest true "Score"
// @Success 200 {object} model.Challenge
// @Router /companies/{id}/challenges/{challengeId}/score [post]
func (h *ChallengeHandler) ReportScore(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	challengeID := c.Params("challengeId")
	m := middleware.GetMembership(c)

	var req model.ReportScoreRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | reporting score=%d for challenge=%s by=%s", traceID, req.Score, challengeID, m.ID)

	challenge, err := h.svc.ReportScore(c.Context(), challengeID, m.ID, req.Score)
	if err != nil {
		if errors.Is(err, model.ErrForbidden) {
			return response.Forbidden(c, "not a participant")
		}
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, challenge)
}

// List godoc
// @Summary List challenges
// @Description List challenges for a company
// @Tags challenges
// @Produce json
// @Param id path string true "Company ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.Response
// @Router /companies/{id}/challenges [get]
func (h *ChallengeHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | listing challenges company=%s", traceID, companyID)

	var statusFilter *model.ChallengeStatus
	if s := c.Query("status"); s != "" {
		st := model.ChallengeStatus(s)
		statusFilter = &st
	}

	challenges, err := h.svc.ListByCompany(c.Context(), companyID, statusFilter)
	if err != nil {
		log.Printf("trace=%s | error listing challenges: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, challenges)
}
