package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/service"
	"github.com/hustlers/motivator-backend/pkg/response"
)

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var req model.CreateTeamRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return response.BadRequest(c, "name is required")
	}

	log.Printf("trace=%s | creating team name=%s company=%s", traceID, req.Name, companyID)

	team, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		log.Printf("trace=%s | error: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Created(c, team)
}

func (h *TeamHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | listing teams company=%s", traceID, companyID)

	teams, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, teams)
}

func (h *TeamHandler) GetByID(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	teamID := c.Params("teamId")

	log.Printf("trace=%s | get team=%s", traceID, teamID)

	team, err := h.svc.GetByID(c.Context(), teamID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "team not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, team)
}

func (h *TeamHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	teamID := c.Params("teamId")

	log.Printf("trace=%s | deleting team=%s", traceID, teamID)

	if err := h.svc.Delete(c.Context(), teamID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "team not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

func (h *TeamHandler) AddMember(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	teamID := c.Params("teamId")

	var req model.AddTeamMemberRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | adding member=%s to team=%s", traceID, req.MembershipID, teamID)

	tm, err := h.svc.AddMember(c.Context(), teamID, req.MembershipID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, tm)
}

func (h *TeamHandler) RemoveMember(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	teamID := c.Params("teamId")
	membershipID := c.Params("membershipId")

	log.Printf("trace=%s | removing member=%s from team=%s", traceID, membershipID, teamID)

	if err := h.svc.RemoveMember(c.Context(), teamID, membershipID); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"removed": true})
}

func (h *TeamHandler) ListMembers(c fiber.Ctx) error {
	teamID := c.Params("teamId")
	members, err := h.svc.ListMembers(c.Context(), teamID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, members)
}

func (h *TeamHandler) CreateBattle(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var req model.CreateTeamBattleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | creating team battle company=%s", traceID, companyID)

	battle, err := h.svc.CreateBattle(c.Context(), companyID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, battle)
}

func (h *TeamHandler) ListBattles(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | listing team battles company=%s", traceID, companyID)

	battles, err := h.svc.ListBattles(c.Context(), companyID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, battles)
}

func (h *TeamHandler) ReportScore(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	battleID := c.Params("battleId")

	var req model.ReportTeamScoreRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | reporting team score battle=%s team=%s score=%d", traceID, battleID, req.TeamID, req.Score)

	battle, err := h.svc.ReportScore(c.Context(), battleID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, battle)
}
