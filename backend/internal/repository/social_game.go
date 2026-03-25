package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type SocialGameRepository interface {
	CreateGame(ctx context.Context, g *model.SocialGame) error
	ListByCompany(ctx context.Context, companyID string) ([]model.SocialGame, error)
	GetByID(ctx context.Context, id string) (*model.SocialGame, error)
	DeleteGame(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status model.GameStatus) error

	AddQuestion(ctx context.Context, q *model.SocialGameQuestion) error
	ListQuestions(ctx context.Context, gameID string) ([]model.SocialGameQuestion, error)
	DeleteQuestion(ctx context.Context, id string) error

	SubmitAnswer(ctx context.Context, a *model.SocialGameAnswer) error
	GetAnswersByGame(ctx context.Context, gameID string) ([]model.SocialGameAnswer, error)

	CreateSubmission(ctx context.Context, s *model.SocialGameSubmission) error
	ListSubmissions(ctx context.Context, gameID string) ([]model.SocialGameSubmission, error)

	CastVote(ctx context.Context, v *model.SocialGameVote) error
	GetVotesBySubmission(ctx context.Context, submissionID string) ([]model.SocialGameVote, error)
	GetVotesByGame(ctx context.Context, gameID string) ([]model.SocialGameVote, error)

	GetGameStats(ctx context.Context, gameID string) (participantCount int, err error)
}

type socialGameRepo struct {
	pool *pgxpool.Pool
}

func NewSocialGameRepository(pool *pgxpool.Pool) SocialGameRepository {
	return &socialGameRepo{pool: pool}
}

func (r *socialGameRepo) CreateGame(ctx context.Context, g *model.SocialGame) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO social_games (company_id, game_type, name, description, config, xp_reward, coin_reward, starts_at, ends_at, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, status, created_at`,
		g.CompanyID, g.GameType, g.Name, g.Description, g.Config, g.XPReward, g.CoinReward, g.StartsAt, g.EndsAt, g.CreatedBy,
	).Scan(&g.ID, &g.Status, &g.CreatedAt)
}

func (r *socialGameRepo) ListByCompany(ctx context.Context, companyID string) ([]model.SocialGame, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, game_type, name, description, status, config, xp_reward, coin_reward,
		        starts_at, ends_at, created_by, created_at, completed_at
		 FROM social_games WHERE company_id = $1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []model.SocialGame
	for rows.Next() {
		var g model.SocialGame
		if err := rows.Scan(&g.ID, &g.CompanyID, &g.GameType, &g.Name, &g.Description, &g.Status, &g.Config,
			&g.XPReward, &g.CoinReward, &g.StartsAt, &g.EndsAt, &g.CreatedBy, &g.CreatedAt, &g.CompletedAt); err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, rows.Err()
}

func (r *socialGameRepo) GetByID(ctx context.Context, id string) (*model.SocialGame, error) {
	g := &model.SocialGame{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, game_type, name, description, status, config, xp_reward, coin_reward,
		        starts_at, ends_at, created_by, created_at, completed_at
		 FROM social_games WHERE id = $1`, id,
	).Scan(&g.ID, &g.CompanyID, &g.GameType, &g.Name, &g.Description, &g.Status, &g.Config,
		&g.XPReward, &g.CoinReward, &g.StartsAt, &g.EndsAt, &g.CreatedBy, &g.CreatedAt, &g.CompletedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return g, err
}

func (r *socialGameRepo) DeleteGame(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM social_games WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *socialGameRepo) UpdateStatus(ctx context.Context, id string, status model.GameStatus) error {
	q := `UPDATE social_games SET status = $2 WHERE id = $1`
	if status == model.GameCompleted {
		q = `UPDATE social_games SET status = $2, completed_at = now() WHERE id = $1`
	}
	_, err := r.pool.Exec(ctx, q, id, status)
	return err
}

// Questions

func (r *socialGameRepo) AddQuestion(ctx context.Context, q *model.SocialGameQuestion) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO social_game_questions (game_id, question, options, correct_index, sort_order)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		q.GameID, q.Question, q.Options, q.CorrectIndex, q.SortOrder,
	).Scan(&q.ID)
}

func (r *socialGameRepo) ListQuestions(ctx context.Context, gameID string) ([]model.SocialGameQuestion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, game_id, question, options, correct_index, sort_order
		 FROM social_game_questions WHERE game_id = $1 ORDER BY sort_order ASC`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []model.SocialGameQuestion
	for rows.Next() {
		var q model.SocialGameQuestion
		if err := rows.Scan(&q.ID, &q.GameID, &q.Question, &q.Options, &q.CorrectIndex, &q.SortOrder); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

func (r *socialGameRepo) DeleteQuestion(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM social_game_questions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

// Answers

func (r *socialGameRepo) SubmitAnswer(ctx context.Context, a *model.SocialGameAnswer) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO social_game_answers (game_id, question_id, member_id, selected_index, is_correct)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, answered_at`,
		a.GameID, a.QuestionID, a.MemberID, a.SelectedIndex, a.IsCorrect,
	).Scan(&a.ID, &a.AnsweredAt)
}

func (r *socialGameRepo) GetAnswersByGame(ctx context.Context, gameID string) ([]model.SocialGameAnswer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, game_id, question_id, member_id, selected_index, is_correct, answered_at
		 FROM social_game_answers WHERE game_id = $1 ORDER BY answered_at ASC`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []model.SocialGameAnswer
	for rows.Next() {
		var a model.SocialGameAnswer
		if err := rows.Scan(&a.ID, &a.GameID, &a.QuestionID, &a.MemberID, &a.SelectedIndex, &a.IsCorrect, &a.AnsweredAt); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, rows.Err()
}

// Submissions

func (r *socialGameRepo) CreateSubmission(ctx context.Context, s *model.SocialGameSubmission) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO social_game_submissions (game_id, member_id, content, statements, lie_index)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, submitted_at`,
		s.GameID, s.MemberID, s.Content, s.Statements, s.LieIndex,
	).Scan(&s.ID, &s.SubmittedAt)
}

func (r *socialGameRepo) ListSubmissions(ctx context.Context, gameID string) ([]model.SocialGameSubmission, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT s.id, s.game_id, s.member_id, s.content, s.statements, s.lie_index, s.submitted_at,
		        (SELECT COUNT(*) FROM social_game_votes WHERE submission_id = s.id)
		 FROM social_game_submissions s WHERE s.game_id = $1 ORDER BY s.submitted_at ASC`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []model.SocialGameSubmission
	for rows.Next() {
		var s model.SocialGameSubmission
		if err := rows.Scan(&s.ID, &s.GameID, &s.MemberID, &s.Content, &s.Statements, &s.LieIndex, &s.SubmittedAt, &s.VoteCount); err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	return submissions, rows.Err()
}

// Votes

func (r *socialGameRepo) CastVote(ctx context.Context, v *model.SocialGameVote) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO social_game_votes (game_id, submission_id, voter_id, vote_value)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, voted_at`,
		v.GameID, v.SubmissionID, v.VoterID, v.VoteValue,
	).Scan(&v.ID, &v.VotedAt)
}

func (r *socialGameRepo) GetVotesBySubmission(ctx context.Context, submissionID string) ([]model.SocialGameVote, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, game_id, submission_id, voter_id, vote_value, voted_at
		 FROM social_game_votes WHERE submission_id = $1 ORDER BY voted_at ASC`, submissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []model.SocialGameVote
	for rows.Next() {
		var v model.SocialGameVote
		if err := rows.Scan(&v.ID, &v.GameID, &v.SubmissionID, &v.VoterID, &v.VoteValue, &v.VotedAt); err != nil {
			return nil, err
		}
		votes = append(votes, v)
	}
	return votes, rows.Err()
}

func (r *socialGameRepo) GetVotesByGame(ctx context.Context, gameID string) ([]model.SocialGameVote, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, game_id, submission_id, voter_id, vote_value, voted_at
		 FROM social_game_votes WHERE game_id = $1 ORDER BY voted_at ASC`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []model.SocialGameVote
	for rows.Next() {
		var v model.SocialGameVote
		if err := rows.Scan(&v.ID, &v.GameID, &v.SubmissionID, &v.VoterID, &v.VoteValue, &v.VotedAt); err != nil {
			return nil, err
		}
		votes = append(votes, v)
	}
	return votes, rows.Err()
}

// Stats

func (r *socialGameRepo) GetGameStats(ctx context.Context, gameID string) (int, error) {
	// Count distinct participants (answers for trivia, submissions for others)
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(DISTINCT member_id) FROM (
			SELECT member_id FROM social_game_answers WHERE game_id = $1
			UNION
			SELECT member_id FROM social_game_submissions WHERE game_id = $1
		) participants`, gameID,
	).Scan(&count)
	return count, err
}
