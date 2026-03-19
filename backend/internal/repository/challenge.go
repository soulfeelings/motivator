package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type ChallengeRepository interface {
	Create(ctx context.Context, c *model.Challenge) error
	GetByID(ctx context.Context, id string) (*model.Challenge, error)
	ListByCompany(ctx context.Context, companyID string, status *model.ChallengeStatus) ([]model.Challenge, error)
	ListByMember(ctx context.Context, membershipID string) ([]model.Challenge, error)
	UpdateStatus(ctx context.Context, id string, status model.ChallengeStatus) error
	UpdateScore(ctx context.Context, id string, isChallenger bool, score int) error
	Complete(ctx context.Context, id string, winnerID string) error
}

type challengeRepo struct {
	pool *pgxpool.Pool
}

func NewChallengeRepository(pool *pgxpool.Pool) ChallengeRepository {
	return &challengeRepo{pool: pool}
}

const challengeCols = `id, company_id, challenger_id, opponent_id, metric, target, wager, status, challenger_score, opponent_score, winner_id, xp_reward, deadline, created_at, completed_at`

func scanChallenge(row pgx.Row, c *model.Challenge) error {
	return row.Scan(&c.ID, &c.CompanyID, &c.ChallengerID, &c.OpponentID, &c.Metric, &c.Target, &c.Wager, &c.Status, &c.ChallengerScore, &c.OpponentScore, &c.WinnerID, &c.XPReward, &c.Deadline, &c.CreatedAt, &c.CompletedAt)
}

func (r *challengeRepo) Create(ctx context.Context, c *model.Challenge) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO challenges (company_id, challenger_id, opponent_id, metric, target, wager, xp_reward)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, status, challenger_score, opponent_score, deadline, created_at`,
		c.CompanyID, c.ChallengerID, c.OpponentID, c.Metric, c.Target, c.Wager, c.XPReward,
	).Scan(&c.ID, &c.Status, &c.ChallengerScore, &c.OpponentScore, &c.Deadline, &c.CreatedAt)
}

func (r *challengeRepo) GetByID(ctx context.Context, id string) (*model.Challenge, error) {
	c := &model.Challenge{}
	err := scanChallenge(r.pool.QueryRow(ctx, `SELECT `+challengeCols+` FROM challenges WHERE id = $1`, id), c)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return c, err
}

func (r *challengeRepo) ListByCompany(ctx context.Context, companyID string, status *model.ChallengeStatus) ([]model.Challenge, error) {
	var rows pgx.Rows
	var err error
	if status != nil {
		rows, err = r.pool.Query(ctx, `SELECT `+challengeCols+` FROM challenges WHERE company_id = $1 AND status = $2 ORDER BY created_at DESC`, companyID, *status)
	} else {
		rows, err = r.pool.Query(ctx, `SELECT `+challengeCols+` FROM challenges WHERE company_id = $1 ORDER BY created_at DESC`, companyID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanChallenges(rows)
}

func (r *challengeRepo) ListByMember(ctx context.Context, membershipID string) ([]model.Challenge, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+challengeCols+` FROM challenges WHERE challenger_id = $1 OR opponent_id = $1 ORDER BY created_at DESC`, membershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanChallenges(rows)
}

func (r *challengeRepo) UpdateStatus(ctx context.Context, id string, status model.ChallengeStatus) error {
	tag, err := r.pool.Exec(ctx, `UPDATE challenges SET status = $2 WHERE id = $1`, id, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *challengeRepo) UpdateScore(ctx context.Context, id string, isChallenger bool, score int) error {
	col := "opponent_score"
	if isChallenger {
		col = "challenger_score"
	}
	_, err := r.pool.Exec(ctx, `UPDATE challenges SET `+col+` = $2 WHERE id = $1`, id, score)
	return err
}

func (r *challengeRepo) Complete(ctx context.Context, id string, winnerID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE challenges SET status = 'completed', winner_id = $2, completed_at = now() WHERE id = $1`,
		id, winnerID)
	return err
}

func scanChallenges(rows pgx.Rows) ([]model.Challenge, error) {
	var challenges []model.Challenge
	for rows.Next() {
		var c model.Challenge
		if err := scanChallenge(rows, &c); err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}
	return challenges, rows.Err()
}
