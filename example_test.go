package user_test

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	user "github.com/mdigger/oauth2-user"
)

func Example() {
	// запускаем сервер с обработчиком
	if err := http.ListenAndServe("localhost:8001", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// запрашиваем токен из заголовка авторизации
			token := r.Header.Get("Authorization")
			if !strings.HasPrefix(token, "Bearer ") {
				w.Header().Set("WWW-Authenticate",
					fmt.Sprintf("Bearer realm=%s", "App Name"))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// оставляем только токен
			token = strings.TrimPrefix(token, "Bearer ")
			// выбираем провайдера для получения информации о пользователе
			var provider user.Provider
			// получаем имя провайдера токена
			service := r.URL.Query().Get("provider")
			switch service {
			case "facebook.com":
				provider = user.Facebook
			case "google.com":
				provider = user.Google
			case "yandex.ru":
				provider = user.Yandex
			default:
				http.Error(w, fmt.Sprintf("unsupported provider %q", service),
					http.StatusUnauthorized)
				return
			}

			// запрашиваем информацию о пользователе
			user, err := provider.Get(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// формируем новый токен уже от нашего сервиса с описанием пользователя
			_ = user
			w.WriteHeader(http.StatusOK)
		})); err != nil {
		log.Fatal(err)
	}
}
