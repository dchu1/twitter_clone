package handlermodels

type CreateUserRequest struct {
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email,omitempty"`
}
