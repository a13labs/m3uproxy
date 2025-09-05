package providers

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserView struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Users struct {
	Users []User `json:"users"`
}
