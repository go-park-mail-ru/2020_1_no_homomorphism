package csrf

type Repository interface {
	Add(token string, expire int64) error
	Check(token string) error
}
