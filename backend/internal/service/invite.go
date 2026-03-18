package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type InviteService struct {
	pool    *pgxpool.Pool
	invites repository.InviteRepository
	members repository.MembershipRepository
}

func NewInviteService(pool *pgxpool.Pool, invites repository.InviteRepository, members repository.MembershipRepository) *InviteService {
	return &InviteService{pool: pool, invites: invites, members: members}
}

func (s *InviteService) Create(ctx context.Context, companyID, invitedBy string, req model.CreateInviteRequest) (*model.Invite, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	invite := &model.Invite{
		CompanyID: companyID,
		Email:     req.Email,
		Role:      req.Role,
		InvitedBy: invitedBy,
		Token:     token,
	}

	if err := s.invites.Create(ctx, invite); err != nil {
		return nil, err
	}
	return invite, nil
}

func (s *InviteService) ListByCompany(ctx context.Context, companyID string, page model.PaginationRequest) ([]model.Invite, int, error) {
	return s.invites.ListByCompany(ctx, companyID, page)
}

func (s *InviteService) Revoke(ctx context.Context, id string) error {
	return s.invites.Revoke(ctx, id)
}

func (s *InviteService) Accept(ctx context.Context, token, userID, email string) (*model.Membership, error) {
	invite, err := s.invites.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if invite.Status != model.InviteStatusPending {
		return nil, model.ErrNotFound
	}

	if time.Now().After(invite.ExpiresAt) {
		return nil, model.ErrInviteExpired
	}

	if invite.Email != email {
		return nil, model.ErrEmailMismatch
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback error: %v", err)
		}
	}()

	membership := &model.Membership{
		UserID:    userID,
		CompanyID: invite.CompanyID,
		Role:      invite.Role,
	}
	if err := s.members.Create(ctx, tx, membership); err != nil {
		return nil, err
	}

	if err := s.invites.MarkAccepted(ctx, tx, invite.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return membership, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
