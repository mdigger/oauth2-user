package user_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mdigger/jwt"
	"github.com/mdigger/log"
	user "github.com/mdigger/oauth2-user"
)

func Example() {
	var appName = "OAuth2Test"
	// задаем конфигурацию для генерируемых нами токенов
	var jwtConfig = &jwt.Config{
		Issuer:  "http://example.com/",
		Created: true,
		Expires: time.Hour,
		Nonce:   jwt.Nonce(8),
		Key:     jwt.NewES256Key(),
	}

	returnToken := func(w http.ResponseWriter, r *http.Request) {
		// запрашиваем авторизационный токен из заголовка авторизации
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			w.Header().Set("WWW-Authenticate",
				fmt.Sprintf("Bearer realm=%s", appName))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// оставляем только токен
		token = strings.TrimPrefix(token, "Bearer ")
		// получаем имя провайдера токена
		var service = r.URL.Query().Get("provider")
		// выбираем провайдера для получения информации о пользователе
		var provider *user.Provider
		switch service {
		case "facebook.com":
			provider = user.Facebook
		case "google.com":
			provider = user.Google
		case "yandex.ru":
			provider = user.Yandex
		}
		if provider == nil {
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
		token, err = jwtConfig.Token(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// отдаем токен в ответ на этот запрос
		w.Header().Set("Content-Type", "application/jwt")
		io.WriteString(w, token)
	}

	// запускаем сервер
	err := http.ListenAndServe("localhost:8001", http.HandlerFunc(returnToken))
	if err != nil {
		log.WithError(err).Error("server error")
	}
}
