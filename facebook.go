package user

import (
	"encoding/json"
	"io"
)

// Facebook описывает провайдера для получения информации о пользователе с
// сервера facebook.
var Facebook = &Provider{
	URL:    "https://graph.facebook.com/me?fields=id,email,name,picture",
	Prefix: "fb:",
	reader: func(r io.Reader) (*User, error) {
		var user = new(struct {
			ID      string `json:"id"`
			User           // стандартный пользователь
			Picture struct {
				Data struct {
					URL string `json:"url"`
				} `json:"data"`
			} `json:"picture"`
		})
		if err := json.NewDecoder(r).Decode(&user); err != nil {
			return nil, err
		}
		user.User.ID = user.ID
		user.User.Picture = user.Picture.Data.URL
		return &user.User, nil
	},
}
