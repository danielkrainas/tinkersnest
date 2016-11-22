package v1

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`

	Salt           []byte `json:"-"`
	HashedPassword string `json:"-"`
}
