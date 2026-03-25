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

type SocialGameHandler struct {
	svc *service.SocialGameService
}

func NewSocialGameHandler(svc *service.SocialGameService) *SocialGameHandler {
	return &SocialGameHandler{svc: svc}
}

// Create godoc
// @Summary Create a social game
// @Description Create a new social game for a company
// @Tags social-games
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.CreateSocialGameRequest true "Game details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games [post]
func (h *SocialGameHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	m := middleware.GetMembership(c)
	var req model.CreateSocialGameRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	log.Printf("trace=%s | creating social game company=%s", traceID, companyID)
	g, err := h.svc.Create(c.Context(), companyID, m.ID, req)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Created(c, g)
}

// List godoc
// @Summary List social games
// @Description List all social games for a company
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games [get]
func (h *SocialGameHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	log.Printf("trace=%s | listing social games company=%s", traceID, companyID)
	games, err := h.svc.List(c.Context(), companyID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, games)
}

// GetByID godoc
// @Summary Get a social game
// @Description Get a social game by ID
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId} [get]
func (h *SocialGameHandler) GetByID(c fiber.Ctx) error {
	gameID := c.Params("gameId")
	g, err := h.svc.GetByID(c.Context(), gameID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "social game not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, g)
}

// Delete godoc
// @Summary Delete a social game
// @Description Delete a social game by ID
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId} [delete]
func (h *SocialGameHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	gameID := c.Params("gameId")
	log.Printf("trace=%s | deleting social game=%s", traceID, gameID)
	if err := h.svc.Delete(c.Context(), gameID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "social game not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

// AddQuestion godoc
// @Summary Add a question to a social game
// @Description Add a trivia question to a social game
// @Tags social-games
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Param body body model.AddQuestionRequest true "Question details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/questions [post]
func (h *SocialGameHandler) AddQuestion(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	gameID := c.Params("gameId")
	var req model.AddQuestionRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	log.Printf("trace=%s | adding question to game=%s", traceID, gameID)
	q, err := h.svc.AddQuestion(c.Context(), gameID, req)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Created(c, q)
}

// ListQuestions godoc
// @Summary List questions for a social game
// @Description List all questions for a social game
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/questions [get]
func (h *SocialGameHandler) ListQuestions(c fiber.Ctx) error {
	gameID := c.Params("gameId")
	questions, err := h.svc.ListQuestions(c.Context(), gameID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, questions)
}

// Launch godoc
// @Summary Launch a social game
// @Description Transition a social game from draft to active
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/launch [post]
func (h *SocialGameHandler) Launch(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	gameID := c.Params("gameId")
	log.Printf("trace=%s | launching social game=%s", traceID, gameID)
	if err := h.svc.Launch(c.Context(), gameID); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, fiber.Map{"launched": true})
}

// StartVoting godoc
// @Summary Start voting on a social game
// @Description Transition a social game from active to voting phase
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/voting [post]
func (h *SocialGameHandler) StartVoting(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	gameID := c.Params("gameId")
	log.Printf("trace=%s | starting voting social game=%s", traceID, gameID)
	if err := h.svc.StartVoting(c.Context(), gameID); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, fiber.Map{"voting": true})
}

// Complete godoc
// @Summary Complete a social game
// @Description Complete a social game and award XP/coins
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/complete [post]
func (h *SocialGameHandler) Complete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	gameID := c.Params("gameId")
	log.Printf("trace=%s | completing social game=%s", traceID, gameID)
	if err := h.svc.Complete(c.Context(), gameID); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"completed": true})
}

// SubmitAnswer godoc
// @Summary Submit a trivia answer
// @Description Submit an answer to a trivia question
// @Tags social-games
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Param body body model.SubmitAnswerRequest true "Answer details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/answer [post]
func (h *SocialGameHandler) SubmitAnswer(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	gameID := c.Params("gameId")
	m := middleware.GetMembership(c)
	var req model.SubmitAnswerRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.QuestionID == "" {
		return response.BadRequest(c, "question_id is required")
	}
	log.Printf("trace=%s | member=%s submitting answer game=%s", traceID, m.ID, gameID)
	a, err := h.svc.SubmitAnswer(c.Context(), gameID, m.ID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, a)
}

// SubmitEntry godoc
// @Summary Submit a game entry
// @Description Submit a photo or two-truths entry
// @Tags social-games
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Param body body model.SubmitEntryRequest true "Entry details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/submit [post]
func (h *SocialGameHandler) SubmitEntry(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	gameID := c.Params("gameId")
	m := middleware.GetMembership(c)
	var req model.SubmitEntryRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	log.Printf("trace=%s | member=%s submitting entry game=%s", traceID, m.ID, gameID)
	sub, err := h.svc.SubmitEntry(c.Context(), gameID, m.ID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, sub)
}

// CastVote godoc
// @Summary Cast a vote
// @Description Cast a vote on a submission
// @Tags social-games
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Param body body model.CastVoteRequest true "Vote details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/vote [post]
func (h *SocialGameHandler) CastVote(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	gameID := c.Params("gameId")
	m := middleware.GetMembership(c)
	var req model.CastVoteRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.SubmissionID == "" {
		return response.BadRequest(c, "submission_id is required")
	}
	log.Printf("trace=%s | member=%s voting game=%s", traceID, m.ID, gameID)
	v, err := h.svc.CastVote(c.Context(), gameID, m.ID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, v)
}

// GetResults godoc
// @Summary Get game results
// @Description Get results including participation rate and leaderboard
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/results [get]
func (h *SocialGameHandler) GetResults(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	gameID := c.Params("gameId")
	log.Printf("trace=%s | getting results game=%s", traceID, gameID)
	results, err := h.svc.GetResults(c.Context(), gameID, companyID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "social game not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, results)
}

// ListSubmissions godoc
// @Summary List game submissions
// @Description List all submissions for a social game
// @Tags social-games
// @Produce json
// @Param id path string true "Company ID"
// @Param gameId path string true "Game ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /companies/{id}/social-games/{gameId}/submissions [get]
func (h *SocialGameHandler) ListSubmissions(c fiber.Ctx) error {
	gameID := c.Params("gameId")
	submissions, err := h.svc.ListSubmissions(c.Context(), gameID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, submissions)
}

