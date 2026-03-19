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

type TournamentHandler struct {
	svc *service.TournamentService
}

func NewTournamentHandler(svc *service.TournamentService) *TournamentHandler {
	return &TournamentHandler{svc: svc}
}

func (h *TournamentHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	var req model.CreateTournamentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	log.Printf("trace=%s | creating tournament name=%s company=%s", traceID, req.Name, companyID)
	t, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, t)
}

func (h *TournamentHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	log.Printf("trace=%s | listing tournaments company=%s", traceID, companyID)
	ts, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, ts)
}

func (h *TournamentHandler) GetByID(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	tournamentID := c.Params("tournamentId")
	log.Printf("trace=%s | get tournament=%s", traceID, tournamentID)
	t, err := h.svc.GetByID(c.Context(), tournamentID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "tournament not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, t)
}

func (h *TournamentHandler) UpdateStatus(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	tournamentID := c.Params("tournamentId")
	var body struct {
		Status string `json:"status"`
	}
	if err := c.Bind().JSON(&body); err != nil {
		return response.BadRequest(c, "invalid request")
	}
	log.Printf("trace=%s | updating tournament=%s status=%s", traceID, tournamentID, body.Status)
	if err := h.svc.UpdateStatus(c.Context(), tournamentID, model.TournamentStatus(body.Status)); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"updated": true})
}

func (h *TournamentHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	tournamentID := c.Params("tournamentId")
	log.Printf("trace=%s | deleting tournament=%s", traceID, tournamentID)
	if err := h.svc.Delete(c.Context(), tournamentID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "tournament not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

func (h *TournamentHandler) Join(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	tournamentID := c.Params("tournamentId")
	m := middleware.GetMembership(c)
	log.Printf("trace=%s | member=%s joining tournament=%s", traceID, m.ID, tournamentID)
	p, err := h.svc.Join(c.Context(), tournamentID, m.ID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, p)
}

func (h *TournamentHandler) Leave(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	tournamentID := c.Params("tournamentId")
	m := middleware.GetMembership(c)
	log.Printf("trace=%s | member=%s leaving tournament=%s", traceID, m.ID, tournamentID)
	if err := h.svc.Leave(c.Context(), tournamentID, m.ID); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"left": true})
}

func (h *TournamentHandler) SubmitScore(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	tournamentID := c.Params("tournamentId")
	m := middleware.GetMembership(c)
	var req model.SubmitScoreRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request")
	}
	log.Printf("trace=%s | member=%s submitting score=%d tournament=%s", traceID, m.ID, req.Score, tournamentID)
	if err := h.svc.SubmitScore(c.Context(), tournamentID, m.ID, req.Score); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"submitted": true})
}

func (h *TournamentHandler) GetStandings(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	tournamentID := c.Params("tournamentId")
	log.Printf("trace=%s | standings tournament=%s", traceID, tournamentID)
	standings, err := h.svc.GetStandings(c.Context(), tournamentID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, standings)
}

func (h *TournamentHandler) Complete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	tournamentID := c.Params("tournamentId")
	log.Printf("trace=%s | completing tournament=%s", traceID, tournamentID)
	if err := h.svc.Complete(c.Context(), tournamentID); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"completed": true})
}
