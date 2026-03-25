package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type SocialGameService struct {
	games   repository.SocialGameRepository
	members repository.MembershipRepository
}

func NewSocialGameService(games repository.SocialGameRepository, members repository.MembershipRepository) *SocialGameService {
	return &SocialGameService{games: games, members: members}
}

func (s *SocialGameService) Create(ctx context.Context, companyID, createdBy string, req model.CreateSocialGameRequest) (*model.SocialGame, error) {
	startsAt := time.Now()
	if req.StartsAt != "" {
		if t, err := time.Parse(time.RFC3339, req.StartsAt); err == nil {
			startsAt = t
		}
	}
	endsAt := startsAt.Add(24 * time.Hour)
	if req.EndsAt != "" {
		if t, err := time.Parse(time.RFC3339, req.EndsAt); err == nil {
			endsAt = t
		}
	}

	cfg := req.Config
	if cfg == nil {
		cfg = json.RawMessage(`{}`)
	}

	g := &model.SocialGame{
		CompanyID:   companyID,
		GameType:    req.GameType,
		Name:        req.Name,
		Description: req.Description,
		Config:      cfg,
		XPReward:    req.XPReward,
		CoinReward:  req.CoinReward,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		CreatedBy:   createdBy,
	}
	if g.Name == "" {
		g.Name = "Social Game"
	}
	if g.XPReward == 0 {
		g.XPReward = 25
	}
	if g.CoinReward == 0 {
		g.CoinReward = 10
	}
	if err := s.games.CreateGame(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *SocialGameService) List(ctx context.Context, companyID string) ([]model.SocialGame, error) {
	return s.games.ListByCompany(ctx, companyID)
}

func (s *SocialGameService) GetByID(ctx context.Context, id string) (*model.SocialGame, error) {
	return s.games.GetByID(ctx, id)
}

func (s *SocialGameService) Delete(ctx context.Context, id string) error {
	return s.games.DeleteGame(ctx, id)
}

// AddQuestion adds a trivia question to a game.
func (s *SocialGameService) AddQuestion(ctx context.Context, gameID string, req model.AddQuestionRequest) (*model.SocialGameQuestion, error) {
	q := &model.SocialGameQuestion{
		GameID:       gameID,
		Question:     req.Question,
		Options:      req.Options,
		CorrectIndex: req.CorrectIndex,
		SortOrder:    req.SortOrder,
	}
	if err := s.games.AddQuestion(ctx, q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *SocialGameService) ListQuestions(ctx context.Context, gameID string) ([]model.SocialGameQuestion, error) {
	return s.games.ListQuestions(ctx, gameID)
}

// Launch transitions a game from draft to active.
func (s *SocialGameService) Launch(ctx context.Context, gameID string) error {
	g, err := s.games.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	if g.Status != model.GameDraft {
		return fmt.Errorf("game is not in draft status")
	}
	return s.games.UpdateStatus(ctx, gameID, model.GameActive)
}

// StartVoting transitions an active game to voting phase (for photo_challenge/two_truths).
func (s *SocialGameService) StartVoting(ctx context.Context, gameID string) error {
	g, err := s.games.GetByID(ctx, gameID)
	if err != nil {
		return err
	}
	if g.Status != model.GameActive {
		return fmt.Errorf("game is not in active status")
	}
	return s.games.UpdateStatus(ctx, gameID, model.GameVoting)
}

// Complete finishes the game, calculates winners, and awards XP/coins.
func (s *SocialGameService) Complete(ctx context.Context, gameID string) error {
	g, err := s.games.GetByID(ctx, gameID)
	if err != nil {
		return err
	}

	// Award participation XP/coins based on game type
	switch g.GameType {
	case model.GameTrivia:
		answers, err := s.games.GetAnswersByGame(ctx, gameID)
		if err != nil {
			return err
		}
		// Track unique participants
		awarded := make(map[string]bool)
		for _, a := range answers {
			if !awarded[a.MemberID] {
				s.members.AwardXP(ctx, a.MemberID, g.XPReward)
				awarded[a.MemberID] = true
			}
		}
	case model.GamePhotoChallenge, model.GameTwoTruths:
		submissions, err := s.games.ListSubmissions(ctx, gameID)
		if err != nil {
			return err
		}
		for _, sub := range submissions {
			s.members.AwardXP(ctx, sub.MemberID, g.XPReward)
		}
	}

	return s.games.UpdateStatus(ctx, gameID, model.GameCompleted)
}

// SubmitAnswer submits a trivia answer for a member.
func (s *SocialGameService) SubmitAnswer(ctx context.Context, gameID, memberID string, req model.SubmitAnswerRequest) (*model.SocialGameAnswer, error) {
	g, err := s.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if g.Status != model.GameActive {
		return nil, fmt.Errorf("game is not active")
	}

	// Look up the question to check correctness
	questions, err := s.games.ListQuestions(ctx, gameID)
	if err != nil {
		return nil, err
	}
	var correct bool
	var found bool
	for _, q := range questions {
		if q.ID == req.QuestionID {
			correct = q.CorrectIndex == req.SelectedIndex
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("question not found in this game")
	}

	a := &model.SocialGameAnswer{
		GameID:        gameID,
		QuestionID:    req.QuestionID,
		MemberID:      memberID,
		SelectedIndex: req.SelectedIndex,
		IsCorrect:     correct,
	}
	if err := s.games.SubmitAnswer(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

// SubmitEntry submits a photo or two-truths entry for a member.
func (s *SocialGameService) SubmitEntry(ctx context.Context, gameID, memberID string, req model.SubmitEntryRequest) (*model.SocialGameSubmission, error) {
	g, err := s.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if g.Status != model.GameActive {
		return nil, fmt.Errorf("game is not active")
	}

	sub := &model.SocialGameSubmission{
		GameID:     gameID,
		MemberID:   memberID,
		Content:    req.Content,
		Statements: req.Statements,
		LieIndex:   req.LieIndex,
	}
	if err := s.games.CreateSubmission(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

// CastVote casts a vote on a submission.
func (s *SocialGameService) CastVote(ctx context.Context, gameID, voterID string, req model.CastVoteRequest) (*model.SocialGameVote, error) {
	g, err := s.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if g.Status != model.GameVoting {
		return nil, fmt.Errorf("game is not in voting phase")
	}

	v := &model.SocialGameVote{
		GameID:       gameID,
		SubmissionID: req.SubmissionID,
		VoterID:      voterID,
		VoteValue:    req.VoteValue,
	}
	if err := s.games.CastVote(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

// ListSubmissions returns all submissions for a game.
func (s *SocialGameService) ListSubmissions(ctx context.Context, gameID string) ([]model.SocialGameSubmission, error) {
	return s.games.ListSubmissions(ctx, gameID)
}

// GetResults returns game results including participation and leaderboard.
func (s *SocialGameService) GetResults(ctx context.Context, gameID, companyID string) (*model.SocialGameResults, error) {
	g, err := s.games.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	participantCount, err := s.games.GetGameStats(ctx, gameID)
	if err != nil {
		return nil, err
	}

	// Get total member count for participation rate
	page := model.PaginationRequest{Page: 1, PerPage: 1}
	_, totalMembers, err := s.members.ListByCompany(ctx, companyID, page)
	if err != nil {
		return nil, err
	}

	var rate float64
	if totalMembers > 0 {
		rate = float64(participantCount) / float64(totalMembers) * 100
	}

	// Build leaderboard based on game type
	var leaderboard []model.LeaderboardEntry

	switch g.GameType {
	case model.GameTrivia:
		answers, err := s.games.GetAnswersByGame(ctx, gameID)
		if err != nil {
			return nil, err
		}
		scores := make(map[string]int)
		for _, a := range answers {
			if a.IsCorrect {
				scores[a.MemberID]++
			}
		}
		for memberID, score := range scores {
			leaderboard = append(leaderboard, model.LeaderboardEntry{MemberID: memberID, Score: score})
		}
	case model.GamePhotoChallenge, model.GameTwoTruths:
		submissions, err := s.games.ListSubmissions(ctx, gameID)
		if err != nil {
			return nil, err
		}
		for _, sub := range submissions {
			leaderboard = append(leaderboard, model.LeaderboardEntry{MemberID: sub.MemberID, Score: sub.VoteCount})
		}
	}

	// Sort leaderboard descending by score
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	var winner *model.LeaderboardEntry
	if len(leaderboard) > 0 {
		winner = &leaderboard[0]
	}

	return &model.SocialGameResults{
		Game:              *g,
		ParticipantCount:  participantCount,
		TotalMembers:      totalMembers,
		ParticipationRate: rate,
		Leaderboard:       leaderboard,
		Winner:            winner,
	}, nil
}
