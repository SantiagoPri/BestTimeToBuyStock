package game_session

type Repository interface {
	Save(*GameSession) error
	FindBySessionID(string) (*GameSession, error)
	Update(*GameSession) error
	FindLeaderboardTop10(page, pageSize int) ([]GameSession, error)
}

type Pagination struct {
	Page     int
	PageSize int
	Total    int64
}
