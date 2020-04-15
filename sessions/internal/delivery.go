package session

type Delivery interface {
	Create(login string) (string, error)
	Delete(sessionID string) error
	Check(sessionID string) (string, error)
}
