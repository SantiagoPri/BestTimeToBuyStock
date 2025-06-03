package gm_session

// Repository defines the interface for GM data storage operations
type Repository interface {
	SaveWeekData(sessionID string, week int, data *GMWeekData) error
	GetWeekData(sessionID string, week int) (*GMWeekData, error)
	ClearSessionData(sessionID string) error // optional, for dev cleanup or retries
}
