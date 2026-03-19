package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/service"
	"github.com/hustlers/motivator-backend/pkg/response"
)

type LeaderboardHandler struct {
	svc *service.MembershipService
}

func NewLeaderboardHandler(svc *service.MembershipService) *LeaderboardHandler {
	return &LeaderboardHandler{svc: svc}
}

// Get godoc
// @Summary Get leaderboard
// @Description Returns members ranked by XP
// @Tags leaderboard
// @Produce json
// @Param id path string true "Company ID"
// @Param limit query int false "Number of entries (default 50)"
// @Success 200 {object} response.Response
// @Router /companies/{id}/leaderboard [get]
func (h *LeaderboardHandler) Get(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	limit := 50
	if v, err := strconv.Atoi(c.Query("limit", "50")); err == nil && v > 0 && v <= 100 {
		limit = v
	}

	log.Printf("trace=%s | fetching leaderboard company=%s limit=%d", traceID, companyID, limit)

	page := model.PaginationRequest{Page: 1, PerPage: limit}
	members, _, err := h.svc.ListByCompany(c.Context(), companyID, page)
	if err != nil {
		log.Printf("trace=%s | error fetching leaderboard: %v", traceID, err)
		return response.InternalError(c)
	}

	type entry struct {
		Rank        int     `json:"rank"`
		MemberID    string  `json:"member_id"`
		UserID      string  `json:"user_id"`
		DisplayName *string `json:"display_name,omitempty"`
		XP          int     `json:"xp"`
		Level       int     `json:"level"`
		Coins       int     `json:"coins"`
		Role        string  `json:"role"`
	}

	entries := make([]entry, len(members))
	for i, m := range members {
		entries[i] = entry{
			Rank:        i + 1,
			MemberID:    m.ID,
			UserID:      m.UserID,
			DisplayName: m.DisplayName,
			XP:          m.XP,
			Level:       m.Level,
			Coins:       m.Coins,
			Role:        string(m.Role),
		}
	}

	return response.Success(c, entries)
}
