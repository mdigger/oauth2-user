# Библиотека для запроса информации об авторизованном пользователе OAuth2

На сегодняшний день в ней описаны настройки для получения минимальной информации 
о пользователе сервисов Google, Facebook и Yandex.

#### Пример использования

```golang
package main

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

var appName = "OAuth2Test"
// задаем конфигурацию для генерируемых нами токенов
var jwtConfig = jwt.Config{
	Issuer:  "http://example.com/",
	Created: true,
	Expires: time.Hour,
	Nonce:   jwt.Nonce(8),
	Key:     jwt.NewES256Key(),
}


func returnToken(w http.ResponseWriter, r *http.Request) {
	// запрашиваем токен из заголовка авторизации
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
	// выбираем провайдера для получения информации о пользователе
	var provider user.Provider
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
	token, err = jwtConfig.Token(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// отдаем токен в ответ на этот запрос
	w.Header().Set("Content-Type", "application/jwt")
	io.WriteString(w, token)
}

func main() {
	// запускаем сервер
	err := http.ListenAndServe("localhost:8001", 
		http.HandlerFunc(returnToken))
	if err != nil {
		log.WithError(err).Error("server error")
	}
}
```
