package session

type Repository interface {
	Create(sId string, value string, expire uint64) error
	Delete(sId string) error
	GetLoginBySessionID(sId string) (string, error)
}
