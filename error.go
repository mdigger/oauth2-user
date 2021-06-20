package user

// Error описывает ошибку, возвращаемую сервером провайдера при запросе
// информации о пользователе.
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// Error возвращает строку с описанием ошибки.
func (e Error) Error() string {
	return e.Message
}
