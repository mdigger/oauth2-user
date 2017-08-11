package user

// User описывает информацию об авторизованном пользователе.
type User struct {
	ID      string `json:"sub"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Picture string `json:"picture,omitempty"`
}
