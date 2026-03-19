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

type RewardHandler struct {
	svc *service.RewardService
}

func NewRewardHandler(svc *service.RewardService) *RewardHandler {
	return &RewardHandler{svc: svc}
}

// Create godoc
// @Summary Create a reward
// @Description Create a reward in the store (admin+ only)
// @Tags rewards
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.CreateRewardRequest true "Reward data"
// @Success 201 {object} model.Reward
// @Router /companies/{id}/rewards [post]
func (h *RewardHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var req model.CreateRewardRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" || req.CostCoins <= 0 {
		return response.BadRequest(c, "name and cost_coins are required")
	}

	log.Printf("trace=%s | creating reward name=%s cost=%d company=%s", traceID, req.Name, req.CostCoins, companyID)

	reward, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		log.Printf("trace=%s | error creating reward: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Created(c, reward)
}

// List godoc
// @Summary List rewards
// @Description List all active rewards in the store
// @Tags rewards
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/rewards [get]
func (h *RewardHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | listing rewards company=%s", traceID, companyID)

	rewards, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		log.Printf("trace=%s | error listing rewards: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, rewards)
}

// Delete godoc
// @Summary Delete a reward
// @Description Deactivate a reward (admin+ only)
// @Tags rewards
// @Param id path string true "Company ID"
// @Param rewardId path string true "Reward ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/rewards/{rewardId} [delete]
func (h *RewardHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	rewardID := c.Params("rewardId")

	log.Printf("trace=%s | deleting reward id=%s", traceID, rewardID)

	if err := h.svc.Delete(c.Context(), rewardID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "reward not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

// Redeem godoc
// @Summary Redeem a reward
// @Description Spend coins to redeem a reward
// @Tags rewards
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.RedeemRequest true "Reward to redeem"
// @Success 200 {object} model.Redemption
// @Router /companies/{id}/rewards/redeem [post]
func (h *RewardHandler) Redeem(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	m := middleware.GetMembership(c)

	var req model.RedeemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.RewardID == "" {
		return response.BadRequest(c, "reward_id is required")
	}

	log.Printf("trace=%s | redeeming reward=%s by member=%s", traceID, req.RewardID, m.ID)

	rd, err := h.svc.Redeem(c.Context(), m.ID, req.RewardID)
	if err != nil {
		log.Printf("trace=%s | error redeeming: %v", traceID, err)
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, rd)
}

// ListRedemptions godoc
// @Summary List all redemptions
// @Description List all reward redemptions for the company (admin+ only)
// @Tags rewards
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/rewards/redemptions [get]
func (h *RewardHandler) ListRedemptions(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | listing redemptions company=%s", traceID, companyID)

	redemptions, err := h.svc.ListRedemptions(c.Context(), companyID)
	if err != nil {
		log.Printf("trace=%s | error listing redemptions: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, redemptions)
}

// FulfillRedemption godoc
// @Summary Fulfill a redemption
// @Description Mark a redemption as fulfilled (admin+ only)
// @Tags rewards
// @Param id path string true "Company ID"
// @Param redemptionId path string true "Redemption ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/rewards/redemptions/{redemptionId}/fulfill [post]
func (h *RewardHandler) FulfillRedemption(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	redemptionID := c.Params("redemptionId")

	log.Printf("trace=%s | fulfilling redemption=%s", traceID, redemptionID)

	if err := h.svc.UpdateRedemptionStatus(c.Context(), redemptionID, model.RedemptionFulfilled); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "redemption not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"fulfilled": true})
}
