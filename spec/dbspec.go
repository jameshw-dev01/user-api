package spec

type DbInterface interface {
	Create(user User) error
	ReadAll() ([]User, error)
	Read(username string) (User, error)
	Update(user User) error
	Delete(user User) error
}
