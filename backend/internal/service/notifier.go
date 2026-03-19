package service

import (
	"context"
	"fmt"

	"github.com/hustlers/motivator-backend/internal/model"
)

// Notifier is an interface so services can send pushes without importing the full NotificationService.
type Notifier interface {
	SendToMember(ctx context.Context, membershipID string, notif model.Notification)
	SendToMembers(ctx context.Context, membershipIDs []string, notif model.Notification)
}

// Notification helpers for common events.

func NotifyAchievementCompleted(ctx context.Context, n Notifier, membershipID, achievementName string, xp, coins int) {
	if n == nil {
		return
	}
	body := fmt.Sprintf("You earned +%d XP", xp)
	if coins > 0 {
		body += fmt.Sprintf(" and +%d coins", coins)
	}
	n.SendToMember(ctx, membershipID, model.Notification{
		Title: fmt.Sprintf("Achievement Unlocked: %s", achievementName),
		Body:  body,
		Data:  map[string]string{"type": "achievement"},
	})
}

func NotifyBadgeAwarded(ctx context.Context, n Notifier, membershipID, badgeName string) {
	if n == nil {
		return
	}
	n.SendToMember(ctx, membershipID, model.Notification{
		Title: "New Badge Earned!",
		Body:  fmt.Sprintf("You earned the \"%s\" badge", badgeName),
		Data:  map[string]string{"type": "badge"},
	})
}

func NotifyChallengeInvite(ctx context.Context, n Notifier, opponentID, metric string, target int) {
	if n == nil {
		return
	}
	n.SendToMember(ctx, opponentID, model.Notification{
		Title: "New Challenge!",
		Body:  fmt.Sprintf("You've been challenged: %s >= %d", metric, target),
		Data:  map[string]string{"type": "challenge_invite"},
	})
}

func NotifyChallengeCompleted(ctx context.Context, n Notifier, winnerID, loserID string, xpReward int) {
	if n == nil {
		return
	}
	n.SendToMember(ctx, winnerID, model.Notification{
		Title: "Challenge Won!",
		Body:  fmt.Sprintf("You won the challenge! +%d XP", xpReward),
		Data:  map[string]string{"type": "challenge_won"},
	})
	n.SendToMember(ctx, loserID, model.Notification{
		Title: "Challenge Complete",
		Body:  "The challenge has ended. Better luck next time!",
		Data:  map[string]string{"type": "challenge_lost"},
	})
}

func NotifyRewardRedeemed(ctx context.Context, n Notifier, membershipID, rewardName string) {
	if n == nil {
		return
	}
	n.SendToMember(ctx, membershipID, model.Notification{
		Title: "Reward Redeemed!",
		Body:  fmt.Sprintf("You redeemed \"%s\". An admin will fulfill it soon.", rewardName),
		Data:  map[string]string{"type": "reward_redeemed"},
	})
}
