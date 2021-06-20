package user

import (
	"encoding/json"
)

// Facebook описывает провайдера для получения информации о пользователе с
// сервера facebook.
var Facebook = Provider{
	URL:    "https://graph.facebook.com/me?fields=id,email,name,picture",
	Prefix: "fb:",
	Decoder: func(dec *json.Decoder) (*User, error) {
		// разбираем данные о пользователе
		var user struct {
			ID      string `json:"id"`
			User           // стандартный пользователь
			Picture struct {
				Data struct {
					URL string `json:"url"`
				} `json:"data"`
			} `json:"picture"`
		}
		if err := dec.Decode(&user); err != nil {
			return nil, err
		}

		// заполняем нестандартные свойства
		user.User.ID = user.ID
		user.User.Picture = user.Picture.Data.URL

		return &user.User, nil
	},
}
