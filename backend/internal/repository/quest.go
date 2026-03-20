package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type QuestRepository interface {
	Create(ctx context.Context, q *model.Quest) error
	GetByID(ctx context.Context, id string) (*model.Quest, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.Quest, error)
	UpdateStatus(ctx context.Context, id string, status model.QuestStatus) error
	Complete(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	CreatePair(ctx context.Context, p *model.QuestPair) error
	GetPairBySender(ctx context.Context, questID, senderID string) (*model.QuestPair, error)
	ListPairs(ctx context.Context, questID string) ([]model.QuestPair, error)
	SendMessage(ctx context.Context, pairID, message string) error
	GetReceivedMessages(ctx context.Context, questID, receiverID string, revealed bool) ([]model.ReceivedMessage, error)
	Vote(ctx context.Context, v *model.QuestVote) error
	GetVoteCounts(ctx context.Context, questID string) (map[string]int, error)
	GetWinnerPair(ctx context.Context, questID string) (*model.QuestPair, error)
}

type questRepo struct {
	pool *pgxpool.Pool
}

func NewQuestRepository(pool *pgxpool.Pool) QuestRepository {
	return &questRepo{pool: pool}
}

func (r *questRepo) Create(ctx context.Context, q *model.Quest) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO quests (company_id, name, description, xp_reward, coin_reward, bonus_xp, bonus_coins, deadline)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, status, created_at`,
		q.CompanyID, q.Name, q.Description, q.XPReward, q.CoinReward, q.BonusXP, q.BonusCoins, q.Deadline,
	).Scan(&q.ID, &q.Status, &q.CreatedAt)
}

func (r *questRepo) GetByID(ctx context.Context, id string) (*model.Quest, error) {
	q := &model.Quest{}
	err := r.pool.QueryRow(ctx,
		`SELECT q.id, q.company_id, q.name, q.description, q.status, q.xp_reward, q.coin_reward, q.bonus_xp, q.bonus_coins,
		        q.deadline, q.reveal_at, q.created_at, q.completed_at,
		        (SELECT COUNT(*) FROM quest_pairs WHERE quest_id = q.id),
		        (SELECT COUNT(*) FROM quest_pairs WHERE quest_id = q.id AND sent_at IS NOT NULL)
		 FROM quests q WHERE q.id = $1`, id,
	).Scan(&q.ID, &q.CompanyID, &q.Name, &q.Description, &q.Status, &q.XPReward, &q.CoinReward, &q.BonusXP, &q.BonusCoins,
		&q.Deadline, &q.RevealAt, &q.CreatedAt, &q.CompletedAt, &q.PairCount, &q.SentCount)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return q, err
}

func (r *questRepo) ListByCompany(ctx context.Context, companyID string) ([]model.Quest, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT q.id, q.company_id, q.name, q.description, q.status, q.xp_reward, q.coin_reward, q.bonus_xp, q.bonus_coins,
		        q.deadline, q.reveal_at, q.created_at, q.completed_at,
		        (SELECT COUNT(*) FROM quest_pairs WHERE quest_id = q.id),
		        (SELECT COUNT(*) FROM quest_pairs WHERE quest_id = q.id AND sent_at IS NOT NULL)
		 FROM quests q WHERE q.company_id = $1 ORDER BY q.created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quests []model.Quest
	for rows.Next() {
		var q model.Quest
		if err := rows.Scan(&q.ID, &q.CompanyID, &q.Name, &q.Description, &q.Status, &q.XPReward, &q.CoinReward, &q.BonusXP, &q.BonusCoins,
			&q.Deadline, &q.RevealAt, &q.CreatedAt, &q.CompletedAt, &q.PairCount, &q.SentCount); err != nil {
			return nil, err
		}
		quests = append(quests, q)
	}
	return quests, rows.Err()
}

func (r *questRepo) UpdateStatus(ctx context.Context, id string, status model.QuestStatus) error {
	q := `UPDATE quests SET status = $2 WHERE id = $1`
	if status == model.QuestRevealed {
		q = `UPDATE quests SET status = $2, reveal_at = now() WHERE id = $1`
	}
	_, err := r.pool.Exec(ctx, q, id, status)
	return err
}

func (r *questRepo) Complete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE quests SET status = 'completed', completed_at = now() WHERE id = $1`, id)
	return err
}

func (r *questRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM quests WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *questRepo) CreatePair(ctx context.Context, p *model.QuestPair) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO quest_pairs (quest_id, sender_id, receiver_id) VALUES ($1, $2, $3) RETURNING id`,
		p.QuestID, p.SenderID, p.ReceiverID,
	).Scan(&p.ID)
}

func (r *questRepo) GetPairBySender(ctx context.Context, questID, senderID string) (*model.QuestPair, error) {
	p := &model.QuestPair{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, quest_id, sender_id, receiver_id, message, sent_at FROM quest_pairs WHERE quest_id = $1 AND sender_id = $2`,
		questID, senderID,
	).Scan(&p.ID, &p.QuestID, &p.SenderID, &p.ReceiverID, &p.Message, &p.SentAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return p, err
}

func (r *questRepo) ListPairs(ctx context.Context, questID string) ([]model.QuestPair, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.id, p.quest_id, p.sender_id, p.receiver_id, p.message, p.sent_at,
		        (SELECT COUNT(*) FROM quest_votes WHERE pair_id = p.id)
		 FROM quest_pairs p WHERE p.quest_id = $1 ORDER BY p.sent_at ASC NULLS LAST`, questID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pairs []model.QuestPair
	for rows.Next() {
		var p model.QuestPair
		if err := rows.Scan(&p.ID, &p.QuestID, &p.SenderID, &p.ReceiverID, &p.Message, &p.SentAt, &p.VoteCount); err != nil {
			return nil, err
		}
		pairs = append(pairs, p)
	}
	return pairs, rows.Err()
}

func (r *questRepo) SendMessage(ctx context.Context, pairID, message string) error {
	_, err := r.pool.Exec(ctx, `UPDATE quest_pairs SET message = $2, sent_at = now() WHERE id = $1`, pairID, message)
	return err
}

func (r *questRepo) GetReceivedMessages(ctx context.Context, questID, receiverID string, revealed bool) ([]model.ReceivedMessage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.id, p.message, p.sent_at, p.sender_id FROM quest_pairs p
		 WHERE p.quest_id = $1 AND p.receiver_id = $2 AND p.sent_at IS NOT NULL`, questID, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.ReceivedMessage
	for rows.Next() {
		var m model.ReceivedMessage
		var senderID string
		if err := rows.Scan(&m.PairID, &m.Message, &m.SentAt, &senderID); err != nil {
			return nil, err
		}
		if revealed {
			m.SenderID = &senderID
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *questRepo) Vote(ctx context.Context, v *model.QuestVote) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO quest_votes (quest_id, pair_id, voter_id) VALUES ($1, $2, $3) RETURNING id`,
		v.QuestID, v.PairID, v.VoterID,
	).Scan(&v.ID)
}

func (r *questRepo) GetVoteCounts(ctx context.Context, questID string) (map[string]int, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pair_id, COUNT(*) FROM quest_votes WHERE quest_id = $1 GROUP BY pair_id`, questID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var pairID string
		var count int
		if err := rows.Scan(&pairID, &count); err != nil {
			return nil, err
		}
		counts[pairID] = count
	}
	return counts, rows.Err()
}

func (r *questRepo) GetWinnerPair(ctx context.Context, questID string) (*model.QuestPair, error) {
	p := &model.QuestPair{}
	err := r.pool.QueryRow(ctx,
		`SELECT p.id, p.quest_id, p.sender_id, p.receiver_id, p.message, p.sent_at,
		        (SELECT COUNT(*) FROM quest_votes WHERE pair_id = p.id) as votes
		 FROM quest_pairs p WHERE p.quest_id = $1 AND p.sent_at IS NOT NULL
		 ORDER BY votes DESC LIMIT 1`, questID,
	).Scan(&p.ID, &p.QuestID, &p.SenderID, &p.ReceiverID, &p.Message, &p.SentAt, &p.VoteCount)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return p, err
}
