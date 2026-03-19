package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

const membershipCols = `id, user_id, company_id, role, display_name, job_title, xp, level, coins, is_active, joined_at, created_at, updated_at`

type MembershipRepository interface {
	Create(ctx context.Context, tx pgx.Tx, m *model.Membership) error
	GetByID(ctx context.Context, id string) (*model.Membership, error)
	GetByUserAndCompany(ctx context.Context, userID, companyID string) (*model.Membership, error)
	ListByCompany(ctx context.Context, companyID string, page model.PaginationRequest) ([]model.Membership, int, error)
	ListByUser(ctx context.Context, userID string) ([]model.Membership, error)
	Update(ctx context.Context, id string, req model.UpdateMembershipRequest) (*model.Membership, error)
	AwardXP(ctx context.Context, id string, amount int) (*model.Membership, error)
	AwardCoins(ctx context.Context, tx pgx.Tx, id string, amount int) error
	Deactivate(ctx context.Context, id string) error
}

type membershipRepo struct {
	pool *pgxpool.Pool
}

func NewMembershipRepository(pool *pgxpool.Pool) MembershipRepository {
	return &membershipRepo{pool: pool}
}

func scanMembership(row pgx.Row, m *model.Membership) error {
	return row.Scan(&m.ID, &m.UserID, &m.CompanyID, &m.Role, &m.DisplayName, &m.JobTitle, &m.XP, &m.Level, &m.Coins, &m.IsActive, &m.JoinedAt, &m.CreatedAt, &m.UpdatedAt)
}

func scanMembershipRows(rows pgx.Rows) ([]model.Membership, error) {
	var members []model.Membership
	for rows.Next() {
		var m model.Membership
		if err := rows.Scan(&m.ID, &m.UserID, &m.CompanyID, &m.Role, &m.DisplayName, &m.JobTitle, &m.XP, &m.Level, &m.Coins, &m.IsActive, &m.JoinedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *membershipRepo) Create(ctx context.Context, tx pgx.Tx, m *model.Membership) error {
	return tx.QueryRow(ctx,
		`INSERT INTO memberships (user_id, company_id, role, display_name, job_title)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, xp, level, coins, is_active, joined_at, created_at, updated_at`,
		m.UserID, m.CompanyID, m.Role, m.DisplayName, m.JobTitle,
	).Scan(&m.ID, &m.XP, &m.Level, &m.Coins, &m.IsActive, &m.JoinedAt, &m.CreatedAt, &m.UpdatedAt)
}

func (r *membershipRepo) GetByID(ctx context.Context, id string) (*model.Membership, error) {
	m := &model.Membership{}
	err := scanMembership(r.pool.QueryRow(ctx,
		`SELECT `+membershipCols+` FROM memberships WHERE id = $1`, id), m)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return m, err
}

func (r *membershipRepo) GetByUserAndCompany(ctx context.Context, userID, companyID string) (*model.Membership, error) {
	m := &model.Membership{}
	err := scanMembership(r.pool.QueryRow(ctx,
		`SELECT `+membershipCols+` FROM memberships WHERE user_id = $1 AND company_id = $2`, userID, companyID), m)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return m, err
}

func (r *membershipRepo) ListByCompany(ctx context.Context, companyID string, page model.PaginationRequest) ([]model.Membership, int, error) {
	page.Normalize()

	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM memberships WHERE company_id = $1 AND is_active = true`, companyID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+membershipCols+`
		 FROM memberships WHERE company_id = $1 AND is_active = true
		 ORDER BY xp DESC, joined_at ASC LIMIT $2 OFFSET $3`,
		companyID, page.PerPage, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	members, err := scanMembershipRows(rows)
	return members, total, err
}

func (r *membershipRepo) ListByUser(ctx context.Context, userID string) ([]model.Membership, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+membershipCols+`
		 FROM memberships WHERE user_id = $1 AND is_active = true ORDER BY joined_at ASC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMembershipRows(rows)
}

func (r *membershipRepo) Update(ctx context.Context, id string, req model.UpdateMembershipRequest) (*model.Membership, error) {
	m := &model.Membership{}
	err := scanMembership(r.pool.QueryRow(ctx,
		`UPDATE memberships SET
			role         = COALESCE($2, role),
			display_name = COALESCE($3, display_name),
			job_title    = COALESCE($4, job_title),
			updated_at   = now()
		 WHERE id = $1
		 RETURNING `+membershipCols,
		id, req.Role, req.DisplayName, req.JobTitle), m)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return m, err
}

func (r *membershipRepo) AwardXP(ctx context.Context, id string, amount int) (*model.Membership, error) {
	m := &model.Membership{}
	// Level up every 100 XP
	err := scanMembership(r.pool.QueryRow(ctx,
		`UPDATE memberships SET
			xp = xp + $2,
			level = 1 + (xp + $2) / 100,
			updated_at = now()
		 WHERE id = $1
		 RETURNING `+membershipCols, id, amount), m)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return m, err
}

func (r *membershipRepo) AwardCoins(ctx context.Context, tx pgx.Tx, id string, amount int) error {
	_, err := tx.Exec(ctx,
		`UPDATE memberships SET coins = coins + $2, updated_at = now() WHERE id = $1`,
		id, amount,
	)
	return err
}

func (r *membershipRepo) Deactivate(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE memberships SET is_active = false, updated_at = now() WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
