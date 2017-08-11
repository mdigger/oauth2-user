package user

import (
	"encoding/json"
	"io"
)

// Google описывает провайдера для получения информации о пользователе с
// сервера Google.
var Google = &Provider{
	URL:    "https://www.googleapis.com/oauth2/v3/userinfo",
	Prefix: "gl:",
	reader: func(r io.Reader) (*User, error) {
		var user = new(User)
		if err := json.NewDecoder(r).Decode(&user); err != nil {
			return nil, err
		}
		return user, nil
	},
}
