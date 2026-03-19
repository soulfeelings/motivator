package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type TournamentRepository interface {
	Create(ctx context.Context, t *model.Tournament) error
	GetByID(ctx context.Context, id string) (*model.Tournament, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.Tournament, error)
	UpdateStatus(ctx context.Context, id string, status model.TournamentStatus) error
	Delete(ctx context.Context, id string) error
	Join(ctx context.Context, p *model.TournamentParticipant) error
	Leave(ctx context.Context, tournamentID, membershipID string) error
	UpdateScore(ctx context.Context, tournamentID, membershipID string, score int) error
	GetStandings(ctx context.Context, tournamentID string) ([]model.TournamentParticipant, error)
	UpdateRanks(ctx context.Context, tournamentID string) error
	Complete(ctx context.Context, id string) error
}

type tournamentRepo struct {
	pool *pgxpool.Pool
}

func NewTournamentRepository(pool *pgxpool.Pool) TournamentRepository {
	return &tournamentRepo{pool: pool}
}

func (r *tournamentRepo) Create(ctx context.Context, t *model.Tournament) error {
	xpJSON, _ := json.Marshal(t.XPPrizes)
	coinJSON, _ := json.Marshal(t.CoinPrizes)
	return r.pool.QueryRow(ctx,
		`INSERT INTO tournaments (company_id, name, description, season, metric, prize_pool, xp_prizes, coin_prizes, max_participants, starts_at, ends_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id, status, created_at`,
		t.CompanyID, t.Name, t.Description, t.Season, t.Metric, t.PrizePool, xpJSON, coinJSON, t.MaxParticipants, t.StartsAt, t.EndsAt,
	).Scan(&t.ID, &t.Status, &t.CreatedAt)
}

func (r *tournamentRepo) GetByID(ctx context.Context, id string) (*model.Tournament, error) {
	t := &model.Tournament{}
	var xpJSON, coinJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT t.id, t.company_id, t.name, t.description, t.season, t.metric, t.prize_pool, t.xp_prizes, t.coin_prizes,
		        t.status, t.max_participants, t.starts_at, t.ends_at, t.created_at, t.completed_at,
		        (SELECT COUNT(*) FROM tournament_participants WHERE tournament_id = t.id)
		 FROM tournaments t WHERE t.id = $1`, id,
	).Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.Season, &t.Metric, &t.PrizePool, &xpJSON, &coinJSON,
		&t.Status, &t.MaxParticipants, &t.StartsAt, &t.EndsAt, &t.CreatedAt, &t.CompletedAt, &t.ParticipantCount)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal(xpJSON, &t.XPPrizes)
	json.Unmarshal(coinJSON, &t.CoinPrizes)
	return t, nil
}

func (r *tournamentRepo) ListByCompany(ctx context.Context, companyID string) ([]model.Tournament, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.company_id, t.name, t.description, t.season, t.metric, t.prize_pool, t.xp_prizes, t.coin_prizes,
		        t.status, t.max_participants, t.starts_at, t.ends_at, t.created_at, t.completed_at,
		        (SELECT COUNT(*) FROM tournament_participants WHERE tournament_id = t.id)
		 FROM tournaments t WHERE t.company_id = $1 ORDER BY t.created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []model.Tournament
	for rows.Next() {
		var t model.Tournament
		var xpJSON, coinJSON []byte
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.Name, &t.Description, &t.Season, &t.Metric, &t.PrizePool, &xpJSON, &coinJSON,
			&t.Status, &t.MaxParticipants, &t.StartsAt, &t.EndsAt, &t.CreatedAt, &t.CompletedAt, &t.ParticipantCount); err != nil {
			return nil, err
		}
		json.Unmarshal(xpJSON, &t.XPPrizes)
		json.Unmarshal(coinJSON, &t.CoinPrizes)
		tournaments = append(tournaments, t)
	}
	return tournaments, rows.Err()
}

func (r *tournamentRepo) UpdateStatus(ctx context.Context, id string, status model.TournamentStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE tournaments SET status = $2 WHERE id = $1`, id, status)
	return err
}

func (r *tournamentRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM tournaments WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *tournamentRepo) Join(ctx context.Context, p *model.TournamentParticipant) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO tournament_participants (tournament_id, membership_id) VALUES ($1, $2)
		 RETURNING id, score, joined_at`, p.TournamentID, p.MembershipID,
	).Scan(&p.ID, &p.Score, &p.JoinedAt)
}

func (r *tournamentRepo) Leave(ctx context.Context, tournamentID, membershipID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM tournament_participants WHERE tournament_id = $1 AND membership_id = $2`, tournamentID, membershipID)
	return err
}

func (r *tournamentRepo) UpdateScore(ctx context.Context, tournamentID, membershipID string, score int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tournament_participants SET score = $3 WHERE tournament_id = $1 AND membership_id = $2`,
		tournamentID, membershipID, score)
	return err
}

func (r *tournamentRepo) GetStandings(ctx context.Context, tournamentID string) ([]model.TournamentParticipant, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, tournament_id, membership_id, score, rank, joined_at
		 FROM tournament_participants WHERE tournament_id = $1 ORDER BY score DESC`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []model.TournamentParticipant
	for rows.Next() {
		var p model.TournamentParticipant
		if err := rows.Scan(&p.ID, &p.TournamentID, &p.MembershipID, &p.Score, &p.Rank, &p.JoinedAt); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}

func (r *tournamentRepo) UpdateRanks(ctx context.Context, tournamentID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tournament_participants SET rank = sub.rank
		 FROM (SELECT id, ROW_NUMBER() OVER (ORDER BY score DESC) as rank
		       FROM tournament_participants WHERE tournament_id = $1) sub
		 WHERE tournament_participants.id = sub.id`, tournamentID)
	return err
}

func (r *tournamentRepo) Complete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE tournaments SET status = 'completed', completed_at = now() WHERE id = $1`, id)
	return err
}
