package handler

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-backend/internal/service"
	"github.com/hustlers/motivator-backend/pkg/response"
)

type AnalyticsHandler struct {
	svc *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// GetDashboard godoc
// @Summary Get analytics dashboard
// @Description Returns comprehensive analytics for a company
// @Tags analytics
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} model.AnalyticsDashboard
// @Router /companies/{id}/analytics [get]
func (h *AnalyticsHandler) GetDashboard(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | fetching analytics dashboard company=%s", traceID, companyID)

	dashboard, err := h.svc.GetDashboard(c.Context(), companyID)
	if err != nil {
		log.Printf("trace=%s | error fetching analytics: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, dashboard)
}
