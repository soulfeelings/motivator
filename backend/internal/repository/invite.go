package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type InviteRepository interface {
	Create(ctx context.Context, invite *model.Invite) error
	GetByToken(ctx context.Context, token string) (*model.Invite, error)
	ListByCompany(ctx context.Context, companyID string, page model.PaginationRequest) ([]model.Invite, int, error)
	UpdateStatus(ctx context.Context, tx pgx.Tx, id string, status model.InviteStatus) error
	MarkAccepted(ctx context.Context, tx pgx.Tx, id string) error
	Revoke(ctx context.Context, id string) error
}

type inviteRepo struct {
	pool *pgxpool.Pool
}

func NewInviteRepository(pool *pgxpool.Pool) InviteRepository {
	return &inviteRepo{pool: pool}
}

func (r *inviteRepo) Create(ctx context.Context, invite *model.Invite) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO invites (company_id, email, role, invited_by, token)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, status, expires_at, created_at`,
		invite.CompanyID, invite.Email, invite.Role, invite.InvitedBy, invite.Token,
	).Scan(&invite.ID, &invite.Status, &invite.ExpiresAt, &invite.CreatedAt)
}

func (r *inviteRepo) GetByToken(ctx context.Context, token string) (*model.Invite, error) {
	i := &model.Invite{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, email, role, status, invited_by, token, expires_at, accepted_at, created_at
		 FROM invites WHERE token = $1`, token,
	).Scan(&i.ID, &i.CompanyID, &i.Email, &i.Role, &i.Status, &i.InvitedBy, &i.Token, &i.ExpiresAt, &i.AcceptedAt, &i.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return i, err
}

func (r *inviteRepo) ListByCompany(ctx context.Context, companyID string, page model.PaginationRequest) ([]model.Invite, int, error) {
	page.Normalize()

	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM invites WHERE company_id = $1`, companyID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, email, role, status, invited_by, token, expires_at, accepted_at, created_at
		 FROM invites WHERE company_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		companyID, page.PerPage, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invites []model.Invite
	for rows.Next() {
		var i model.Invite
		if err := rows.Scan(&i.ID, &i.CompanyID, &i.Email, &i.Role, &i.Status, &i.InvitedBy, &i.Token, &i.ExpiresAt, &i.AcceptedAt, &i.CreatedAt); err != nil {
			return nil, 0, err
		}
		invites = append(invites, i)
	}
	return invites, total, rows.Err()
}

func (r *inviteRepo) UpdateStatus(ctx context.Context, tx pgx.Tx, id string, status model.InviteStatus) error {
	_, err := tx.Exec(ctx,
		`UPDATE invites SET status = $2 WHERE id = $1`, id, status,
	)
	return err
}

func (r *inviteRepo) MarkAccepted(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx,
		`UPDATE invites SET status = 'accepted', accepted_at = now() WHERE id = $1`, id,
	)
	return err
}

func (r *inviteRepo) Revoke(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE invites SET status = 'revoked' WHERE id = $1 AND status = 'pending'`, id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
