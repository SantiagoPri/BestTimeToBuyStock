package game_session

// GameSessionTx represents a transaction for game session operations
type GameSessionTx interface {
	Commit() error
	Rollback() error
	GetSession() *GameSession
	Update(*GameSession) error
}

// Repository represents the game session repository interface
type Repository interface {
	Save(*GameSession) error
	FindBySessionID(string) (*GameSession, error)
	FindLeaderboardTop10(page, pageSize int) ([]GameSession, error)
	BeginTransaction(sessionID string) (GameSessionTx, error)
}

type Pagination struct {
	Page     int
	PageSize int
	Total    int64
}
