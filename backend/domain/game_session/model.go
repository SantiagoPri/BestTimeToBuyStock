package game_session

type GameSessionStatus string

const (
	StatusStarting GameSessionStatus = "starting"
	StatusWeek1    GameSessionStatus = "week1"
	StatusWeek2    GameSessionStatus = "week2"
	StatusWeek3    GameSessionStatus = "week3"
	StatusWeek4    GameSessionStatus = "week4"
	StatusWeek5    GameSessionStatus = "week5"
	StatusFinished GameSessionStatus = "finished"
	StatusExpired  GameSessionStatus = "expired"
)

func (s GameSessionStatus) String() string {
	return string(s)
}

func (s GameSessionStatus) IsFinished() bool {
	return s == StatusFinished || s == StatusExpired
}

type HoldingInfo struct {
	Quantity   int     `json:"quantity"`
	TotalSpent float64 `json:"total_spent"`
}

type SessionMetadata struct {
	Holdings map[string]HoldingInfo `json:"holdings"`
}

type GameSession struct {
	ID        int64             `json:"id"`
	SessionID string            `json:"session_id"`
	Username  string            `json:"username"`
	Cash      float64           `json:"cash"`
	Status    GameSessionStatus `json:"status"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	Metadata  *SessionMetadata  `json:"metadata,omitempty"`
}
