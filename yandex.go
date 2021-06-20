package user

import (
	"encoding/json"
)

// Yandex описывает провайдера для получения информации о пользователе с
// сервера Yandex.
var Yandex = Provider{
	URL:    "https://login.yandex.ru/info",
	Prefix: "ya:",
	Decoder: func(dec *json.Decoder) (*User, error) {
		var yaUser struct {
			ID           string `json:"id"`
			Name         string `json:"display_name"`
			Email        string `json:"default_email"`
			Picture      string `json:"default_avatar_id"`
			PictureEmpty bool   `json:"is_avatar_empty"`
		}
		if err := dec.Decode(&yaUser); err != nil {
			return nil, err
		}

		user := User{
			ID:    yaUser.ID,
			Name:  yaUser.Name,
			Email: yaUser.Email,
		}

		if !yaUser.PictureEmpty {
			user.Picture = "https://avatars.yandex.net/get-yapic/" +
				yaUser.Picture + "/islands-200"
		}

		return &user, nil
	},
}
