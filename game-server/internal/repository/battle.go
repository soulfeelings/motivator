package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-game/internal/model"
)

type BattleRepository interface {
	Create(ctx context.Context, b *model.Battle) error
	GetByID(ctx context.Context, id string) (*model.Battle, error)
	ListByBase(ctx context.Context, baseID string) ([]model.Battle, error)
}

type battleRepo struct {
	pool *pgxpool.Pool
}

func NewBattleRepository(pool *pgxpool.Pool) BattleRepository {
	return &battleRepo{pool: pool}
}

func (r *battleRepo) Create(ctx context.Context, b *model.Battle) error {
	alJSON, _ := json.Marshal(b.AttackerLost)
	dlJSON, _ := json.Marshal(b.DefenderLost)
	replayJSON, _ := json.Marshal(b.ReplayData)

	return r.pool.QueryRow(ctx,
		`INSERT INTO battles (attacker_id, defender_id, winner_id, attacker_lost, defender_lost, replay_data, coins_won, xp_won)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, fought_at`,
		b.AttackerID, b.DefenderID, b.WinnerID, alJSON, dlJSON, replayJSON, b.CoinsWon, b.XPWon,
	).Scan(&b.ID, &b.FoughtAt)
}

func (r *battleRepo) GetByID(ctx context.Context, id string) (*model.Battle, error) {
	b := &model.Battle{}
	var alJSON, dlJSON, replayJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, attacker_id, defender_id, winner_id, attacker_lost, defender_lost, replay_data, coins_won, xp_won, fought_at
		 FROM battles WHERE id = $1`, id,
	).Scan(&b.ID, &b.AttackerID, &b.DefenderID, &b.WinnerID, &alJSON, &dlJSON, &replayJSON, &b.CoinsWon, &b.XPWon, &b.FoughtAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(alJSON, &b.AttackerLost)
	json.Unmarshal(dlJSON, &b.DefenderLost)
	json.Unmarshal(replayJSON, &b.ReplayData)
	return b, nil
}

func (r *battleRepo) ListByBase(ctx context.Context, baseID string) ([]model.Battle, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, attacker_id, defender_id, winner_id, attacker_lost, defender_lost, replay_data, coins_won, xp_won, fought_at
		 FROM battles WHERE attacker_id = $1 OR defender_id = $1 ORDER BY fought_at DESC LIMIT 20`, baseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []model.Battle
	for rows.Next() {
		var b model.Battle
		var alJSON, dlJSON, replayJSON []byte
		if err := rows.Scan(&b.ID, &b.AttackerID, &b.DefenderID, &b.WinnerID, &alJSON, &dlJSON, &replayJSON, &b.CoinsWon, &b.XPWon, &b.FoughtAt); err != nil {
			return nil, err
		}
		json.Unmarshal(alJSON, &b.AttackerLost)
		json.Unmarshal(dlJSON, &b.DefenderLost)
		json.Unmarshal(replayJSON, &b.ReplayData)
		battles = append(battles, b)
	}
	return battles, rows.Err()
}
