package game_session

import (
	"time"
)

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

type GameSession struct {
	SessionID string
	Username  string
	Cash      float64
	Status    GameSessionStatus
	CreatedAt time.Time
}
