package model

type AnalyticsOverview struct {
	TotalMembers       int `json:"total_members"`
	ActiveMembers      int `json:"active_members"`
	TotalXPAwarded     int `json:"total_xp_awarded"`
	TotalCoinsAwarded  int `json:"total_coins_awarded"`
	TotalCoinsSpent    int `json:"total_coins_spent"`
	TotalBadgesAwarded int `json:"total_badges_awarded"`
	TotalAchievements  int `json:"total_achievements_completed"`
	TotalChallenges    int `json:"total_challenges"`
	TotalRedemptions   int `json:"total_redemptions"`
}

type TopPerformer struct {
	MembershipID string  `json:"membership_id"`
	DisplayName  *string `json:"display_name,omitempty"`
	XP           int     `json:"xp"`
	Level        int     `json:"level"`
	Badges       int     `json:"badges"`
	Achievements int     `json:"achievements"`
}

type AchievementStat struct {
	AchievementID string `json:"achievement_id"`
	Name          string `json:"name"`
	Metric        string `json:"metric"`
	Completions   int    `json:"completions"`
}

type ChallengeStat struct {
	TotalChallenges int `json:"total_challenges"`
	Completed       int `json:"completed"`
	Active          int `json:"active"`
	Pending         int `json:"pending"`
	AvgXPReward     int `json:"avg_xp_reward"`
}

type RewardStat struct {
	RewardID       string `json:"reward_id"`
	Name           string `json:"name"`
	CostCoins      int    `json:"cost_coins"`
	TotalRedeemed  int    `json:"total_redeemed"`
	TotalCoinsSpent int   `json:"total_coins_spent"`
}

type XPDistribution struct {
	Level int `json:"level"`
	Count int `json:"count"`
}

type AnalyticsDashboard struct {
	Overview       AnalyticsOverview  `json:"overview"`
	TopPerformers  []TopPerformer     `json:"top_performers"`
	AchievementStats []AchievementStat `json:"achievement_stats"`
	ChallengeStats ChallengeStat      `json:"challenge_stats"`
	RewardStats    []RewardStat       `json:"reward_stats"`
	XPDistribution []XPDistribution   `json:"xp_distribution"`
}
