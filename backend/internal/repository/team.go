package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type TeamRepository interface {
	Create(ctx context.Context, t *model.Team) error
	GetByID(ctx context.Context, id string) (*model.Team, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.Team, error)
	Delete(ctx context.Context, id string) error
	AddMember(ctx context.Context, tm *model.TeamMember) error
	RemoveMember(ctx context.Context, teamID, membershipID string) error
	ListMembers(ctx context.Context, teamID string) ([]model.TeamMember, error)
	GetMemberTeam(ctx context.Context, membershipID string) (*model.Team, error)
	CreateBattle(ctx context.Context, b *model.TeamBattle) error
	GetBattle(ctx context.Context, id string) (*model.TeamBattle, error)
	ListBattles(ctx context.Context, companyID string) ([]model.TeamBattle, error)
	UpdateBattleScore(ctx context.Context, id string, isTeamA bool, score int) error
	CompleteBattle(ctx context.Context, id string, winnerID string) error
}

type teamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) TeamRepository {
	return &teamRepo{pool: pool}
}

func (r *teamRepo) Create(ctx context.Context, t *model.Team) error {
	if t.Color == "" {
		t.Color = "#8b5cf6"
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO teams (company_id, name, description, color) VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		t.CompanyID, t.Name, t.Description, t.Color,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *teamRepo) GetByID(ctx context.Context, id string) (*model.Team, error) {
	t := &model.Team{}
	err := r.pool.QueryRow(ctx,
		`SELECT t.id, t.company_id, t.name, t.description, t.color, t.created_at,
		        (SELECT COUNT(*) FROM team_members WHERE team_id = t.id)
		 FROM teams t WHERE t.id = $1`, id,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.Color, &t.CreatedAt, &t.MemberCount)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return t, err
}

func (r *teamRepo) ListByCompany(ctx context.Context, companyID string) ([]model.Team, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.company_id, t.name, t.description, t.color, t.created_at,
		        (SELECT COUNT(*) FROM team_members WHERE team_id = t.id)
		 FROM teams t WHERE t.company_id = $1 ORDER BY t.created_at ASC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []model.Team
	for rows.Next() {
		var t model.Team
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.Color, &t.CreatedAt, &t.MemberCount); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *teamRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM teams WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *teamRepo) AddMember(ctx context.Context, tm *model.TeamMember) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO team_members (team_id, membership_id) VALUES ($1, $2)
		 RETURNING id, joined_at`, tm.TeamID, tm.MembershipID,
	).Scan(&tm.ID, &tm.JoinedAt)
}

func (r *teamRepo) RemoveMember(ctx context.Context, teamID, membershipID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM team_members WHERE team_id = $1 AND membership_id = $2`, teamID, membershipID)
	return err
}

func (r *teamRepo) ListMembers(ctx context.Context, teamID string) ([]model.TeamMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, team_id, membership_id, joined_at FROM team_members WHERE team_id = $1`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.TeamMember
	for rows.Next() {
		var tm model.TeamMember
		if err := rows.Scan(&tm.ID, &tm.TeamID, &tm.MembershipID, &tm.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, tm)
	}
	return members, rows.Err()
}

func (r *teamRepo) GetMemberTeam(ctx context.Context, membershipID string) (*model.Team, error) {
	t := &model.Team{}
	err := r.pool.QueryRow(ctx,
		`SELECT t.id, t.company_id, t.name, t.description, t.color, t.created_at, 0
		 FROM teams t JOIN team_members tm ON tm.team_id = t.id
		 WHERE tm.membership_id = $1 LIMIT 1`, membershipID,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.Color, &t.CreatedAt, &t.MemberCount)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return t, err
}

func (r *teamRepo) CreateBattle(ctx context.Context, b *model.TeamBattle) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO team_battles (company_id, team_a_id, team_b_id, metric, target, xp_reward, coin_reward)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, status, team_a_score, team_b_score, deadline, created_at`,
		b.CompanyID, b.TeamAID, b.TeamBID, b.Metric, b.Target, b.XPReward, b.CoinReward,
	).Scan(&b.ID, &b.Status, &b.TeamAScore, &b.TeamBScore, &b.Deadline, &b.CreatedAt)
}

func (r *teamRepo) GetBattle(ctx context.Context, id string) (*model.TeamBattle, error) {
	b := &model.TeamBattle{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, team_a_id, team_b_id, metric, target, team_a_score, team_b_score, status, winner_id, xp_reward, coin_reward, deadline, created_at, completed_at
		 FROM team_battles WHERE id = $1`, id,
	).Scan(&b.ID, &b.CompanyID, &b.TeamAID, &b.TeamBID, &b.Metric, &b.Target, &b.TeamAScore, &b.TeamBScore, &b.Status, &b.WinnerID, &b.XPReward, &b.CoinReward, &b.Deadline, &b.CreatedAt, &b.CompletedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return b, err
}

func (r *teamRepo) ListBattles(ctx context.Context, companyID string) ([]model.TeamBattle, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, team_a_id, team_b_id, metric, target, team_a_score, team_b_score, status, winner_id, xp_reward, coin_reward, deadline, created_at, completed_at
		 FROM team_battles WHERE company_id = $1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []model.TeamBattle
	for rows.Next() {
		var b model.TeamBattle
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.TeamAID, &b.TeamBID, &b.Metric, &b.Target, &b.TeamAScore, &b.TeamBScore, &b.Status, &b.WinnerID, &b.XPReward, &b.CoinReward, &b.Deadline, &b.CreatedAt, &b.CompletedAt); err != nil {
			return nil, err
		}
		battles = append(battles, b)
	}
	return battles, rows.Err()
}

func (r *teamRepo) UpdateBattleScore(ctx context.Context, id string, isTeamA bool, score int) error {
	col := "team_b_score"
	if isTeamA {
		col = "team_a_score"
	}
	_, err := r.pool.Exec(ctx, `UPDATE team_battles SET `+col+` = $2 WHERE id = $1`, id, score)
	return err
}

func (r *teamRepo) CompleteBattle(ctx context.Context, id string, winnerID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE team_battles SET status = 'completed', winner_id = $2, completed_at = now() WHERE id = $1`, id, winnerID)
	return err
}
