package user

// Google описывает провайдера для получения информации о пользователе с
// сервера Google.
var Google = Provider{
	URL:    "https://www.googleapis.com/oauth2/v3/userinfo",
	Prefix: "gl:",
}
