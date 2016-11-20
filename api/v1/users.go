package v1

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`

	Salt           int    `json:"-"`
	HashedPassword string `json:"-"`
}
